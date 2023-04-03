package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	openai "github.com/sashabaranov/go-openai"
)

type ApiController struct {
	beego.Controller
}

var token = beego.AppConfig.String("token")
var proxy = beego.AppConfig.String("proxy")

// ai config
var config openai.ClientConfig

// 用户信息
type UserBox struct {
	InQuestion bool
	Client     *openai.Client
	Time       time.Time
}

var userBox = make(map[string]*UserBox)

func init() {
	beego.Debug("token", token)
	beego.Debug("proxy", proxy)

	//定时器
	go timeOut()

	config = openai.DefaultConfig(token)
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}
}

type ChatData struct {
	Id  string `json:"id"`
	Msg string `json:"msg"`
}

type Response struct {
	Code int      `json:"code"`
	Data ChatData `json:"data"`
}

// @Title Chat
// @Description Chat by openai
// @Param	body		body 	ChatData	true		"body for ChatData content"
// @Success 200 {int} Response
// @Failure 403 body is empty
// @router / [post]
func (api *ApiController) Chat() {
	var data ChatData
	err := json.Unmarshal(api.Ctx.Input.RequestBody, &data)

	if err != nil {
		api.Data["json"] = &Response{
			Code: -1,
			Data: ChatData{
				Id:  data.Id,
				Msg: err.Error(),
			},
		}
		api.ServeJSON()
		return
	}

	var client *openai.Client

	user, ok := userBox[data.Id]
	if ok {
		//用户存在对话
		if user.InQuestion {
			//用户正在处于对话中，不应该让他开始新的对话
			api.Data["json"] = &Response{
				Code: -1,
				Data: ChatData{
					Id:  data.Id,
					Msg: "等下啦，正在思考之前的问题呢~",
				},
			}
			api.ServeJSON()
			return
		}
		user.InQuestion = true
		client = user.Client
	} else {
		//用户不存在对话
		user = new(UserBox)
		user.Client = openai.NewClientWithConfig(config)
		user.InQuestion = true
		userBox[data.Id] = user
		client = user.Client
	}

	user.Time = time.Now()

	beego.Debug(data.Id, data.Msg)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    data.Id,
					Content: data.Msg,
				},
			},
		},
	)

	if err != nil {
		userBox[data.Id].InQuestion = false
		beego.Error("ChatCompletion error: %v\n", err)
		api.Data["json"] = &Response{
			Code: -1,
			Data: ChatData{
				Id:  data.Id,
				Msg: err.Error(),
			},
		}
		api.ServeJSON()
		return
	}

	msgText := ""
	px := ""
	for i, ccc := range resp.Choices {
		if len(ccc.Message.Content) > 0 {
			if i != 0 {
				px += "\n"
			}
			msgText += ccc.Message.Content + px
		}
	}

	api.Data["json"] = &Response{
		Code: 0,
		Data: ChatData{
			Id:  data.Id,
			Msg: msgText,
		},
	}

	userBox[data.Id].InQuestion = false
	api.ServeJSON()
}

// 定时器，超时
func timeOut() {
	timeTicker := time.NewTicker(60 * time.Second)
	for {
		<-timeTicker.C
		for id, user := range userBox {
			if user.InQuestion && user.Time.Add(5*time.Second).Before(time.Now()) {
				beego.Debug("清除超时的用户", id)
				user.InQuestion = false
			}
		}
		continue
	}
}

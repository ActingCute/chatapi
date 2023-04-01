/*
 * @Author: ActingCute酱 rem486@qq.com
 * @Date: 2023-03-29 20:40:08
 * @LastEditors: ActingCute酱 rem486@qq.com
 * @LastEditTime: 2023-03-29 21:21:13
 * @FilePath: \chatgpt\routers\router.go
 * @Description: 说明
 */
// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"chatgpt/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/api",
			beego.NSInclude(
				&controllers.ApiController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

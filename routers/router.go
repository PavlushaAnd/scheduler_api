// @APIVersion 1.0.0
// @Title Scheduler Test API
// @Description API fot scheduler application
// @Contact pavelkiwiandrosov.00@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /v1

package routers

import (
	"scheduler_api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/task", &controllers.TaskController{}, "get:GetAll"),
		beego.NSRouter("/user", &controllers.UserController{}, "get:GetAll"),
		beego.NSRouter("/task", &controllers.TaskController{}, "post:Post"),
		beego.NSRouter("/task/:task_code", &controllers.TaskController{}, "get:Get"),
		beego.NSRouter("/taskUpd/:task_code", &controllers.TaskController{}, "post:Put"),
		beego.NSRouter("/taskDel/:task_code", &controllers.TaskController{}, "delete:Delete"),
	)
	beego.AddNamespace(ns)
}

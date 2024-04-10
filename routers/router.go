// @version 1.0.0
// @Title Scheduler Test API
// @Description API fot scheduler application
// @Contact pavelkiwiandrosov@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/v1
// @schemes http https

package routers

import (
	"scheduler_api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		//endpoints for task
		beego.NSRouter("/task", &controllers.TaskController{}, "get:GetAll"),
		beego.NSRouter("/task", &controllers.TaskController{}, "post:Post"),
		beego.NSRouter("/task/:task_code", &controllers.TaskController{}, "get:Get"),
		beego.NSRouter("/task/taskUpd/:task_code", &controllers.TaskController{}, "post:Put"),
		beego.NSRouter("/task/taskDel/:task_code", &controllers.TaskController{}, "delete:Delete"),
		beego.NSRouter("/task/taskRecUpd/:task_code", &controllers.TaskController{}, "post:PutCascade"),
		beego.NSRouter("/task/taskRecDel/:task_code", &controllers.TaskController{}, "delete:DeleteCascade"),
		//endpoints for user
		beego.NSRouter("/user/userlist", &controllers.CoreController{}, "get:GetUserList"),
		beego.NSRouter("/user/login", &controllers.CoreController{}, "post:Login"),
		beego.NSRouter("/user/addorupd", &controllers.CoreController{}, "post:AddOrUpdateUser"),
		beego.NSRouter("/user/delete", &controllers.CoreController{}, "delete:DeleteUser"),
		beego.NSRouter("/user/updpasswd", &controllers.CoreController{}, "post:ModifyPassword"),
		beego.NSRouter("/user/rstpasswd", &controllers.CoreController{}, "post:ResetPassword"),
		//endpoints for rooms
		beego.NSRouter("/room", &controllers.RoomController{}, "get:GetRoomList"),
		beego.NSRouter("/room", &controllers.RoomController{}, "post:PostAndUpdRoom"),
		beego.NSRouter("/room", &controllers.RoomController{}, "delete:DeleteRoom"),
		//endpoints for projects
		beego.NSRouter("/project", &controllers.ProjectController{}, "get:GetProjectList"),
		beego.NSRouter("/project", &controllers.ProjectController{}, "post:PostAndUpdProject"),
		beego.NSRouter("/project", &controllers.ProjectController{}, "delete:DeleteProject"),
		//endpoints for clients
		beego.NSRouter("/client", &controllers.ClientController{}, "get:GetClientList"),
		beego.NSRouter("/client", &controllers.ClientController{}, "post:PostAndUpdClient"),
		beego.NSRouter("/client", &controllers.ClientController{}, "delete:DeleteClient"),
		//endpoints for positions
		beego.NSRouter("/position", &controllers.PositionController{}, "get:GetPositionList"),
		beego.NSRouter("/position", &controllers.PositionController{}, "post:PostAndUpdPosition"),
		beego.NSRouter("/position", &controllers.PositionController{}, "delete:DeletePosition"),
	)
	beego.AddNamespace(ns)
}

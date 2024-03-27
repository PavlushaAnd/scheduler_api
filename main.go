package main

import (
	"fmt"
	"scheduler_api/models"
	_ "scheduler_api/routers"

	_ "time/tzdata"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	//DB connection
	orm.Debug = true
	//conn := "root:WhisperingW@ves22@tcp(127.0.0.1:3306)/schedulerdbdev?charset=utf8&parseTime=true&loc=Local"
	conn := "rooty:WhisperingW@ves22@tcp(192.168.69.52:3306)/schedulerdb?charset=utf8&parseTime=true&loc=Local"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", conn)
	if err != nil {
		errors.New(fmt.Sprintf("connect to database failed, err: %v", err))
		return
	}
	orm.RegisterModel(new(models.Task))
	//orm.RunSyncdb("default", true, true)
	orm.RunSyncdb("default", false, true)
	//CORS permissions
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET ", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Access-Control-Allow-Origin", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Access-Control-Allow-Headers"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

package main

import (
	"fmt"
	"scheduler_api/models"
	_ "scheduler_api/routers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	//DB connection
	orm.Debug = true
	conn := "rooty:WhisperingW@ves22@tcp(127.0.0.1:3306)/schedulerdb?charset=utf8&parseTime=true&loc=Local"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	fmt.Println("1")
	err := orm.RegisterDataBase("default", "mysql", conn)
	fmt.Println("2")
	if err != nil {
		errors.New(fmt.Sprintf("connect to database failed, err: %v", err))
		return
	}
	orm.RegisterModel(new(models.Task))
	orm.RegisterModel(new(models.User))
	orm.RunSyncdb("default", true, true)

	//CORS permitions
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET ", "OPTIONS"},
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

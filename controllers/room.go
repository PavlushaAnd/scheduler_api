package controllers

import (
	"encoding/json"
	"fmt"
	"scheduler_api/core"
	"scheduler_api/logger"
	"scheduler_api/models"
	"scheduler_api/utils"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type RoomController struct {
	core.Core
}

type RoomView struct {
	Id       int    `json:"room_id"`
	Name     string `json:"room_name"`
	Sequence int    `json:"room_sequence"`
}

// swagger comments
func (c *RoomController) PostAndUpdRoom() {
	c.RequireLogin()

	roomDetailStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", roomDetailStr)
	d := &RoomView{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	room := &models.Room{
		Name:         d.Name,
		Sequence:     d.Sequence,
		LastModified: time.Now(),
		Version:      0,
	}

	if d.Id == 0 {
		err = models.InsertRoom(room, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		room.Id = d.Id
		err = models.UpdateRoom(room, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success"}
	c.ServeJSON()
}

// swagger comment
func (c *RoomController) GetRoomList() {
	c.RequireLogin()

	o := orm.NewOrmUsingDB("default")
	roomList, err := models.ListRoom("", o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	roomView := make([]*RoomView, 0)
	for _, v := range roomList {
		roomView = append(roomView, &RoomView{
			Name:     v.Name,
			Sequence: v.Sequence,
			Id:       v.Id,
		})
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: roomView}
	c.ServeJSON()
}

// swagger comment
func (c *RoomController) DeleteRoom() {
	c.RequireLogin()

	delRoom := c.GetString("room_name")

	o := orm.NewOrmUsingDB("default")
	room, err := models.GetRoom(delRoom, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find room %s, err: - %s", delRoom, err.Error())}
		c.ServeJSON()
		return
	}
	err = models.DeleteRoom(room, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot delete room %s, err: - %s", delRoom, err.Error())}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success"}
	c.ServeJSON()
}

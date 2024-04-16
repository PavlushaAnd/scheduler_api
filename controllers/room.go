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
	Inactive bool   `json:"room_inactive"`
	Sequence int    `json:"room_sequence"`
}

// swagger comments
// @Title add or update room
// @tags rooms
// @Description add or update room
// @Param	roomDetail		body		RoomView	true		"room detail"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /room [post]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string false "With the bearer in front"
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

	if d.Id == 0 {
		room := &models.Room{
			Name:         d.Name,
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			CreatorCode:  c.CurrentUserDetail.UserCode,
			CreatedAt:    time.Now(),
			LastModified: time.Now(),
		}

		err = models.InsertRoom(room, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		room := &models.Room{
			Name:         d.Name,
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			EditorCode:   c.CurrentUserDetail.UserCode,
			LastModified: time.Now(),
		}

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
// @Title Get room list
// @tags rooms
// @Description get room list
// @Success 200 {object} utils.JSONStruct{data=RoomView}
// @Failure 400
// @router /room [get]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
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
			Inactive: v.Inactive,
			Id:       v.Id,
		})
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: roomView}
	c.ServeJSON()
}

// swagger comment
// @Title  delete room
// @tags rooms
// @Description delete room
// @Param	room_name		query		string	true		"room name"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /room [delete]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *RoomController) DeleteRoom() {
	c.RequireLogin()

	delRoom := c.GetString("room_name")

	o := orm.NewOrmUsingDB("default")
	room, err := models.GetRoom(delRoom, o)
	if (err != nil) || (room == nil) {
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

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

type PositionController struct {
	core.Core
}

type PositionView struct {
	Id          int    `json:"position_id"`
	Description string `json:"position_description"`
	Name        string `json:"position_name"`
	Code        string `json:"position_code"`
	Inactive    bool   `json:"position_inactive"`
	Sequence    int    `json:"position_sequence"`
}

// swagger comments
func (c *PositionController) PostAndUpdPosition() {
	c.RequireLogin()

	positionDetailStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", positionDetailStr)
	d := &PositionView{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	if d.Id == 0 {
		position := &models.Position{
			Name:         d.Name,
			Description:  d.Description,
			Code:         d.Code,
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			CreatorCode:  c.CurrentUserDetail.UserCode,
			CreatedAt:    time.Now(),
			LastModified: time.Now(),
		}

		err = models.InsertPosition(position, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		position := &models.Position{
			Name:         d.Name,
			Code:         d.Code,
			Description:  d.Description,
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			EditorCode:   c.CurrentUserDetail.UserCode,
			LastModified: time.Now(),
		}

		position.Id = d.Id
		err = models.UpdatePosition(position, o)
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
func (c *PositionController) GetPositionList() {
	c.RequireLogin()

	o := orm.NewOrmUsingDB("default")
	positionList, err := models.ListPosition("", o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	positionView := make([]*PositionView, 0)
	for _, v := range positionList {
		positionView = append(positionView, &PositionView{
			Name:        v.Name,
			Description: v.Description,
			Code:        v.Code,
			Sequence:    v.Sequence,
			Inactive:    v.Inactive,
			Id:          v.Id,
		})
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: positionView}
	c.ServeJSON()
}

// swagger comment
func (c *PositionController) DeletePosition() {
	c.RequireLogin()

	delPosition := c.GetString("position_code")

	o := orm.NewOrmUsingDB("default")
	position, err := models.GetPosition(delPosition, o)
	if (err != nil) || (position == nil) {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find position %s, err: - %s", delPosition, err.Error())}
		c.ServeJSON()
		return
	}
	err = models.DeletePosition(position, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot delete position %s, err: - %s", delPosition, err.Error())}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success"}
	c.ServeJSON()
}

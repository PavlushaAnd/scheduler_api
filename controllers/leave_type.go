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

type LeaveTypeController struct {
	core.Core
}

type LeaveTypeView struct {
	Id              int    `json:"leave_type_id"`
	Name            string `json:"leave_type_name"`
	Code            string `json:"leave_type_code"`
	ColorText       string `json:"leave_type_color_text"`
	ColorBackground string `json:"leave_type_color_background"`
	Inactive        bool   `json:"leave_type_inactive"`
	Sequence        int    `json:"leave_type_sequence"`
}

// swagger comments
// @Title add or update leave type
// @tags leave types
// @Description add or update leave type
// @Param	leaveTypeDetail		body		LeaveTypeView	true		"leave type detail"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /leave_type [post]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string false "With the bearer in front"
func (c *LeaveTypeController) PostAndUpdLeaveType() {
	c.RequireLogin()

	leaveTypeDetailStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", leaveTypeDetailStr)
	d := &LeaveTypeView{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	if d.Id == 0 {
		leave_type := &models.LeaveType{
			Name:            d.Name,
			Code:            d.Code,
			Inactive:        d.Inactive,
			Sequence:        d.Sequence,
			ColorText:       d.ColorText,
			ColorBackground: d.ColorBackground,
			CreatorCode:     c.CurrentUserDetail.UserCode,
			CreatedAt:       time.Now(),
			LastModified:    time.Now(),
		}

		err = models.InsertLeaveType(leave_type, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		leave_type := &models.LeaveType{
			Name:            d.Name,
			Code:            d.Code,
			Inactive:        d.Inactive,
			Sequence:        d.Sequence,
			ColorText:       d.ColorText,
			ColorBackground: d.ColorBackground,
			EditorCode:      c.CurrentUserDetail.UserCode,
			LastModified:    time.Now(),
		}

		leave_type.Id = d.Id
		err = models.UpdateLeaveType(leave_type, o)
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
// @Title Get leave type list
// @tags leave types
// @Description get leave type list
// @Param	filter_inactive	query	bool	false	"hide inactive leave types"
// @Success 200 {object} utils.JSONStruct{data=LeaveTypeView}
// @Failure 400
// @router /leave_type [get]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *LeaveTypeController) GetLeaveTypeList() {
	c.RequireLogin()
	filterInactive, _ := c.GetBool("filter_inactive")

	o := orm.NewOrmUsingDB("default")
	leaveTypeList, err := models.ListLeaveType("", filterInactive, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	leaveTypeView := make([]*LeaveTypeView, 0)
	for _, v := range leaveTypeList {
		leaveTypeView = append(leaveTypeView, &LeaveTypeView{
			Name:            v.Name,
			Code:            v.Code,
			ColorBackground: v.ColorBackground,
			ColorText:       v.ColorText,
			Sequence:        v.Sequence,
			Inactive:        v.Inactive,
			Id:              v.Id,
		})
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: leaveTypeView}
	c.ServeJSON()
}

// swagger comment
// @Title  delete leave type
// @tags leave types
// @Description delete leave types
// @Param	leave_type_code		query		string	true		"leave type code"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /leave_type [delete]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *LeaveTypeController) DeleteLeaveType() {
	c.RequireLogin()

	delLeaveType := c.GetString("leave_type_code")

	o := orm.NewOrmUsingDB("default")
	leave_type, err := models.GetLeaveType(delLeaveType, o)
	if (err != nil) || (leave_type == nil) {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find leave_type %s, err: - %s", delLeaveType, err.Error())}
		c.ServeJSON()
		return
	}
	err = models.DeleteLeaveType(leave_type, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot delete leave_type %s, err: - %s", delLeaveType, err.Error())}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success"}
	c.ServeJSON()
}

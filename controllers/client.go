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

type ClientController struct {
	core.Core
}

type ClientView struct {
	Id       int    `json:"client_id"`
	Name     string `json:"client_name"`
	Code     string `json:"client_code"`
	Inactive bool   `json:"client_inactive"`
	Sequence int    `json:"client_sequence"`
}

// swagger comments
// @Title add or update client
// @tags clients
// @Description add or update client
// @Param	clientDetail		body		ClientView	true		"client detail"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /client [post]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string false "With the bearer in front"
func (c *ClientController) PostAndUpdClient() {
	c.RequireLogin()

	clientDetailStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", clientDetailStr)
	d := &ClientView{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	if d.Id == 0 {
		client := &models.Client{
			Name:         d.Name,
			Code:         d.Code,
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			CreatorCode:  c.CurrentUserDetail.UserCode,
			CreatedAt:    time.Now(),
			LastModified: time.Now(),
		}

		err = models.InsertClient(client, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		client := &models.Client{
			Name:         d.Name,
			Code:         d.Code,
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			EditorCode:   c.CurrentUserDetail.UserCode,
			LastModified: time.Now(),
		}

		client.Id = d.Id
		err = models.UpdateClient(client, o)
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
// @Title Get client list
// @tags clients
// @Description get client list
// @Param	filter_inactive	query	bool	false	"hide inactive clients"
// @Success 200 {object} utils.JSONStruct{data=ClientView}
// @Failure 400
// @router /client [get]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *ClientController) GetClientList() {
	c.RequireLogin()
	filterInactive, _ := c.GetBool("filter_inactive")

	o := orm.NewOrmUsingDB("default")
	clientList, err := models.ListClient("", filterInactive, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	clientView := make([]*ClientView, 0)
	for _, v := range clientList {
		clientView = append(clientView, &ClientView{
			Name:     v.Name,
			Code:     v.Code,
			Sequence: v.Sequence,
			Inactive: v.Inactive,
			Id:       v.Id,
		})
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: clientView}
	c.ServeJSON()
}

// swagger comment
// @Title  delete client
// @tags clients
// @Description delete client
// @Param	client_code		query		string	true		"client code"
// @Success 200 {object} utils.JSONStruct
// @Failure 400
// @router /client [delete]
// @Security ApiKeyAuth
// @SecurityDefinition BearerAuth api_key Authorization header with JWT token
// @Param Authorization header string true "With the bearer in front"
func (c *ClientController) DeleteClient() {
	c.RequireLogin()

	delClient := c.GetString("client_code")

	o := orm.NewOrmUsingDB("default")
	client, err := models.GetClient(delClient, o)
	if (err != nil) || (client == nil) {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find client %s, err: - %s", delClient, err.Error())}
		c.ServeJSON()
		return
	}
	err = models.DeleteClient(client, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot delete client %s, err: - %s", delClient, err.Error())}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success"}
	c.ServeJSON()
}

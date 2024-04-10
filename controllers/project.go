package controllers

import (
	"encoding/json"
	"fmt"
	"scheduler_api/core"
	"scheduler_api/logger"
	"scheduler_api/models"
	"scheduler_api/utils"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type ProjectController struct {
	core.Core
}

type ProjectView struct {
	Id         int    `json:"project_id"`
	Name       string `json:"project_name"`
	Inactive   bool   `json:"project_inactive"`
	Sequence   int    `json:"project_sequence"`
	ClientCode string `json:"client_code"`
}

// swagger comments
func (c *ProjectController) PostAndUpdProject() {
	c.RequireLogin()

	projectDetailStr := string(c.Ctx.Input.RequestBody)
	logger.D("docdetail json:", projectDetailStr)
	d := &ProjectView{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, d)
	if err != nil {
		logger.E("json.Unmarshal failed, err", err)
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "Request body is not a valid json"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrmUsingDB("default")

	if d.Id == 0 {
		project := &models.Project{
			Name:         fmt.Sprintf("%s_%s", d.ClientCode, d.Name),
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			CreatorCode:  c.CurrentUserDetail.UserCode,
			CreatedAt:    time.Now(),
			LastModified: time.Now(),
		}
		err = models.InsertProject(project, o)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("error on orm using - %s", err.Error())}
			c.ServeJSON()
			return
		}
	} else {
		project := &models.Project{
			Name:         fmt.Sprintf("%s_%s", d.ClientCode, d.Name),
			Inactive:     d.Inactive,
			Sequence:     d.Sequence,
			EditorCode:   c.CurrentUserDetail.UserCode,
			LastModified: time.Now(),
		}
		project.Id = d.Id
		err = models.UpdateProject(project, o)
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
func (c *ProjectController) GetProjectList() {
	c.RequireLogin()

	o := orm.NewOrmUsingDB("default")
	projectList, err := models.ListProject("", o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	projectView := make([]*ProjectView, 0)
	for _, v := range projectList {
		projectName := strings.Split(v.Name, "_")
		projectView = append(projectView, &ProjectView{
			Name:       projectName[1],
			ClientCode: projectName[0],
			Inactive:   v.Inactive,
			Sequence:   v.Sequence,
			Id:         v.Id,
		})
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success", Data: projectView}
	c.ServeJSON()
}

// swagger comment
func (c *ProjectController) DeleteProject() {
	c.RequireLogin()

	delProject := fmt.Sprintf("%s_%s", c.GetString("client_code"), c.GetString("project_name"))

	o := orm.NewOrmUsingDB("default")
	project, err := models.GetProject(delProject, o)
	if (err != nil) || (project == nil) {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot find project %s, err: - %s", delProject, err.Error())}
		c.ServeJSON()
		return
	}
	err = models.DeleteProject(project, o)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: fmt.Sprintf("Cannot delete project %s, err: - %s", delProject, err.Error())}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success"}
	c.ServeJSON()
}

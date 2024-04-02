package controllers

import (
	"encoding/json"
	"fmt"
	"scheduler_api/core"
	"scheduler_api/models"
	"scheduler_api/utils"

	"github.com/beego/beego/v2/client/orm"
)

// Operations about Tasks
type TaskController struct {
	core.Core
}

// @Title CreateTask
// @Description create single/recurrent Task
// @Param	body		body 	models.FTask	true		"body for user content"
// @Success 200  success post!
// @Failure 403 error message
// @router / [post]
func (c *TaskController) Post() {
	c.RequireLogin()
	o := orm.NewOrm()

	var task *models.FTask
	json.Unmarshal(c.Ctx.Input.RequestBody, &task)
	//forcing non admin users to schedule only tasks for them
	if (c.CurrentUserDetail.Role != "admin") && (c.CurrentUserDetail.UserCode != task.UserCode) {
		task.UserCode = c.CurrentUserDetail.UserCode
	}
	_, addErr := models.AddTask(o, task)
	if addErr != nil {
		c.Data["json"] = addErr.Error()
	} else {
		c.Data["json"] = "success post!"
	}
	c.ServeJSON()
}

// @Title GetAllTasks
// @Description get all Tasks
// @Success 200 {object} models.FTask
// @router / [get]
func (c *TaskController) GetAll() {
	c.RequireLogin()
	var res []*models.FTask
	tasks, err := models.GetAllTasks()
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	if c.CurrentUserDetail.Role == "admin" {
		c.Data["json"] = tasks
		c.ServeJSON()
	} else {
		for _, task := range tasks {
			if task.UserCode == c.CurrentUserDetail.UserCode {
				res = append(res, task)
			}
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
}

// @Title GetTask
// @Description get task by task_code
// @Param	task_code		path 	string	true		"The key for Task"
// @Success 200 {object} models.FTask
// @Failure 403 {task_code} is empty
// @router /:task_code [get]
func (c *TaskController) Get() {
	c.RequireLogin()
	//permissions should be discussed
	if c.CurrentUserDetail.Role != "admin" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - permission denied"}
		c.ServeJSON()
		return
	}
	tid := c.GetString(":task_code")
	if tid != "" {
		task, err := models.GetTask(tid)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = task
		}
	}
	c.ServeJSON()
}

// @Title UpdateTask
// @Description update the task
// @Param	task_code		path 	string	true		"The task_code you want to update"
// @Param	body		body 	models.FTask	true		"body for task content"
// @Success 200 {object} models.FTask
// @Failure 403 error message
// @router /taskUpd/:task_code [post]
func (c *TaskController) Put() {
	c.RequireLogin()
	tid := c.GetString(":task_code")
	if tid != "" {
		var task models.FTask
		json.Unmarshal(c.Ctx.Input.RequestBody, &task)
		if (c.CurrentUserDetail.Role != "admin") && (c.CurrentUserDetail.UserCode != task.UserCode) {
			task.UserCode = c.CurrentUserDetail.UserCode
		}
		err := models.UpdateTask(tid, &task)
		if err != nil {
			c.Data["json"] = err.Error()
		}
		c.ServeJSON()
	}
}

// @Title DeleteTask
// @Description delete the task
// @Param	task_code		path 	string	true		"The task_code you want to delete"
// @Success 200 {task_code} delete success!
// @Failure 403 {task_code} is empty
// @router /taskDel/:task_code [delete]
func (c *TaskController) Delete() {
	c.RequireLogin()
	if c.CurrentUserDetail.Role != "admin" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "error - permission denied"}
		c.ServeJSON()
		return
	}
	tid := c.GetString(":task_code")
	smth, err := models.DeleteTask(tid)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		if smth {
			c.Data["json"] = fmt.Sprintf("%v delete success!", tid)
		} else {
			c.Data["json"] = fmt.Sprintf("%v is empty", tid)
		}
	}
	c.ServeJSON()
}

// @Title DeleteCascadeTask
// @Description delete recurrence by Task
// @Param	task_code		path 	string	true		"The task_code you want to delete"
// @Success 200 {task_code} delete success!
// @Failure 403 {task_code} is empty
// @router /taskRecDel/:task_code [delete]
func (c *TaskController) DeleteCascade() {
	c.RequireLogin()
	tid := c.GetString(":task_code")
	smth, err := models.CascadeDeleteRecurrentTask(tid)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		if smth {
			c.Data["json"] = fmt.Sprintf("%v delete success!", tid)
		} else {
			c.Data["json"] = fmt.Sprintf("%v is empty", tid)
		}
	}
	c.ServeJSON()
}

// @Title UpdateCascadeTask
// @Description update recurrence by Tasks (can receive FTask but will update only Title, Description and Location)
// @Param	task_code		path 	string	true		"The task_code you want to update"
// @Param	body		body 	models.FTask	true		"body for task content"
// @Success 200 {object} models.FTask
// @Failure 403 error message
// @router /taskRecUpd/:task_code [post]
func (c *TaskController) PutCascade() {
	c.RequireLogin()
	tid := c.GetString(":task_code")
	if tid != "" {
		var task models.FTask
		json.Unmarshal(c.Ctx.Input.RequestBody, &task)
		if (c.CurrentUserDetail.Role != "admin") && (c.CurrentUserDetail.UserCode != task.UserCode) {
			task.UserCode = c.CurrentUserDetail.UserCode
		}
		tt, err := models.CascadeUpdateRecurrentTask(tid, &task)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = tt
		}
	}
	c.Data["json"] = fmt.Sprint("no task code in request")
	c.ServeJSON()
}

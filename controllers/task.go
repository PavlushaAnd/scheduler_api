package controllers

import (
	"encoding/json"
	"scheduler_api/models"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Tasks
type TaskController struct {
	beego.Controller
}

// @Title PostTask
// @Description post task
// @Param	body		body 	models.Task	true		"body for user content"
// @Success 200 {string} models.Task.Id
// @Failure 403 body is empty
// @router / [post]
func (t *TaskController) Post() {
	var task models.FTask
	json.Unmarshal(t.Ctx.Input.RequestBody, &task)
	t_code, _ := models.AddTask(task)
	t.Data["json"] = map[string]string{"task_code": t_code}
	t.ServeJSON()
}

// @Title GetAllTasks
// @Description get all Tasks
// @Success 200 {object} models.Task
// @router / [get]
func (t *TaskController) GetAll() {
	tasks, err := models.GetAllTasks()
	if err != nil {
		t.Data["json"] = err.Error()
	} else {
		t.Data["json"] = tasks
	}
	t.ServeJSON()
}

// @Title GetTask
// @Description get user by tid
// @Param	tid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Task
// @Failure 403 :tid is empty
// @router /:tid [get]
func (t *TaskController) Get() {
	tid := t.GetString(":tid")
	if tid != "" {
		task, err := models.GetTask(tid)
		if err != nil {
			t.Data["json"] = err.Error()
		} else {
			t.Data["json"] = task
		}
	}
	t.ServeJSON()
}

// @Title UpdateTask
// @Description update the task
// @Param	tid		path 	string	true		"The tid you want to update"
// @Param	body		body 	models.Task	true		"body for task content"
// @Success 200 {object} models.Task
// @Failure 403 :tid is not int
// @router /:tid [put]
func (t *TaskController) Put() {
	tid := t.GetString(":tid")
	if tid != "" {
		var task models.Task
		json.Unmarshal(t.Ctx.Input.RequestBody, &task)
		tt, err := models.UpdateTask(tid, &task)
		if err != nil {
			t.Data["json"] = err.Error()
		} else {
			t.Data["json"] = tt
		}
	}
	t.ServeJSON()
}

// @Title DeleteTask
// @Description delete the task
// @Param	tid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 tid is empty
// @router /:tid [delete]
func (t *TaskController) Delete() {
	tid := t.GetString(":tid")
	models.DeleteTask(tid)
	t.Data["json"] = "delete success!"
	t.ServeJSON()
}

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
// @router /task [post]
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
// @router /task [get]
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
// @Description get task by task_code
// @Param	task_code		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Task
// @Failure 403 {task_code} is empty
// @router /task/:task_code [get]
func (t *TaskController) Get() {
	tid := t.GetString(":task_code")
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
// @Param	task_code		path 	string	true		"The task_code you want to update"
// @Param	body		body 	models.Task	true		"body for task content"
// @Success 200 {object} models.Task
// @Failure 403 {task_code} is not int
// @router /taskUpd/:task_code [put]
func (t *TaskController) Put() {
	tid := t.GetString(":task_code")
	if tid != "" {
		var task models.FTask
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
// @Param	task_code		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 {task_code} is empty
// @router /taskDel/:task_code [delete]
func (t *TaskController) Delete() {
	tid := t.GetString(":task_code")
	models.DeleteTask(tid)
	t.Data["json"] = "delete success!"
	t.ServeJSON()
}

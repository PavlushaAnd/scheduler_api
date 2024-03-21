package controllers

import (
	"encoding/json"
	"fmt"
	"scheduler_api/models"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Tasks
type TaskController struct {
	beego.Controller
}

// @Title CreateTask
// @Description create single/recurrent Task
// @Param	body		body 	models.FTask	true		"body for user content"
// @Success 200  success post!
// @Failure 403 error message
// @router / [post]
func (t *TaskController) Post() {
	var task *models.FTask
	json.Unmarshal(t.Ctx.Input.RequestBody, &task)
	_, addErr := models.AddTask(task)
	if addErr != nil {
		t.Data["json"] = addErr.Error()
	} else {
		t.Data["json"] = "success post!"
	}

	t.ServeJSON()
}

// @Title GetAllTasks
// @Description get all Tasks
// @Success 200 {object} models.FTask
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
// @Description get task by task_code
// @Param	task_code		path 	string	true		"The key for Task"
// @Success 200 {object} models.FTask
// @Failure 403 {task_code} is empty
// @router /:task_code [get]
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
// @Param	body		body 	models.FTask	true		"body for task content"
// @Success 200 {object} models.FTask
// @Failure 403 error message
// @router /taskUpd/:task_code [post]
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
// @Param	task_code		path 	string	true		"The task_code you want to delete"
// @Success 200 {task_code} delete success!
// @Failure 403 {task_code} is empty
// @router /taskDel/:task_code [delete]
func (t *TaskController) Delete() {
	tid := t.GetString(":task_code")
	smth, err := models.DeleteTask(tid)
	if err != nil {
		t.Data["json"] = err.Error()
	} else {
		if smth {
			t.Data["json"] = fmt.Sprintf("%v delete success!", tid)
		} else {
			t.Data["json"] = fmt.Sprintf("%v is empty", tid)
		}
	}
	t.ServeJSON()
}

// @Title DeleteCascadeTask
// @Description delete recurrence by Task
// @Param	task_code		path 	string	true		"The task_code you want to delete"
// @Success 200 {task_code} delete success!
// @Failure 403 {task_code} is empty
// @router /taskRecDel/:task_code [delete]
func (t *TaskController) DeleteCascade() {
	tid := t.GetString(":task_code")
	smth, err := models.CascadeDeleteRecurrentTask(tid)
	if err != nil {
		t.Data["json"] = err.Error()
	} else {
		if smth {
			t.Data["json"] = fmt.Sprintf("%v delete success!", tid)
		} else {
			t.Data["json"] = fmt.Sprintf("%v is empty", tid)
		}
	}
	t.ServeJSON()
}

// @Title UpdateCascadeTask
// @Description update recurrence by Tasks (can receive FTask but will update only Title, Description and Location)
// @Param	task_code		path 	string	true		"The task_code you want to update"
// @Param	body		body 	models.FTask	true		"body for task content"
// @Success 200 {object} models.FTask
// @Failure 403 error message
// @router /taskRecUpd/:task_code [post]
func (t *TaskController) PutCascade() {
	tid := t.GetString(":task_code")
	if tid != "" {
		var task models.FTask
		json.Unmarshal(t.Ctx.Input.RequestBody, &task)
		tt, err := models.CascadeUpdateRecurrentTask(tid, &task)
		if err != nil {
			t.Data["json"] = err.Error()
		} else {
			t.Data["json"] = tt
		}
	}
	t.ServeJSON()
}

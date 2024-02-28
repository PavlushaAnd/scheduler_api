package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

var (
	TaskList map[string]*Task
)

func init() {
	TaskList = make(map[string]*Task)
	t := Task{"user_11111", "task1", "simple description", "default location", "01.01.2024 17:25", "01.01.2024 18:20"}
	TaskList["user_11111"] = &t
}

type Task struct {
	Task_code   string
	Title       string
	Description string
	Location    string
	StartDate   string `json:"StartDate" orm:"auto_now_add;type(datetime)"`
	EndDate     string `json:"EndDate" orm:"auto_now;type(datetime)"`
}

func AddTask(t Task) (string, error) {
	o := orm.NewOrm()

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	TaskList[t.Task_code] = &t
	_, insertErr := o.Insert(&t)
	if insertErr != nil {
		return "", errors.New("failed to insert task to database")
	}

	return t.Task_code, nil
}

func GetTask(tid string) (u *Task, err error) {
	if t, ok := TaskList[tid]; ok {
		return t, nil
	}
	return nil, errors.New("Task not exists")
}

func GetAllTasks() map[string]*Task {
	return TaskList
}

func UpdateTask(tid string, tt *Task) (a *Task, err error) {
	if t, ok := TaskList[tid]; ok {
		if tt.Title != "" {
			t.Title = tt.Title
		}
		if tt.Description != "" {
			t.Description = tt.Description
		}
		if tt.Location != "" {
			t.Location = tt.Location
		}
		if tt.StartDate != "" {
			t.StartDate = tt.StartDate
		}
		if tt.EndDate != "" {
			t.EndDate = tt.EndDate
		}
		return t, nil
	}
	return nil, errors.New("Task Not Exist")
}

func DeleteTask(tid string) {
	delete(TaskList, tid)
}

package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

var (
	TaskList map[string]*Task
)

func init() {
	TaskList = make(map[string]*Task)
	t := FTask{"user_11111", "task1", "simple description", "default location", "2024.01.01 17:00", "2024.01.01 20:00"}
	tb := ConvertToBackend(t)
	TaskList["user_11111"] = &tb
}

type Task struct {
	Task_code   string
	Title       string
	Description string
	Location    string
	StartDate   time.Time `json:"StartDate" orm:"auto_now_add ;type(datetime)"`
	EndDate     time.Time `json:"EndDate" orm:"auto_now; type(datetime)"`
}

type FTask struct {
	Task_code   string
	Title       string
	Description string
	Location    string
	StartDate   string
	EndDate     string
}

func AddTask(t FTask) (string, error) {
	o := orm.NewOrm()

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	tb := ConvertToBackend(t)
	TaskList[t.Task_code] = &tb
	_, insertErr := o.Insert(&tb)
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
		if tt.StartDate.IsZero() {
			t.StartDate = tt.StartDate
		}
		if tt.EndDate.IsZero() {
			t.EndDate = tt.EndDate
		}
		return t, nil
	}
	return nil, errors.New("Task Not Exist")
}

func DeleteTask(tid string) {
	delete(TaskList, tid)
}

const customLayout = "2006.01.02 15:04"

func ConvertToBackend(t FTask) Task {
	var res Task
	startDate, err := time.ParseInLocation(customLayout, t.StartDate, time.Local)
	if err != nil {
		errors.New(fmt.Sprintf("Error parsing StartDate:", err))
	}
	res.StartDate = startDate

	endDate, err := time.ParseInLocation(customLayout, t.EndDate, time.Local)
	if err != nil {
		errors.New(fmt.Sprintf("Error parsing EndDate:", err))
	}
	res.EndDate = endDate
	res.Title = t.Title
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.Location = t.Location
	return res
}

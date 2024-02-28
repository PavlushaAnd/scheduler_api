package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	TaskList map[string]*Task
)

func init() {
	TaskList = make(map[string]*Task)
	now := time.Now()
	t := Task{"user_11111", "task1", "simple description", "default location", now, now}
	TaskList["user_11111"] = &t
}

type Task struct {
	Id          string
	Title       string
	Description string
	Location    string
	StartDate   time.Time
	EndDate     time.Time
}

func AddTask(t Task) string {
	t.Id = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	TaskList[t.Id] = &t
	return t.Id
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
		/* if tt.StartDate != nil {
			t.StartDate = tt.StartDate
		}
		if tt.EndDate != nil {
			t.EndDate = tt.EndDate
		} */
		return t, nil
	}
	return nil, errors.New("Task Not Exist")
}

func DeleteTask(tid string) {
	delete(TaskList, tid)
}

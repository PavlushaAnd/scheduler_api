package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Task struct {
	Id          int       `orm:"column(id)"`
	Task_code   string    `orm:"column(task_code)"`
	Title       string    `orm:"column(title)"`
	Description string    `orm:"column(description); null"`
	Location    string    `orm:"column(location)"`
	Repeatable  bool      `orm:"column(repeatable)"`
	StartDate   time.Time `json:"StartDate" orm:"auto_now_add ;type(datetime)"`
	EndDate     time.Time `json:"EndDate" orm:"auto_now; type(datetime)"`
}

type FTask struct {
	Task_code   string
	Title       string
	Description string
	Location    string
	Repeatable  string
	StartDate   string
	EndDate     string
}

func AddTask(t *FTask) (string, error) {
	o := orm.NewOrm()

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	tb, err := ConvertToBackend(t)
	if err != nil {
		return "", err
	}
	_, insertErr := o.Insert(tb)

	if insertErr != nil {
		return "", errors.New("failed to insert task to database")
	}

	return t.Task_code, nil
}

func GetTask(tid string) (*Task, error) {
	o := orm.NewOrm()

	// Init task with code
	t := &Task{Task_code: tid}

	// Read from database
	o.QueryTable("task").Filter("task_code", tid).One(t)
	if t.StartDate.IsZero() {
		return nil, errors.New("Task not exist")
	}
	return t, nil
}

func GetAllTasks() ([]*FTask, error) {
	// New ORM object
	o := orm.NewOrm()

	var t []*Task

	count, e := o.QueryTable(new(Task)).All(&t)
	if e != nil {
		return nil, e
	}

	if count <= 0 {
		return nil, errors.New("nothing found")
	}

	var tf []*FTask
	for _, v := range t {
		tf = append(tf, ConvertToFrontend(v))
	}
	return tf, nil
}

func UpdateTask(tid string, tt *FTask) (a *FTask, err error) {
	o := orm.NewOrm()

	changeTask, convertErr := ConvertToBackend(tt)
	if convertErr != nil {
		return nil, convertErr
	}
	updTask := new(Task)
	err = o.QueryTable("task").Filter("task_code", tid).One(updTask)
	if err == orm.ErrNoRows {
		return nil, fmt.Errorf("item with ID %v not found", tid)
	} else if err != nil {
		return nil, err
	}
	updTask.Title = changeTask.Title
	updTask.Description = changeTask.Description
	updTask.StartDate = changeTask.StartDate
	updTask.EndDate = changeTask.EndDate
	updTask.Location = changeTask.Location
	_, err = o.Update(updTask)
	if err != nil {
		return nil, err
	}
	res := ConvertToFrontend(updTask)
	return res, nil
}

func DeleteTask(tid string) (bool, error) {
	o := orm.NewOrm()

	i, err := o.QueryTable("task").Filter("task_code", tid).Delete()
	if err != nil {
		return false, errors.New("deletion problem")
	}
	if i != 0 {

		return true, nil
	}

	return false, nil
}

const customLayout = "2006.01.02 15:04"

func ConvertToBackend(t *FTask) (*Task, error) {
	res := new(Task)
	startDate, err := time.ParseInLocation(time.RFC3339Nano, t.StartDate, time.Local)
	if err != nil {
		return nil, errors.New("error parsing start_date")
	}
	res.StartDate = startDate

	endDate, err := time.ParseInLocation(time.RFC3339Nano, t.EndDate, time.Local)
	if err != nil {
		return nil, errors.New("error parsing end_date")
	}
	res.EndDate = endDate
	res.Title = t.Title
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.Location = t.Location
	return res, nil
}

func ConvertToFrontend(t *Task) *FTask {
	res := new(FTask)
	startDate := t.StartDate.Format(customLayout)
	endDate := t.EndDate.Format(customLayout)
	res.Title = t.Title
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.Location = t.Location
	res.EndDate = endDate
	res.StartDate = startDate
	return res
}

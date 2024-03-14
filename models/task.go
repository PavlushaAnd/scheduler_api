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
	Repeatable  string    `orm:"column(repeatable)"`
	StartDate   time.Time `json:"StartDate" orm:"type(datetime)"`
	EndDate     time.Time `json:"EndDate" orm:"type(datetime)"`
	RecEndDate  time.Time `json:"RecEndDate" orm:"type(datetime)"`
}

type FTask struct {
	Task_code   string
	Title       string
	Description string
	Location    string
	Repeatable  string
	StartDate   string
	EndDate     string
	RecEndDate  string
}

func AddTask(t *FTask) ([]string, error) {
	o := orm.NewOrm()

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	tb, err := ConvertTaskToBackend(t)
	if err != nil {
		return nil, err
	}
	if t.Repeatable != "" {
		rt, recErr := Recurrence(tb)
		if recErr != nil {
			return nil, recErr
		}
		for i := range rt {
			_, insertErr := o.Insert(rt[i])
			if insertErr != nil {
				return nil, errors.New("failed to insert task to database")
			}
		}
	} else {
		_, insertErr := o.Insert(tb)
		if insertErr != nil {
			return nil, errors.New("failed to insert task to database")
		}

	}
	return []string{t.Task_code}, nil
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
		tf = append(tf, ConvertTaskToFrontend(v))
	}
	return tf, nil
}

func UpdateTask(tid string, tt *FTask) (a *FTask, err error) {
	o := orm.NewOrm()

	changeTask, convertErr := ConvertTaskToBackend(tt)
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
	updTask.Repeatable = changeTask.Repeatable
	updTask.Description = changeTask.Description
	updTask.StartDate = changeTask.StartDate
	updTask.EndDate = changeTask.EndDate
	updTask.Location = changeTask.Location
	_, err = o.Update(updTask)
	if err != nil {
		return nil, err
	}
	res := ConvertTaskToFrontend(updTask)
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

func ConvertTaskToBackend(t *FTask) (*Task, error) {
	res := new(Task)
	startDate, err := time.ParseInLocation(time.RFC3339Nano, t.StartDate, time.UTC)
	if err != nil {
		return nil, errors.New("error parsing start_date")
	}
	res.StartDate = startDate

	endDate, err := time.ParseInLocation(time.RFC3339Nano, t.EndDate, time.UTC)
	if err != nil {
		return nil, errors.New("error parsing end_date")
	}
	res.EndDate = endDate

	recEndDate, err := time.ParseInLocation(time.RFC3339Nano, t.RecEndDate, time.UTC)
	if err != nil {
		return nil, errors.New("error parsing rec_end_date")
	}
	res.RecEndDate = recEndDate
	res.Title = t.Title
	res.Repeatable = t.Repeatable
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.Location = t.Location
	return res, nil
}

func ConvertTaskToFrontend(t *Task) *FTask {
	res := new(FTask)
	startDate := t.StartDate.Format(customLayout)
	endDate := t.EndDate.Format(customLayout)
	recEndDate := t.RecEndDate.Format(customLayout)
	res.RecEndDate = recEndDate
	res.Title = t.Title
	res.Repeatable = t.Repeatable
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.Location = t.Location
	res.EndDate = endDate
	res.StartDate = startDate
	return res
}

func Recurrence(t *Task) (recTaskList []*Task, recError error) {
	var (
		recStartDate, recEndDate time.Time
	)
	switch t.Repeatable {
	case "FREQ=DAILY":
		for i := 0; t.StartDate.AddDate(0, 0, i).Before(t.RecEndDate); i++ {
			recStartDate = t.StartDate.AddDate(0, 0, i)
			recEndDate = t.EndDate.AddDate(0, 0, i)

			task := *t
			task.StartDate = recStartDate
			task.EndDate = recEndDate
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
		return
	case "FREQ=WEEKLY":
		for i := 0; t.StartDate.AddDate(0, 0, 7*i).Before(t.RecEndDate); i++ {
			recStartDate = t.StartDate.AddDate(0, 0, 7*i)
			recEndDate = t.EndDate.AddDate(0, 0, 7*i)

			task := *t
			task.StartDate = recStartDate
			task.EndDate = recEndDate
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
		return
	case "FREQ=MONTHLY":
		for i := 0; t.StartDate.AddDate(0, i, 0).Before(t.RecEndDate); i++ {
			recStartDate = t.StartDate.AddDate(0, i, 0)
			recEndDate = t.EndDate.AddDate(0, i, 0)

			task := *t
			task.StartDate = recStartDate
			task.EndDate = recEndDate
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
		return
	default:
		return nil, errors.New("wrong recurrence format")
	}
}

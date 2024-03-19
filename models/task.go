package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Task struct {
	Id           int       `orm:"column(id)"`
	Task_code    string    `orm:"column(task_code)"`
	Title        string    `orm:"column(title)"`
	Description  string    `orm:"column(description); null"`
	Location     string    `orm:"column(location)"`
	Repeatable   string    `orm:"column(repeatable)"`
	StartDate    time.Time `orm:"type(datetime)"`
	EndDate      time.Time `orm:"type(datetime)"`
	RecEndDate   time.Time `orm:"type(datetime); nul"`
	RecStartDate time.Time `orm:"type(datetime); nul"`
	Version      int       `orm:"version"`
	LastModified time.Time
}

type FTask struct {
	Task_code    string
	Title        string
	Description  string
	Location     string
	Repeatable   string
	StartDate    string
	EndDate      string
	RecEndDate   string
	RecStartDate string
}

func AddTask(t *FTask) ([]string, error) {
	o := orm.NewOrm()

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	tb, err := ConvertTaskToBackend(t)
	if err != nil {
		return nil, err
	}
	if t.Repeatable != "" {
		rt, recErr := CreateRecurrence(tb)
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
		tb.LastModified = time.Now()
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

func UpdateTask(tid string, tt *FTask) (res *FTask, err error) {
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
	updTask.Description = changeTask.Description
	updTask.StartDate = changeTask.StartDate
	updTask.EndDate = changeTask.EndDate
	updTask.Location = changeTask.Location
	updTask.Repeatable = changeTask.Repeatable
	updTask.LastModified = time.Now()
	if count, _ := o.QueryTable("task").Filter("task_code", tid).Filter("version", updTask.Version).Count(); count != 0 {
		updTask.Version++
		_, err = o.Update(updTask)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("concurrency error")
	}

	res = ConvertTaskToFrontend(updTask)
	return res, nil
}

func CascadeUpdateRecurrentTask(tid string, changeTask *FTask) (res *FTask, err error) {
	o := orm.NewOrm()

	updTask := new(Task)
	err = o.QueryTable("task").Filter("task_code", tid).One(updTask)
	if err == orm.ErrNoRows {
		return nil, fmt.Errorf("item with ID %v not found", tid)
	} else if err != nil {
		return nil, err
	}
	tasks := updTask.recurrentCascadeTaskCodeParser(o)

	for _, v := range tasks {
		o.QueryTable("task").Filter("task_code", v).One(updTask)
		updTask.Title = changeTask.Title
		updTask.Description = changeTask.Description
		updTask.Location = changeTask.Location
		updTask.LastModified = time.Now()
		if count, _ := o.QueryTable("task").Filter("task_code", updTask.Task_code).Filter("version", updTask.Version).Count(); count != 0 {
			updTask.Version++
			_, err = o.Update(updTask)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("concurrency error")
		}
	}

	res = ConvertTaskToFrontend(updTask)
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

func CascadeDeleteRecurrentTask(tid string) (bool, error) {
	o := orm.NewOrm()

	var deletedItems int64
	delTask := new(Task)
	err := o.QueryTable("task").Filter("task_code", tid).One(delTask)
	if err != nil {
		return false, errors.New("deletion problem")
	}
	tasks := delTask.recurrentCascadeTaskCodeParser(o)
	for _, v := range tasks {
		i, err := o.QueryTable("task").Filter("task_code", v).Delete()
		deletedItems += i
		if err != nil {
			return false, errors.New("deletion problem")
		}
	}

	if deletedItems != 0 {
		return true, nil
	}
	return false, nil
}

func (t Task) recurrentCascadeTaskCodeParser(o orm.Ormer) (res []string) {
	var (
		tasks []*Task
	)
	o.QueryTable("task").Filter("rec_start_date", t.RecStartDate).Filter("rec_end_date", t.RecEndDate).All(&tasks)
	for _, v := range tasks {
		if v.Repeatable == t.Repeatable {
			res = append(res, v.Task_code)
		}
	}
	return
}

const customLayout = "2006.01.02 15:04"

func ConvertTaskToBackend(t *FTask) (*Task, error) {
	res := new(Task)
	startDate, err := time.ParseInLocation(time.RFC3339Nano, t.StartDate, time.UTC)
	if err != nil {
		return nil, errors.New("error parsing start_date; wrong date")
	}
	res.StartDate = startDate

	endDate, err := time.ParseInLocation(time.RFC3339Nano, t.EndDate, time.UTC)
	if err != nil {
		return nil, errors.New("error parsing end_date; wrong date")
	}
	res.EndDate = endDate

	if t.RecEndDate != "" {
		recEndDate, err := time.ParseInLocation(time.RFC3339Nano, t.RecEndDate, time.UTC)
		if err != nil {
			return nil, errors.New("error parsing rec_end_date; wrong date")
		}
		res.RecEndDate = recEndDate
		res.RecStartDate = startDate
	}

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
	res.RecStartDate = t.RecStartDate.Format(customLayout)
	return res
}

func CreateRecurrence(t *Task) (recTaskList []*Task, recError error) {
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
			task.LastModified = time.Now()
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
	case "FREQ=WEEKLY":
		for i := 0; t.StartDate.AddDate(0, 0, 7*i).Before(t.RecEndDate); i++ {
			recStartDate = t.StartDate.AddDate(0, 0, 7*i)
			recEndDate = t.EndDate.AddDate(0, 0, 7*i)

			task := *t
			task.StartDate = recStartDate
			task.EndDate = recEndDate
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			task.LastModified = time.Now()
			recTaskList = append(recTaskList, &task)
		}
	case "FREQ=MONTHLY":
		for i := 0; t.StartDate.AddDate(0, i, 0).Before(t.RecEndDate); i++ {
			recStartDate = t.StartDate.AddDate(0, i, 0)
			recEndDate = t.EndDate.AddDate(0, i, 0)

			task := *t
			task.StartDate = recStartDate
			task.EndDate = recEndDate
			task.LastModified = time.Now()
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
	case "FREQ=YEARLY":
		for i := 0; t.StartDate.AddDate(i, i, 0).Before(t.RecEndDate); i++ {
			recStartDate = t.StartDate.AddDate(i, 0, 0)
			recEndDate = t.EndDate.AddDate(i, 0, 0)

			task := *t
			task.StartDate = recStartDate
			task.EndDate = recEndDate
			task.LastModified = time.Now()
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
	default:
		return nil, errors.New("wrong recurrence format")
	}
	return
}

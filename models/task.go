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
	RecEndDate   time.Time `orm:"type(datetime); null"`
	RecStartDate time.Time `orm:"type(datetime); null"`
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
	Hours        string
	RecEndDate   string
	RecStartDate string
}

func AddTask(t *FTask) (string, error) {
	o := orm.NewOrm()

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	tb, err := ConvertTaskToBackend(t)
	if err != nil {
		return "", err
	}
	if t.Repeatable != "" {
		rt, recErr := CreateRecurrence(tb)
		if recErr != nil {
			return "", recErr
		}
		tb.LastModified = time.Now()
		_, insertErr := o.Insert(tb)
		if insertErr != nil {
			return "", insertErr
		}
		for i := range rt {
			if (rt[i].StartDate != tb.StartDate) && (rt[i].EndDate != tb.EndDate) {
				_, insertErr := o.Insert(rt[i])
				if insertErr != nil {
					return "", insertErr
				}
			}
		}
	} else {
		tb.LastModified = time.Now()
		_, insertErr := o.Insert(tb)
		if insertErr != nil {
			return "", insertErr
		}

	}
	return tb.Task_code, nil
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
	updTask.LastModified = time.Now()
	if (changeTask.Repeatable != "") && (!changeTask.RecEndDate.IsZero()) {
		updTask.Repeatable = changeTask.Repeatable
		updTask.RecEndDate = changeTask.RecEndDate
		updTask.RecStartDate = changeTask.StartDate
	}
	if count, _ := o.QueryTable("task").Filter("task_code", tid).Filter("version", updTask.Version).Count(); count != 0 {
		updTask.Version++

		_, err = o.Update(updTask)
		if err != nil {
			return nil, err
		}
		//creating recurrent tasks and deleting entry duplicate
		if (changeTask.Repeatable != "") && (!changeTask.RecEndDate.IsZero()) {
			tid, _ := AddTask(tt)
			DeleteTask(tid)
		}

	} else {
		return nil, errors.New("concurrency error")

	}

	res = ConvertTaskToFrontend(updTask)
	return res, nil
}

func CascadeUpdateRecurrentTask(tid string, tt *FTask) (res *FTask, err error) {
	o := orm.NewOrm()

	updTask := new(Task)
	changeTask, convertErr := ConvertTaskToBackend(tt)
	if convertErr != nil {
		return nil, convertErr
	}
	err = o.QueryTable("task").Filter("task_code", tid).One(updTask)
	if err == orm.ErrNoRows {
		return nil, fmt.Errorf("item with ID %v not found", tid)
	} else if err != nil {
		return nil, err
	}

	sTimeDelta := changeTask.StartDate.Sub(updTask.StartDate)
	eTimeDelta := changeTask.EndDate.Sub(updTask.EndDate)
	tasks := updTask.recurrentCascadeTaskCodeParser(o)

	for _, v := range tasks {
		o.QueryTable("task").Filter("task_code", v).One(updTask)
		updTask.Title = changeTask.Title
		updTask.Description = changeTask.Description
		updTask.Location = changeTask.Location
		updTask.LastModified = time.Now()
		if sTimeDelta != 0 {
			updTask.StartDate = updTask.StartDate.Add(sTimeDelta)
			//avoiding dst
			updTask.StartDate = time.Date(updTask.StartDate.Year(), updTask.StartDate.Month(), updTask.StartDate.Day(), changeTask.StartDate.Hour(), updTask.StartDate.Minute(), 0, 0, updTask.StartDate.Location())
		}
		if eTimeDelta != 0 {
			updTask.EndDate = updTask.EndDate.Add(eTimeDelta)
			//avoiding dst
			updTask.EndDate = time.Date(updTask.EndDate.Year(), updTask.EndDate.Month(), updTask.EndDate.Day(), changeTask.EndDate.Hour(), updTask.EndDate.Minute(), 0, 0, updTask.EndDate.Location())
		}

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
		if (v.Repeatable == t.Repeatable) && (v.Title == t.Title) {
			res = append(res, v.Task_code)
		}
	}
	return
}

const customLayout = "2006.01.02 15:04"

func ConvertTaskToBackend(t *FTask) (*Task, error) {
	res := new(Task)
	loc, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse(time.RFC3339Nano, t.StartDate)
	if err != nil {
		return nil, errors.New("error parsing start_date; wrong date")
	}
	res.StartDate = startDate.In(loc)

	if t.Hours != "" {
		duration, err := strconv.ParseFloat(t.Hours, 32)
		if err != nil {
			return nil, err
		}
		res.EndDate = res.StartDate.Add(time.Duration(duration * float64(time.Hour)))
	} else {
		endDate, err := time.Parse(time.RFC3339Nano, t.EndDate)
		if err != nil {
			return nil, errors.New("error parsing end_date; wrong date")
		}
		res.EndDate = endDate.In(loc)
	}

	if t.RecEndDate != "" {
		recEndDate, err := time.Parse(time.RFC3339Nano, t.RecEndDate)
		if err != nil {
			return nil, errors.New("error parsing rec_end_date; wrong date")
		}
		res.RecEndDate = recEndDate.In(loc)
		res.RecStartDate = startDate.In(loc)
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
	duration := t.EndDate.Sub(t.StartDate).Hours()

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
	res.Hours = fmt.Sprintf("%.2f", duration)
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
		for i := 0; t.StartDate.AddDate(i, 0, 0).Before(t.RecEndDate); i++ {
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

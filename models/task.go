package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Task struct {
	Id           int       `orm:"column(id)"`
	Task_code    string    `orm:"column(task_code)"`
	Title        string    `orm:"column(title)"`
	UserCode     string    `orm:"column(user_code)"`
	RoomName     string    `orm:"column(room_name)"`
	ProjectName  string    `orm:"column(project_name)"`
	Description  string    `orm:"column(description); null"`
	Repeatable   string    `orm:"column(repeatable)"`
	StartDate    time.Time `orm:"type(datetime); column(start_date)"`
	EndDate      time.Time `orm:"type(datetime); column(end_date)"`
	RecEndDate   time.Time `orm:"type(datetime); null; column(rec_end_date)"`
	RecStartDate time.Time `orm:"type(datetime); null; column(rec_start_date)"`
	Version      int       `orm:"version"`
	LastModified time.Time `orm:"column(last_modified)"`
	CreatedAt    time.Time `orm:"column(created_at)"`
	CreatorCode  string    `orm:"column(creator_code)"`
	EditorCode   string    `orm:"column(editor_code)"`
}

type FTask struct {
	Task_code    string
	Title        string
	Description  string
	UserCode     string
	RoomName     string
	ProjectName  string
	ClientCode   string
	Repeatable   string
	StartDate    string
	EndDate      string
	Hours        string
	RecEndDate   string
	RecStartDate string
}

func init() {
	orm.RegisterModel(new(Task))
}

func AddTask(o orm.Ormer, t *Task) (string, error) {

	err := CheckDependencies(t, o)
	if err != nil {
		return "", err
	}

	t.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if t.Repeatable != "" {
		rt, recErr := CreateRecurrence(t, t.CreatorCode)
		if recErr != nil {
			return "", recErr
		}
		_, insertErr := o.Insert(t)
		if insertErr != nil {
			return "", insertErr
		}
		for i := range rt {
			if (rt[i].StartDate != t.StartDate) && (rt[i].EndDate != t.EndDate) {
				_, insertErr := o.Insert(rt[i])
				if insertErr != nil {
					return "", insertErr
				}
			}
		}
	} else {
		t.LastModified = time.Now()
		_, insertErr := o.Insert(t)
		if insertErr != nil {
			return "", insertErr
		}

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
		tf = append(tf, ConvertTaskToFrontend(v))
	}
	return tf, nil
}

func UpdateTask(tid string, tt *FTask, code string) error {
	o := orm.NewOrm()

	changeTask, err := ConvertTaskToBackend(tt)
	if err != nil {
		return err
	}
	err = CheckDependencies(changeTask, o)
	if err != nil {
		return err
	}
	updTask := new(Task)
	err = o.QueryTable("task").Filter("task_code", tid).One(updTask)
	if err == orm.ErrNoRows {
		return fmt.Errorf("item with ID %v not found", tid)
	} else if err != nil {
		return err
	}
	updTask.Title = changeTask.Title
	updTask.Description = changeTask.Description
	updTask.StartDate = changeTask.StartDate
	updTask.EndDate = changeTask.EndDate
	updTask.ProjectName = changeTask.ProjectName
	updTask.UserCode = changeTask.UserCode
	updTask.RoomName = changeTask.RoomName
	updTask.LastModified = time.Now()
	updTask.EditorCode = code
	if (changeTask.Repeatable != "") && (!changeTask.RecEndDate.IsZero()) {
		updTask.Repeatable = changeTask.Repeatable
		updTask.RecEndDate = changeTask.RecEndDate
		updTask.RecStartDate = changeTask.StartDate
	}
	if count, _ := o.QueryTable("task").Filter("task_code", tid).Filter("version", updTask.Version).Count(); count != 0 {
		updTask.Version++

		_, err = o.Update(updTask)
		if err != nil {
			return err
		}
		//creating recurrent tasks and deleting entry duplicate
		if (changeTask.Repeatable != "") && (!changeTask.RecEndDate.IsZero()) {
			tid, _ := AddTask(o, updTask)
			DeleteTask(tid)
		}

	} else {
		return errors.New("concurrency error")

	}

	return nil
}

func CascadeUpdateRecurrentTask(tid string, tt *FTask, code string) (res *FTask, err error) {
	o := orm.NewOrm()

	updTask := new(Task)
	changeTask, convertErr := ConvertTaskToBackend(tt)
	if convertErr != nil {
		return nil, convertErr
	}
	err = CheckDependencies(changeTask, o)
	if err != nil {
		return nil, err
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
	sh := updTask.RecStartDate.Hour()
	eh := updTask.RecEndDate.Hour()
	for _, task := range tasks {
		o.QueryTable("task").Filter("task_code", task).One(updTask)
		updTask.Title = changeTask.Title
		updTask.Description = changeTask.Description
		updTask.ProjectName = changeTask.ProjectName
		updTask.UserCode = changeTask.UserCode
		updTask.RoomName = changeTask.RoomName
		updTask.LastModified = time.Now()
		updTask.EditorCode = code

		if sTimeDelta.Minutes() != 0 {
			updTask.StartDate = updTask.StartDate.Add(sTimeDelta)

			updTask.RecStartDate = updTask.RecStartDate.Add(sTimeDelta)
			//avoiding dst
			updTask.StartDate = time.Date(updTask.StartDate.Year(), updTask.StartDate.Month(), updTask.StartDate.Day(), changeTask.StartDate.Hour(), updTask.StartDate.Minute(), 0, 0, updTask.StartDate.Location())
			updTask.RecStartDate = time.Date(updTask.RecStartDate.Year(), updTask.RecStartDate.Month(), updTask.RecStartDate.Day(), sh, updTask.RecStartDate.Minute(), 0, 0, updTask.RecStartDate.Location())
		}
		if eTimeDelta.Minutes() != 0 {
			updTask.EndDate = updTask.EndDate.Add(eTimeDelta)

			updTask.RecEndDate = updTask.RecEndDate.Add(eTimeDelta)
			//avoiding dst
			updTask.EndDate = time.Date(updTask.EndDate.Year(), updTask.EndDate.Month(), updTask.EndDate.Day(), changeTask.EndDate.Hour(), updTask.EndDate.Minute(), 0, 0, updTask.EndDate.Location())
			updTask.RecEndDate = time.Date(updTask.RecEndDate.Year(), updTask.RecEndDate.Month(), updTask.RecEndDate.Day(), eh, updTask.RecEndDate.Minute(), 0, 0, updTask.RecEndDate.Location())
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
	res.RoomName = t.RoomName
	res.UserCode = t.UserCode
	res.Repeatable = t.Repeatable
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.ProjectName = t.ProjectName
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
	res.RoomName = t.RoomName
	res.UserCode = t.UserCode
	res.Repeatable = t.Repeatable
	res.Description = t.Description
	res.Task_code = t.Task_code
	res.EndDate = endDate
	res.StartDate = startDate
	res.RecStartDate = t.RecStartDate.Format(customLayout)
	res.Hours = fmt.Sprintf("%.2f", duration)

	projectName := strings.Split(t.ProjectName, "_")
	res.ClientCode = projectName[0]
	res.ProjectName = projectName[1]

	return res
}

func CreateRecurrence(t *Task, code string) (recTaskList []*Task, recError error) {
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
			task.EditorCode = code
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
			task.EditorCode = code
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
			task.EditorCode = code
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
			task.EditorCode = code
			task.Task_code = "task_" + strconv.FormatInt(time.Now().UnixNano()+int64(i), 10)
			recTaskList = append(recTaskList, &task)
		}
	default:
		return nil, errors.New("wrong recurrence format")
	}
	return
}

func CheckDependencies(t *Task, o orm.Ormer) error {
	user := User{}
	err := o.QueryTable("user").Filter("user_code", t.UserCode).One(&user)
	if err != nil {
		return fmt.Errorf("user %s not exist in the database", t.UserCode)
	}
	if user.Inactive {
		return fmt.Errorf("user %s is inactive", t.UserCode)
	}
	room := Room{}
	err = o.QueryTable("room").Filter("name", t.RoomName).One(&room)
	if err != nil {
		return fmt.Errorf("room %s not exist in the database", t.RoomName)
	}
	if room.Inactive {
		return fmt.Errorf("room %s is inactive", t.RoomName)
	}
	if strings.IndexByte(t.ProjectName, '_') == -1 {
		return fmt.Errorf("wrong project name format")
	}
	projectName := strings.Split(t.ProjectName, "_")
	clientCode := projectName[0]
	client := Client{}
	err = o.QueryTable("client").Filter("code", clientCode).One(&client)
	if err != nil {
		return fmt.Errorf("client %s not exist in the database", clientCode)
	}
	if client.Inactive {
		return fmt.Errorf("client %s is inactive", clientCode)
	}
	project := Project{}
	err = o.QueryTable("project").Filter("name", t.ProjectName).One(&project)
	if err != nil {
		return fmt.Errorf("project %s not exist in the database", t.ProjectName)
	}
	if project.Inactive {
		return fmt.Errorf("project %s is inactive", t.ProjectName)
	}

	return nil
}

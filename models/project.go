package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Project struct {
	Id           int       `orm:"column(id);auto"`
	Name         string    `orm:"column(name);unique"`
	Inactive     bool      `orm:"column(inactive)"`
	Sequence     int       `orm:"column(sequence)"`
	CreatorCode  string    `orm:"column(creator_code)"`
	EditorCode   string    `orm:"column(editor_code)"`
	CreatedAt    time.Time `orm:"column(created_at)"`
	LastModified time.Time `orm:"column(last_modified)"`
}

func init() {
	orm.RegisterModel(new(Project))
}

func InsertProject(project *Project, o orm.Ormer) error {
	_, err := o.Insert(project)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProject(project *Project, o orm.Ormer) error {
	_, err := o.Update(project, "name", "inactive", "sequence", "editor_code", "last_modified")
	if err != nil {
		return err
	}

	return nil
}

func DeleteProject(project *Project, o orm.Ormer) error {
	_, err := o.Delete(project)
	if err != nil {
		return err
	}

	return nil
}

func ListProject(projectName string, filterInactive bool, o orm.Ormer) ([]*Project, error) {
	var projects []*Project
	qs := o.QueryTable(new(Project))

	if filterInactive {
		qs = qs.Filter("inactive", 0)
	}

	if projectName != "" {
		qs = qs.Filter("name__icontains", projectName)
	}

	_, err := qs.All(&projects)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return projects, nil
}

func GetProject(projectName string, o orm.Ormer) (*Project, error) {
	project := Project{}

	qs := o.QueryTable(new(Project))

	err := qs.Filter("name", projectName).One(&project)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

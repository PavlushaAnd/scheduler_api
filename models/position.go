package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Position struct {
	Id           int       `orm:"column(id);auto"`
	Name         string    `orm:"column(name)"`
	Description  string    `orm:"column(description)"`
	Code         string    `orm:"column(code);unique"`
	Inactive     bool      `orm:"column(inactive)"`
	Sequence     int       `orm:"column(sequence)"`
	CreatorCode  string    `orm:"column(creator_code)"`
	EditorCode   string    `orm:"column(editor_code)"`
	CreatedAt    time.Time `orm:"column(created_at)"`
	LastModified time.Time `orm:"column(last_modified)"`
}

func init() {
	orm.RegisterModel(new(Position))
}

func InsertPosition(position *Position, o orm.Ormer) error {
	_, err := o.Insert(position)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePosition(position *Position, o orm.Ormer) error {
	_, err := o.Update(position, "name", "description", "code", "inactive", "sequence", "editor_code", "last_modified")
	if err != nil {
		return err
	}

	return nil
}

func DeletePosition(position *Position, o orm.Ormer) error {
	_, err := o.Delete(position)
	if err != nil {
		return err
	}

	return nil
}

func ListPosition(positionCode string, o orm.Ormer) ([]*Position, error) {
	var position []*Position
	qs := o.QueryTable(new(Position))

	if positionCode != "" {
		qs = qs.Filter("name__icontains", positionCode)
	}

	_, err := qs.All(&position)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return position, nil
}

func GetPosition(positionCode string, o orm.Ormer) (*Position, error) {
	position := Position{}

	qs := o.QueryTable(new(Position))

	err := qs.Filter("code", positionCode).One(&position)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &position, nil
}

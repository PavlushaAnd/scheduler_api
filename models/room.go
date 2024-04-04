package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Room struct {
	Id           int       `orm:"column(id);auto"`
	Name         string    `orm:"column(name);unique"`
	Sequence     int       `orm:"column(sequence)"`
	LastModified time.Time `orm:"column(last_modified)"`
	Version      int       `orm:"version"`
}

func init() {
	orm.RegisterModel(new(Room))
}

func InsertRoom(room *Room, o orm.Ormer) error {
	_, err := o.Insert(room)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRoom(room *Room, o orm.Ormer) error {
	_, err := o.Update(room)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoom(room *Room, o orm.Ormer) error {
	_, err := o.Delete(room)
	if err != nil {
		return err
	}

	return nil
}

func ListRoom(roomName string, o orm.Ormer) ([]*Room, error) {
	var rooms []*Room
	qs := o.QueryTable(new(Room))

	if roomName != "" {
		qs = qs.Filter("name__icontains", roomName)
	}

	_, err := qs.All(&rooms)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil
}

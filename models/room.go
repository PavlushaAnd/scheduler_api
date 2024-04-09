package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Room struct {
	Id           int       `orm:"column(id);auto"`
	Name         string    `orm:"column(name);unique"`
	Sequence     int       `orm:"column(sequence)"`
	Inactive     bool      `orm:"column(inactive)"`
	CreatorCode  string    `orm:"column(creator_code)"`
	EditorCode   string    `orm:"column(editor_code)"`
	CreatedAt    time.Time `orm:"column(created_at)"`
	DeletedAt    time.Time `orm:"column(deleted_at);null"`
	LastModified time.Time `orm:"column(last_modified)"`
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
	_, err := o.Update(room, "name", "inactive", "sequence", "editor_code", "last_modified")
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

func GetRoom(roomName string, o orm.Ormer) (*Room, error) {
	room := Room{}

	qs := o.QueryTable(new(Room))

	err := qs.Filter("name", roomName).One(&room)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

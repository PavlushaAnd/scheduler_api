package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Client struct {
	Id           int       `orm:"column(id);auto"`
	Name         string    `orm:"column(name)"`
	Code         string    `orm:"column(code);unique"`
	Inactive     bool      `orm:"column(inactive)"`
	Sequence     int       `orm:"column(sequence)"`
	CreatorCode  string    `orm:"column(creator_code)"`
	EditorCode   string    `orm:"column(editor_code)"`
	CreatedAt    time.Time `orm:"column(created_at)"`
	LastModified time.Time `orm:"column(last_modified)"`
}

func init() {
	orm.RegisterModel(new(Client))
}

func InsertClient(client *Client, o orm.Ormer) error {
	_, err := o.Insert(client)
	if err != nil {
		return err
	}

	return nil
}

func UpdateClient(client *Client, o orm.Ormer) error {
	_, err := o.Update(client, "name", "code", "inactive", "sequence", "editor_code", "last_modified")
	if err != nil {
		return err
	}

	return nil
}

func DeleteClient(client *Client, o orm.Ormer) error {
	_, err := o.Delete(client)
	if err != nil {
		return err
	}

	return nil
}

func ListClient(clientName string, o orm.Ormer) ([]*Client, error) {
	var client []*Client
	qs := o.QueryTable(new(Client))

	if clientName != "" {
		qs = qs.Filter("name__icontains", clientName)
	}

	_, err := qs.All(&client)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return client, nil
}

func GetClient(clientCode string, o orm.Ormer) (*Client, error) {
	client := Client{}

	qs := o.QueryTable(new(Client))

	err := qs.Filter("code", clientCode).One(&client)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

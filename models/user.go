package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id                int    `orm:"column(id);auto"`
	UserCode          string `orm:"column(user_code)"`
	UserName          string `orm:"column(user_name)"`
	Inactive          bool   `orm:"column(inactive)"`
	PhoneNo           string `orm:"column(phone_no)"`
	EmailAddress      string `orm:"column(email_address)"`
	HasUploadedPage   bool   `orm:"column(has_uploaded_page)"`
	HasRecognisedPage bool   `orm:"column(has_recognised_page)"`
	HasConfirmedPage  bool   `orm:"column(has_confirmed_page)"`
	HasPostedPage     bool   `orm:"column(has_posted_page)"`
	Password          string `orm:"column(password)"`
	Role              string `orm:"column(role)"`
	ColorText         string `orm:"column(color_text)"`
	ColorBackground   string `orm:"column(color_background)"`
}

func (t *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

func GetUser(useCode string, o orm.Ormer) (*User, error) {
	user := User{}

	qs := o.QueryTable(new(User)) //.Filter("inactive", 0)
	err := qs.Filter("user_code", useCode).One(&user)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func InsertUser(user *User, o orm.Ormer) error {
	_, err := o.Insert(user)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(user *User, o orm.Ormer) error {
	_, err := o.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserWithoutPwd(user *User, o orm.Ormer) error {
	_, err := o.Update(user, "user_name", "phone_no", "email_address", "inactive", "role", "has_uploaded_page", "has_recognised_page", "has_confirmed_page", "has_posted_page")
	if err != nil {
		return err
	}

	return nil
}

func ListUser(userName string, page, pageSize int, o orm.Ormer) ([]*User, int, error) {
	var users []*User

	qs := o.QueryTable(new(User)) //.Filter("inactive", 0)

	if userName != "" {
		qs = qs.Filter("user_name__icontains", userName)
	}

	count, err := qs.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("error on counting user - %s", err.Error())
	}
	if count == 0 {
		return nil, 0, nil
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		limit := pageSize
		qs = qs.Limit(limit).Offset(offset)
	}

	_, err = qs.All(&users)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, int(count), nil
		}
		return nil, int(count), fmt.Errorf("error on listing user - %s", err.Error())
	}

	return users, int(count), nil
}

func DelUser(user *User, o orm.Ormer) error {
	_, err := o.Delete(user)
	if err != nil {
		return err
	}

	return nil
}

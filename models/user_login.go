package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserLoginInfo struct {
	Id             int       `orm:"column(id);auto"`
	UserCode       string    `orm:"column(user_code)"`
	LastLoginTime  time.Time `orm:"column(last_login_time)"`
	LastLoginIp    string    `orm:"column(last_login_ip)"`
	LastLoginToken string    `orm:"column(last_login_token)"`
}

func (t *UserLoginInfo) TableName() string {
	return "app_user_login"
}

func init() {
	orm.RegisterModel(new(UserLoginInfo))
}

func GetUserLoginInfo(userCode string, o orm.Ormer) (*UserLoginInfo, error) {
	loginInfo := UserLoginInfo{}

	qs := o.QueryTable(new(UserLoginInfo))
	err := qs.Filter("user_code", userCode).One(&loginInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &loginInfo, nil
}

func InsertUserLoginInfo(loginInfo *UserLoginInfo, o orm.Ormer) error {
	_, err := o.Insert(loginInfo)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserLoginInfo(loginInfo *UserLoginInfo, o orm.Ormer) error {
	_, err := o.Update(loginInfo)
	if err != nil {
		return err
	}

	return nil
}

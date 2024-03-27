package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type LoginLog struct {
	Id        int       `orm:"column(id);auto"`
	LoginUser string    `orm:"column(login_user)"`
	LoginIp   string    `orm:"column(login_ip)"`
	LoginTime time.Time `orm:"column(login_time)"`
}

func (t *LoginLog) TableName() string {
	return "app_login_log"
}

func init() {
	orm.RegisterModel(new(LoginLog))
}

func InsertLoginLog(loginInfo *LoginLog, o orm.Ormer) error {
	_, err := o.Insert(loginInfo)
	if err != nil {
		return err
	}

	return nil
}

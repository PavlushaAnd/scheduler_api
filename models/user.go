package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id       int    `orm:"column(id)"`
	UserCode string `orm:"column(user_code)"`
	Username string `orm:"column(user_name)"`
	Password string `orm:"column(password)"`
}

func AddUser(u *User) (string, error) {
	o := orm.NewOrm()

	u.UserCode = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	_, insertErr := o.Insert(u)

	if insertErr != nil {
		return "", errors.New("failed to insert user to database")
	}

	return u.UserCode, nil
}

/* func GetUser(uid string) (u *User, err error) {
	if u, ok := UserList[uid]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}

func GetAllUsers() map[string]*User {
	return UserList
} */

/*
	 func UpdateUser(uid string, uu *User) (a *User, err error) {
		if u, ok := UserList[uid]; ok {
			if uu.Username != "" {
				u.Username = uu.Username
			}
			if uu.Password != "" {
				u.Password = uu.Password
			}
			if uu.Profile.Age != 0 {
				u.Profile.Age = uu.Profile.Age
			}
			if uu.Profile.Address != "" {
				u.Profile.Address = uu.Profile.Address
			}
			if uu.Profile.Gender != "" {
				u.Profile.Gender = uu.Profile.Gender
			}
			if uu.Profile.Email != "" {
				u.Profile.Email = uu.Profile.Email
			}
			return u, nil
		}
		return nil, errors.New("User Not Exist")
	}
*/
func Login(username, password string) (bool, error) {
	o := orm.NewOrm()

	q := o.QueryTable("user").Filter("user_name", username)
	if q.Exist() {
		if q.Filter("password", password).Exist() {
			return true, nil
		}
		return false, fmt.Errorf("wrong password")
	} else {
		return false, fmt.Errorf("user %v not found", username)
	}
}

/* func DeleteUser(uid string) {
	delete(UserList, uid)
} */

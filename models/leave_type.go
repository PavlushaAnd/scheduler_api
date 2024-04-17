package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type LeaveType struct {
	Id              int       `orm:"column(id);auto"`
	Name            string    `orm:"column(name)"`
	Code            string    `orm:"column(code);unique"`
	Inactive        bool      `orm:"column(inactive)"`
	Sequence        int       `orm:"column(sequence)"`
	ColorText       string    `orm:"column(color_text)"`
	ColorBackground string    `orm:"column(color_background)"`
	CreatorCode     string    `orm:"column(creator_code)"`
	EditorCode      string    `orm:"column(editor_code)"`
	CreatedAt       time.Time `orm:"column(created_at)"`
	LastModified    time.Time `orm:"column(last_modified)"`
}

func init() {
	orm.RegisterModel(new(LeaveType))
}

func InsertLeaveType(leave_type *LeaveType, o orm.Ormer) error {
	_, err := o.Insert(leave_type)
	if err != nil {
		return err
	}

	return nil
}

func UpdateLeaveType(leave_type *LeaveType, o orm.Ormer) error {
	_, err := o.Update(leave_type, "name", "code", "inactive", "sequence", "editor_code", "color_background", "color_text", "last_modified")
	if err != nil {
		return err
	}

	return nil
}

func DeleteLeaveType(leave_type *LeaveType, o orm.Ormer) error {
	_, err := o.Delete(leave_type)
	if err != nil {
		return err
	}

	return nil
}

func ListLeaveType(leaveTypeCode string, filterInactive bool, o orm.Ormer) ([]*LeaveType, error) {
	var leaveType []*LeaveType
	qs := o.QueryTable(new(LeaveType))

	if filterInactive {
		qs = qs.Filter("inactive", 0)
	}

	if leaveTypeCode != "" {
		qs = qs.Filter("name__icontains", leaveTypeCode)
	}

	_, err := qs.All(&leaveType)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return leaveType, nil
}

func GetLeaveType(leaveTypeCode string, o orm.Ormer) (*LeaveType, error) {
	leaveType := LeaveType{}

	qs := o.QueryTable(new(LeaveType))

	err := qs.Filter("code", leaveTypeCode).One(&leaveType)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &leaveType, nil
}

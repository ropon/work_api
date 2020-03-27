package models

import (
	"time"
	"workApi/dblayer"
)

type User struct {
	ID         uint      `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	PhoneNum   string    `json:"phoneNum" gorm:"column:phoneNum"`
	MachineId  string    `json:"machineId" gorm:"column:machineId"`
	ExpiryTime time.Time `json:"expiryTime" gorm:"column:expiryTime"`
	RealName   string    `json:"realName" gorm:"column:realName"`
	NickName   string    `json:"nickName" gorm:"column:nickName"`
}

// TableName 指定表名
func (u *User) TableName() string {
	return "user"
}

//通过手机号和机器码查询
func GetUser(user *User, phoneNum, machineId string) (err error) {
	err = dblayer.DB.Where("phoneNum = ? AND machineId = ?", phoneNum, machineId).First(user).Error
	if err != nil {
		return
	}
	return nil
}

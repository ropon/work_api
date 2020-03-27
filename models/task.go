package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"workApi/dblayer"
)

type TaskLog struct {
	ID         uint      `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	XiuJin     float64   `json:"xiuJin" gorm:"column:xiuJin"`
	XueBei     float64   `json:"xueBei" gorm:"column:xueBei"`
	Note       string    `json:"note" gorm:"column:note"`
	CreateTime time.Time `json:"createTime" gorm:"column:createTime"`
	UID        uint      `json:"-" gorm:"column:uid"`
}

// TableName 指定表名
func (t *TaskLog) TableName() string {
	return "taskLog"
}

// 新增时记录创建时间
func (t *TaskLog) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("CreateTime", time.Now())
	return nil
}

// 创建一条任务日志
func CreateTaskLog(t *TaskLog) (err error) {
	err = dblayer.DB.Create(t).Error
	if err != nil {
		return err
	}
	return nil
}

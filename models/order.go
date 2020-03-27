package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"workApi/dblayer"
)

type Order struct {
	ID         uint      `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	OrderId    string    `json:"orderId" gorm:"column:orderId"`
	OrderPrice float64   `json:"orderPrice" gorm:"column:orderPrice"`
	OrderSum   int64     `json:"orderSum" gorm:"column:orderSum"`
	CreateTime time.Time `json:"createTime" gorm:"column:createTime"`
	UID        uint      `json:"-" gorm:"column:uid"`
}

// TableName 指定表名
func (o *Order) TableName() string {
	return "order"
}

// 新增时记录创建时间
func (o *Order) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("CreateTime", time.Now())
	return nil
}

// 创建一条订单记录
func CreateOrder(o *Order) (err error) {
	err = dblayer.DB.Create(o).Error
	if err != nil {
		return err
	}
	return nil
}

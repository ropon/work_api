package models

import (
	"context"
	"fmt"
	"github.com/ropon/work_api/conf"
	"github.com/ropon/work_api/utils"
	"time"
)

type Service struct {
	Id            uint      `json:"id" form:"id" gorm:"primary_key,AUTO_INCREMENT"`
	SvcName       string    `json:"svc_name" form:"svc_name" gorm:"column:svc_name" sql:"unique;not null"`      //服务名称
	SvcType       string    `json:"svc_type" form:"svc_type" gorm:"column:svc_type" sql:"not null"`             //服务类型
	AuthorEmail   string    `json:"author_email" form:"author_email" gorm:"column:author_email" sql:"not null"` //服务负责人
	Ports         string    `json:"ports" form:"ports" gorm:"column:ports" sql:"unique;not null"`               //服务端口
	CreateTimeStr string    `json:"create_time" gorm:"-"`
	UpdateTimeStr string    `json:"update_time" gorm:"-"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time" sql:"type:datetime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time" sql:"type:datetime"`
}

type ServiceList []*Service

func (s *Service) TableName() string {
	return "service"
}

// FormatTime 特殊处理时间
func (s *Service) FormatTime() {
	s.CreateTimeStr = utils.FormatTime(s.CreateTime)
	s.UpdateTimeStr = utils.FormatTime(s.UpdateTime)
}

// Create 增(post /service)
func (s *Service) Create() (err error) {
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Create(s).Error
	return
}

// Delete 删(delete /service/:svc_id)
func (s *Service) Delete() (err error) {
	err = conf.MysqlDb.Delete(s).Error
	return
}

// Update 改(put /service/:svc_id)/全部
func (s *Service) Update() (err error) {
	s.UpdateTime = time.Now()
	err = conf.MysqlDb.Save(s).Error
	return
}

// Patch 改(patch /service/:svc_id)/部分
func (s *Service) Patch(v interface{}) (err error) {
	tmp := v.(*Service)
	tmp.UpdateTime = time.Now()
	err = conf.MysqlDb.Model(s).Updates(tmp).Error
	return
}

// Get 查(get /service/:svc_id)一个
func (s *Service) Get() (err error) {
	err = conf.MysqlDb.Where("id = ?", s.Id).Find(s).Error
	return
}

// GetByName 根据名称查询一个
func (s *Service) GetByName() (err error) {
	err = conf.MysqlDb.Where("svc_name = ?", s.SvcName).Find(s).Error
	return
}

// List 查(get /service)多个
func (s *Service) List(ctx context.Context, PageSize, PageNum int64) (list ServiceList, count int64, err error) {
	sp, _ := utils.ExtractChildSpan("db:get services", ctx)
	defer sp.Finish()
	list = make(ServiceList, 0)
	//默认精确匹配
	db := conf.MysqlDb.Where(s)
	//可以自定义查询
	if s.SvcName != "" {
		db = conf.MysqlDb.Where("svc_name like ?", fmt.Sprintf(`%%%s%%`, s.SvcName))
	}
	if s.Ports != "" {
		db = conf.MysqlDb.Where("ports like ?", fmt.Sprintf(`%%%s%%`, s.Ports))
	}
	//获取条件匹配后总记录数
	//err = db.Find(&ServiceList{}).Count(&count).Error
	//if err != nil {
	//	return nil, count, err
	//}
	offset, limit := utils.GetOffsetAndLimit(PageSize, PageNum)
	err = db.Model(s).Count(&count).Offset(offset).Limit(limit).Find(&list).Error
	return list, count, err
}

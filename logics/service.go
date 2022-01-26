package logics

import (
	"context"
	"fmt"
	"github.com/ropon/work_api/models"
)

type BaseData struct {
	UserEmail string `json:"user_email" form:"user_email"`
	OpsAdmin  string `json:"ops_admin" form:"ops_admin"`
}

type CUServiceReq struct {
	SvcName     string `json:"svc_name" form:"svc_name" binding:"required"`
	SvcType     string `json:"svc_type" form:"svc_type" binding:"required"`
	AuthorEmail string `json:"author_email" form:"author_email" binding:"required"`
	Ports       string `json:"ports" form:"ports" binding:"required"`
	BaseData
}

type ServiceReq struct {
	models.Service
	PageSize int64 `json:"page_size" form:"page_size"`
	PageNum  int64 `json:"page_num" form:"page_num"`
}

type ServiceRes struct {
	TotalCount  int64              `json:"total_count"`
	ServiceList models.ServiceList `json:"service_list"`
}

func (bp *BaseData) Init(userEmail, opsAdmin string) {
	if bp.UserEmail == "" {
		bp.UserEmail = userEmail
	}
	if bp.OpsAdmin == "" {
		bp.OpsAdmin = opsAdmin
	}
}

func initService(req *CUServiceReq) *models.Service {
	return &models.Service{
		SvcName:     req.SvcName,
		SvcType:     req.SvcType,
		AuthorEmail: req.AuthorEmail,
		Ports:       req.Ports,
	}
}

//创建服务返回服务详情
func CreateService(req *CUServiceReq) (*models.Service, error) {
	s := initService(req)
	err := s.GetByName()
	if err == nil && s.Id != 0 {
		return nil, fmt.Errorf("名称:%s对应服务已存在", s.SvcName)
	}

	err = s.Create()
	if err != nil {
		return nil, err
	}

	s.FormatTime()
	return s, nil
}

//通过服务ID删除指定服务
func DeleteService(id uint) error {
	do := &DbObj{
		Id:  id,
		Obj: &models.Service{Id: id},
	}
	return do.delete()
}

//通过服务ID更新指定服务全部信息
func UpdateService(id uint, req *CUServiceReq) (*models.Service, error) {
	do := &DbObj{
		Id:  id,
		Obj: &models.Service{Id: id},
	}
	if err := do.get(); err != nil {
		return nil, err
	}

	s := initService(req)
	s.Id = id
	s.CreateTime = do.Obj.(*models.Service).CreateTime

	return s, do.update(s)
}

//通过服务ID更新指定服务部分信息
func PatchUpdateService(id uint, req *ServiceReq) (interface{}, error) {
	s := req.Service
	do := &DbObj{
		Id:  id,
		Obj: &models.Service{Id: id},
	}
	if err := do.patch(&s); err != nil {
		return nil, err
	}
	return do.Obj, nil
}

//获取服务列表
func GetServices(ctx context.Context, req *ServiceReq) (*ServiceRes, error) {
	s := req.Service
	sl, count, err := s.List(ctx, req.PageSize, req.PageNum)
	if err != nil {
		return nil, err
	}
	for _, s := range sl {
		s.FormatTime()
	}
	res := &ServiceRes{
		TotalCount:  count,
		ServiceList: sl,
	}
	return res, nil
}

//获取单个服务详情
func GetService(id uint) (interface{}, error) {
	do := &DbObj{
		Id:  id,
		Obj: &models.Service{Id: id},
	}
	return do.Obj, do.get()
}

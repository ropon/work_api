package logics

import (
	"fmt"
)

type Obj interface {
	Get() (err error)
	Delete() (err error)
	Update() (err error)
	Patch(v interface{}) (err error)
	FormatTime()
}

type DbObj struct {
	Id uint
	Obj
}

func (do *DbObj) get() error {
	err := do.Obj.Get()
	if err != nil {
		return fmt.Errorf("ID:%d对应记录不存在", do.Id)
	}
	do.Obj.FormatTime()
	return nil
}

func (do *DbObj) delete() error {
	if err := do.get(); err != nil {
		return err
	}
	return do.Delete()
}

func (do *DbObj) update(v Obj) error {
	if err := do.get(); err != nil {
		return err
	}
	err := v.Update()
	if err != nil {
		return err
	}
	v.FormatTime()
	return nil
}

func (do *DbObj) patch(v Obj) error {
	if err := do.get(); err != nil {
		return err
	}
	err := do.Patch(v)
	if err != nil {
		return err
	}
	v.FormatTime()
	return nil
}

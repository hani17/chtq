package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var (
	x *xorm.Engine

	tables []interface{}
)

func init() {
	tables = append(tables, new(User))
}

func getEngine() (*xorm.Engine, error) {
	datasourcename := "root:test@/chtq"
	return xorm.NewEngine("mysql", datasourcename)
}

func setEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Fail to connect to database %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	//x.ShowSQL(true)
	return nil
}

func NewEngine() (err error) {
	if err = setEngine(); err != nil {
		return err
	}
	x.Sync(tables...)
	return nil
}

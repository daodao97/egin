package db

import (
	"database/sql"

	"github.com/daodao97/egin/utils/config"
)

type ModelConf struct {
	Driver     string
	Table      string
	FakeDelete bool
	PrimaryKey string
	Connection string
	FakeDelKey string
}

type Model interface {
	FindById(id int, selectFields []string, binding interface{}) error
	DeleteById(id int) error
	UpdateById(id int, record Record) error
	SelectCount(filter Filter) (int, error)
	Select(filter Filter, attr Attr, binding interface{}) (err error)
	Insert(record Record) (lastId int64, affected int64, err error)
	Update(filter Filter, record Record) (lastId int64, affected int64, err error)
	Delete(filter Filter) (lastId int64, affected int64, err error)
	DB() *sql.DB
	BeforeSelect(function func(m Model, filter Filter, attr Attr))
	AfterSelect(function func(m Model, result []map[string]interface{}, err error) []map[string]interface{})
	BeforeInsert(function func(m Model, record Record) Record)
	AfterInsert(function func(m Model, lastId int64, record Record, err error))
	BeforeUpdate(function func(m Model, filter Filter, record Record) (Filter, Record))
	AfterUpdate(function func(m Model, lastId int64, record Record, err error))
	BeforeDelete(function func(m Model, filter Filter))
	AfterDelete(function func(m Model, lastId int64, affected int64, err error))
	GetTableName() string
	GetConnectionName() string
	GetConf() config.Database
	GetLastSql() string
}

type DateColumn struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DelColumn struct {
	IsDeleted int `json:"is_deleted"`
}

const (
	StatusOff = iota
	StatusOn
)

const (
	StatusTrue  = true
	StatusFalse = false
)

const FakeDelKey = "is_deleted"

type StatusColumn struct {
	Status int `json:"status" comment:"状态 0 禁用, 1 启用"`
}

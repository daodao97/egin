package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 从池子中捞出一个db对象, 注意, 这里并非 mysql 连接池
func getDBInPool(key string) (*sql.DB, bool) {
	val, ok := pool.Load(key)
	return val.(*sql.DB), ok
}

func NewModel(conf ModelConf) Model {
	if conf.Driver == "" {
		conf.Driver = "mysql"
	}
	if conf.Connection == "" {
		conf.Connection = "default"
	}
	if conf.FakeDelKey == "" {
		conf.FakeDelKey = "is_deleted"
	}
	if conf.PrimaryKey == "" {
		conf.PrimaryKey = "id"
	}
	currentDb, ok := dbConf[conf.Connection]
	if !ok {
		logger.Error(fmt.Sprintf("database %s not found", conf.Connection))
	}

	db, ok := getDBInPool(conf.Connection)
	if !ok {
		panic("not found db conn")
	}
	return &baseModel{
		conf:       currentDb,
		Driver:     conf.Driver,
		Table:      conf.Table,
		Connection: conf.Connection,
		FakeDelete: conf.FakeDelete,
		FakeDelKey: conf.FakeDelKey,
		PrimaryKey: conf.PrimaryKey,
		db:         db,
	}
}

// TODO 主从库的支持
// Model 的基础封装
type baseModel struct {
	conf         Database
	Driver       string
	Table        string
	PrimaryKey   string
	Connection   string
	db           *sql.DB
	Entity       interface{}
	LastSql      string
	FakeDelete   bool
	FakeDelKey   string
	ActionType   string
	Args         interface{}
	beforeSelect []func(m Model, filter Filter, attr Attr)
	afterSelect  []func(m Model, result []map[string]interface{}, err error) []map[string]interface{}
	beforeUpdate []func(m Model, filter Filter, record Record) (Filter, Record)
	afterUpdate  []func(m Model, lastId int64, record Record, err error)
	beforeInsert []func(m Model, record Record) Record
	afterInsert  []func(m Model, lastId int64, record Record, err error)
	beforeDelete []func(m Model, filter Filter)
	afterDelete  []func(m Model, lastId int64, affected int64, err error)
}

func (m *baseModel) FindById(id int, selectFields []string, binding interface{}) error {
	var b []interface{}
	err := m.Select(Filter{"id": id}, Attr{Select: selectFields}, &b)
	if err != nil {
		return err
	}
	if len(b) != 1 {
		return errors.New("not found")
	}
	r, _ := json.Marshal(b[0])
	err = json.Unmarshal(r, &binding)
	if err != nil {
		return err
	}
	return nil
}

func (m *baseModel) UpdateById(id int, record Record) error {
	_, _, err := m.Update(Filter{"id": id}, record)
	return err
}

func (m *baseModel) DeleteById(id int) error {
	_, _, err := m.Delete(Filter{"id": id})
	return err
}

func (m *baseModel) SelectCount(filter Filter) (int, error) {
	defer timeCost()(m)

	sqlWhere, args1 := filterToQuery(filter)

	_sql := strings.Trim(fmt.Sprintf("select count(*) as count from %s %s", m.Table, sqlWhere), " ")
	m.LastSql = _sql
	m.ActionType = "read"

	result, err := Query(m.db, _sql, args1...)

	if err != nil {
		return 0, err
	}
	count, ok := result[0]["count"]
	if !ok {
		return 0, errors.New("select count error")
	}
	countInt, err := strconv.Atoi(*count.(*string))
	if err != nil {
		return 0, err
	}

	return countInt, nil
}

// 数据库查询
// filter {sex:0, class:{in:[1,2]}}
// attr {Select:[id,name], OrderBy:"id desc"}
// sql: select `id`, `name` from ${table} where `sex` = ? and class in (?, ?)
// result: 查询结果, 错误信息
func (m *baseModel) Select(filter Filter, attr Attr, binding interface{}) (err error) {
	defer timeCost()(m)

	for _, beforeHook := range m.beforeSelect {
		beforeHook(m, filter, attr)
	}

	var args []interface{}
	if m.FakeDelete {
		filter[m.FakeDelKey] = 0
	}

	sqlWhere, args1 := filterToQuery(filter)
	sqlField := attrToSelectQuery(attr)
	sqlAttr, args2 := attrToQuery(attr)

	_sql := strings.Trim(fmt.Sprintf("select %s from %s %s %s", sqlField, m.Table, sqlWhere, sqlAttr), " ")
	args1 = append(args1, args2...)
	args = append(args, args1...)
	m.LastSql = _sql
	m.ActionType = "read"
	m.Args = args

	result, err := Query(m.db, _sql, args...)

	if err != nil {
		return err
	}

	for _, afterHook := range m.afterSelect {
		result = afterHook(m, result, err)
	}

	t, err := json.Marshal(result)
	if err != nil {
		return err
	}

	if string(t) == "null" {
		t = []byte("[]")
	}

	err = json.Unmarshal(t, binding)
	if err != nil {
		return err
	}

	return nil
}

// Record: {name:"Joke", sex:1}
// sql: insert into ${table} (`name`, `sex`) values (?, ?)
// return: 最后的记录id, 受影响行数, 错误信息
func (m *baseModel) Insert(record Record) (lastId int64, affected int64, err error) {
	defer timeCost()(m)
	delete(record, m.PrimaryKey)

	for _, beforeHook := range m.beforeInsert {
		record = beforeHook(m, record)
	}

	sqlInsert, args := insertRecordToQuery(record)
	_sql := fmt.Sprintf("insert into %s %s", m.Table, sqlInsert)
	m.LastSql = _sql
	m.ActionType = "create"
	m.Args = record

	lastId, affected, err = exec(m.db, _sql, args...)

	for _, afterHook := range m.afterInsert {
		afterHook(m, lastId, record, err)
	}

	return lastId, affected, err
}

// filter: {id: 1}
// Record: {name:"Joke", sex:1}
// sql: update ${table} set `name` = ?, `sex` = ? where `id` = ?
// return: 最后的记录id, 受影响行数, 错误信息
func (m *baseModel) Update(filter Filter, record Record) (lastId int64, affected int64, err error) {
	defer timeCost()(m)
	delete(record, m.PrimaryKey)

	for _, beforeHook := range m.beforeUpdate {
		filter, record = beforeHook(m, filter, record)
	}

	sqlUpdate, argsUpdate := updateRecordToQuery(record)
	sqlWhere, argsWhere := filterToQuery(filter)
	_sql := fmt.Sprintf("update `%s` set %s %s", m.Table, sqlUpdate, sqlWhere)
	args := append(argsUpdate, argsWhere...)
	m.LastSql = _sql
	m.ActionType = "update"
	m.Args = record

	lastId, affected, err = exec(m.db, _sql, args...)

	if pk, ok := filter[m.PrimaryKey]; ok {
		lastId = int64(pk.(int))
	}

	for _, afterHook := range m.afterUpdate {
		afterHook(m, lastId, record, err)
	}

	return lastId, affected, err
}

// filter: {id: 1}
// Record: {name:"Joke", sex:1}
// 物理删除sql: delete from ${table} where `id` = ?
// 伪删除sql: update ${table} set ${FakeDelKey} = 1 where `id` = ?
// return: 最后的记录id, 受影响行数, 错误信息
func (m *baseModel) Delete(filter Filter) (lastId int64, affected int64, err error) {
	defer timeCost()(m)

	for _, beforeHook := range m.beforeDelete {
		beforeHook(m, filter)
	}

	if m.FakeDelete {
		return m.Update(filter, Record{m.FakeDelKey: 1})
	}

	sqlWhere, args := filterToQuery(filter)

	_sql := fmt.Sprintf("delete from %s %s", m.Table, sqlWhere)
	m.LastSql = _sql
	m.ActionType = "delete"

	lastId, affected, err = exec(m.db, _sql, args...)

	for _, afterHook := range m.afterDelete {
		afterHook(m, lastId, affected, err)
	}

	return lastId, affected, err
}

// 获取 db原生对象, 可以执行原生sql语句等更多操作
func (m *baseModel) DB() *sql.DB {
	return m.db
}

func (m *baseModel) BeforeSelect(function func(m Model, filter Filter, attr Attr)) {
	m.beforeSelect = append(m.beforeSelect, function)
}

func (m *baseModel) AfterSelect(function func(m Model, result []map[string]interface{}, err error) []map[string]interface{}) {
	m.afterSelect = append(m.afterSelect, function)
}

func (m *baseModel) BeforeInsert(function func(m Model, record Record) Record) {
	m.beforeInsert = append(m.beforeInsert, function)
}

func (m *baseModel) AfterInsert(function func(m Model, lastId int64, record Record, err error)) {
	m.afterInsert = append(m.afterInsert, function)
}

func (m *baseModel) BeforeUpdate(function func(m Model, filter Filter, record Record) (Filter, Record)) {
	m.beforeUpdate = append(m.beforeUpdate, function)
}

func (m *baseModel) AfterUpdate(function func(m Model, lastId int64, record Record, err error)) {
	m.afterUpdate = append(m.afterUpdate, function)
}

func (m *baseModel) BeforeDelete(function func(m Model, filter Filter)) {
	m.beforeDelete = append(m.beforeDelete, function)
}

func (m *baseModel) AfterDelete(function func(m Model, lastId int64, affected int64, err error)) {
	m.afterDelete = append(m.afterDelete, function)
}

func (m *baseModel) GetTableName() string {
	return m.Table
}

func (m *baseModel) GetConnectionName() string {
	return m.Connection
}

func (m *baseModel) GetConf() Database {
	return m.conf
}

func (m *baseModel) GetLastSql() string {
	return m.LastSql
}

// 记录 sql, 耗时 等信息
func timeCost() func(m *baseModel) {
	start := time.Now()
	return func(m *baseModel) {
		tc := time.Since(start)
		logger.Info(
			fmt.Sprintf("use time %v", tc),
			map[string]interface{}{
				"table":      m.Table,
				"connection": m.Connection,
				"sql":        m.LastSql,
				"type":       m.ActionType,
				"ums":        fmt.Sprintf("%v", tc),
				"args":       m.Args,
			})
	}
}

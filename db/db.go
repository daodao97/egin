package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Host     string
	Port     int
	User     string
	Passwd   string
	Database string
	Driver   string
	Options  map[string]string
	Pool     struct {
		MaxOpenConns int
		MaxIdleConns int
	}
}

type Databases map[string]Database

type Logger interface {
	Info(message interface{}, content ...interface{})
	Error(message interface{}, content ...interface{})
}

var (
	pool   sync.Map
	dbConf Databases
	logger Logger
	once   sync.Once
)

func Init(_dbConf Databases, _logger Logger) {
	once.Do(func() {
		dbConf = _dbConf
		for key, conf := range dbConf {
			db := makeDb(conf)
			pool.Store(key, db)
		}
		logger = _logger
		logger.Info("db init")
		go clearStmt()
	})
}

// 生成 原生 DB 对象
func makeDb(conf Database) *sql.DB {
	server := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Passwd, server, conf.Database)
	driver := conf.Driver
	if driver == "" {
		driver = "mysql"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(fmt.Sprintf("failed Connection database: %s", err))
	}

	// 设置数据库连接池最大连接数
	var MaxOpenConns int
	if conf.Pool.MaxOpenConns == 0 {
		MaxOpenConns = 100
	} else {
		MaxOpenConns = conf.Pool.MaxOpenConns
	}
	db.SetMaxOpenConns(MaxOpenConns)

	// 连接池最大允许的空闲连接数
	// 如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	var MaxIdleConns int
	if conf.Pool.MaxIdleConns == 0 {
		MaxIdleConns = 20
	} else {
		MaxIdleConns = conf.Pool.MaxIdleConns
	}
	db.SetMaxIdleConns(MaxIdleConns)
	return db
}

// 一般用Prepared Statements和Exec()完成INSERT, UPDATE, DELETE操作
func exec(db *sql.DB, _sql string, args ...interface{}) (int64, int64, error) {

	tx, err := db.Begin()
	if err != nil {
		logger.Error(err)
		return 0, 0, err
	}
	var flag bool
	addr := &flag
	defer func(flag *bool, errMsg *error) {
		if *flag {
			return
		}
		err := tx.Rollback()
		if err != nil {
			logger.Error(fmt.Sprintf("db.exec.rollback fail: %s", err), map[string]interface{}{
				"sql":  _sql,
				"args": args,
				"msg":  errMsg,
			})
		} else {
			logger.Info("db.exec.rollback", map[string]interface{}{
				"sql":  _sql,
				"args": args,
				"msg":  errMsg,
			})
		}
	}(addr, &err)

	stmt, err := tx.Prepare(_sql)
	if err != nil {
		logger.Error("db.exec.prepare fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return 0, 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		logger.Error("db.exec.exec fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("db.exec.commit fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return 0, 0, err
	}
	flag = true

	lastId, err := res.LastInsertId()
	if err != nil {
		logger.Error("db.exec.lastId fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return 0, 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logger.Error("db.exec.affected fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return 0, 0, err
	}

	return lastId, affected, nil
}

func Query(db *sql.DB, _sql string, args ...interface{}) (result []map[string]interface{}, err error) {
	stmt, err := makeStmt(db, _sql)
	if err != nil {
		logger.Error("db.query.prepare fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return result, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		logger.Error("db.query.query fail", map[string]interface{}{
			"sql":  _sql,
			"args": args,
			"msg":  err,
		})
		return result, err
	}

	return rows2SliceMap(rows)
}

// 将 sql.Rows 的结果转换为 map
// 注意这里所有的value 均为string
// TODO 是否可以根据 rows.ColumnTypes() databasesType 做类型转换?
func rows2SliceMap(rows *sql.Rows) (list []map[string]interface{}, err error) {
	// 字段名称
	columns, _ := rows.Columns()
	// 多少个字段
	length := len(columns)
	for rows.Next() {
		var dest []interface{}

		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			continue
		}
		for _, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "CHAR":
				dest = append(dest, new(string))
			case "INT", "TINYINT":
				dest = append(dest, new(int))
			case "TEXT":
				dest = append(dest, new(sql.NullString))
			default:
				dest = append(dest, new(string))
			}
		}
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("scan error", err)
			continue
		}
		if dest == nil {
			continue
		}
		// 每一行
		row := make(map[string]interface{})
		for i := 0; i < length; i++ {
			if val, ok := dest[i].(*sql.NullString); ok {
				row[columns[i]] = val.String
			} else {
				row[columns[i]] = dest[i]
			}
		}
		list = append(list, row)
	}
	if err := rows.Err(); err != nil {
		return list, err
	}
	return list, nil
}

var statements sync.Map

type stmtData struct {
	stmt     *sql.Stmt
	createAt time.Time
}

func makeStmt(db *sql.DB, sqlStr string) (stmt *sql.Stmt, err error) {
	value, exist := statements.Load(sqlStr)
	if exist {
		return value.(stmtData).stmt, nil
	}

	stmt, err = db.Prepare(sqlStr)

	if err != nil {
		return stmt, err
	}

	statements.Store(sqlStr, stmtData{
		stmt:     stmt,
		createAt: time.Now(),
	})

	return stmt, nil
}

func clearStmt() {
	for range time.Tick(time.Second) {
		statements.Range(func(key, value interface{}) bool {
			if time.Now().Unix()-value.(stmtData).createAt.Unix() > 60 {
				statements.Delete(key)
			}
			return true
		})
	}
}

type db struct {
	db     *sql.DB
	result interface{}
}

func (d *db) Conn(connection string) *db {
	db, ok := getDBInPool(connection)
	if ok {
		d.db = db
	}
	return d
}

func (d *db) Query(sql string, args ...interface{}) *db {
	result, err := Query(d.db, sql, args...)
	if err != nil {
		panic(err)
	}
	d.result = result
	return d
}

func (d *db) Bind(bind interface{}) error {
	t, err := json.Marshal(d.result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(t, bind)
	if err != nil {
		return err
	}

	return nil
}

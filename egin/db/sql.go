package db

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/daodao97/egin/lib"
)

// sql 片段的构造

const LogicStrAnd = "and"
const LogicStrOr = "or"

const LogicIntAnd = 1
const LogicIntOr = 0

type Filter map[string]interface{}

type Attr struct {
	Select  []string
	Limit   int
	Offset  int
	OrderBy string
	GroupBy string
	Having  string
}

type Record map[string]interface{}

// Record {"field":"val", "field2":"val2"}
// return arg1: (`field`, `field2`) values (?, ?)
// return arg2: ["val", "val2"]
func insertRecordToQuery(record Record) (string, []interface{}) {
	var fields []string
	var value []string
	var args []interface{}
	for k, v := range record {
		fields = append(fields, k)
		value = append(value, "?")
		args = append(args, v)
	}

	_sql := fmt.Sprintf(
		"(%s) values (%s)",
		fmt.Sprintf("`%s`", strings.Join(fields, "`,`")),
		fmt.Sprintf("%s", strings.Join(value, ",")),
	)

	return _sql, args
}

// Record {"field":"val", "field2":"val2"}
// return arg1: set `field` = ?, `field2` = ?
// return arg2: ["val", "val2"]
func updateRecordToQuery(record Record) (string, []interface{}) {
	var scopes []string
	var args []interface{}
	for k, v := range record {
		scopes = append(scopes, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	return strings.Join(scopes, ", "), args
}

// field {"field":"val", "field2":{in:[1,3]}}
// return arg1: where `field` = ? and `field2` in (?, ?)
// return arg2: ["val", "val2"]
func filterToQuery(filter Filter) (string, []interface{}) {
	var sql string
	var args []interface{}
	if len(filter) == 0 {
		return sql, args
	}
	logic := filter["__logic"]
	if logic == nil {
		logic = LogicStrAnd
	}
	var scopes []string

	for k, v := range filter {
		switch v.(type) {
		case []int:
			pl := []string{}
			for _, val := range v.([]int) {
				pl = append(pl, strconv.Itoa(val))
			}
			scopes = append(scopes, fmt.Sprintf("`%s` in (%s)", k, strings.Join(pl, ",")))
		case []string:
			scopes = append(scopes, fmt.Sprintf("`%s` in (%s)", k, strings.Join(v.([]string), ",")))
		case map[string]string:
			var _scope []string
			var _logic = LogicStrAnd
			for op, val := range v.(map[string]string) {
				_, found := lib.Find([]string{">", ">=", "<", "<=", "=", "<>", "!=", "like"}, op)
				if found {
					_scope = append(_scope, fmt.Sprintf("`%s` %s ?", k, op))
				}
				if op == "find_in_set" {
					_scope = append(_scope, fmt.Sprintf("find_in_set(`%s`, ?)", k))
				}
				if op == "__logic" {
					_logic = val
				}
				args = append(args, val)
			}
			scopes = append(scopes, fmt.Sprintf("(%s)", strings.Join(_scope, fmt.Sprintf(" %s ", _logic))))
		case map[string]int:
			var _scope []string
			var _logic = LogicIntAnd
			for op, val := range v.(map[string]int) {
				_, found := lib.Find([]string{">", ">=", "<", "<=", "=", "<>", "!=", "like"}, op)
				if found {
					_scope = append(_scope, fmt.Sprintf("`%s` %s ?", k, op))
				}
				if op == "__logic" {
					_logic = val
				}
				args = append(args, val)
			}
			var logicStr string
			if _logic == LogicIntAnd {
				logicStr = LogicStrAnd
			} else if _logic == LogicIntOr {
				logicStr = LogicStrOr
			}
			scopes = append(scopes, fmt.Sprintf("(%s)", strings.Join(_scope, fmt.Sprintf(" %s ", logicStr))))
		default:
			args = append(args, v)
			scopes = append(scopes, fmt.Sprintf("`%s` = ?", k))
		}
	}
	sql = fmt.Sprintf("where %s", strings.Join(scopes, fmt.Sprintf(" %s ", logic.(string))))
	return sql, args
}

// attr: {OrderBy:"id desc",Limit: 1}
// return arg1: order by ? limit ?
// return arg2: ["id desc", 1]
func attrToQuery(attr Attr) (string, []interface{}) {
	var sql string
	var scopes []string
	var args []interface{}
	// https://stackoverflow.com/questions/30867337/golang-order-by-issue-with-mysql
	// 占位符 ? 只能用在 insert 或 条件中, order by/ group by 不支持
	valid := regexp.MustCompile("^[A-Za-z0-9_ `]+$")
	if attr.OrderBy != "" && valid.MatchString(attr.OrderBy) {
		scopes = append(scopes, fmt.Sprintf("order by %s", attr.OrderBy))
	}
	if attr.GroupBy != "" && valid.MatchString(attr.GroupBy) {
		scopes = append(scopes, fmt.Sprintf("group by %s", attr.GroupBy))
	}
	if attr.Limit != 0 {
		args = append(args, attr.Limit)
		scopes = append(scopes, "limit ?")
	}
	if attr.Offset != 0 {
		args = append(args, attr.Offset)
		scopes = append(scopes, "offset ?")
	}
	sql = strings.Join(scopes, " ")
	return sql, args
}

// attr: {Select:[id,name]}
// return arg: `id`, `name`
func attrToSelectQuery(attr Attr) string {
	if len(attr.Select) == 0 {
		return "*"
	}
	return fmt.Sprintf("`%s`", strings.Join(attr.Select, "`,`"))
}

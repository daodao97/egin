package gen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/daodao97/egin/db"
	"github.com/daodao97/egin/lib"

	"github.com/daodao97/egin-tools/parser"
)

func Gen(data map[string]interface{}, tpl string) (string, error) {
	var buf = bytes.NewBufferString("")

	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", errors.Wrapf(err, "template init err")
	}

	err = t.Execute(buf, data)
	if err != nil {
		return "", errors.Wrapf(err, "template data err")
	}

	// fmt.Println(string(buf.Bytes()))

	byteData, err := format.Source(buf.Bytes())
	if err != nil {
		return "", errors.Wrapf(err, "format err \n %s", string(buf.Bytes()))
	}
	return string(byteData), nil
}

func apiMethod(doc []string) (method string, path string, err error) {
	if doc == nil {
		return "", "", errors.New("not api")
	}
	reg := regexp.MustCompile("^@(Any|Get|Put|Post|Delete)Api\\s?([a-zA-Z0-9_/:*]+)?\\s?.*")
	matched := reg.FindStringSubmatch(doc[0])
	if len(matched) != 3 {
		return "", "", errors.New("not api")
	}
	return strings.ToUpper(matched[1]), matched[2], nil
}

func apiMiddleware(doc []string) []string {
	reg := regexp.MustCompile(`@Middleware .*`)
	for _, v := range doc {
		if matched := reg.MatchString(v); matched {
			return strings.Split(v, " ")[1:]
		}
	}
	return []string{}
}

func MakeRouteHandle(entity string, info parser.StructFunc) (code string, err error) {
	method, path, err := apiMethod(info.Doc)
	if err != nil {
		return "", errors.New("not api")
	}

	res := regexp.MustCompile(":([a-zA-Z0-9*])+")
	pathArgs := res.FindAllString(path, -1)
	for i, v := range pathArgs {
		pathArgs[i] = strings.TrimPrefix(v, ":")
	}

	method = strings.ToUpper(method)

	if method == "ANY" {
		method = "Any"
	}

	args := map[string]interface{}{
		"method":     method,
		"entity":     entity,
		"path":       path,
		"funcName":   info.Name,
		"moduleName": ModuleName(),
		"pathArgs":   pathArgs,
	}

	args["middleware"] = ""
	middleware := apiMiddleware(info.Doc)
	if len(middleware) > 0 {
		for i, v := range middleware {
			if !strings.HasSuffix(v, ")") {
				middleware[i] = fmt.Sprintf("middleware.%s()", v)
			} else {
				middleware[i] = fmt.Sprintf("middleware.%s", v)
			}
		}
		args["middleware"] = ", " + strings.Join(middleware, ",")
	}

	tpl := SimpleApi
	if len(info.Params) == 3 {
		args["paramsStruct"] = info.Params[len(info.Params)-1].Type
		tpl = ApiWithParam
	}
	if len(info.Params) == 2 {
		if info.Params[1].Type == "int" {
			tpl = SimpleApi
		} else {
			args["paramsStruct"] = info.Params[1].Type
			tpl = ApiWithParam
		}
	}
	if info.ResultCount == 0 {
		tpl = ApiReturnVoid
	}

	return Gen(args, tpl)
}

func MakeRouteFile(structInfo []parser.StructInfo, varsInfo []parser.VarInfo) {
	for _, v := range structInfo {
		entity := v.Name
		var handles []string
		for _, f := range v.Funcs {
			handle, err := MakeRouteHandle(v.Name, f)
			if err != nil {
				fmt.Println(err)
				continue
			}
			handles = append(handles, handle)
		}
		if len(handles) == 0 {
			continue
		}

		argsR := map[string]interface{}{
			"strconv":                true,
			"entity":                 entity,
			"handles":                handles,
			"hasCustomValidateFuncs": false,
			"moduleName":             ModuleName(),
		}

		var customValidateVarsName []string
		for _, v := range varsInfo {
			if v.Type == "utils.CustomValidateFunc" {
				argsR["hasCustomValidateFuncs"] = true
				customValidateVarsName = append(customValidateVarsName, v.Name)
			}
		}
		argsR["customValidateFuncs"] = customValidateVarsName

		tpl, err := Gen(argsR, RouteFile)
		if err != nil {
			fmt.Println("gen config/routes/***.go error", err)
			os.Exit(1)
		}
		_ = ioutil.WriteFile(fmt.Sprintf("config/routes/%s.go", lib.ToSnakeCase(entity)), []byte(tpl), os.FileMode(0644))
	}
}

func MakeRouteExport() {
	var funcs []parser.FuncInfo
	lib.RecursiveDir("config/routes", func(filePath string) {
		fmt.Println(filePath)
		funcInfo, err := parser.FileFunInfo(filePath)
		if err != nil {
			panic(err)
		}
		funcs = append(funcs, funcInfo...)
	})
	var list []string
	for _, v := range funcs {
		var s string
		if strings.HasSuffix(v.Name, "Router") {
			s = "routes." + v.Name + "(r)"
		} else {
			s = "routes." + v.Name + "()"
		}
		list = append(list, s)
	}
	args := map[string]interface{}{
		"list":       sort.StringSlice(list),
		"moduleName": ModuleName(),
	}
	tpl, err := Gen(args, RouteExport)
	if err != nil {
		fmt.Println("gen config/route.go error", err)
		os.Exit(1)
	}
	_ = ioutil.WriteFile("config/routes.go", []byte(tpl), os.FileMode(0644))
}

type TableField struct {
	Name    string
	Field   string
	Type    string
	Comment string
}

func getType(str string) string {
	switch str {
	case "varchar", "char", "tinytext", "datetime", "text", "longtext", "timestamp":
		return "string"
	case "tinyint":
		return "int"
	default:
		return str
	}
}

func MakeModel(connection string, databases string, table TableInfo) {
	mysqlDb, ok := db.GetDBInPool(connection)
	if !ok {
		fmt.Println("get pool error")
		return
	}
	rows, err := mysqlDb.Query("select `COLUMN_NAME`, `DATA_TYPE`, `COLUMN_COMMENT` from information_schema.COLUMNS where `TABLE_SCHEMA` = ? and `TABLE_NAME` = ? order by ORDINAL_POSITION", databases, table.Name)
	if err != nil {
		fmt.Println("table schema fail ", err)
		return
	}

	var fieldList []interface{}

	var fakeDel bool

	for rows.Next() {
		dest := []interface{}{
			new(string),
			new(string),
			new(string),
		}
		rows.Scan(dest...)

		row := TableField{
			Name:    lib.ToCamelCase(*dest[0].(*string)),
			Field:   *dest[0].(*string),
			Type:    getType(*dest[1].(*string)),
			Comment: getType(*dest[2].(*string)),
		}

		fieldList = append(fieldList, row)

		if row.Field == "is_deleted" {
			fakeDel = true
		}
	}

	args := map[string]interface{}{
		"entityName":   lib.ToCamelCase(table.Name),
		"table":        table.Name,
		"connection":   connection,
		"fieldList":    fieldList,
		"backquote":    "`",
		"fakeDel":      fakeDel,
		"tableComment": table.Comment,
		"moduleName":   ModuleName(),
	}

	tpl, err := Gen(args, Entity)
	if err != nil {
		fmt.Println("gen model/***.go error", err)
		os.Exit(1)
	}
	_ = ioutil.WriteFile(fmt.Sprintf("model/%s.go", table.Name), []byte(tpl), os.FileMode(0644))
}

type TableInfo struct {
	Name    string
	Comment string
}

func GetDbAllTable(connection string, databases string) (result []TableInfo) {
	mysqlDb, ok := db.GetDBInPool(connection)
	if !ok {
		fmt.Println("get pool error")
		return
	}
	rows, err := mysqlDb.Query("select table_name,table_comment from information_schema.tables where table_schema=? and table_type='base table';", databases)
	if err != nil {
		fmt.Println("table schema fail ", err)
		return
	}

	for rows.Next() {
		dest := []interface{}{
			new(string),
			new(string),
		}
		rows.Scan(dest...)
		result = append(result, TableInfo{
			Name:    *dest[0].(*string),
			Comment: *dest[1].(*string),
		})
	}

	return result
}

func GetTableInfo(connection string, databases string, table string) (result TableInfo) {
	mysqlDb, ok := db.GetDBInPool(connection)
	if !ok {
		fmt.Println("get pool error")
		return
	}
	rows, err := mysqlDb.Query("select table_name,table_comment from information_schema.tables where table_schema=? and table_type='base table' and table_name = ?;", databases, table)
	if err != nil {
		fmt.Println("table schema fail ", err)
		return
	}

	for rows.Next() {
		dest := []interface{}{
			new(string),
			new(string),
		}
		rows.Scan(dest...)
		result.Name = *dest[0].(*string)
		result.Comment = *dest[1].(*string)
	}

	return result
}

func ModuleName() string {
	fi, err := os.Open("go.mod")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	a, _, c := br.ReadLine()
	if c == io.EOF {
		return ""
	}

	return strings.Replace(string(a), "module ", "", 1)
}

func MakeController(tableName string, desc string) {
	table := lib.ToCamelCase(tableName)
	args := map[string]interface{}{
		"table":      table,
		"tableName":  tableName,
		"EntityDesc": desc,
	}
	tpl, err := Gen(args, ctrlTpl)
	if err != nil {
		fmt.Println("gen controller/***.go error", err)
		os.Exit(1)
	}
	_ = ioutil.WriteFile(fmt.Sprintf("controller/%s.go", tableName), []byte(tpl), os.FileMode(0644))

}

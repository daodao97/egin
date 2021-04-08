package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// 反射调用结构体方法
func Invoke(any interface{}, name string, args ...interface{}) ([]reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	var _result []reflect.Value
	notExist := method == reflect.Value{}
	if notExist {
		return _result, fmt.Errorf("Method %s not found ", name)
	}
	methodType := method.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return _result, fmt.Errorf("Method %s must have minimum %d params. Have %d ", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return _result, fmt.Errorf("Method %s must have %d params. Have %d ", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return _result, fmt.Errorf("Method %s. Param[%d] must be %s. Have %s ", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return _result, fmt.Errorf("Method %s. Param[%d] must be %s. Have %s ", name, i, inType, argType)
		}
	}
	return method.Call(in), nil
}

// 获取结构体中字段的名称
func GetStructFieldsName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

// 获取结构体中Tag的值，如果没有tag则返回字段值
func GetStructTags(structName interface{}) map[string]map[string]string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make(map[string]map[string]string)
	for i := 0; i < fieldNum; i++ {
		fieldName := t.Field(i).Name
		tagStr := string(t.Field(i).Tag)
		if tagStr != "" {
			tokens := strings.Split(tagStr, " ")
			part := make(map[string]string)
			for i := range tokens {
				tagInfo := strings.Split(strings.Replace(tokens[i], "\"", "", -1), ":")
				if len(tagInfo) > 1 {
					part[tagInfo[0]] = tagInfo[1]
				}
			}
			result[fieldName] = part
		}
	}
	return result
}

func WaitGo(number int, fu func()) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(index int) {
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func UpdateStructByTagMap(result interface{}, tagName string, tagMap map[string]interface{}) error {
	t := reflect.TypeOf(result)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("result have to be a pointer")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("result pointer not struct")
	}
	v := reflect.ValueOf(result).Elem()
	fieldNum := v.NumField()
	for i := 0; i < fieldNum; i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		f := v.FieldByName(fieldInfo.Name)
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		value, ok := tagMap[tag]
		if !ok {
			continue
		}

		valueRealType := reflect.TypeOf(value).Kind()
		targetType := f.Kind()
		if valueRealType == targetType {
			f.Set(reflect.ValueOf(value))
			continue
		}
		expr := fmt.Sprintf("%s-to-%s", valueRealType, targetType)
		switch expr {
		case "int-to-string":
			f.SetString(strconv.Itoa(value.(int)))
		case "string-to-string":
			f.SetString(value.(string))
		case "string-to-int":
			_v, _ := strconv.Atoi(value.(string))
			f.SetInt(int64(_v))
		case "float64-to-int":
			f.SetInt(int64(value.(float64)))
			// TODO more case
		}
	}
	return nil
}

// 获取协程id, 十分低效, 生产环境不要使用
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// 根据标签名获取结构体中的值, 返回 field : tagVale 的map
func GetStructAllTag(obj interface{}, tagName string) (map[string]string, error) {
	s := reflect.TypeOf(obj)

	info := make(map[string]string)

	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	if s.Kind() != reflect.Struct {
		return info, errors.New("arg1 is not a struct")
	}

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		info[f.Name] = f.Tag.Get(tagName)
	}

	return info, nil
}

func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

var re = regexp.MustCompile("([_\\-])([a-zA-Z]+)")

func ToCamelCase(str string) string {
	camel := re.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)

	return camel
}

func RecursiveDir(pathname string, f func(filePath string)) {
	rd, err := ioutil.ReadDir(pathname)
	if err == nil {
		for _, fi := range rd {
			if fi.IsDir() {
				fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
				RecursiveDir(pathname+fi.Name()+"\\", f)
			} else {
				f(pathname + "/" + fi.Name())
			}
		}
	}
}

func RecursiveData(data interface{}, handler func(tree string, val interface{})) {
	switch concreteVal := data.(type) {
	case map[string]interface{}:
		RecursiveMap(concreteVal, "", handler)
	case []interface{}:
		RecursiveArray(concreteVal, "", handler)
	default:
		handler("", data)
	}
}

func RecursiveMap(aMap map[string]interface{}, preTree string, handler func(tree string, val interface{})) {
	for key, val := range aMap {
		tree := fmt.Sprintf("%s.%s", preTree, key)
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			RecursiveMap(val.(map[string]interface{}), tree, handler)
		case []interface{}:
			RecursiveArray(val.([]interface{}), tree, handler)
		default:
			handler(tree, concreteVal)
		}
	}
}

func RecursiveArray(anArray []interface{}, preTree string, handler func(tree string, val interface{})) {
	for i, val := range anArray {
		tree := fmt.Sprintf("%s.%d", preTree, i)
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			RecursiveMap(val.(map[string]interface{}), tree, handler)
		case []interface{}:
			RecursiveArray(val.([]interface{}), tree, handler)
		default:
			handler(tree, concreteVal)
		}
	}
}

func StructToMap(data interface{}) map[string]interface{} {
	var b map[string]interface{}
	tmp, _ := json.Marshal(data)
	_ = json.Unmarshal(tmp, &b)
	return b
}

func BindToStruct(data interface{}, bind interface{}) {
	tmp, _ := json.Marshal(data)
	_ = json.Unmarshal(tmp, bind)
}

func Ternary(who interface{}, right interface{}, wrong interface{}) interface{} {
	switch who.(type) {
	case string:
		if len(who.(string)) > 0 {
			return right
		} else {
			return wrong
		}
	case int:
		if who.(int) > 0 {
			return right
		} else {
			return wrong
		}
	case bool:
		if who.(bool) {
			return right
		} else {
			return wrong
		}
	}

	return nil
}

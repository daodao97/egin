package db

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

func getModel() Model {
	return NewModel(ModelConf{
		Connection: "default",
		Table:      "data_version",
	})
}

func getLastVersionData(tableName string, dbName string) (preData map[string]interface{}) {
	var pre []map[string]interface{}
	_ = getModel().Select(
		Filter{"table_name": tableName, "db_name": dbName},
		Attr{Select: []string{"data"}, Limit: 1, OrderBy: "id desc"},
		&pre)
	if len(pre) != 1 {
		return preData
	}
	_ = json.Unmarshal([]byte(pre[0]["data"].(string)), &preData)
	return preData
}

func getDiffFields(preData map[string]interface{}, newData map[string]interface{}) (diffKey []string) {
	for k, v := range newData {
		pre, ok := preData[k]
		if !ok {
			diffKey = append(diffKey, k)
			continue
		}
		if v != pre {
			diffKey = append(diffKey, k)
			continue
		}
	}

	return diffKey
}

func DataVersion(m Model, lastId int64, record Record, err error) {
	versionModel := getModel()
	dataStr, _ := json.Marshal(record)
	preData := getLastVersionData(m.GetTableName(), m.GetConnectionName())
	_, _, _ = versionModel.Insert(Record{
		"version_id":    "",
		"db_name":       m.GetConf().Database,
		"table_name":    m.GetTableName(),
		"pk":            lastId,
		"data":          string(dataStr),
		"user_id":       0,
		"modify_fields": strings.Join(getDiffFields(preData, record), ","),
	})

	fmt.Println(runtime.Stack([]byte(""), true))
}

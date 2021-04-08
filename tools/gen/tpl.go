package gen

const SimpleApi = `
r.{{ .method }}("{{ .path }}", func(ctx *gin.Context) {
	{{- range $index, $value := .pathArgs -}} 
		{{ if eq $value "id"}}
		{{$value}}, _ := strconv.Atoi(ctx.Param("{{$value}}"))
		{{ else }}
		{{$value}} := ctx.Param("{{$value}}")
		{{ end }}
	{{- end -}}
	result, code, err := ctrl.{{ .funcName }}(ctx{{- range $index, $value := .pathArgs }}, {{$value}}{{- end }})
	egin.Response(ctx, result, code, err)
}{{ .middleware }})
`

const ApiWithParam = `
r.{{ .method }}("{{ .path }}", func () func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params controller.{{ .paramsStruct }}
		errs := utils.Validated(ctx, &params)
		if errs != nil {
			egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
			return
		}
		{{- range $index, $value := .pathArgs }} 
			{{ if eq $value "id"}}
			{{$value}}, _ := strconv.Atoi(ctx.Param("{{$value}}"))
			{{ else }}
			{{$value}} := ctx.Param("{{$value}}")
			{{ end }}
		{{- end }}
		result, code, err := ctrl.{{ .funcName }}(ctx{{- range $index, $value := .pathArgs }}, {{$value}}{{- end }}, params)
		egin.Response(ctx, result, code, err)
	}
}(){{ .middleware }})
`

const ApiReturnVoid = `
r.{{ .method }}("{{ .path }}", func(ctx *gin.Context) {
	ctrl.{{ .funcName }}(ctx)
})
`

const RouteFile = `
// ****************************
// 该文件为系统生成, 请勿更改
// ****************************
package routes

import (
	{{ if .strconv }}"strconv"{{ end }}
	"strings"

	"github.com/daodao97/egin"
	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/middleware"
	"github.com/daodao97/egin/utils"
	"github.com/gin-gonic/gin"

	"{{ .moduleName }}/controller"
)

func Reg{{ .entity }}Router(r *gin.Engine) {
	ctrl := controller.{{ .entity }}{}
	{{ range $index, $value := .handles }} {{$value}} {{ end }}
}

{{ if .hasCustomValidateFuncs }}
func RegUserCustomValidateFunc() {
	utils.RegCustomValidateFuncs([]utils.CustomValidateFunc{
		{{ range $index, $value := .customValidateFuncs }} controller.{{$value}}, {{ end }}
	})
}
{{ end }}
`

const RouteExport = `
// ****************************
// 该文件为系统生成, 请勿更改
// ****************************
package config

import (
	"github.com/gin-gonic/gin"

	"{{ .moduleName }}/config/routes"
)

func RegRouter(r *gin.Engine) {
	{{ range $index, $value := .list }}
		{{ $value }}
	{{ end }}
}
`

const Entity = `
package model

import (
	"github.com/daodao97/egin/db"
)

{{ if .tableComment }} // {{ .tableComment }} {{ end }}
type {{ .entityName }}Entity struct {
	{{- range $index, $value := .fieldList }}
		{{ $value.Name }} {{ $value.Type }} {{ $.backquote }}json:"{{ $value.Field }}"{{ if $value.Comment }} comment:"{{ $value.Comment }}"{{ end }}{{ $.backquote }}
	{{- end }}
}

type {{ .entityName }}Model struct {
	db.BaseModel
}

var {{ .entityName }} {{ .entityName }}Model

func init() {
	{{ .entityName }} = *New{{ .entityName }}Model()
}

func New{{ .entityName }}Model() *{{ .entityName }}Model {
	return &{{ .entityName }}Model{
		BaseModel: db.BaseModel{
			Connection: "{{ .connection }}",
			Table:      "{{ .table }}",
			{{- if .fakeDel }}
			FakeDelete: true,
			FakeDelKey: "is_deleted",
			{{- end }}
		},
	}
}
`

const ctrlTpl = `
package controller

import (
	"encoding/json"

	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/db"
	"github.com/gin-gonic/gin"

	"oms/model"
)

// @Controller {{ .EntityDesc }} 
type {{ .table }} struct {
}

type {{ .table }}Filter struct {
	Id int
}

// @GetApi /{{ .tableName }}
// @Summary 列表接口
// @Desc 列表接口 维护者: 刀刀
// @Params {{ .table }}Filter
// @Response
// @Middleware IpLimiter
func (u {{ .table }}) Get(c *gin.Context, params {{ .table }}Filter) (interface{}, consts.ErrCode, error) {
	filter := db.Filter{}
	if params.Id != 0 {
		filter["id"] = params.Id
	}
	var result []model.{{ .table }}Entity
	err := model.{{ .table }}.Get(filter, db.Attr{
		Select:  []string{"id"},
		OrderBy: "id desc",
	}, &result)
	return result, 0, err
}

type {{ .table }}Form struct {
}

// @PostApi /{{ .tableName }}
// @Summary 创建
// @Desc 创建接口 维护者: 刀刀
// @Params {{ .table }}Form 接口参数所对应的结构体
// @Response
func (u {{ .table }}) Post(c *gin.Context, params {{ .table }}Form) (interface{}, consts.ErrCode, error) {
	var paramsMap db.Record
	tmp, _ := json.Marshal(params)
	err := json.Unmarshal(tmp, &paramsMap)
	lastId, _, err := model.{{ .table }}.Insert(paramsMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	var result struct {
		LastId int64 
	}
	result.LastId = lastId
	return result, code, err
}

// @PutApi /{{ .tableName }}/:id
// @Summary 更新接口
// @Desc 维护者: 刀刀
// @Params {{ .table }}Form 接口参数所对应的结构体
// @Response
func (u {{ .table }}) Put(c *gin.Context, id int, params {{ .table }}Form) (interface{}, consts.ErrCode, error) {
	var paramsMap db.Record
	tmp, _ := json.Marshal(params)
	err := json.Unmarshal(tmp, &paramsMap)
	_, affected, err := model.{{ .table }}.Update(db.Filter{"id": id}, paramsMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	return affected, code, err
}

// @DeleteApi /{{ .tableName }}/:id
// @Summary 删除
// @Desc 维护者: 刀刀
// @Response
func (u {{ .table }}) Delete(c *gin.Context, id int) (interface{}, consts.ErrCode, error) {
	_, affected, err := model.{{ .table }}.Delete(db.Filter{
		"id": id,
	})
	return affected, 0, err
}
`

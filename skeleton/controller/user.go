package controller

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/daodao97/egin/egin/cache/redis"
	"github.com/daodao97/egin/egin/consts"
	"github.com/daodao97/egin/egin/db"
	"github.com/daodao97/egin/egin/utils"

	"skeleton/model"
)

type BaseApi struct {
	name string
}

// @Controller 用户管理 这里是简介
type User struct {
	BaseApi
}

type ParamsValidate struct {
	CheckIn  time.Time `form:"check_in" time_format:"2006-01-02" json:"check_in" binding:"required,bookabledate" label:"输入时间" in:"query"`
	CheckOut time.Time `form:"check_out" json:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02" label:"输出时间" in:"query"`
}

var Bookabledate = utils.CustomValidateFunc{
	Handle: func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			today := time.Now()
			if today.After(date) {
				return false
			}
		}
		return true
	},
	TagName: "bookabledate",
	Message: "{0}不能早于当前时间或{1}格式错误!",
}

// @GetApi /user 若无默认为全小写的方法名
// @Summary 用户列表接口
// @Desc 接口简介, 若无则为空 维护者: 刀刀
// @Params ParamsValidate 接口参数所对应的结构体
// @Response
// @Middleware IpLimiter
func (u User) Get(c *gin.Context, params ParamsValidate) (interface{}, consts.ErrCode, error) {

	var result []model.UserEntity
	user := model.User
	err := user.Get(db.Filter{
		"id": map[string]int{
			">": 20,
		},
	}, db.Attr{
		Select:  []string{"realname", "id", "username", "password"},
		OrderBy: "id desc",
	}, &result)

	redis := redis.NewDefault()
	setV, _ := json.Marshal([]int{1, 2, 4})
	err = redis.Set("egin:test", setV, 0)
	_cache, err := redis.Get("egin:test")

	return []interface{}{result, params, _cache}, 0, err
}

// @PostApi /user 若无默认为全小写的方法名
// @Summary 创建用户
// @Desc 这个接口支持你创建一个用户 维护者: 刀刀
// @Params ParamsValidate 接口参数所对应的结构体
// @Response
func (u User) Post(c *gin.Context) (interface{}, consts.ErrCode, error) {
	user := model.User
	result, _, err := user.Insert(db.Record{
		"username": "test33333",
		"realname": "你好",
		"password": "cool",
	})
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	return []interface{}{result}, code, err
}

// @PutApi /user
// @Summary 更新用户信息
// @Desc 维护者: 刀刀
// @Params ParamsValidate 接口参数所对应的结构体
// @Response
func (u User) Put(c *gin.Context) (interface{}, consts.ErrCode, error) {
	user := model.User
	_, affected, err := user.Update(
		db.Filter{
			"id": 13,
		},
		db.Record{
			"username": "test12",
		})
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	return affected, code, err
}

// @DeleteApi /user
// @Summary 删除用户
// @Desc 维护者: 刀刀
// @Params ParamsValidate 接口参数所对应的结构体
// @Response
func (u User) Delete(c *gin.Context) (interface{}, consts.ErrCode, error) {
	user := model.User
	_, affected, err := user.Delete(db.Filter{
		"id": 22,
	})
	return affected, 0, err
}

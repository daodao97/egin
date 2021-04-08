package utils

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/daodao97/egin/lib"
)

var Validate *validator.Validate

var Trans ut.Translator

func init() {
	InitValidator()
}

func InitValidator() {
	uni := ut.New(zh.New())
	Trans, _ = uni.GetTranslator("zh")
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册翻译器
		_ = zh_translations.RegisterDefaultTranslations(validate, Trans)
		// 注册一个函数，获取struct tag里自定义的label作为字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})

		Validate = validate
	}
}

type CustomValidateFunc struct {
	Handle  validator.Func
	TagName string
	Message string
}

func RegCustomValidateFunc(customValidate CustomValidateFunc) {
	// 注册自定义函数
	_ = Validate.RegisterValidation(customValidate.TagName, customValidate.Handle)
	// 根据提供的标记注册翻译
	_ = Validate.RegisterTranslation(customValidate.TagName, Trans, func(ut ut.Translator) error {
		return ut.Add(customValidate.TagName, customValidate.Message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(customValidate.TagName, fe.Field(), fe.Field())
		return t
	})
}

func RegCustomValidateFuncs(funs []CustomValidateFunc) {
	for _, f := range funs {
		RegCustomValidateFunc(f)
	}
}

// 替换错误信息中的 struct field 为 label
func TransErr(obj interface{}, validateError validator.ValidationErrors) (map[string]string, error) {
	labelMap, err := lib.GetStructAllTag(obj, "label")
	if err != nil {
		return nil, err
	}

	vErr := validateError.Translate(Trans)
	errMessage := make(map[string]string)
	for k, v := range vErr {
		for key, value := range labelMap {
			v = strings.Replace(v, key, value, -1)
		}
		ks := strings.Split(k, ".")
		_k := ks[1]
		errMessage[_k] = v
	}

	return errMessage, nil
}

func Validated(ctx *gin.Context, params interface{}) (errMsg []string) {
	err := ctx.ShouldBind(params)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			return []string{err.Error()}
		case validator.ValidationErrors:
			errs, _ := TransErr(params, err.(validator.ValidationErrors))
			for _, msg := range errs {
				errMsg = append(errMsg, msg)
			}
			return errMsg
		}
	}
	return nil
}

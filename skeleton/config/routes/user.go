// ****************************
// 该文件为系统生成, 请勿随意更改
// ****************************
package routes

import (
	"strings"

	"github.com/daodao97/egin/egin"
	"github.com/daodao97/egin/egin/consts"
	"github.com/daodao97/egin/egin/middleware"
	"github.com/daodao97/egin/egin/utils"
	"github.com/gin-gonic/gin"

	"skeleton/controller"
)

func RegUserRouter(r *gin.Engine) {

	r.Handle("GET", "/user", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.ParamsValidate
			errs := utils.Validated(ctx, &params)
			if errs != nil {
				egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
				return
			}
			result, code, err := controller.User{}.Get(ctx, params)
			egin.Response(ctx, result, code, err)
		}
	}(), middleware.IpLimiter())

	r.Handle("POST", "/user", func(ctx *gin.Context) {
		result, code, err := controller.User{}.Post(ctx)
		egin.Response(ctx, result, code, err)
	})

	r.Handle("PUT", "/user", func(ctx *gin.Context) {
		result, code, err := controller.User{}.Put(ctx)
		egin.Response(ctx, result, code, err)
	})

	r.Handle("DELETE", "/user", func(ctx *gin.Context) {
		result, code, err := controller.User{}.Delete(ctx)
		egin.Response(ctx, result, code, err)
	})

}

func RegUserCustomValidateFunc() {
	utils.RegCustomValidateFuncs([]utils.CustomValidateFunc{
		controller.Bookabledate,
	})
}

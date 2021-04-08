// ****************************
// 该文件为系统生成, 请勿随意更改
// ****************************
package config

import (
	"github.com/gin-gonic/gin"

	"skeleton/config/routes"
)

func RegRouter(r *gin.Engine) {

	routes.RegUserRouter(r)

	routes.RegUserCustomValidateFunc()

}

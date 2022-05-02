package web

import (
	"github.com/gin-gonic/gin"
)


func Start() {
	r := gin.Default()

	Router(r)
	Api(r)

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run()
}

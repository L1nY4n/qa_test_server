package web

import (
	"net/http"
	"qa_test_server/model"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {

	// 指定模板的位置
	r.LoadHTMLGlob("templates/*.html")
	// 静态文件映射
	r.StaticFS("/assets", http.Dir("templates/assets"))
	r.Static("favicon.ico",".templates/favicon.ico")

	// 根路径加载 index 模板，web 页面的入口
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})




	r.GET("/test", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "hello.html", gin.H{ // H是一个开箱即用的map
			"message": "hello world!",
		})
	})

	r.GET("/dy", func(c *gin.Context) {

		c.HTML(http.StatusOK, "dynamic.html", gin.H{

			"title": "lf-test1",
			//"desc":    "描述",
			//"Content": "内容",
		})
	})

	r.GET("/index", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "hello",
			"desc":    "描述",
			"Content": "内容",
		})
	})
	r.GET("/info", func(c *gin.Context) {

		c.JSON(http.StatusOK, model.Dev_capture_packed{})



	})

	//fmt.Printf("receive from client, data: %v\n", string(buf[:10]))
	r.GET("/hello", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{
			"test": 1,
		})
	})
}


func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}


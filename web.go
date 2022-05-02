package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mon_packed struct {
	Title   string
	Desc    string
	Content string
}

func main() {
	go tcp_server()
	// 创建一个默认的路由引擎
	r := gin.Default()
	// GET：请求方式；/hello：请求的路径
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.LoadHTMLGlob("templates/*")
	r.GET("/test", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(200, "hello.html", gin.H{ // H是一个开箱即用的map
			"message": "hello world!",
		})
	})

	r.GET("/json", func(c *gin.Context) {
		a := mon_packed{
			Title:   "测试标题",
			Desc:    "描述",
			Content: "内容",
		}
		c.JSON(200, a)

	})

	r.GET("/jsonp", func(c *gin.Context) {
		a := mon_packed{
			Title:   "jsonp",
			Desc:    "描述",
			Content: "内容",
		}
		c.JSON(200, a)

	})

	r.GET("/index", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", gin.H{

			"title":   string(buf[:100]),
			"desc":    "描述",
			"Content": "内容",
		})
		//fmt.Printf("receive from client, data: %v\n", string(buf[:10]))
		r.GET("/hello", func(c *gin.Context) {

			c.HTML(http.StatusOK, "hello.html", gin.H{
				"test": 1,
			})
		})
	})

	r.GET("/dy", func(c *gin.Context) {

		c.HTML(http.StatusOK, "dynamic.html", gin.H{

			"title": "lf-test1",
			//"desc":    "描述",
			//"Content": "内容",
		})
	})
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run()

}

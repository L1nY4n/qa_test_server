package web

import (
	"net/http"
	"qa_test_server/device"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {

	// 指定模板的位置
	r.LoadHTMLGlob("templates/*.html")
	// 静态文件映射
	r.StaticFS("/assets", http.Dir("assets"))

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

		c.JSON(http.StatusOK, device.Dev_cap)
		// c.HTML(http.StatusOK, "device.html", gin.H{
		// 	"capture":  (string)(device.Dev_cap.Sys_para.Pro_info.SN[:20]),
		// 	"sys_mon":  device.Dev_cap.Sys_mon,
		// 	"sys_para": device.Dev_cap.Sys_para,
		// })

	})

	//fmt.Printf("receive from client, data: %v\n", string(buf[:10]))
	r.GET("/hello", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{
			"test": 1,
		})
	})
}

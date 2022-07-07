package web

import (
	"net/http"
	"qa_test_server/device"
	"sort"
	"strings"
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

	device_route := r.Group("/device")
	{
		device_route.GET("/list", devices)
	}


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

		//fmt.Printf("receive from client, data: %v\n", string(buf[:10]))
		r.GET("/hello", func(c *gin.Context) {
			c.HTML(http.StatusOK, "hello.html", gin.H{
				"test": 1,
			})
		})
	})
}


func devices(c *gin.Context) {
	manager := &device.ManagerGlabal
	list := manager.List()
	// sn 排序
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i].Sn, list[j].Sn) < 0
	})
	c.JSON(200, gin.H{
		"data":    list,
		"success": true,
	})
}
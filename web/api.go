package web

import (
	"fmt"
	"math/rand"
	"qa_test_server/manager"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//  web 服务对外提供的接口
func Api(r *gin.Engine) {

	// 路由分组管理

	// 系统
	system_route := r.Group("/system")
	{
		system_route.GET("/info", sysInfo)
	}

	// 设备
	device_route := r.Group("/device")
	{
		device_route.GET("/list", deviceList)
		device_route.GET("/info/:id", deviceInfo)
		device_route.GET("/randomData/:current", randomData)
	}

}

func sysInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"ver":    11,
		"author": "lf",
	})
}

// 设备列表
func deviceList(c *gin.Context) {

		manager := &manager.ManagerGlabal
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

// 查询单个设备的信息
func deviceInfo(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, gin.H{
		"id":   id,
		"info": "info",
	})
}

func randomData(c *gin.Context) {
	curr := c.Param("current")
	fmt.Printf("curr %v",curr)
	num, _ := strconv.Atoi(curr)
	i := rand.Intn(1000) -500
	time_str := time.Now().Format("2006/01/02 15:04:05")

	c.JSON(200, gin.H{
		"name": time_str,
		"data": gin.H{
			"time":  time_str,
			"value": num + i,
		},
	})
}

package device

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User3 struct {
	ID       int
	Name     string
	Age      int
	Value    uint64
	Birthday time.Time
}

var Dev_cap = Dev_capture_packed{
	Cap_if: Capture_info{Model: "unkonw1"},
}

func Sql_op() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:L1nFen9.com@tcp(localhost:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{

		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("连接服务器失败")
	}

	fmt.Println("连接服务器成功")
	db.Debug().AutoMigrate(&Dev_capture_packed{})
	//db.Debug().Create(&Dev_cap)
	//user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

}

package db

import (
	"fmt"
	"qa_test_server/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Sql_op() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:L1nFen9.com@tcp(localhost:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{

		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("连接服务器失败")
	}

	fmt.Println("连接服务器成功")
	DB.Debug().AutoMigrate(model.Nano_Dev_capture_packed{})
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	fmt.Println("#################迁移成功#############")

}

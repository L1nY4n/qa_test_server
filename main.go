package main

import (
	"qa_test_server/db"
	"qa_test_server/elog"
	"qa_test_server/tcpserver"
	"qa_test_server/web"
)

func main() {
	elog.RedirectStderr()
	db.Sql_op()
	// 启动 tcp server
	go tcpserver.Tcpserver()
	// 启动web 服务
	web.Start()
}

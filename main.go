package main

import (
	"log"

	"qa_test_server/config"
	"qa_test_server/db"
	"qa_test_server/elog"
	"qa_test_server/tcpserver"
	"qa_test_server/web"
)

func main() {
	cfg := config.Load()

	if err := elog.RedirectStderr(); err != nil {
		log.Printf("redirect stderr failed: %v", err)
	}

	if err := db.Init(cfg.DBDSN, cfg.DBAutoMigrate); err != nil {
		log.Printf("db init failed: %v", err)
	}

	go tcpserver.Tcpserver(cfg.TCPAddr, cfg.ProxyFromPort, cfg.ProxyToPort)
	web.Start(cfg)
}

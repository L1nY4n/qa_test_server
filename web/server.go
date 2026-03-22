package web

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"

	"qa_test_server/config"
	"qa_test_server/manager"
	"qa_test_server/model"
)

var appConfig config.AppConfig

func Start(cfg config.AppConfig) {
	appConfig = cfg

	if err := manager.UserManagerGlobal.Init(
		cfg.AuthSecret,
		cfg.AuthTokenTTL,
		cfg.DefaultAdminUsername,
		cfg.DefaultAdminPassword,
	); err != nil {
		log.Printf("init user manager failed: %v", err)
	}
	if err := manager.HistoryManagerGlobal.Init(
		&manager.ManagerGlabal,
		cfg.HistorySampleEvery,
		cfg.HistoryRetention,
	); err != nil {
		log.Printf("init history manager failed: %v", err)
	}
	if err := manager.ParamChangeManagerGlobal.Init(cfg.ChangeLogRetention); err != nil {
		log.Printf("init param-change manager failed: %v", err)
	}
	if err := manager.DecryptManagerGlobal.Init(); err != nil {
		log.Printf("init decrypt manager failed: %v", err)
	}
	manager.VirtualDeviceManagerGlobal.SetBroadcastHook(func(device model.Device) bool {
		data, err := json.Marshal(device)
		if err != nil {
			return false
		}
		return WsManager.TryGroupbroadcast("device_upload", data)
	})

	r := gin.Default()
	r.Use(Cors())
	_ = r.SetTrustedProxies(nil)

	Router(r)
	Api(r)
	WsRoute(r)

	go WsManager.Start()

	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Printf("web server exited: %v", err)
	}
}

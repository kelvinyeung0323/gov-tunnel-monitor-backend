package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/config"
	"github.com/ky/gov-tunnel-monitor-backend/base/db"
	"github.com/ky/gov-tunnel-monitor-backend/base/web"
	"github.com/ky/gov-tunnel-monitor-backend/base/websocket"
	"github.com/ky/gov-tunnel-monitor-backend/handlers"
	"github.com/ky/gov-tunnel-monitor-backend/handlers/monitor"
	"github.com/ky/gov-tunnel-monitor-backend/routes"
)

func main() {
	InitApp()
	r := gin.Default()
	//统一异常处理
	r.Use(web.ErrorHandleMiddleware)
	r.Use(db.TransactionMiddleware)
	routes.InitRoute(r)
	r.Run(":" + config.AppConfig.Port)
}

var MonHub *monitor.MonitorHub

func InitApp() {
	//初始化配置
	config.InitConfig()
	//初始化数据库
	db.InitDB()
	//初始化SQL模板
	db.InitSqlTemplate()
	//TODO: 初始化路由
	//初始化webSocket
	websocket.InitWebSocketHub()
	websocket.WSHub.MessageHanler = handlers.HandlerWebsocketRecievedMessage
	websocket.WSHub.Run()

	//初始化监控
	monitor.InitMonitor()
	monitor.Hub.ProcessData = handlers.ProcessMonitorData
	//TODO:获取所有设备
	// devices := repo.GetAllDevices(nil)
	// monitor.Hub.LoadMonitorList(devices)
	monitor.Hub.Start()
}

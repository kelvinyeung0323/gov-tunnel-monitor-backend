package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/websocket"
	"github.com/ky/gov-tunnel-monitor-backend/handlers/monitor"
	repo "github.com/ky/gov-tunnel-monitor-backend/repositories"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

//处理监控数据，单独拆出来是为了避免循环引用
func ProcessMonitorData(d *monitor.Device) error {
	pushMonitorData(d)
	persistMonitorData(d)
	return nil
}

//推送状态信息
func pushMonitorData(device *monitor.Device) {
	msg := &websocket.Message{}
	msg.MsgType = websocket.WS_BROCAST

	data := fmt.Sprintf("%v", device.Data)
	now := types.Time(time.Now())

	msg.Body = types.MonitorData{
		DeviceId:     device.DeviceId,
		DeviceName:   device.DeviceName,
		DeviceStatus: device.Status,
		DeviceType:   device.Type,
		Data:         &data,
		DataTime:     &now,
	}
	websocket.WSHub.Send(msg)
}

//持久化数据
func persistMonitorData(device *monitor.Device) {
	msg := &websocket.Message{}
	msg.MsgType = websocket.WS_BROCAST

	val := fmt.Sprintf("%v", device.Data)
	now := types.Time(time.Now())

	data := types.MonitorData{
		DeviceId:     device.DeviceId,
		DeviceName:   device.DeviceName,
		DeviceStatus: device.Status,
		DeviceType:   device.Type,
		Data:         &val,
		DataTime:     &now,
	}

	//TODO：更新设备最新监控值
	err := repo.UpdateDeviceStatus(nil, &device.Device)
	if err != nil {
		log.Printf("error:持久化设备状态数据%v", err)
	}
	repo.PersistMonitorData(&data)
	if err != nil {
		log.Printf("error:持久化监控数据%v", err)
	}
}

func AddToMonitor(deviceId int) error {
	device := repo.GetDeviceById(nil, deviceId)
	if device == nil {
		return fmt.Errorf("没有对应的设备ID")
	}
	return nil
}
func AddToMonitor2(device *types.Device) error {
	monitor.Hub.AddToMonitor(*device)
	return nil
}
func RemoveFromMonitor(deviceId int) error {
	monitor.Hub.RemoveFromMonitor(deviceId)
	return nil
}

func QueryMonitorData(ctx *gin.Context, form *types.MonitorDataQueryForm) ([]types.MonitorData, error) {
	return repo.QueryMonitorData(ctx, form), nil
}

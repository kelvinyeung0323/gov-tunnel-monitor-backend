package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/db"
	"github.com/ky/gov-tunnel-monitor-backend/base/web"
	repo "github.com/ky/gov-tunnel-monitor-backend/repositories"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

//查询设备
func QueryDevices(ctx *gin.Context, queryForm *types.DeviceQueryForm) []types.Device {
	return repo.QueryDevices(ctx, queryForm)
}

//根据ID获取设备
func GetDeviceById(ctx *gin.Context, deviceId int) *types.Device {
	device := repo.GetDeviceById(ctx, deviceId)
	if device == nil {
		web.BizErr("设备不存在")
	}
	return device
}

//创建设备
func CreateDevice(ctx *gin.Context, form *types.DeviceCreatedForm) *types.Device {
	//开户事务
	db.BeginTx(ctx)
	//TODO检查字段
	device := &types.Device{}
	device.DeviceId = form.DeviceId
	device.DeviceName = form.DeviceName
	device.Description = form.Description
	device.Ip = form.Ip
	device.Port = form.Port
	device.SetupAddr = form.SetupAddr
	device.Type = form.Type
	device.Status = types.D_OFFLINE
	device.Rtsp = form.Rtsp
	repo.CreateDevice(ctx, device)
	//提交事务
	db.CommitTx(ctx)
	return device
}

//更新设备
//返回更新后的设备
func UpdateDevice(ctx *gin.Context, device *types.Device) *types.Device {
	db.BeginTx(ctx)
	if device.DeviceId == nil {
		web.BizErr("设备不存在.")
	}

	oDvice := repo.GetDeviceById(ctx, *device.DeviceId)
	if oDvice == nil {
		web.BizErr("设备不存在.")
	}

	repo.UpdateDevice(ctx, device)
	device = repo.GetDeviceById(ctx, *device.DeviceId)
	db.CommitTx(ctx)
	return device
}

//删除设备
func DelDevice(ctx *gin.Context, deviceId int) *types.Device {
	db.BeginTx(ctx)
	db.CommitTx(ctx)
	return nil
}

//删除多个设备
//返回删除后的设务列表
func DelDevices(ctx *gin.Context, deviceIds []int) []types.Device {
	db.BeginTx(ctx)
	repo.DeleteDevices(ctx, deviceIds)
	db.CommitTx(ctx)
	return nil
}

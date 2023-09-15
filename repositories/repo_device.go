package repo

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/db"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

func GetDeviceById(ctx *gin.Context, deviceId int) *types.Device {
	u := &types.Device{}
	sql := db.GetSqlString("getDeviceById", deviceId)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(ctx).Get(u, sql)
	if err != nil {
		log.Printf("error:%v\n", err)
		return nil
	}
	return u
}

func QueryDevices(ctx *gin.Context, queryForm *types.DeviceQueryForm) []types.Device {

	devices := []types.Device{}
	sql := db.GetSqlString("queryDevices", queryForm)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(ctx).Select(&devices, sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		panic(err)
	}

	return devices
}

func CreateDevice(ctx *gin.Context, device *types.Device) {

	sql := db.GetSqlString("createDevice", device)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		panic(err)
	}

	if lastId, e := r.LastInsertId(); e == nil {
		deviceId := int(lastId)
		device.DeviceId = &deviceId
	} else {
		panic(e)
	}

}

func UpdateDevice(ctx *gin.Context, device *types.Device) {
	sql := db.GetSqlString("updateDevice", device)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		panic(err)
	}
	cnt, _ := r.RowsAffected()
	log.Printf("Rows Affected: %v", cnt)
}

func UpdateDeviceStatus(ctx *gin.Context, device *types.Device) error {
	sql := db.GetSqlString("updateDeviceStatus", device)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(ctx).Exec(sql)
	if err != nil {
		return err
	}
	cnt, _ := r.RowsAffected()
	log.Printf("Rows Affected: %v", cnt)
	return nil
}

func DeleteDevices(ctx *gin.Context, deviceIds []int) {
	sql := db.GetSqlString("deleteDevices", deviceIds)
	log.Printf("SQL is: %s\n", sql)
	// r, err := db.GetConn(ctx).Exec(sql)
	// if err != nil {
	// 	panic(err)
	// }
	// cnt, _ := r.RowsAffected()
	// log.Printf("Rows Affected: %v", cnt)
}

func GetAllDevices(ctx *gin.Context) []types.Device {
	devices := []types.Device{}
	sql := db.GetSqlString("getAllDevices", nil)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(nil).Select(&devices, sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		panic(err)
	}

	return devices
}

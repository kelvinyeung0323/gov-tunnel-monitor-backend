package routes

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/web"
	"github.com/ky/gov-tunnel-monitor-backend/handlers"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

func registerDeviceRoute(r *gin.RouterGroup) {
	//查询设备
	r.GET("/devices", func(ctx *gin.Context) {

		form := types.DeviceQueryForm{}
		if err := ctx.ShouldBindQuery(&form); err != nil {
			web.Err(web.VALID_ERROR)

		}
		devices := handlers.QueryDevices(ctx, &form)
		web.ReturnOK(ctx, &devices)

	})

	//根据用户ID获取用户
	r.GET("/device/:deviceId", func(ctx *gin.Context) {
		s := ctx.Param("deviceId")
		deviceId, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("error:%v\n", err)
			web.Err(web.VALID_ERROR)
		}

		device := handlers.GetDeviceById(ctx, deviceId)
		web.ReturnOK(ctx, device)
	})

	//创建设备
	r.POST("/device", func(ctx *gin.Context) {
		form := &types.DeviceCreatedForm{}
		if err := ctx.ShouldBind(&form); err != nil {
			log.Printf("error:%v", err)
			web.Err(web.VALID_ERROR)
		}

		device := handlers.CreateDevice(ctx, form)
		web.ReturnOK(ctx, device)

	})

	//更新用户
	r.PUT("/device", func(ctx *gin.Context) {
		device := &types.Device{}
		if err := ctx.ShouldBind(&device); err != nil {
			log.Printf("error:%v", err)
			web.Err(web.VALID_ERROR)
		}

		device = handlers.UpdateDevice(ctx, device)
		web.ReturnOKWithMsg(ctx, device, "用户修改成功！")
	})

	//删除用户
	r.DELETE("/device/:deviceId", func(ctx *gin.Context) {
		s := ctx.Param("deviceId")
		deviceId, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("error:%v\n", err)
			web.Err(web.VALID_ERROR)
		}

		device := handlers.DelDevice(ctx, deviceId)
		web.ReturnOKWithMsg(ctx, device, "用户已删除！")
	})

	//删除多个设备
	r.DELETE("/devices", func(ctx *gin.Context) {
		queryStr := ctx.Query("deviceIds")
		if queryStr == "" {
			web.Err(web.VALID_ERROR)
		}

		arr := strings.Split(queryStr, ",")
		ids := []int{}
		for _, v := range arr {
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("error: %v", err)
				web.Err(web.VALID_ERROR)
			}
			ids = append(ids, i)
		}
		device := handlers.DelDevices(ctx, ids)
		web.ReturnOKWithMsg(ctx, device, "设备已删除！")
	})
}

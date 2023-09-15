package repo

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ky/gov-tunnel-monitor-backend/base/db"
	"github.com/ky/gov-tunnel-monitor-backend/types"
)

func QueryMonitorData(ctx *gin.Context, queryForm *types.MonitorDataQueryForm) []types.MonitorData {

	data := []types.MonitorData{}
	sql := db.GetSqlString("queryMonitorData", queryForm)
	log.Printf("SQL is: %s\n", sql)
	err := db.GetConn(ctx).Select(&data, sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		panic(err)
	}

	return data
}

func PersistMonitorData(data *types.MonitorData) error {
	sql := db.GetSqlString("persistMonitorData", data)
	log.Printf("SQL is: %s\n", sql)
	r, err := db.GetConn(nil).Exec(sql)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return err
	}

	if lastId, e := r.LastInsertId(); e == nil {
		dataId := int(lastId)
		data.Id = dataId
	} else {
		return err
	}
	return nil
}

func ReportMonitorData(ctx *gin.Context, device *types.Device) {
	//TODO:
}

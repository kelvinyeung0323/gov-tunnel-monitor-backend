package types

//监控数据
type MonitorData struct {
	Id           int              `json:"id" db:"id"`
	DeviceId     *int             `json:"deviceId" db:"device_id"`
	DeviceName   *string          `json:"deviceName" db:"device_name"`
	DeviceType   DeviceTypeEnum   `json:"deviceType" db:"device_type"`
	DeviceStatus DeviceStatusEnum `json:"deviceStatus" db:"device_status"`
	OriginData   []byte           `json:"-" db:"origin_data"`
	Data         *string          `json:"data" db:"data"`
	Remark       *string          `json:"remark" db:"remark"`
	DataTime     *Time            `json:"dataTime" db:"data_time"`
}

type MonitorDataQueryForm struct {
	DeviceType   DeviceTypeEnum   `json:"deviceType"`
	DeviceName   *string          `json:"deviceName" `
	DeviceStatus DeviceStatusEnum `json:"deviceStatus"`
	StartTime    *Time            `json:"startTime" `
	EndTime      *Time            `json:"endTime" `
}

type MonitorDataReport struct {
	//TODO:
}

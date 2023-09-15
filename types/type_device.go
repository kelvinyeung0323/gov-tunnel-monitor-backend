package types

/**
* 设备类型
**/
type DeviceTypeEnum int

const (
	D_WATER_LEVELER  DeviceTypeEnum = iota //超声波液位仪
	D_WATER_RULER                          //电子水尺
	D_MSG_BOARD                            //可变情报板
	D_ALARM                                //声光报警器
	D_BARRIER_GATE                         //道闸
	D_TRAFFIC_LIGHT                        //机动车信号灯
	D_LANE_INDICATOR                       //车道指示器
	D_CCTV                                 //报像机
	D_TILT                                 //倾角变送器
)

func (d DeviceTypeEnum) String() string {
	switch d {
	case D_WATER_LEVELER:
		return "超声波液位仪"
	case D_WATER_RULER:
		return "电子水尺"
	case D_MSG_BOARD:
		return "可变情报板"
	case D_ALARM:
		return "声光报警器"
	case D_BARRIER_GATE:
		return "道闸"
	case D_TRAFFIC_LIGHT:
		return "机动车信号灯"
	case D_LANE_INDICATOR:
		return "车道指示器"
	case D_CCTV:
		return "报像机"
	case D_TILT:
		return "倾角变送器"
	default:
		return "未知"

	}
}

/**
* 设备状态
**/
type DeviceStatusEnum int

const (
	D_NORMAL DeviceStatusEnum = iota
	D_OFFLINE
	D_FAULT
	D_UNKNOW
)

func (s DeviceStatusEnum) String() string {
	switch s {
	case D_NORMAL:
		return "正常"
	case D_OFFLINE:
		return "未启用"
	case D_FAULT:
		return "故障"
	case D_UNKNOW:
		return "未知"
	default:
		return "未知"
	}
}

/**
* 设备
**/
type Device struct {
	DeviceId    *int             `json:"deviceId" db:"device_id"`
	DeviceName  *string          `json:"deviceName" db:"device_name"`
	Ip          *string          `json:"ip" db:"ip"`
	Port        *string          `json:"port" db:"port"`
	Type        DeviceTypeEnum   `json:"type" db:"type"`
	Status      DeviceStatusEnum `json:"status" db:"status"`
	Description *string          `json:"description" db:"description"`
	SetupAddr   *string          `json:"setupAddr" db:"setup_addr"`
	FaultCause  *string          `json:"faultCause" db:"fault_cause"`
	RefreshedAt *Time            `json:"refreshedAt" db:"refreshed_at" time_format:"2006-01-02 15:04:05" time_utc:"1"` //状态刷新时间
	CreatedAt   *Time            `json:"createdAt" db:"created_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	CreatedBy   *string          `json:"createdBy" db:"created_by"`
	UpdatedAt   *Time            `json:"updated_at" db:"updated_at" time_format:"2006-01-02 15:04:05" time_utc:"1"`
	UpdatedBy   *string          `json:"updatedBy" db:"updated_by"`
	Rtsp        *string          `json:"rtsp" db:"rtsp"`
	AddrCode    *byte            `json:"-" db:"addr_code"` //地址码
}

type DeviceQueryForm struct {
	Type       DeviceTypeEnum `json:"type" db:"type"`
	DeviceName *string        `json:"deviceName" db:"device_name"`
}

type DeviceCreatedForm struct {
	DeviceId    *int             `json:"deviceId" db:"device_id"`
	DeviceName  *string          `json:"deviceName" db:"device_name"`
	Ip          *string          `json:"ip" db:"ip"`
	Port        *string          `json:"port" db:"port"`
	Type        DeviceTypeEnum   `json:"type" db:"type"`
	Status      DeviceStatusEnum `json:"status" db:"status"`
	Description *string          `json:"description" db:"description"`
	SetupAddr   *string          `json:"setupAddr" db:"setup_addr"`
	Rtsp        *string          `json:"rtsp" db:"rtsp"`
}

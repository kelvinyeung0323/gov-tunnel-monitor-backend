package monitor

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/ky/gov-tunnel-monitor-backend/types"
	"github.com/ky/gov-tunnel-monitor-backend/utils/typeconv"
)

type Device struct {
	types.Device
	hub  *MonitorHub
	conn net.Conn
	Data any
}

//连接设备
func (d *Device) Dial() {
	if d.conn == nil || IsClosed(d.conn) {
		//两秒超时
		addr := *d.Ip + ":" + *d.Port
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)

		if err != nil {
			log.Printf("error: 连接异常:addr:%v\n error:%v\n", addr, err)
			d.Status = types.D_FAULT
		}
		d.conn = conn
		d.Status = types.D_NORMAL
	}
	d.hub.bus <- d
}

//检查连接是否关闭
func IsClosed(conn net.Conn) bool {
	deadline := time.Now().Add(time.Second * 1)
	if err := conn.SetReadDeadline(deadline); err != nil {
		return true
	}
	buf := []byte{}
	_, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			return true
		} else {
			return false
		}
	}
	log.Printf("debug:cheeck connection: %v", buf)
	return false

}

func RefreshWaterRulerData(d *Device) {
	data, err := d.GetWaterRulerData()
	if err != nil {
		log.Printf("error:%v", err)
		return
	}

	v := typeconv.BytesToInt16(data[3:4])
	d.Data = v
	d.hub.bus <- d
}

func RefreshWaterLevelerData(d *Device) {
	data, err := d.GetWaterLevelerData()
	if err != nil {
		log.Printf("error:%v", err)
		return
	}

	v := typeconv.BytesToFloat32(data[3:6]) * 100
	d.Data = v
	d.hub.bus <- d
}

func RefreshTiltData(d *Device) {
	x, err := d.GetTiltData(AXIS_X)
	if err != nil {
		log.Printf("error:%v", err)
	}
	d1 := typeconv.BytesToFloat32(x[3:4])

	y, err := d.GetTiltData(AXIS_Y)
	if err != nil {
		log.Printf("error:%v", err)
	}
	d2 := typeconv.BytesToFloat32(y[3:4])

	d.Data = map[string]float32{"x": d1, "y": d2}
	d.hub.bus <- d
}

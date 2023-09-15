package monitor

import (
	"fmt"
)

//获取电子水尺或超声波液位器监控数据
func (d *Device) GetWaterRulerData() ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("电子水尺地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := WaterRulerCmd(*d.AddrCode)
	if err != nil {
		return nil, fmt.Errorf("获取电子水尺指令错误:%v", err)
	}
	return send(d, cmd)
}

//获取超声波液位器监控数据
func (d *Device) GetWaterLevelerData() ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("超声波液位器地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := WaterLevelerCmd(*d.AddrCode)
	if err != nil {
		return nil, fmt.Errorf("获取超声波液位器指令错误:%v", err)
	}
	return send(d, cmd)
}

//获取倾角变送器监控数据
func (d *Device) GetTiltData(axisType AxisType) ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("倾角变送器地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := TiltCmd(*d.AddrCode, axisType)
	if err != nil {
		return nil, fmt.Errorf("获取倾角变送器指令错误:axis:%v,err:%v", axisType, err)
	}
	return send(d, cmd)
}

func (d *Device) OpenAlarm(content string, vol int) ([]byte, error) {
	cmd, err := AlarmCmd(&content, vol)
	if err != nil {
		return nil, fmt.Errorf("获取声光报警器指令错误:%v", err)
	}
	return send(d, cmd)
}
func (d *Device) CloseAlarm() ([]byte, error) {

	cmd, err := AlarmCloseCmd()
	if err != nil {
		return nil, fmt.Errorf("获取声光报警器指令关闭指令错误::%v", err)
	}
	return send(d, cmd)
}

func (d *Device) OpenBarrierGate() ([]byte, error) {
	cmd := BarrierGateOpenCmd()
	return send(d, cmd)
}
func (d *Device) CloseBarrierGate() ([]byte, error) {
	cmd := BarrierGateCloseCmd()
	return send(d, cmd)
}
func (d *Device) StopBarrierGate() ([]byte, error) {
	cmd := BarrierGateStopCmd()
	return send(d, cmd)
}
func (d *Device) OpenLaneIndicator() ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("车道指示器地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := SwitchCmd(*d.AddrCode, CT_OPEN)
	if err != nil {
		return nil, fmt.Errorf("获取车道指示器打开指令错误:%v", err)
	}
	return send(d, cmd)
}
func (d *Device) CloseLaneIndicator() ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("车道指示器地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := SwitchCmd(*d.AddrCode, CT_CLOSE)
	if err != nil {
		return nil, fmt.Errorf("获取车道指示器关闭指令错误:%v", err)
	}
	return send(d, cmd)
}
func (d *Device) OpenTrafficLight() ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("交通灯地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := SwitchCmd(*d.AddrCode, CT_OPEN)
	if err != nil {
		return nil, fmt.Errorf("获取交通灯打开指令错误:%v", err)
	}
	return send(d, cmd)
}
func (d *Device) CloseTrafficLight() ([]byte, error) {
	if d.AddrCode == nil {
		return nil, fmt.Errorf("交通灯地址码为空,id:%v", d.DeviceId)
	}
	cmd, err := SwitchCmd(*d.AddrCode, CT_CLOSE)
	if err != nil {
		return nil, fmt.Errorf("获取交通灯关闭指令错误:%v", err)
	}
	return send(d, cmd)
}

func send(d *Device, cmd []byte) ([]byte, error) {
	if IsClosed(d.conn) {
		//重连
		d.Dial()
	}
	_, err := d.conn.Write(cmd)
	if err != nil {
		return nil, fmt.Errorf("发送指令失败！%v", err)
	}
	data := []byte{}
	_, err = d.conn.Read(data)
	if err != err {
		return nil, fmt.Errorf("读取电子水尺返回数据失败:%v", err)
	}
	return data, nil
}

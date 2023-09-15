package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ky/gov-tunnel-monitor-backend/types"
)

type MonitorHub struct {
	Devices     map[string]*Device
	bus         chan *Device
	stop_chan   chan byte
	isStopped   bool
	ProcessData func(d *Device) error
}

func New() (*MonitorHub, error) {
	myHub := &MonitorHub{
		Devices:   make(map[string]*Device),
		bus:       make(chan *Device),
		stop_chan: make(chan byte),
		isStopped: true,
	}

	return myHub, nil
}

var Hub *MonitorHub

func InitMonitor() {
	Hub, _ = New()
}

//加载监控设备
func (m *MonitorHub) LoadMonitorList(devices []types.Device) {
	for _, v := range devices {
		device := &Device{Device: v, hub: m}
		m.Devices[strconv.Itoa(*device.DeviceId)] = device
	}
}

func (m *MonitorHub) AddToMonitor(d types.Device) {
	device := &Device{Device: d, hub: m}
	m.Devices[strconv.Itoa(*device.DeviceId)] = device
}

func (m *MonitorHub) RemoveFromMonitor(deviceId int) {
	delete(m.Devices, strconv.Itoa(deviceId))
}

//启动轮询监控
func (hub *MonitorHub) Start() {
	//TODO:可重入
	if !hub.isStopped {
		return
	}
	hub.isStopped = false
	hub.stop_chan = make(chan byte)

	go func() {

		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-hub.stop_chan:
				log.Println("monitor has stopped.")
				hub.isStopped = true
				return
			case device := <-hub.bus:
				if hub.ProcessData != nil {
					go hub.ProcessData(device)
				}
			case <-ticker.C:
				//刷新所有设备的连接状态
				go RefreshConnStatus(hub)
			}
		}
	}()

}

//停止轮询监控
func (hub *MonitorHub) Stop() {
	close(hub.stop_chan)
}

func (hub *MonitorHub) OperationHandler(oper *Operation) error {

	device, ok := hub.Devices[oper.DeviceId]
	if !ok {
		return fmt.Errorf("error: 设备不在监控列表中")
	}
	var err error
	switch oper.Cmd {
	case "open":
		switch device.Type {
		case types.D_ALARM:
			_, err = device.OpenAlarm(oper.Content, oper.Vol)

		case types.D_LANE_INDICATOR:
			_, err = device.OpenLaneIndicator()

		case types.D_TRAFFIC_LIGHT:
			_, err = device.OpenTrafficLight()

		case types.D_BARRIER_GATE:
			_, err = device.OpenBarrierGate()

		default:
			return fmt.Errorf("打不到对应的操作:%v", oper)
		}
	case "close":
		switch device.Type {
		case types.D_ALARM:
			_, err = device.CloseAlarm()

		case types.D_LANE_INDICATOR:
			_, err = device.CloseLaneIndicator()

		case types.D_TRAFFIC_LIGHT:
			_, err = device.CloseTrafficLight()

		case types.D_BARRIER_GATE:
			_, err = device.CloseBarrierGate()

		default:
			err = fmt.Errorf("打不到对应的操作:%v", oper)
		}

	case "stop":
		switch device.Type {
		case types.D_BARRIER_GATE:
			_, err := device.StopBarrierGate()
			return err
		default:
			return fmt.Errorf("打不到对应的操作:%v", oper)
		}
	default:
		return fmt.Errorf("打不到对应的操作:%v", oper)
	}
	go RefreshConnStatus(hub)
	return err
}

func RefreshConnStatus(hub *MonitorHub) {
	for _, v := range hub.Devices {
		//轮询连接状态
		go v.Dial()
		//获取电子水尺水位数据
		if v.Type == types.D_WATER_RULER {
			go RefreshWaterRulerData(v)

		}
		//获取超声波液位器水位数据
		if v.Type == types.D_WATER_LEVELER {
			go RefreshWaterLevelerData(v)
		}
		//获取倾角传感器数据
		if v.Type == types.D_TILT {
			go RefreshTiltData(v)
		}

	}
}

type Operation struct {
	DeviceId string
	Cmd      string
	Content  string
	Vol      int //音量
}

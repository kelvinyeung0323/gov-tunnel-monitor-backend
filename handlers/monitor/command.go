package monitor

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ky/gov-tunnel-monitor-backend/utils/typeconv"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type CmdType int

const (
	CT_OPEN  CmdType = 1
	CT_CLOSE CmdType = 2
)

type AxisType int

const (
	AXIS_X = 1
	AXIS_Y = 2
	AXIS_Z = 3
)

//信号灯指令/车道指示器指令获取   写开关量输出状态(得电)
func SwitchCmd(addr byte, cmdType CmdType) ([]byte, error) {
	//从机地址+功能码（写开关量输出状态）+开关量地址
	cmd := []byte{0x01, 0x05, addr}
	//待写入的开关量状态(高位在前  FF00 – 开关量输出得电     0000 – 开关量输出失电
	switch cmdType {
	case CT_OPEN:
		cmd = append(cmd, 0xFF, 0x00)
	case CT_CLOSE:
		cmd = append(cmd, 0x00, 0x00)
	}

	c, err := Crc(cmd)
	if err != nil {
		return nil, err
	}
	cmd = append(cmd, c...)
	//CHUNKOU_INDEX0
	cmd = append([]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06}, cmd...)
	return cmd, nil
}

//声光报警器关闭指令
func AlarmCloseCmd() ([]byte, error) {
	cmd := []byte{0x01, 0x06, 0x04, 0x0E, 0x00, 0x00}
	crc, err := Crc(cmd)
	if err != nil {
		return nil, err
	}
	cmd = append(cmd, crc...)
	return cmd, nil
}

//声光报警器指令
//content 为播报的指定内容
func AlarmCmd(content *string, vol int) ([]byte, error) {

	gbk, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(*content))
	if err != nil {
		return nil, fmt.Errorf("获取声光报警器指令，GBK转换：%v", err)
	}
	cmd := []byte{0x01, 0x1C, 0x00, 0x03}
	//音量
	volBytes := typeconv.Int16ToBytes(int16(vol))
	cmd = append(cmd, volBytes...)
	cmd = append(cmd, 0xB1, 0x02, 0x00, 0x01, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01)
	cmd = append(cmd, gbk...)
	crc, err := Crc(cmd)
	if err != nil {
		return nil, fmt.Errorf("获取声光报警器指令循环冗余校验码错误：%v", err)
	}
	cmd = append(cmd, crc...)
	return cmd, nil
}
func BarrierGateOpenCmd() []byte {
	return []byte{0xA6, 0x01, 0x01, 0x81, 0x81}
}
func BarrierGateCloseCmd() []byte {
	return []byte{0xA6, 0x01, 0x01, 0x80, 0x80}
}

func BarrierGateStopCmd() []byte {
	return []byte{0xA6, 0x01, 0x01, 0x82, 0x82}
}

//电子水尺问询指令
//返回水位值
//addr:地址码
func WaterRulerCmd(addr byte) ([]byte, error) {
	//地址码 + 功能码03 +超始地址0000+ +数据长度0006 + 校验码
	cmd := []byte{addr, 0x03, 0x00, 0x00, 0x00, 0x06}
	c, err := Crc(cmd)
	if err != nil {
		return nil, err
	}
	cmd = append(cmd, c...)
	return cmd, nil
}

//超声波液位器问询指令
func WaterLevelerCmd(addr byte) ([]byte, error) {
	//地址码 + 功能码03 +超始地址0000+ +数据长度0002 + 校验码
	cmd := []byte{addr, 0x03, 0x00, 0x00, 0x00, 0x02}
	c, err := Crc(cmd)
	if err != nil {
		return nil, err
	}
	cmd = append(cmd, c...)
	return cmd, nil
}

//倾角变送器问询指令
func TiltCmd(addr byte, axis AxisType) ([]byte, error) {
	//地址码 + 功能码03 +
	cmd := []byte{addr, 0x03}
	//寄存器地址  00  00   x轴角度
	//0001 H  y 角度
	switch axis {
	case AXIS_X:
		cmd = append(cmd, 0x00, 0x00)
	case AXIS_Y:
		cmd = append(cmd, 0x00, 0x01)
	}
	//数据长度
	cmd = append(cmd, 0x00, 0x01)
	//检验码
	c, err := Crc(cmd)
	if err != nil {
		return nil, err
	}
	cmd = append(cmd, c...)
	return cmd, nil
}

//循环冗余校检码
func Crc(data []byte) ([]byte, error) {
	// CRC寄存器全为1
	CRC := 0x0000ffff
	// 多项式校验值
	POLYNOMIAL := 0x0000a001
	dataLen := len(data)
	for i := 0; i < dataLen; i++ {
		CRC ^= (int(data[i]) & 0x000000ff)
		for j := 0; j < 8; j++ {
			if (CRC & 0x00000001) != 0 {
				CRC >>= 1
				CRC ^= POLYNOMIAL
			} else {
				CRC >>= 1
			}
		}
	}

	//高低位转换
	// high := (CRC << 16) & 0xffffffff
	// low := CRC >> 16
	// newVal := low | high
	bytesBuf := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuf, binary.LittleEndian, uint16(CRC))
	return bytesBuf.Bytes(), err
}

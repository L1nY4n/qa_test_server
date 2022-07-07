package device

import (
	"strings"
	"time"
)

type Device struct {
	Sn   string
	Name string

	Packet       Dev_capture_packed
	Last_rx_time time.Time
	// todo 添加其他需要呈现的字段，从 Dev_capture_packed 结构中解码出来利于呈现
}

// 从 原始数据中解析出设备的信息内容
func Decode(packet Dev_capture_packed) Device {

	name := strings.Trim(string(packet.Cap_info.Name[:]), "\x00")
	sn := strings.Trim(string(packet.Sys_para.Pro_info.SN[:]), "\x00")
	return Device{
		Sn:           sn,
		Name:         name,
		Last_rx_time: time.Now(),
		Packet:       packet,
	}
}

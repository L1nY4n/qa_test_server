package model

import (
	"fmt"
	"strings"
	"time"
)

type Device struct {
	Sn   string
	Name string

	Packet       Femto_msg_packed
	Last_rx_time time.Time
	// todo 添加其他需要呈现的字段，从 Dev_capture_packed 结构中解码出来利于呈现
}

// 从 原始数据中解析出设备的信息内容
func Femto_Decode(packet Femto_msg_packed) Device {

	name := strings.Trim(string(packet.Femto_holding_reg.Laser_para.Laser_info.Model[:]), "\x00")
	sn := strings.Trim(string(packet.Femto_holding_reg.Laser_para.Laser_info.SN[:]), "\x00")
	fmt.Printf("##sn==%s\r\n\n", sn)
	return Device{
		Sn:           sn,
		Name:         name,
		Last_rx_time: time.Now(),
		Packet:       packet,
	}
}

/*
	消息体结构信息
	1.采集板板卡信息
	2.所连接的设备类型
*/

package model

/*板卡信息*/
//数据来源于上位机
//数据也有可能来源于硬件采集板卡
//设备信息

type Capture_info struct {
	Time uint64    `gorm:"autoCreateTime:nano"` //时间戳
	Name [30]uint8 `gorm:"-"`
	Bate [20]uint8 `gorm:"-"`
}

// type Dev_info struct {
// 	ID                 uint32
// 	Model              uint8 //设备型号
// 	Hardware_bate      uint8 //版本信息
// 	Fpga_software_bate uint16
// 	Mcu_software_bate  uint8
// 	Nano_mon           Nano_laser_mon `gorm:"embedded"` //监控结构体
// }

//消息体结构(下位机组包形式)
//采集设备+激光器类型

type Nano_Dev_capture_packed struct {
	Cap_info Capture_info               `json:"采集设备信息" gorm:"embedded"` //老化板卡
	Sys_para Nano_Sys_init_param_struct `json:"系统参数" gorm:"embedded"`   //老化板卡
	Sys_mon  Nano_laser_mon             `json:"系统监控" gorm:"embedded"`   //采集数据

}

type Pico_Dev_capture_packed struct {
	Cap_info Capture_info               `json:"采集设备信息" gorm:"embedded"` //老化板卡
	Sys_para Pico_Sys_init_param_struct `json:"系统参数" gorm:"embedded"`   //老化板卡
	Sys_mon  Pico_laser_mon             `json:"系统监控" gorm:"embedded"`   //采集数据

}

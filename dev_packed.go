/*
	消息体结构信息
	1.采集板板卡信息
	2.所连接的设备类型


*/

package main

/*板卡信息*/
//数据来源于上位机
//数据也有可能来源于硬件采集板卡

type Capture_info struct {
	Model         string //采集卡型号，老化卡或上位机上传
	Hardware_bate byte   //硬件版本
	Software_bate byte   //软件版本
}

/*设备信息*/
type Dev_info struct {
	Model              string //设备型号
	Hardware_bate      byte   //版本信息
	Fpga_software_bate uint16
	Mcu_software_bate  byte
	Capture_data       [512]byte //数据采集包
}

//消息体结构(下位机组包形式)
//采集设备+激光器类型
type Dev_capture_packed struct {
	cap_info Capture_info //老化板卡
	dev_info Dev_info     //采集数据

}

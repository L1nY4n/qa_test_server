/*
	消息体结构信息
	1.采集板板卡信息
	2.所连接的设备类型
*/

package device

/*板卡信息*/
//数据来源于上位机
//数据也有可能来源于硬件采集板卡

type Capture_info struct {
	Model         string //采集卡型号，老化卡或上位机上传
	Hardware_bate byte   //硬件版本
	Software_bate byte   //软件版本
}

//#########纳秒设备监控结构体###############//

//设备的种子开关状态
type Dev_Laser_status struct {
	laser_mode   byte
	laser_status byte
	seed_on      byte
	ld1_on       byte
	ld2_on       byte
	fpga_ap_on   byte
	fpga_apm1_on byte
	fpga_apm2_on byte
	fpga_aom_on  byte
}

//电流监控参数
type Dev_cur_mon struct {
	ld_cur   [2]float32
	pre_cur  [2]float32
	main_cur [4]float32
	seed_cur float32
	tcm_cur  float32
}

/*电压监控 结构体*/
type Dev_vol_mon struct {
	vp12_vol float32
	vn12_vol float32
	pre_vol  [2]float32
	main_vol [4]float32
	n3v3_vol float32
	f3v3_vol float32
	f2v5_vol float32
	ld_vol   [2]float32
}

/*温度监控*/
type Dev_temperature_mon struct {
	ld_temp         [2]float32
	tcm_temp        [2]float32
	tcm_module_temp float32
	seed_temp       float32
	pre_pump_temp   [2]float32
	main_pump_temp  [4]float32
	plate_temp      float32
	hdc_temp        float32
	hdc_humi        float32
	hdc_dew_point   float32
}

/*种子模块*/
type Dev_seed_mon struct {
	enable         byte
	alarm          byte
	trig_mode      byte
	pulse_width    uint16
	pulse_freq     uint16
	fpga_read_freq uint16
	supply_vol     uint16
	power          uint16
	mudule_temp    float32
	ld_temp        float32
}

/*温控模块*/
type Dev_tcm_mon struct {
	sw          byte
	actual_temp float32
	adjust_temp float32
	rms         float32
	senser_type float32
	control_p   float32
	control_i   float32
	control_d   float32
	alarm       byte
	TCM_Ready   byte
}

/*告警监控模块*/
type Dev_alarm_mon struct {
	Alarm_flag_reg_now     uint64
	Alarm_flag_reg_active  uint64
	Alarm_flag_reg_history uint64
}

/*fpga寄存器值回读*/
type Dev_fpga_mon struct {
	pre_cur_set         float32
	main_cur_set        float32
	aom_vol_set         float32
	ap_freq             float32
	ap_duty             float32
	apm_freq            float32
	apm_duty            float32
	aom_freq            float32
	aom_duty            float32
	fpga_apm_pwm_sw     uint16
	fpga_aom_delay      uint16
	ext_power_level_bit uint16
}

type dev_time struct {
	year   byte
	month  byte
	day    byte
	hour   byte
	miunte byte
	second byte
}

type Dev_time_mon struct {
	sys_time dev_time
}

type Nano_laser_mon struct {
	headA     byte
	headB     byte
	mon_laser Dev_Laser_status
	mon_cur   Dev_cur_mon
	mon_vol   Dev_vol_mon
	mon_temp  Dev_temperature_mon
	mon_tcm   Dev_tcm_mon
	mon_seed  Dev_seed_mon
	mon_alarm Dev_alarm_mon
	mon_fpga  Dev_fpga_mon
	mon_time  dev_time
}

type Dev_info struct {
	Model              string //设备型号
	Hardware_bate      byte   //版本信息
	Fpga_software_bate uint16
	Mcu_software_bate  byte
	nano_              Nano_laser_mon //监控结构体
}

//消息体结构(下位机组包形式)
//采集设备+激光器类型

type Dev_capture_packed struct {
	cap_info Capture_info //老化板卡
	dev_info Dev_info     //采集数据

}

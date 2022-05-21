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
	Model         string `gorm:"primaryKey"` //采集卡型号，老化卡或上位机上传
	Hardware_bate byte   //硬件版本
	Software_bate byte   //软件版本
}

//#########纳秒设备监控结构体###############//

//设备的种子开关状态
type Dev_Laser_status struct {
	Laser_mode   byte
	Laser_status byte
	Seed_on      byte
	Ld1_on       byte
	Ld2_on       byte
	Fpga_ap_on   byte
	Fpga_apm1_on byte
	Fpga_apm2_on byte
	Fpga_aom_on  byte
}

//电流监控参数
type Dev_cur_mon struct {
	Ld_cur   [2]float32 `gorm:"serializer:json"`
	Pre_cur  [2]float32 `gorm:"serializer:json"`
	Main_cur [4]float32 `gorm:"serializer:json"`
	Seed_cur float32
	Tcm_cur  float32
}

/*电压监控 结构体*/
type Dev_vol_mon struct {
	Vp12_vol float32
	Vn12_vol float32
	Pre_vol  [2]float32 `gorm:"serializer:json"`
	Main_vol [4]float32 `gorm:"serializer:json"`
	N3v3_vol float32
	F3v3_vol float32
	F2v5_vol float32
	Ld_vol   [2]float32 `gorm:"serializer:json"`
}

/*温度监控*/
type Dev_temperature_mon struct {
	Ld_temp         [2]float32 `gorm:"serializer:json"`
	Tcm_temp        [2]float32 `gorm:"serializer:json"`
	Tcm_module_temp float32
	Seed_temp       float32
	Pre_pump_temp   [2]float32 `gorm:"serializer:json"`
	Main_pump_temp  [4]float32 `gorm:"serializer:json"`
	Plate_temp      float32
	Hdc_temp        float32
	Hdc_humi        float32
	Hdc_dew_point   float32
}

/*种子模块*/
type Dev_seed_mon struct {
	Enable         byte
	Alarm          byte
	Trig_mode      byte
	Pulse_width    uint16
	Pulse_freq     uint16
	Fpga_read_freq uint16
	Supply_vol     uint16
	Power          uint16
	Mudule_temp    float32
	Ld_temp        float32
}

/*温控模块*/
type Dev_tcm_mon struct {
	SW          byte
	Actual_temp float32
	Adjust_temp float32
	Rms         float32
	Senser_type float32
	Control_p   float32
	Control_i   float32
	Control_d   float32
	Alarm       byte
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
	Pre_cur_set         float32
	Main_cur_set        float32
	Aom_vol_set         float32
	Ap_freq             float32
	Ap_duty             float32
	Apm_freq            float32
	Apm_duty            float32
	Aom_freq            float32
	Aom_duty            float32
	Fpga_apm_pwm_sw     uint16
	Fpga_aom_delay      uint16
	Ext_power_level_bit uint16
}

type Dev_time struct {
	Year   byte
	Month  byte
	Day    byte
	Hour   byte
	Miunte byte
	Second byte
}

type Dev_time_mon struct {
	Sys_time Dev_time
}

type Nano_laser_mon struct {
	HeadA     byte
	HeadB     byte
	Mon_laser Dev_Laser_status    `gorm:"embedded"`
	Mon_cur   Dev_cur_mon         `gorm:"embedded"`
	Mon_vol   Dev_vol_mon         `gorm:"embedded"`
	Mon_temp  Dev_temperature_mon `gorm:"embedded"`
	Mon_tcm   Dev_tcm_mon         `gorm:"embedded"`
	Mon_seed  Dev_seed_mon        `gorm:"embedded"`
	Mon_alarm Dev_alarm_mon       `gorm:"embedded"`
	Mon_fpga  Dev_fpga_mon        `gorm:"embedded"`
	Mon_time  Dev_time            `gorm:"embedded"`
}

type Dev_info struct {
	ID                 int
	Model              string //设备型号
	Hardware_bate      byte   //版本信息
	Fpga_software_bate uint16
	Mcu_software_bate  byte
	Nano_mon           Nano_laser_mon `gorm:"embedded"` //监控结构体
}

//消息体结构(下位机组包形式)
//采集设备+激光器类型

type Dev_capture_packed struct {
	Cap_if Capture_info `gorm:"embedded"` //老化板卡
	Dev_if Dev_info     `gorm:"embedded"` //采集数据
}

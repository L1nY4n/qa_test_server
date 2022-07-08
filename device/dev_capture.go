/*
	消息体结构信息
	1.采集板板卡信息
	2.所连接的设备类型
*/

package device

/*板卡信息*/
//数据来源于上位机
//数据也有可能来源于硬件采集板卡
//设备信息

type Capture_info struct {
	Time uint64    `gorm:"autoCreateTime:nano"` //时间戳
	Name [30]uint8 `gorm:"-"`
	Bate [20]uint8 `gorm:"-"`
}

//#########纳秒设备监控结构体###############//

//设备的种子开关状态
type Dev_Laser_status struct {
	Laser_mode   uint8 `json:"工作模式"`
	Laser_status uint8 `json:"激光器状态"`
	Seed_on      uint8 `json:"种子开关"`
	Ld1_on       uint8 `json:"单模1开关"`
	Ld2_on       uint8 `json:"单模2开关"`
	Fpga_ap_on   uint8 `json:"预放开关"`
	Fpga_apm1_on uint8 `json:"主放1开关"`
	Fpga_apm2_on uint8 `json:"主放2开关"`
	Fpga_apm3_on uint8 `json:"主放3开关"`
	Fpga_apm4_on uint8 `json:"主放4开关"`
	Fpga_aom_on  uint8 `json:"声光开关"`
}

//电流监控参数
type Dev_cur_mon struct {
	Ld_1cur        float32 `json:"单模1电流"`
	Ld_2cur        float32 `json:"单模2电流"`
	Pre1_cur       float32 `json:"预放1电流"`
	Pre2_cur       float32 `json:"预放2电流"`
	Main1_cur      float32 `json:"主放1电流"`
	Main2_cur      float32 `json:"主放2电流"`
	Main3_cur      float32 `json:"主放3电流"`
	Main4_cur      float32 `json:"主放4电流"`
	Main_cur_limit float32 `json:"主放电流限制"`
	Seed_cur       float32 `json:"种子电流"`
	Tcm_cur        float32 `json:"温控电流"`
}

/*种子模块*/
type Dev_seed_mon struct {
	Enable         uint8     `json:"种子开关"`
	Alarm          uint8     `json:"告警值"`
	Trig_mode      uint8     `json:"触发模式"`
	Pulse_width    uint16    `json:"脉冲宽度"`
	Pulse_freq     uint16    `json:"脉冲设置频率"`
	Select_freq    uint16    `json:"选单频率"`
	Fpga_read_freq uint16    `json:"测频值"`
	Supply_vol     uint16    `json:"供电电压"`
	Power          uint16    `json:"功率"`
	Mudule_temp    float32   `json:"模块温度"`
	Seed_ld_temp   float32   `json:"种子ld温度"`
	Para_reserved  [20]uint8 `json:"预留参数"`
}

/*电压监控 结构体*/
type Dev_vol_mon struct {
	Vp12_vol  float32 `json:"vp电压"`
	Vn12_vol  float32 `json:"vn电压"`
	Pre1_vol  float32 `json:"预放1电压"`
	Pre2_vol  float32 `json:"预放2电压"`
	Main1_vol float32 `json:"主放1电压"`
	Main2_vol float32 `json:"主放2电压"`
	Main3_vol float32 `json:"主放3电压"`
	Main4_vol float32 `json:"主放4电压"`
	N3v3_vol  float32 `json:"常规3.3V"`
	F3v3_vol  float32 `json:"FPGA3.3V"`
	F1v25_vol float32 `json:"FPGA2.5V"`
	Ld1_vol   float32 `json:"单模1电压"`
	Ld2_vol   float32 `json:"单模2电压"`
	Vcc12     float32 `json:"VCC12"`
	Vcc24     float32 `json:"VCC24"`
}

/*温度监控*/
type Dev_temperature_mon struct {
	Ld1_temp        float32 `json:"单模1温度"`
	Ld2_temp        float32 `json:"单模2温度"`
	Tcm_temp        float32 `json:"温控温度"`
	Tcm_module_temp float32 `json:"温控模块温度"`
	Seed_temp       float32 `json:"种子温度"`
	Pre1_temp       float32 `json:"预放1温度"`
	Pre2_temp       float32 `json:"预放2温度"`
	Main1_temp      float32 `json:"主放1温度"`
	Main2_temp      float32 `json:"主放2温度"`
	Main3_temp      float32 `json:"主放3温度"`
	Main4_temp      float32 `json:"主放4温度"`
	Plate_temp      float32 `json:"底盘温度"`
	Hdc_temp        float32 `json:"温湿度传感器-温度"`
	Hdc_humi        float32 `json:"温湿度传感器-湿度"`
	Hdc_dew_point   float32 `json:"露点值"`
}

/*温控模块*/
type Dev_tcm_mon struct {
	SW          uint8   `json:"开关"`
	Actual_temp float32 `json:"实际温度"`
	Adjust_temp float32 `json:"设置温度"`
	Rms         float32 `json:"误差"`
	Senser_type float32 `json:"传感器类型"`
	Control_p   float32 `json:"控制-P"`
	Control_i   float32 `json:"控制-I"`
	Control_d   float32 `json:"控制-D"`
	Alarm       uint8   `json:"告警"`
	TCM_Ready   uint8   `json:"准备信号"`
}

/*告警监控模块*/
type Dev_alarm_mon struct {
	Alarm_flag_reg_now     uint64 `json:"当前告警"`
	Alarm_flag_reg_active  uint64 `json:"生效告警"`
	Alarm_flag_reg_history uint64 `json:"历史告警"`
}

/*fpga寄存器值回读*/
type Dev_fpga_mon struct {
	Pre_cur_set         float32 `json:"预放电流设置"`
	Main_cur_set        float32 `json:"主放电流设置"`
	Aom_vol_set         float32 `json:"aom设置"`
	Ap_freq             float32 `json:"预放频率"`
	Ap_duty             float32 `json:"预放占空比"`
	Apm_freq            float32 `json:"主放频率"`
	Apm_duty            float32 `json:"主放占空比"`
	Aom_freq            float32 `json:"声光频率"`
	Aom_duty            float32 `json:"声光占空比"`
	Fpga_apm_pwm_sw     uint16  `json:"pwm设置"`
	Fpga_aom_delay      uint16  `json:"声光延时"`
	Ext_io_status       uint16  `json:"外控io1"`
	Ext_power_level_bit uint16  `json:"外控io2"`
}

type Dev_time struct {
	Year   uint8 `json:"年"`
	Month  uint8 `json:"月"`
	Day    uint8 `json:"日"`
	Hour   uint8 `json:"时"`
	Miunte uint8 `json:"分"`
	Second uint8 `json:"秒"`
}

type Dev_time_mon struct {
	Sys_time     Dev_time `gorm:"embedded" json:"设备时间"`
	Sys_run_time uint32   `gorm:"-" json:"运行时间"`
}

type Nano_laser_mon struct {
	HeadA     uint8
	HeadB     uint8
	Mon_laser Dev_Laser_status    `json:"激光器状态监测" gorm:"embedded"`
	Mon_cur   Dev_cur_mon         `json:"激光器电流监测" gorm:"embedded"`
	Mon_vol   Dev_vol_mon         `json:"激光器电压监测" gorm:"embedded"`
	Mon_temp  Dev_temperature_mon `json:"激光器温度监测" gorm:"embedded"`
	Mon_tcm   [3]Dev_tcm_mon      `json:"激光器温控板监测" gorm:"embedded"`
	Mon_seed  Dev_seed_mon        `json:"激光器种子模块监测" gorm:"embedded"`
	Mon_alarm Dev_alarm_mon       `json:"激光器告警监测" gorm:"embedded"`
	Mon_fpga  Dev_fpga_mon        `json:"激光器FPGA寄存器监测" gorm:"embedded"`
	Mon_time  Dev_time            `json:"激光器时间监测" gorm:"embedded"`
}

//------------------------------------存储结构体-------------------------------------//
type Product_info struct {
	Model           [20]byte   `json:"型号,string" gorm:"-"` //sn号码
	SN              [20]byte   `gorm:"-"`                  //sn号码
	PN              [20]uint8  `gorm:"-"`                  //sn号码
	Hardware_bate   [5]uint8   `gorm:"-"`                  //硬件版本号
	Optics_bate     [15]uint8  `gorm:"-"`                  //光学版本
	Mcu_vision      [20]uint8  `json:"mcu版本"`              //mcu版本号
	Fpga_vision     [10]uint16 `json:"fpga版本"`             //fpga 版本号
	Storage_cnt     uint16     `json:"存储次数"`
	Record_time     Dev_time   `json:"记录次数" gorm:"embedded"` //此表更新时间
	Activation_time Dev_time   `json:"激活时间" gorm:"embedded"` //激活时间
	Deadline_time   Dev_time   `json:"失效时间" gorm:"embedded"` //失效时间
}

//标志位

type Sys_conf_check struct {
	IsLock          uint8  `json:"锁机"`   //锁机标志位
	Unlock_time     uint16 `json:"锁机次数"` //解锁次数
	Fiber_alarm_cnt uint8  `json:"光纤告警"` //
	Res1            uint8  `json:"预留"`   //
	IsActive        uint8  `json:"激活状态"`
	IsOverdue       uint8  `json:"超期"`  //超期
	IsOverdue_op    uint8  `json:"升级"`  //锁机操作
	IsFactory       uint8  `json:"初始化"` //出厂
}

//种子参数
type Seed_param struct {
	Seed_board_type uint8     `json:"种子板类型"`
	Width           uint16    `json:"脉宽"`
	Freq            uint16    `json:"频率"`
	Burst_number    uint16    `json:"串个数"`
	Burst_interval  uint8     `json:"串间隔"`
	Burst_wareform  [20]uint8 `json:"波形参数"`
	Para_reserved   [20]uint8 `json:"预留"`
	Power_level     uint16    `json:"功率"`

	Trigger_mode uint8 `json:"触发模式"`
	//Temp         float32 `json:"温度"`
}

//FPGA 参数
//     [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Ansi, Pack = 1)]
//     struct fpga_init_param
// {
//         public uint8 adc_start;
//         public uint8 work_mode;
//         public uint pwm_freq;
//         public uint pwm_duty;
//         public uint int_enable;
//         public uint8 ole_red_mode;
//         public uint8 ole_red_on;
//         public uint k;
//         public uint8 aom_current_level;
//         public uint8 apm_current_level;
//         public uint8 ap_current;
//         public uint8 aom_apm_set;
//         public uint8 sw_ap;
//         public uint8 sw_apm;
//         [MarshalAs(UnmanagedType.ByValArray, SizeConst = 4)]
//         public uint[] pwm_ramp;
//         public uint8 seed_dac_set;
//         public uint lvds_period_hs05;
//         public uint lvds_period_dac;
//         public uint fiber_aom_ttl;
//         public uint space_aom_ttl;
//         public uint alarm_filter_time;
//         public uint alarm_en;
//         public uint alarm_clear;
//         public uint alarm_out;
//         public uint alarm_lat;
//         public uint fiber_break_delay;

//         public uint pd_freq1_thr_h;
//         public uint pd_freq1_thr_l;
//         public uint pd_freq2_thr_h;
//         public uint pd_freq2_thr_l;
//         public uint pd_freq3_thr_h;
//         public uint pd_freq3_thr_l;
// };

type PD_param struct {
	Freq_ref    [3]float32 `json:"测频阈值" gorm:"-"`
	Protect_ref [3]float32 `json:"保护阈值" gorm:"-"`
}

type Temp_param struct {
	Ld_temp   [2]float32 `gorm:"-" json:"单模设定温度" `
	Tcm_temp  [3]float32 `gorm:"-" json:"tcm设定温度"`
	Seed_temp float32    `json:"种子温度"`
}

type PID struct {
	En        float32 `json:"使能"`
	Set       float32 `json:"板内pid设定值"`
	Actual    float32 `json:"板内pid实际值"`
	Err       float32 `json:"板内pid误差"`
	Err_last  float32 `json:"板内pid历史误差"`
	Kp        float32 `json:"板内pid值-p"`
	Ki        float32 `json:"板内pid值-i"`
	Kd        float32 `json:"板内pid值-d"`
	Kout      float32 `json:"板内pid值-out"`
	Voltage   float32 `json:"板内pid值-电压"`
	Integral  float32 `json:"板内pid值-积分值"`
	Limit_max float32 `json:"板内pid值-最大限定"`
	Limit_min float32 `json:"板内pid值-最小限定"`
}

//温控模块控制

type Tcm_control_param struct {
	Tcm_control_p [3]float32 `json:"温控模块-p"`
	Tcm_control_i [3]float32 `json:"温控模块-i"`
	Tcm_control_d [3]float32 `json:"温控模块-d"`
	Pid           PID        `json:"板内温控模块" gorm:"embedded"`
}

type Cur_param struct {
	Seed_cur    float32    `json:"种子电流"`
	Single_cur  [2]float32 `json:"单模电流" gorm:"-"`
	Ap_cur      float32    `json:"预放电流"`
	Ap_thr_h    float32    `json:"预放电流高阈值"` //预放高阈值
	Ap_thr_l    float32    `json:"预放电流低阈值"` //预放低阈值
	Apm_max_cur float32    `json:"主放默认电流"`
}

type Mon_limit_param struct {
	//阈值

	Vol_limit_max [10]float32 `json:"电压限定高" gorm:"-"`

	Vol_limit_min [10]float32 `json:"电压限定低" gorm:"-"`

	Cur_limit_max [15]float32 `json:"电流限定高" gorm:"-"`

	Cur_limit_min [15]float32 `json:"电流限定低" gorm:"-"`

	Temp_limit_max [10]float32 `json:"温度限定高" gorm:"-"`

	Temp_limit_min [10]float32 `json:"温度限定低" gorm:"-"`

	//四组校准k值

	//Power_k    [16]float32 `json:"三次项校准参数" gorm:"-"`
	Fpga_thr   [150]uint8 `json:"一次项校准参数" gorm:"-"`
	Alarm_mask uint64     `json:"告警屏蔽位"`
}

//存储合集

type Sys_init_param_struct struct {
	HeadA      uint8          `json:"包头1"`
	HeadB      uint8          `json:"包头2"`
	Pro_info   Product_info   `json:"产品信息"  gorm:"embedded"`
	Seed_param Seed_param     `json:"种子信息"  gorm:"embedded"`
	Sys_cc     Sys_conf_check `json:"设备标志位"  gorm:"embedded"`
	//public fpga_init_param fpga_init_p;
	Pd_param  PD_param          `json:"pd参数"  gorm:"embedded"`
	Temp      Temp_param        `json:"温度参数"  gorm:"embedded"`
	Tcm       Tcm_control_param `json:"温控参数"  gorm:"embedded"`
	Cur       Cur_param         `json:"电流参数"  gorm:"embedded"`
	Mon_limit Mon_limit_param   `json:"监控限定值"  gorm:"embedded"`
	CRCH      uint8             `json:"校验1"  gorm:"embedded"`
	CRCL      uint8             `json:"校验2"`
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

type Dev_capture_packed struct {
	Cap_info Capture_info          `json:"采集设备信息" gorm:"embedded"` //老化板卡
	Sys_para Sys_init_param_struct `json:"系统参数" gorm:"embedded"`   //老化板卡
	Sys_mon  Nano_laser_mon        `json:"系统监控" gorm:"embedded"`   //采集数据

}

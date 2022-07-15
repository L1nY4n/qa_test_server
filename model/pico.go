package model

//皮秒解析结构体
type Pico_Product_info struct {
	Model           [20]uint8 //设备信息
	SN              [20]uint8 //产品sn号
	PN              [20]uint8 //PN信息
	Hardware_bate   uint8     //硬件版本号
	Optics_bate     uint8     //光学方案
	Mcu_vision      [20]uint8 //muc软件版本
	Fpga_vision     uint16    //fpga版本
	Storage_cnt     uint16    //存储参数保存
	Record_time     Dev_time  //此表更新时间
	Activation_time Dev_time  //激活时间
	Deadline_time   Dev_time  //失效时间
}

//配置相关
type Pico_Sys_conf_check struct {
	IsLock      uint8  //锁机状态
	Unlock_time uint32 //解锁次数
	IsActive    uint8  //是否激活
	IsOverdue   uint8  //过期
	IsUpdate    uint8  //升级标志位
	IsFactory   uint8  //恢复出厂设置标志位

}

//种子参数
type Pico_Seed_param struct {
	Seed_board_type uint8  //种类板类型
	Width           uint16 //脉冲宽度
	Freq            uint16 //脉冲频率
	Burst_number    uint16
	Burst_interval  uint8     //脉冲间隔
	Burst_wareform  [20]uint8 //脉冲波形
	Power_level     uint16    //功率
	Trigger_mode    uint8     //触发方式
	//float temp;    //温度
}

//PD 阈值参数
type Pico_PD_param struct {
	Freq_ref    [3]float32 //3路测频阈值设定
	Protect_ref [3]float32 //3路保护阈值
}

////温度设定参数
type Pico_Temp_param struct {
	Ld_temp   [2]float32 //单模泵浦温度
	Tcm_temp  [3]float32 //温控模块控制
	Seed_temp float32    //种子温度控制
}

//温控电路参数
type Pico_Tcm_control_param struct {
	Tcm_control_p [3]float32
	Tcm_control_i [3]float32
	Tcm_control_d [3]float32
	Pid           PID //温控模块温度

}

//电流设定参数
type Pico_Cur_param struct {
	//		__packed struct{
	//		float cur[9];
	//	} cur_sort;
	Seed_cur       float32
	Single_def_cur [2]float32 //单模电流限定值
	Ap_def_cur     [2]float32 //预放电流值
	Apm_def_cur    [4]float32 //四路ld电流
	Apm_ramp_speed float32
}

type Pico_Mon_limit_param struct {
	//阈值
	Vol_limit_max [10]float32
	Vol_limit_min [10]float32
	//0-2作为曲线校准apm索引
	//3是ap电流限制
	Cur_limit_max [15]float32
	Cur_limit_min [15]float32
	//温度限制
	Temp_limit_max [10]float32
	Temp_limit_min [10]float32
	//四组校准k值,内部软件使用
	//	float power_k[16];
	Fpga_thr [110]uint8
	//警告掩膜
	Alarm_mask uint64
}

type Pico_Motor_para struct {
	Model       [2]uint8  //电机滑台型号
	Aoto_change [2]uint8  //自动换点
	Reset_flag  [2]uint8  //上电复位
	Change_duty [2]uint16 //换点周期,单位分钟

	Pos     [2]uint32 //坐标位置
	Pos_max [2]uint32
	Step    [2]uint32 //移动转化步进值

}

type Pico_Sys_init_param_struct struct {
	//包头参数
	HeadA uint8
	HeadB uint8

	Pro_info Pico_Product_info
	//种子板默认参数
	Seed_param     Pico_Seed_param
	Sys_conf_check Pico_Sys_conf_check

	//pd 阈值设定
	Pd_param Pico_PD_param
	// 温度设定
	Temp_param Pico_Temp_param
	//温控参数
	Tcm_para Pico_Tcm_control_param
	//默认电流参数
	Cur_param Pico_Cur_param
	Mon_limit Pico_Mon_limit_param
	CRCH      uint8
	CRCL      uint8
}

//监控结构体

////-----------------监控量缓存结构体--------------//
type Pico_Mon_Laser_Status_Struct struct {
	Laser_mode   uint8 //工作模式
	Laser_status uint8 //激光状态
	Seed_on      uint8 //种子开光状态
	Ld1_on       uint8
	Ld2_on       uint8
	Fpga_ap_on   uint8
	Fpga_apm1_on uint8
	Fpga_apm2_on uint8
	Fpga_apm3_on uint8
	Fpga_apm4_on uint8
	Fpga_aom_on  uint8
}

//电流监控 单位是A
type Pico_Mon_Cur_Struct struct {
	Single_pump_cur [2]float32
	pre_pump_cur    [2]float32
	main_pump_cur   [4]float32
	main_cur_limit  float32
	seed_cur        float32
	tcm_cur         float32
}

//电压监控，单位是v
type Pico_Mon_Vol_Struct struct {
	//从adc128s102采集得来
	vp12v_vol float32
	vn12v_vol float32
	pre_vol   [2]float32
	main_vol  [4]float32
	//使用片内adc采集而来
	normal_3v3_vol float32
	fpga_3v3_vol   float32
	fpga_1v2_vol   float32
	single_vol     [2]float32
	vcc12v         float32
	vcc24v         float32
}

//温度监控
type Pico_Mon_Temp_Struct struct {
	Single_temp     [2]float32 //单模温度;
	Tcm_temp        float32    //板内温控模块温度
	Tcm_module_temp float32    //温控模块
	Seed_temp       float32
	Pre_pump_temp   [2]float32
	Main_pump_temp  [4]float32
	Plate_temp      float32

	//温湿度传感器
	Hdc_temp  float32
	Hdc_humi  float32
	Dew_point float32
}

//-----种子模块监控
type Pico_Mon_Seed_Struct struct {
	SW                      uint8
	mode                    uint8
	seed_type               uint8
	work_status             uint16
	seed_mode_locked_status uint16
	burst                   uint16
	select_freq             uint32 //选频
	freq                    uint32
	fpga_read_freq          uint16 //测频回读
	power                   uint16
}

//温度控制模块
type Pico_Mon_TCM_Struct struct {
	SW          uint8
	Actual_temp float32
	Adjust_temp float32
	Rms         float32
	Senser_type float32
	Control_p   float32
	Control_i   float32
	Control_d   float32
	Alarm       uint8
	TCM_Ready   uint8
}

type Pico_Mon_ALARM_Struct struct {
	//警报标志寄存器
	//警告出现时,设置为1,警告清除设置为0
	Alarm_flag_reg_now     uint64 //实时监测监控寄存器,会随着系统自动变化,
	Alarm_flag_reg_active  uint64 //实际报错触发用到的寄存器，支持屏蔽设置，支持清除操作
	Alarm_flag_reg_history uint64
}

//FPGA监控结构体
type Pico_Mon_FPGA_Struct struct {
	Pre_cur_set         float32 //预放电流设置
	Main_cur_set        float32 //主放电流设置
	Apm_cur             [4]float32
	Apm_cur_dest        [4]float32
	main_cur_ramp_speed float32
	Aom_vol_set         float32 //声光电压设置
	//预放频率
	Ap_freq float32
	//预放占空比
	Ap_duty float32
	//主放频率
	Apm_freq float32
	//主放占空比
	Apm_duty float32
	//声光频率
	Aom_freq float32
	//声光占空比
	Aom_duty float32

	Fpga_apm_pwm_sw     uint16
	Fpga_aom_delay      uint16
	Ext_io_status       uint16
	Ext_power_level_bit uint16
}

type Pico_Time_Struct struct {
	Year    uint8
	Month   uint8
	Day     uint8
	Hour    uint8
	Miuntes uint8
	Second  uint8
}

type Pico_Mon_Time_Struct struct {

	//激活日期
	//	Time_Struct activation_time;
	//系统时间
	Sys_time Pico_Time_Struct
	//失效时间
	//Time_Struct deadline_time;
	Sys_runtime uint32
}

type Pico_Motor_Mon_Struct struct {
	//复位标志位,回到0坐标值
	Reset_flag  [2]uint8  //是否归0
	Init_flag   [2]uint8  //是否初始化
	Cur_pos     [2]uint32 //当前位置
	Dest_pos    [2]uint32 //目标位置
	Pos_max     [2]uint32
	Total_range [2]uint32 //
	Step        [2]uint16 //步进值

}

type Pico_laser_mon struct {
	HeadA     uint8
	HeadB     uint8
	Mon_laser Pico_Mon_Laser_Status_Struct
	Mon_cur   Pico_Mon_Cur_Struct
	Mon_vol   Pico_Mon_Vol_Struct
	Mon_temp  Pico_Mon_Temp_Struct
	Mon_tcm   [3]Pico_Mon_TCM_Struct
	Mon_seed  Pico_Mon_Seed_Struct
	Mon_alarm Pico_Mon_ALARM_Struct
	Mon_fpga  Pico_Mon_FPGA_Struct
	Mon_time  Pico_Mon_Time_Struct
	Mon_motor Pico_Motor_Mon_Struct
}

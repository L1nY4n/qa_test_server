package model

//飞秒数据结构使用modbus和数据通讯

const USER_PARA_OFFECT = 1500

//输入寄存器
const (
	//[Description("激光器状态")]
	MOD_LASER_STATUS = iota /*激光器状态显示*/
	//[Description("在位检测0")]
	MOD_ADDR_ONLINE_0 /*模块在位检测IO回读0-15*/
	//[Description("在位检测1")]
	MOD_ADDR_ONLINE_1 /*模块在位检测IO回读16-31*/

	//[Description("外控8bit信号监测")]
	MOD_ADDR_EXT_DATA /*外控8bit信号监测*/
	//[Description("外控锁存信号监测")]
	MOD_ADDR_EXT_LATCH /*外控锁存信号监测*/
	//[Description("水冷机告警")]
	MOD_ADDR_EXT_ALARM_COOLER /*水冷机告警*/
	//[Description("测试IO监测")]
	MOD_ADDR_EXT_TEST /*测试IO监测*/
	//[Description("测试预留1信号监测")]
	MOD_ADDR_EXT_RESERVE1 /*测试预留1信号监测*/
	//[Description("测试预留2信号监测")]
	MOD_ADDR_EXT_RESERVE2 /*测试预留2信号监测*/
	//[Description("外控GATE信号监测")]
	MOD_ADDR_EXT_GATE /*外控GATE信号监测*/
	//[Description("外控PWM信号监测")]
	MOD_ADDR_EXT_PWM /*外控PWM信号监测*/
	//[Description("外控PRR信号监测")]
	MOD_ADDR_EXT_PRR /*外控PRR信号监测*/
	//[Description("外控TRIG信号监测")]
	MOD_ADDR_EXT_TRIG /*外控TRIG信号监测*/
	//[Description("外控SYNC信号监测")]
	MOD_ADDR_EXT_SYNC /*外控SYNC信号监测*/
	//[Description("外控模拟转换8bit监测")]
	MOD_ADDR_ANOLOG_DATA /*外控模拟转换8bit监测*/
	//[Description("水冷机监测频率")]
	MOD_ADDR_WATER_FLOW_FREQ
	//[Description("水流量监测")]
	MOD_ADDR_WATER_FLOW

	/*版本信息*/
	//[Description("电路硬件版本")]
	MOD_ADDR_HARDWARE_BATE /*电路硬件版本*/
	//[Description("boot版本")]
	MOD_ADDR_MCU_BOOT_BATE /*boot版本*/
	//[Description("app版本0")]
	MOD_ADDR_MCU_APP_BATE_0 /*app版本*/
	//[Description("app版本1")]
	MOD_ADDR_MCU_APP_BATE_1
	//[Description("app版本2")]
	MOD_ADDR_MCU_APP_BATE_2
	//[Description("app版本3")]
	MOD_ADDR_MCU_APP_BATE_3
	//[Description("app版本4")]
	MOD_ADDR_MCU_APP_BATE_4
	//[Description("app版本5")]
	MOD_ADDR_MCU_APP_BATE_5

	//[Description("FPGA 版本0")]
	MOD_ADDR_FPGA_BATE_0
	//[Description("FPGA 版本1")]
	MOD_ADDR_FPGA_BATE_1
	//[Description("FPGA 版本2")]
	MOD_ADDR_FPGA_BATE_2
	//[Description("FPGA 版本3")]
	MOD_ADDR_FPGA_BATE_3
	//[Description("FPGA 版本4")]
	MOD_ADDR_FPGA_BATE_4
	//[Description("FPGA 版本5")]
	MOD_ADDR_FPGA_BATE_5

	//[Description("PD0检测值")]
	MOD_ADDR_PD0
	//[Description("PD1检测值")]
	MOD_ADDR_PD1
	//[Description("PD2检测值")]
	MOD_ADDR_PD2
	//[Description("PD3检测值")]
	MOD_ADDR_PD3

	/*种子状态*/
	//[Description("种子状态1")]
	MOD_ADDR_DIGI_SEED_STATUS1
	//[Description("种子状态2")]
	MOD_ADDR_DIGI_SEED_STATUS2
	//[Description("种子状态3")]
	MOD_ADDR_DIGI_SEED_STATUS3
	//[Description("种子状态4")]
	MOD_ADDR_DIGI_SEED_STATUS4
	//[Description("种子状态5")]
	MOD_ADDR_DIGI_SEED_STATUS5
	//[Description("种子状态6")]
	MOD_ADDR_DIGI_SEED_STATUS6

	//[Description("温度点0监控")]
	MOD_ADDR_TEMP0
	//[Description("温度点1监控")]
	MOD_ADDR_TEMP1
	//[Description("温度点2监控")]
	MOD_ADDR_TEMP2
	//[Description("温度点3监控")]
	MOD_ADDR_TEMP3
	//[Description("温度点4监控")]
	MOD_ADDR_TEMP4
	//[Description("温度点5监控")]
	MOD_ADDR_TEMP5
	//[Description("温度点6监控")]
	MOD_ADDR_TEMP6
	//[Description("温度点7监控")]
	MOD_ADDR_TEMP7
	//[Description("温度点8监控")]
	MOD_ADDR_TEMP8
	//[Description("温度点9监控")]
	MOD_ADDR_TEMP9
	//[Description("温度点10监控")]
	MOD_ADDR_TEMP10
	//[Description("温度点11监控")]
	MOD_ADDR_TEMP11
	//[Description("温度点12监控")]
	MOD_ADDR_TEMP12
	//[Description("温度点13监控")]
	MOD_ADDR_TEMP13
	//[Description("温度点14监控")]
	MOD_ADDR_TEMP14
	//[Description("温度点15监控")]
	MOD_ADDR_TEMP15
	//[Description("温度点16监控")]
	MOD_ADDR_TEMP16
	//[Description("温度点17监控")]
	MOD_ADDR_TEMP17
	//[Description("温度点18监控")]
	MOD_ADDR_TEMP18
	//[Description("温度点19监控")]
	MOD_ADDR_TEMP19

	/*电流值回读*/
	//[Description("PUMP0实际开关状态")]
	MOD_ADDR_PUMP0_ON
	//[Description("PUMP0实际电流回读")]
	MOD_ADDR_PUMP0_ACTUAL_CUR
	//[Description("PUMP0_FPGA寄存器电流")]
	MOD_ADDR_PUMP0_FPGA_CUR

	//[Description("PUMP1实际开关状态")]
	MOD_ADDR_PUMP1_ON
	//[Description("PUMP1实际电流回读")]
	MOD_ADDR_PUMP1_ACTUAL_CUR
	//[Description("PUMP0_FPGA寄存器电流")]
	MOD_ADDR_PUMP1_FPGA_CUR

	//[Description("PUMP2实际开关状态")]
	MOD_ADDR_PUMP2_ON
	//[Description("PUMP2实际电流回读")]
	MOD_ADDR_PUMP2_ACTUAL_CUR
	//[Description("PUMP2_FPGA寄存器电流")]
	MOD_ADDR_PUMP2_FPGA_CUR

	//[Description("PUMP3实际开关状态")]
	MOD_ADDR_PUMP3_ON
	//[Description("PUMP3实际电流回读")]
	MOD_ADDR_PUMP3_ACTUAL_CUR
	//[Description("PUMP3_FPGA寄存器电流")]
	MOD_ADDR_PUMP3_FPGA_CUR

	//[Description("PUMP4实际开关状态")]
	MOD_ADDR_PUMP4_ON
	//[Description("PUMP4实际电流回读")]
	MOD_ADDR_PUMP4_ACTUAL_CUR
	//[Description("PUMP4_FPGA寄存器电流")]
	MOD_ADDR_PUMP4_FPGA_CUR

	//[Description("PUMP5实际开关状态")]
	MOD_ADDR_PUMP5_ON
	//[Description("PUMP5实际电流回读")]
	MOD_ADDR_PUMP5_ACTUAL_CUR
	//[Description("PUMP5_FPGA寄存器电流")]
	MOD_ADDR_PUMP5_FPGA_CUR

	//[Description("PUMP6实际开关状态")]
	MOD_ADDR_PUMP6_ON
	//[Description("PUMP6实际电流回读")]
	MOD_ADDR_PUMP6_ACTUAL_CUR
	//[Description("PUMP6_FPGA寄存器电流")]
	MOD_ADDR_PUMP6_FPGA_CUR

	//[Description("PUMP7实际开关状态")]
	MOD_ADDR_PUMP7_ON
	//[Description("PUMP7实际电流回读")]
	MOD_ADDR_PUMP7_ACTUAL_CUR
	//[Description("PUMP7_FPGA寄存器电流")]
	MOD_ADDR_PUMP7_FPGA_CUR

	//[Description("PUMP8实际开关状态")]
	MOD_ADDR_PUMP8_ON
	//[Description("PUMP8实际电流回读")]
	MOD_ADDR_PUMP8_ACTUAL_CUR
	//[Description("PUMP8_FPGA寄存器电流")]
	MOD_ADDR_PUMP8_FPGA_CUR

	//[Description("PUMP9实际开关状态")]
	MOD_ADDR_PUMP9_ON
	//[Description("PUMP9实际电流回读")]
	MOD_ADDR_PUMP9_ACTUAL_CUR
	//[Description("PUMP9_FPGA寄存器电流")]
	MOD_ADDR_PUMP9_FPGA_CUR

	//[Description("PUMP10实际开关状态")]
	MOD_ADDR_PUMP10_ON
	//[Description("PUMP10实际电流回读")]
	MOD_ADDR_PUMP10_ACTUAL_CUR
	//[Description("PUMP10_FPGA寄存器电流")]
	MOD_ADDR_PUMP10_FPGA_CUR

	//[Description("PUMP11实际开关状态")]
	MOD_ADDR_PUMP11_ON
	//[Description("PUMP11实际电流回读")]
	MOD_ADDR_PUMP11_ACTUAL_CUR
	//[Description("PUMP11_FPGA寄存器电流")]
	MOD_ADDR_PUMP11_FPGA_CUR

	//[Description("PUMP12实际开关状态")]
	MOD_ADDR_PUMP12_ON
	//[Description("PUMP12实际电流回读")]
	MOD_ADDR_PUMP12_ACTUAL_CUR
	//[Description("PUMP12_FPGA寄存器电流")]
	MOD_ADDR_PUMP12_FPGA_CUR

	//[Description("PUMP13实际开关状态")]
	MOD_ADDR_PUMP13_ON
	//[Description("PUMP13实际电流回读")]
	MOD_ADDR_PUMP13_ACTUAL_CUR
	//[Description("PUMP13_FPGA寄存器电流")]
	MOD_ADDR_PUMP13_FPGA_CUR

	//[Description("PUMP14实际开关状态")]
	MOD_ADDR_PUMP14_ON
	//[Description("PUMP14实际电流回读")]
	MOD_ADDR_PUMP14_ACTUAL_CUR
	//[Description("PUMP14_FPGA寄存器电流")]
	MOD_ADDR_PUMP14_FPGA_CUR

	/*电压值回读*/
	//[Description("电压点0回读")]
	MOD_ADDR_VOL0
	//[Description("电压点1回读")]
	MOD_ADDR_VOL1
	//[Description("电压点2回读")]
	MOD_ADDR_VOL2
	//[Description("电压点3回读")]
	MOD_ADDR_VOL3
	//[Description("电压点4回读")]
	MOD_ADDR_VOL4
	//[Description("电压点5回读")]
	MOD_ADDR_VOL5
	//[Description("电压点6回读")]
	MOD_ADDR_VOL6
	//[Description("电压点7回读")]
	MOD_ADDR_VOL7
	//[Description("电压点8回读")]
	MOD_ADDR_VOL8
	//[Description("电压点9回读")]
	MOD_ADDR_VOL9
	//[Description("电压点10回读")]
	MOD_ADDR_VOL10
	//[Description("电压点11回读")]
	MOD_ADDR_VOL11
	//[Description("电压点12回读")]
	MOD_ADDR_VOL12
	//[Description("电压点13回读")]
	MOD_ADDR_VOL13
	//[Description("电压点14回读")]
	MOD_ADDR_VOL14
	//[Description("电压点15回读")]
	MOD_ADDR_VOL15
	//[Description("电压点16回读")]
	MOD_ADDR_VOL16
	//[Description("电压点17回读")]
	MOD_ADDR_VOL17
	//[Description("电压点18回读")]
	MOD_ADDR_VOL18
	//[Description("电压点19回读")]
	MOD_ADDR_VOL19

	/*电机实际坐标点*/
	//[Description("电机1实际位置点")]
	MOD_ADDR_MOTOR0_ACTUAL_POS
	//[Description("电机2实际位置点")]
	MOD_ADDR_MOTOR1_ACTUAL_POS
	// [Description("电机3实际位置点")]
	MOD_ADDR_MOTOR2_ACTUAL_POS
	// [Description("电机4实际位置点")]
	MOD_ADDR_MOTOR3_ACTUAL_POS

	//  [Description("温湿度模块0温度")]
	MOD_ADDR_TH0_MDDULE_TEMP
	//  [Description("温湿度模块0湿度")]
	MOD_ADDR_TH0_MDDULE_HUMI
	//  [Description("温湿度模块1温度")]
	MOD_ADDR_TH1_MDDULE_TEMP
	//  [Description("温湿度模块1湿度")]
	MOD_ADDR_TH1_MDDULE_HUMI

	//[Description("数字温控0当前温度")]
	MOD_ADDR_DIGI_TCM0_ACTUAL_TEMP
	//[Description("数字温控0告警寄存器")]
	MOD_ADDR_DIGI_TCM0_ALARM_REG
	//[Description("数字温控1当前温度")]
	MOD_ADDR_DIGI_TCM1_ACTUAL_TEMP
	//[Description("数字温控1告警寄存器")]
	MOD_ADDR_DIGI_TCM1_ALARM_REG
	//[Description("数字温控2当前温度")]
	MOD_ADDR_DIGI_TCM2_ACTUAL_TEMP
	//[Description("数字温控2告警寄存器")]
	MOD_ADDR_DIGI_TCM2_ALARM_REG
	//[Description("数字温控3当前温度")]
	MOD_ADDR_DIGI_TCM3_ACTUAL_TEMP
	//[Description("数字温控3告警寄存器")]
	MOD_ADDR_DIGI_TCM3_ALARM_REG
	//[Description("数字温控4当前温度")]
	MOD_ADDR_DIGI_TCM4_ACTUAL_TEMP
	//[Description("数字温控4告警寄存器")]
	MOD_ADDR_DIGI_TCM4_ALARM_REG
	//[Description("数字温控5当前温度")]
	MOD_ADDR_DIGI_TCM5_ACTUAL_TEMP
	//[Description("数字温控5告警寄存器")]
	MOD_ADDR_DIGI_TCM5_ALARM_REG

	//[Description("系统过期")]
	MOD_ADDR_LASER_OVERDUE

	/*FPGA调试回读*/
	//[Description("FPGA调试值地址回读")]
	MOD_ADDR_FPGA_DEBUG_ADDR_ECHO //
	//[Description("FPGA调试值回读_0")]
	MOD_ADDR_FPGA_DEBUG_VAL_0
	//[Description("FPGA调试值回读_1")]
	MOD_ADDR_FPGA_DEBUG_VAL_1
	//[Description("FPGA调试次数统计")]
	MOD_ADDR_FPGA_DEBUG_CNT

	/*#############时间定义#############*/

	/*激光器泵浦工作时间(秒)*/
	//[Description("泵浦工作时长0")]
	MOD_ADDR_LASER_PUMP_WORK_TIME_0
	//[Description("泵浦工作时长1")]
	MOD_ADDR_LASER_PUMP_WORK_TIME_1

	/*激光器出光时间(秒)*/
	//[Description("激光器出光时长0")]
	MOD_ADDR_LASER_EMISSON_TIME_0
	//[Description("激光器出光时长1")]
	MOD_ADDR_LASER_EMISSON_TIME_1

	/*激光器此次上电时间(秒)*/
	//[Description("激光器上电时长0")]
	MOD_ADDR_LASER_UPTIME_0
	//[Description("激光器上电时长1")]
	MOD_ADDR_LASER_UPTIME_1

	/*激光器总运行时间(秒)*/
	//[Description("激光器总运行时长0")]
	MOD_ADDR_LASER_TOTAL_UPTIME_0
	//[Description("激光器总运行时长1")]
	MOD_ADDR_LASER_TOTAL_UPTIME_1

	/*激光器系统时间*/
	//[Description("激光器系统时间--年&月")]
	MOD_ADDR_LASER_SYSTIME_YEAR_MON
	//[Description("激光器系统时间--日&时")]
	MOD_ADDR_LASER_SYSTIME_DAY_HOR
	//[Description("激光器系统时间--分&秒")]
	MOD_ADDR_LASER_SYSTIME_MINU_SEC

	/*激光器当前告警*/
	//[Description("激光器当前告警0")]
	MOD_ADDR_LASER_ALARM_NOW_0
	//[Description("激光器当前告警1")]
	MOD_ADDR_LASER_ALARM_NOW_1
	//[Description("激光器当前告警2")]
	MOD_ADDR_LASER_ALARM_NOW_2
	//[Description("激光器当前告警3")]
	MOD_ADDR_LASER_ALARM_NOW_3
	//[Description("激光器当前告警4")]
	MOD_ADDR_LASER_ALARM_NOW_4
	//[Description("激光器当前告警5")]
	MOD_ADDR_LASER_ALARM_NOW_5
	//[Description("激光器当前告警6")]
	MOD_ADDR_LASER_ALARM_NOW_6
	//[Description("激光器当前告警7")]
	MOD_ADDR_LASER_ALARM_NOW_7
	//[Description("激光器当前告警8")]
	MOD_ADDR_LASER_ALARM_NOW_8
	//[Description("激光器当前告警9")]
	MOD_ADDR_LASER_ALARM_NOW_9
	//[Description("激光器当前告警10")]
	MOD_ADDR_LASER_ALARM_NOW_10
	//[Description("激光器当前告警11")]
	MOD_ADDR_LASER_ALARM_NOW_11

	/*激光器历史告警*/
	//[Description("激光器历史告警0")]
	MOD_ADDR_LASER_ALARM_HISTORY_0
	//[Description("激光器历史告警1")]
	MOD_ADDR_LASER_ALARM_HISTORY_1
	//[Description("激光器历史告警2")]
	MOD_ADDR_LASER_ALARM_HISTORY_2
	//[Description("激光器历史告警3")]
	MOD_ADDR_LASER_ALARM_HISTORY_3
	//[Description("激光器历史告警4")]
	MOD_ADDR_LASER_ALARM_HISTORY_4
	//[Description("激光器历史告警5")]
	MOD_ADDR_LASER_ALARM_HISTORY_5
	//[Description("激光器历史告警6")]
	MOD_ADDR_LASER_ALARM_HISTORY_6
	//[Description("激光器历史告警7")]
	MOD_ADDR_LASER_ALARM_HISTORY_7
	//[Description("激光器历史告警8")]
	MOD_ADDR_LASER_ALARM_HISTORY_8
	//[Description("激光器历史告警9")]
	MOD_ADDR_LASER_ALARM_HISTORY_9
	//[Description("激光器历史告警10")]
	MOD_ADDR_LASER_ALARM_HISTORY_10
	//[Description("激光器历史告警11")]
	MOD_ADDR_LASER_ALARM_HISTORY_11

	/*输入寄存器结尾*/

	MOD_ADDR_INPUT_END
)

// MODBUS_INPUT_REGISTERS;

const (

	//[Description("参数表头")]
	MOD_ADDR_PARA_HEAD = iota

	/*######设备型号############*/
	//[Description("设备型号_0")]
	MOD_ADDR_LASER_MODEL_0
	//[Description("设备型号_1")]
	MOD_ADDR_LASER_MODEL_1
	//[Description("设备型号_2")]
	MOD_ADDR_LASER_MODEL_2
	//[Description("设备型号_3")]
	MOD_ADDR_LASER_MODEL_3
	//[Description("设备型号_4")]
	MOD_ADDR_LASER_MODEL_4
	//[Description("设备型号_5")]
	MOD_ADDR_LASER_MODEL_5
	//[Description("设备型号_6")]
	MOD_ADDR_LASER_MODEL_6
	//[Description("设备型号_7")]
	MOD_ADDR_LASER_MODEL_7
	//[Description("设备型号_8")]
	MOD_ADDR_LASER_MODEL_8
	//[Description("设备型号_9")]
	MOD_ADDR_LASER_MODEL_9

	/*######设备SN############*/
	//[Description("设备SN_0")]
	MOD_ADDR_LASER_SN_0
	//[Description("设备SN_1")]
	MOD_ADDR_LASER_SN_1
	//[Description("设备SN_2")]
	MOD_ADDR_LASER_SN_2
	//[Description("设备SN_3")]
	MOD_ADDR_LASER_SN_3
	//[Description("设备SN_4")]
	MOD_ADDR_LASER_SN_4
	//[Description("设备SN_5")]
	MOD_ADDR_LASER_SN_5
	//[Description("设备SN_6")]
	MOD_ADDR_LASER_SN_6
	//[Description("设备SN_7")]
	MOD_ADDR_LASER_SN_7
	//[Description("设备SN_8")]
	MOD_ADDR_LASER_SN_8
	//[Description("设备SN_9")]
	MOD_ADDR_LASER_SN_9

	/*######设备PN############*/
	//[Description("设备PN_0")]
	MOD_ADDR_LASER_PN_0
	//[Description("设备PN_1")]
	MOD_ADDR_LASER_PN_1
	//[Description("设备PN_2")]
	MOD_ADDR_LASER_PN_2
	//[Description("设备PN_3")]
	MOD_ADDR_LASER_PN_3
	//[Description("设备PN_4")]
	MOD_ADDR_LASER_PN_4
	//[Description("设备PN_5")]
	MOD_ADDR_LASER_PN_5
	//[Description("设备PN_6")]
	MOD_ADDR_LASER_PN_6
	//[Description("设备PN_7")]
	MOD_ADDR_LASER_PN_7
	//[Description("设备PN_8")]
	MOD_ADDR_LASER_PN_8
	//[Description("设备PN_9")]
	MOD_ADDR_LASER_PN_9

	//[Description("激光系列")]
	MOD_ADDR_LASER_SERIES
	//[Description("激光功率输出")]
	MOD_ADDR_LASER_POWER_LELVE

	/*###########esp32模块###########*/
	//[Description("ESP32_使能")]
	MOD_ADDR_ESP_EN
	//[Description("ESP32_模式")]
	MOD_ADDR_ESP_MODE
	//[Description("ESP32_WIFI_SSID0")]
	MOD_ADDR_ESP_SSID_0
	//[Description("ESP32_WIFI_SSID1")]
	MOD_ADDR_ESP_SSID_1
	//[Description("ESP32_WIFI_SSID2")]
	MOD_ADDR_ESP_SSID_2
	//[Description("ESP32_WIFI_SSID3")]
	MOD_ADDR_ESP_SSID_3
	//[Description("ESP32_WIFI_SSID4")]
	MOD_ADDR_ESP_SSID_4
	//[Description("ESP32_WIFI_SSID5")]
	MOD_ADDR_ESP_SSID_5
	//[Description("ESP32_WIFI_SSID6")]
	MOD_ADDR_ESP_SSID_6
	//[Description("ESP32_WIFI_SSID7")]
	MOD_ADDR_ESP_SSID_7
	//[Description("ESP32_WIFI_SSID8")]
	MOD_ADDR_ESP_SSID_8
	//[Description("ESP32_WIFI_SSID9")]
	MOD_ADDR_ESP_SSID_9
	//[Description("ESP32_WIFI_SSID10")]
	MOD_ADDR_ESP_SSID_10
	//[Description("ESP32_WIFI_SSID11")]
	MOD_ADDR_ESP_SSID_11
	//[Description("ESP32_WIFI_SSID12")]
	MOD_ADDR_ESP_SSID_12
	//[Description("ESP32_WIFI_SSID13")]
	MOD_ADDR_ESP_SSID_13
	//[Description("ESP32_WIFI_SSID14")]
	MOD_ADDR_ESP_SSID_14

	//[Description("ESP32_IP地址0")]
	MOD_ADDR_ESP_WIFI_IP_0
	//[Description("ESP32_IP地址1")]
	MOD_ADDR_ESP_WIFI_IP_1
	//[Description("ESP32_IP地址2")]
	MOD_ADDR_ESP_WIFI_IP_2
	//[Description("ESP32_IP地址3")]
	MOD_ADDR_ESP_WIFI_IP_3

	//[Description("ESP32_网口使能")]
	MOD_ADDR_ETH_EN
	//[Description("ESP32_网口模式")]
	MOD_ADDR_ETH_MODE
	//[Description("ESP32_网口IP0")]
	MOD_ADDR_ETH_IP_0
	//[Description("ESP32_网口IP1")]
	MOD_ADDR_ETH_IP_1
	//[Description("ESP32_网口IP2")]
	MOD_ADDR_ETH_IP_2
	//[Description("ESP32_网口IP3")]
	MOD_ADDR_ETH_IP_3

	/*###########温湿度模块###########*/
	//[Description("温湿度模块0使能")]
	MOD_ADDR_TH0_MDDULE_EN
	//[Description("湿度模块0温度上限")]
	MOD_ADDR_TH0_MDDULE_TEMP_MAX
	//[Description("湿度模块0温度下限")]
	MOD_ADDR_TH0_MDDULE_TEMP_MIN
	//[Description("湿度模块0湿度上限")]
	MOD_ADDR_TH0_MDDULE_HUMI_MAX
	//[Description("湿度模块0湿度下限")]
	MOD_ADDR_TH0_MDDULE_HUMI_MIN

	//[Description("温湿度模块1使能")]
	MOD_ADDR_TH1_MDDULE_EN
	//[Description("湿度模块1温度上限")]
	MOD_ADDR_TH1_MDDULE_TEMP_MAX
	//[Description("湿度模块1温度下限")]
	MOD_ADDR_TH1_MDDULE_TEMP_MIN
	//[Description("湿度模块1湿度上限")]
	MOD_ADDR_TH1_MDDULE_HUMI_MAX
	//[Description("湿度模块1湿度下限")]
	MOD_ADDR_TH1_MDDULE_HUMI_MIN

	/*###########气泵模块###########*/
	//[Description("气泵0使能")]
	MOD_ADDR_MODULE_AIR_PUMP0_MODE
	//[Description("气泵0开启阈值")]
	MOD_ADDR_MODULE_AIR_PUMP0_OPEN_TH
	//[Description("气泵0关闭阈值")]
	MOD_ADDR_MODULE_AIR_PUMP0_CLOSE_TH

	//[Description("气泵1使能")]
	MOD_ADDR_MODULE_AIR_PUMP1_MODE
	//[Description("气泵1开启阈值")]
	MOD_ADDR_MODULE_AIR_PUMP1_OPEN_TH
	//[Description("气泵1关闭阈值")]
	MOD_ADDR_MODULE_AIR_PUMP1_CLOSE_TH

	/*###########种子参数定义#############*/

	//[Description("种子使能")]
	MOD_ADDR_SEED_EN
	//[Description("种子类型")]
	MOD_ADDR_SEED_TYPE
	//[Description("种子开关")]
	MOD_ADDR_SEED_SW
	//[Description("种子参数1")]
	MOD_ADDR_SEED_PARA1
	//[Description("种子参数2")]
	MOD_ADDR_SEED_PARA2
	//[Description("种子参数3")]
	MOD_ADDR_SEED_PARA3
	//[Description("种子参数4")]
	MOD_ADDR_SEED_PARA4
	//[Description("种子参数5")]
	MOD_ADDR_SEED_PARA5

	/*###########泵浦驱动参数定义#############*/

	//[Description("泵浦0使能")]
	MOD_ADDR_PUMP0_EN /*泵浦使能*/
	//[Description("泵浦0开关")]
	MOD_ADDR_PUMP0_SW /*泵浦开关*/
	//[Description("泵浦0优先级设置")]
	MOD_ADDR_PUMP0_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦0目标电流值")]
	MOD_ADDR_PUMP0_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦0补偿系数")]
	MOD_ADDR_PUMP0_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦0电流设置比例系数")]
	MOD_ADDR_PUMP0_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦0电流速率调节ma/s")]
	MOD_ADDR_PUMP0_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦0电流上限")]
	MOD_ADDR_PUMP0_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦0电流回读adc通道绑定")]
	MOD_ADDR_PUMP0_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦0电流回读补偿值")]
	MOD_ADDR_PUMP0_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦0电流回读系数值")]
	MOD_ADDR_PUMP0_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦0回读电流错误阈值")]
	MOD_ADDR_PUMP0_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦0电流采样滤波次数")]
	MOD_ADDR_PUMP0_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦1使能")]
	MOD_ADDR_PUMP1_EN /*泵浦使能*/
	//[Description("泵浦1开关")]
	MOD_ADDR_PUMP1_SW /*泵浦开关*/
	//[Description("泵浦1优先级设置")]
	MOD_ADDR_PUMP1_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦1目标电流值")]
	MOD_ADDR_PUMP1_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦1补偿系数")]
	MOD_ADDR_PUMP1_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦1电流设置比例系数")]
	MOD_ADDR_PUMP1_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦1电流速率调节ma/s")]
	MOD_ADDR_PUMP1_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦1电流上限")]
	MOD_ADDR_PUMP1_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦1电流回读adc通道绑定")]
	MOD_ADDR_PUMP1_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦1电流回读补偿值")]
	MOD_ADDR_PUMP1_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦1电流回读系数值")]
	MOD_ADDR_PUMP1_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦1回读电流错误阈值")]
	MOD_ADDR_PUMP1_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦1电流采样滤波次数")]
	MOD_ADDR_PUMP1_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦2使能")]
	MOD_ADDR_PUMP2_EN /*泵浦使能*/
	//[Description("泵浦2开关")]
	MOD_ADDR_PUMP2_SW /*泵浦开关*/
	//[Description("泵浦2优先级设置")]
	MOD_ADDR_PUMP2_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦2目标电流值")]
	MOD_ADDR_PUMP2_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦2补偿系数")]
	MOD_ADDR_PUMP2_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦2电流设置比例系数")]
	MOD_ADDR_PUMP2_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦2电流速率调节ma/s")]
	MOD_ADDR_PUMP2_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦2电流上限")]
	MOD_ADDR_PUMP2_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦2电流回读adc通道绑定")]
	MOD_ADDR_PUMP2_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦2电流回读补偿值")]
	MOD_ADDR_PUMP2_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦2电流回读系数值")]
	MOD_ADDR_PUMP2_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦2回读电流错误阈值")]
	MOD_ADDR_PUMP2_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦2电流采样滤波次数")]
	MOD_ADDR_PUMP2_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦3使能")]
	MOD_ADDR_PUMP3_EN /*泵浦使能*/
	//[Description("泵浦3开关")]
	MOD_ADDR_PUMP3_SW /*泵浦开关*/
	//[Description("泵浦3优先级设置")]
	MOD_ADDR_PUMP3_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦3目标电流值")]
	MOD_ADDR_PUMP3_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦3补偿系数")]
	MOD_ADDR_PUMP3_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦3电流设置比例系数")]
	MOD_ADDR_PUMP3_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦3电流速率调节ma/s")]
	MOD_ADDR_PUMP3_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦3电流上限")]
	MOD_ADDR_PUMP3_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦3电流回读adc通道绑定")]
	MOD_ADDR_PUMP3_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦3电流回读补偿值")]
	MOD_ADDR_PUMP3_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦3电流回读系数值")]
	MOD_ADDR_PUMP3_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦3回读电流错误阈值")]
	MOD_ADDR_PUMP3_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦3电流采样滤波次数")]
	MOD_ADDR_PUMP3_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦4使能")]
	MOD_ADDR_PUMP4_EN /*泵浦使能*/
	//[Description("泵浦4开关")]
	MOD_ADDR_PUMP4_SW /*泵浦开关*/
	//[Description("泵浦4优先级设置")]
	MOD_ADDR_PUMP4_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦4目标电流值")]
	MOD_ADDR_PUMP4_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦4补偿系数")]
	MOD_ADDR_PUMP4_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦4电流设置比例系数")]
	MOD_ADDR_PUMP4_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦4电流速率调节ma/s")]
	MOD_ADDR_PUMP4_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦4电流上限")]
	MOD_ADDR_PUMP4_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦4电流回读adc通道绑定")]
	MOD_ADDR_PUMP4_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦4电流回读补偿值")]
	MOD_ADDR_PUMP4_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦4电流回读系数值")]
	MOD_ADDR_PUMP4_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦4回读电流错误阈值")]
	MOD_ADDR_PUMP4_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦4电流采样滤波次数")]
	MOD_ADDR_PUMP4_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦5使能")]
	MOD_ADDR_PUMP5_EN /*泵浦使能*/
	//[Description("泵浦5开关")]
	MOD_ADDR_PUMP5_SW /*泵浦开关*/
	//[Description("泵浦5优先级设置")]
	MOD_ADDR_PUMP5_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦5目标电流值")]
	MOD_ADDR_PUMP5_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦5补偿系数")]
	MOD_ADDR_PUMP5_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦5电流设置比例系数")]
	MOD_ADDR_PUMP5_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦5电流速率调节ma/s")]
	MOD_ADDR_PUMP5_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦5电流上限")]
	MOD_ADDR_PUMP5_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦5电流回读adc通道绑定")]
	MOD_ADDR_PUMP5_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦5电流回读补偿值")]
	MOD_ADDR_PUMP5_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦5电流回读系数值")]
	MOD_ADDR_PUMP5_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦5回读电流错误阈值")]
	MOD_ADDR_PUMP5_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦5电流采样滤波次数")]
	MOD_ADDR_PUMP5_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦6使能")]
	MOD_ADDR_PUMP6_EN /*泵浦使能*/
	//[Description("泵浦6开关")]
	MOD_ADDR_PUMP6_SW /*泵浦开关*/
	//[Description("泵浦6优先级设置")]
	MOD_ADDR_PUMP6_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦6目标电流值")]
	MOD_ADDR_PUMP6_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦6补偿系数")]
	MOD_ADDR_PUMP6_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦6电流设置比例系数")]
	MOD_ADDR_PUMP6_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦6电流速率调节ma/s")]
	MOD_ADDR_PUMP6_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦6电流上限")]
	MOD_ADDR_PUMP6_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦6电流回读adc通道绑定")]
	MOD_ADDR_PUMP6_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦6电流回读补偿值")]
	MOD_ADDR_PUMP6_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦6电流回读系数值")]
	MOD_ADDR_PUMP6_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦6回读电流错误阈值")]
	MOD_ADDR_PUMP6_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦6电流采样滤波次数")]
	MOD_ADDR_PUMP6_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦7使能")]
	MOD_ADDR_PUMP7_EN /*泵浦使能*/
	//[Description("泵浦7开关")]
	MOD_ADDR_PUMP7_SW /*泵浦开关*/
	//[Description("泵浦7优先级设置")]
	MOD_ADDR_PUMP7_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦7目标电流值")]
	MOD_ADDR_PUMP7_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦7补偿系数")]
	MOD_ADDR_PUMP7_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦7电流设置比例系数")]
	MOD_ADDR_PUMP7_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦7电流速率调节ma/s")]
	MOD_ADDR_PUMP7_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦7电流上限")]
	MOD_ADDR_PUMP7_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦7电流回读adc通道绑定")]
	MOD_ADDR_PUMP7_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦7电流回读补偿值")]
	MOD_ADDR_PUMP7_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦7电流回读系数值")]
	MOD_ADDR_PUMP7_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦7回读电流错误阈值")]
	MOD_ADDR_PUMP7_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦7电流采样滤波次数")]
	MOD_ADDR_PUMP7_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦8使能")]
	MOD_ADDR_PUMP8_EN /*泵浦使能*/
	//[Description("泵浦8开关")]
	MOD_ADDR_PUMP8_SW /*泵浦开关*/
	//[Description("泵浦8优先级设置")]
	MOD_ADDR_PUMP8_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦8目标电流值")]
	MOD_ADDR_PUMP8_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦8补偿系数")]
	MOD_ADDR_PUMP8_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦8电流设置比例系数")]
	MOD_ADDR_PUMP8_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦8电流速率调节ma/s")]
	MOD_ADDR_PUMP8_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦8电流上限")]
	MOD_ADDR_PUMP8_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦8电流回读adc通道绑定")]
	MOD_ADDR_PUMP8_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦8电流回读补偿值")]
	MOD_ADDR_PUMP8_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦8电流回读系数值")]
	MOD_ADDR_PUMP8_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦8回读电流错误阈值")]
	MOD_ADDR_PUMP8_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦8电流采样滤波次数")]
	MOD_ADDR_PUMP8_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦9使能")]
	MOD_ADDR_PUMP9_EN /*泵浦使能*/
	//[Description("泵浦9开关")]
	MOD_ADDR_PUMP9_SW /*泵浦开关*/
	//[Description("泵浦9优先级设置")]
	MOD_ADDR_PUMP9_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦9目标电流值")]
	MOD_ADDR_PUMP9_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦9补偿系数")]
	MOD_ADDR_PUMP9_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦9电流设置比例系数")]
	MOD_ADDR_PUMP9_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦9电流速率调节ma/s")]
	MOD_ADDR_PUMP9_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦9电流上限")]
	MOD_ADDR_PUMP9_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦9电流回读adc通道绑定")]
	MOD_ADDR_PUMP9_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦9电流回读补偿值")]
	MOD_ADDR_PUMP9_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦9电流回读系数值")]
	MOD_ADDR_PUMP9_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦9回读电流错误阈值")]
	MOD_ADDR_PUMP9_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦9电流采样滤波次数")]
	MOD_ADDR_PUMP9_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦10使能")]
	MOD_ADDR_PUMP10_EN /*泵浦使能*/
	//[Description("泵浦10开关")]
	MOD_ADDR_PUMP10_SW /*泵浦开关*/
	//[Description("泵浦10优先级设置")]
	MOD_ADDR_PUMP10_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦10目标电流值")]
	MOD_ADDR_PUMP10_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦10补偿系数")]
	MOD_ADDR_PUMP10_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦10电流设置比例系数")]
	MOD_ADDR_PUMP10_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦10电流速率调节ma/s")]
	MOD_ADDR_PUMP10_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦10电流上限")]
	MOD_ADDR_PUMP10_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦10电流回读adc通道绑定")]
	MOD_ADDR_PUMP10_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦10电流回读补偿值")]
	MOD_ADDR_PUMP10_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦10电流回读系数值")]
	MOD_ADDR_PUMP10_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦10回读电流错误阈值")]
	MOD_ADDR_PUMP10_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦10电流采样滤波次数")]
	MOD_ADDR_PUMP10_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦11使能")]
	MOD_ADDR_PUMP11_EN /*泵浦使能*/
	//[Description("泵浦11开关")]
	MOD_ADDR_PUMP11_SW /*泵浦开关*/
	//[Description("泵浦11优先级设置")]
	MOD_ADDR_PUMP11_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦11目标电流值")]
	MOD_ADDR_PUMP11_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦11补偿系数")]
	MOD_ADDR_PUMP11_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦11电流设置比例系数")]
	MOD_ADDR_PUMP11_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦11电流速率调节ma/s")]
	MOD_ADDR_PUMP11_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/

	/*监控参数配置*/
	//[Description("泵浦11电流回读adc通道绑定")]
	MOD_ADDR_PUMP11_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦11电流回读补偿值")]
	MOD_ADDR_PUMP11_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦11电流回读系数值")]
	MOD_ADDR_PUMP11_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦11回读电流错误阈值")]
	MOD_ADDR_PUMP11_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦11电流采样滤波次数")]
	MOD_ADDR_PUMP11_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦12使能")]
	MOD_ADDR_PUMP12_EN /*泵浦使能*/
	//[Description("泵浦12开关")]
	MOD_ADDR_PUMP12_SW /*泵浦开关*/
	//[Description("泵浦12优先级设置")]
	MOD_ADDR_PUMP12_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦12目标电流值")]
	MOD_ADDR_PUMP12_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦12补偿系数")]
	MOD_ADDR_PUMP12_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦12电流设置比例系数")]
	MOD_ADDR_PUMP12_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦12电流速率调节ma/s")]
	MOD_ADDR_PUMP12_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦12电流上限")]
	MOD_ADDR_PUMP12_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦12电流回读adc通道绑定")]
	MOD_ADDR_PUMP12_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦12电流回读补偿值")]
	MOD_ADDR_PUMP12_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦12电流回读系数值")]
	MOD_ADDR_PUMP12_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦12回读电流错误阈值")]
	MOD_ADDR_PUMP12_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦12电流采样滤波次数")]
	MOD_ADDR_PUMP12_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦13使能")]
	MOD_ADDR_PUMP13_EN /*泵浦使能*/
	//[Description("泵浦13开关")]
	MOD_ADDR_PUMP13_SW /*泵浦开关*/
	//[Description("泵浦13优先级设置")]
	MOD_ADDR_PUMP13_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦13目标电流值")]
	MOD_ADDR_PUMP13_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦13补偿系数")]
	MOD_ADDR_PUMP13_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦13电流设置比例系数")]
	MOD_ADDR_PUMP13_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦13电流速率调节ma/s")]
	MOD_ADDR_PUMP13_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦13电流上限")]
	MOD_ADDR_PUMP13_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦13电流回读adc通道绑定")]
	MOD_ADDR_PUMP13_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦13电流回读补偿值")]
	MOD_ADDR_PUMP13_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦13电流回读系数值")]
	MOD_ADDR_PUMP13_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦13回读电流错误阈值")]
	MOD_ADDR_PUMP13_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦13电流采样滤波次数")]
	MOD_ADDR_PUMP13_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	//[Description("泵浦14使能")]
	MOD_ADDR_PUMP14_EN /*泵浦使能*/
	//[Description("泵浦14开关")]
	MOD_ADDR_PUMP14_SW /*泵浦开关*/
	//[Description("泵浦14优先级设置")]
	MOD_ADDR_PUMP14_SET_PRIORITY /*泵浦优先级*/
	//[Description("泵浦14目标电流值")]
	MOD_ADDR_PUMP14_SET_DEST_CUR_VALUE /*目标电流*/
	//[Description("泵浦14补偿系数")]
	MOD_ADDR_PUMP14_SET_COMPENSATION_VALUE /*电流设置补偿值*/
	//[Description("泵浦14电流设置比例系数")]
	MOD_ADDR_PUMP14_SET_COEFFICIENT_VALUE /*电流比例系数值*/
	//[Description("泵浦14电流速率调节ma/s")]
	MOD_ADDR_PUMP14_SET_CUR_SPEED_VALUE /*电流速率调节ma/s*/
	//[Description("泵浦14电流上限")]
	MOD_ADDR_PUMP14_SET_CUR_MAX /*电流设置上限*/

	/*监控参数配置*/
	//[Description("泵浦14电流回读adc通道绑定")]
	MOD_ADDR_PUMP14_MON_CUR_CH_BIND /*回读电流adc通道*/
	//[Description("泵浦14电流回读补偿值")]
	MOD_ADDR_PUMP14_MON_CUR_COMPENSATION_VALUE /*回读电流补偿值*/
	//[Description("泵浦14电流回读系数值")]
	MOD_ADDR_PUMP14_MON_CUR_COEFFICIENT_VALUE /*回读电流系数值*/
	//[Description("泵浦14回读电流错误阈值")]
	MOD_ADDR_PUMP14_MON_CUR_ERR_THR /*回读电流错误阈值*/
	//[Description("泵浦14电流采样滤波次数")]
	MOD_ADDR_PUMP14_MON_CUR_FILTER_TIME /*电流采样滤波次数*/

	/*############电机驱动定义##############*/

	//[Description("电机0使能")]
	MOD_ADDR_MOTOR0_EN
	//[Description("电机0驱动频率(Khz")]                                                         //
	MOD_ADDR_MOTOR0_FREQ
	//[Description("电机0目标坐标")]
	MOD_ADDR_MOTOR0_DEST_POS
	//[Description("电机0目标坐标")]
	MOD_ADDR_MOTOR0_POWER_UP_RESET
	//[Description("电机0坐标表0")]
	MOD_ADDR_MOTOR0_POS_TABLE_0
	//[Description("电机0坐标表1")]
	MOD_ADDR_MOTOR0_POS_TABLE_1
	//[Description("电机0坐标表2")]
	MOD_ADDR_MOTOR0_POS_TABLE_2
	//[Description("电机0坐标表3")]
	MOD_ADDR_MOTOR0_POS_TABLE_3
	//[Description("电机0坐标表4")]
	MOD_ADDR_MOTOR0_POS_TABLE_4
	//[Description("电机0坐标表5")]
	MOD_ADDR_MOTOR0_POS_TABLE_5
	//[Description("电机0坐标表6")]
	MOD_ADDR_MOTOR0_POS_TABLE_6
	//[Description("电机0坐标表7")]
	MOD_ADDR_MOTOR0_POS_TABLE_7
	//[Description("电机0坐标表8")]
	MOD_ADDR_MOTOR0_POS_TABLE_8
	//[Description("电机0坐标表9")]
	MOD_ADDR_MOTOR0_POS_TABLE_9
	//[Description("电机0坐标表10")]
	MOD_ADDR_MOTOR0_POS_TABLE_10
	//[Description("电机0坐标表11")]
	MOD_ADDR_MOTOR0_POS_TABLE_11
	//[Description("电机0坐标表12")]
	MOD_ADDR_MOTOR0_POS_TABLE_12
	//[Description("电机0坐标表13")]
	MOD_ADDR_MOTOR0_POS_TABLE_13
	//[Description("电机0坐标表14")]
	MOD_ADDR_MOTOR0_POS_TABLE_14
	//[Description("电机0坐标表15")]
	MOD_ADDR_MOTOR0_POS_TABLE_15
	//[Description("电机0坐标表16")]
	MOD_ADDR_MOTOR0_POS_TABLE_16
	//[Description("电机0坐标表17")]
	MOD_ADDR_MOTOR0_POS_TABLE_17
	//[Description("电机0坐标表18")]
	MOD_ADDR_MOTOR0_POS_TABLE_18
	//[Description("电机0坐标表19")]
	MOD_ADDR_MOTOR0_POS_TABLE_19

	//[Description("电机1使能")]
	MOD_ADDR_MOTOR1_EN
	//[Description("电机1驱动频率(Khz")]                                                         //
	MOD_ADDR_MOTOR1_FREQ
	//[Description("电机1目标坐标")]
	MOD_ADDR_MOTOR1_DEST_POS
	//[Description("电机1目标坐标")]
	MOD_ADDR_MOTOR1_POWER_UP_RESET
	//[Description("电机1坐标表0")]
	MOD_ADDR_MOTOR1_POS_TABLE_0
	//[Description("电机1坐标表1")]
	MOD_ADDR_MOTOR1_POS_TABLE_1
	//[Description("电机1坐标表2")]
	MOD_ADDR_MOTOR1_POS_TABLE_2
	//[Description("电机1坐标表3")]
	MOD_ADDR_MOTOR1_POS_TABLE_3
	//[Description("电机1坐标表4")]
	MOD_ADDR_MOTOR1_POS_TABLE_4
	//[Description("电机1坐标表5")]
	MOD_ADDR_MOTOR1_POS_TABLE_5
	//[Description("电机1坐标表6")]
	MOD_ADDR_MOTOR1_POS_TABLE_6
	//[Description("电机1坐标表7")]
	MOD_ADDR_MOTOR1_POS_TABLE_7
	//[Description("电机1坐标表8")]
	MOD_ADDR_MOTOR1_POS_TABLE_8
	//[Description("电机1坐标表9")]
	MOD_ADDR_MOTOR1_POS_TABLE_9
	//[Description("电机1坐标表10")]
	MOD_ADDR_MOTOR1_POS_TABLE_10
	//[Description("电机1坐标表11")]
	MOD_ADDR_MOTOR1_POS_TABLE_11
	//[Description("电机1坐标表12")]
	MOD_ADDR_MOTOR1_POS_TABLE_12
	//[Description("电机1坐标表13")]
	MOD_ADDR_MOTOR1_POS_TABLE_13
	//[Description("电机1坐标表14")]
	MOD_ADDR_MOTOR1_POS_TABLE_14
	//[Description("电机1坐标表15")]
	MOD_ADDR_MOTOR1_POS_TABLE_15
	//[Description("电机1坐标表16")]
	MOD_ADDR_MOTOR1_POS_TABLE_16
	//[Description("电机1坐标表17")]
	MOD_ADDR_MOTOR1_POS_TABLE_17
	//[Description("电机1坐标表18")]
	MOD_ADDR_MOTOR1_POS_TABLE_18
	//[Description("电机1坐标表19")]
	MOD_ADDR_MOTOR1_POS_TABLE_19

	//[Description("电机2使能")]
	MOD_ADDR_MOTOR2_EN
	//[Description("电机2驱动频率(Khz")]                                                         //
	MOD_ADDR_MOTOR2_FREQ
	//[Description("电机2目标坐标")]
	MOD_ADDR_MOTOR2_DEST_POS
	//[Description("电机2目标坐标")]
	MOD_ADDR_MOTOR2_POWER_UP_RESET
	//[Description("电机2坐标表0")]
	MOD_ADDR_MOTOR2_POS_TABLE_0
	//[Description("电机2坐标表1")]
	MOD_ADDR_MOTOR2_POS_TABLE_1
	//[Description("电机2坐标表2")]
	MOD_ADDR_MOTOR2_POS_TABLE_2
	//[Description("电机2坐标表3")]
	MOD_ADDR_MOTOR2_POS_TABLE_3
	//[Description("电机2坐标表4")]
	MOD_ADDR_MOTOR2_POS_TABLE_4
	//[Description("电机2坐标表5")]
	MOD_ADDR_MOTOR2_POS_TABLE_5
	//[Description("电机2坐标表6")]
	MOD_ADDR_MOTOR2_POS_TABLE_6
	//[Description("电机2坐标表7")]
	MOD_ADDR_MOTOR2_POS_TABLE_7
	//[Description("电机2坐标表8")]
	MOD_ADDR_MOTOR2_POS_TABLE_8
	//[Description("电机2坐标表9")]
	MOD_ADDR_MOTOR2_POS_TABLE_9
	//[Description("电机2坐标表10")]
	MOD_ADDR_MOTOR2_POS_TABLE_10
	//[Description("电机2坐标表11")]
	MOD_ADDR_MOTOR2_POS_TABLE_11
	//[Description("电机2坐标表12")]
	MOD_ADDR_MOTOR2_POS_TABLE_12
	//[Description("电机2坐标表13")]
	MOD_ADDR_MOTOR2_POS_TABLE_13
	//[Description("电机2坐标表14")]
	MOD_ADDR_MOTOR2_POS_TABLE_14
	//[Description("电机2坐标表15")]
	MOD_ADDR_MOTOR2_POS_TABLE_15
	//[Description("电机2坐标表16")]
	MOD_ADDR_MOTOR2_POS_TABLE_16
	//[Description("电机2坐标表17")]
	MOD_ADDR_MOTOR2_POS_TABLE_17
	//[Description("电机2坐标表18")]
	MOD_ADDR_MOTOR2_POS_TABLE_18
	//[Description("电机2坐标表19")]
	MOD_ADDR_MOTOR2_POS_TABLE_19

	//[Description("电机3使能")]
	MOD_ADDR_MOTOR3_EN
	//[Description("电机3驱动频率(Khz")]                                                         //
	MOD_ADDR_MOTOR3_FREQ
	//[Description("电机3目标坐标")]
	MOD_ADDR_MOTOR3_DEST_POS
	//[Description("电机3目标坐标")]
	MOD_ADDR_MOTOR3_POWER_UP_RESET
	//[Description("电机3坐标表0")]
	MOD_ADDR_MOTOR3_POS_TABLE_0
	//[Description("电机3坐标表1")]
	MOD_ADDR_MOTOR3_POS_TABLE_1
	//[Description("电机3坐标表3")]
	MOD_ADDR_MOTOR3_POS_TABLE_2
	//[Description("电机3坐标表3")]
	MOD_ADDR_MOTOR3_POS_TABLE_3
	//[Description("电机3坐标表4")]
	MOD_ADDR_MOTOR3_POS_TABLE_4
	//[Description("电机3坐标表5")]
	MOD_ADDR_MOTOR3_POS_TABLE_5
	//[Description("电机3坐标表6")]
	MOD_ADDR_MOTOR3_POS_TABLE_6
	//[Description("电机3坐标表7")]
	MOD_ADDR_MOTOR3_POS_TABLE_7
	//[Description("电机3坐标表8")]
	MOD_ADDR_MOTOR3_POS_TABLE_8
	//[Description("电机3坐标表9")]
	MOD_ADDR_MOTOR3_POS_TABLE_9
	//[Description("电机3坐标表10")]
	MOD_ADDR_MOTOR3_POS_TABLE_10
	//[Description("电机3坐标表11")]
	MOD_ADDR_MOTOR3_POS_TABLE_11
	//[Description("电机3坐标表13")]
	MOD_ADDR_MOTOR3_POS_TABLE_12
	//[Description("电机3坐标表13")]
	MOD_ADDR_MOTOR3_POS_TABLE_13
	//[Description("电机3坐标表14")]
	MOD_ADDR_MOTOR3_POS_TABLE_14
	//[Description("电机3坐标表15")]
	MOD_ADDR_MOTOR3_POS_TABLE_15
	//[Description("电机3坐标表16")]
	MOD_ADDR_MOTOR3_POS_TABLE_16
	//[Description("电机3坐标表17")]
	MOD_ADDR_MOTOR3_POS_TABLE_17
	//[Description("电机3坐标表18")]
	MOD_ADDR_MOTOR3_POS_TABLE_18
	//[Description("电机3坐标表19")]
	MOD_ADDR_MOTOR3_POS_TABLE_19

	/*数字温控模块*/

	//[Description("数字温控模块0使能")]
	MOD_ADDR_DIGI_TCM0_EN
	//[Description("数字温控模块0开关")]
	MOD_ADDR_DIGI_TCM0_SW
	//[Description("数字温控模块0类型")]
	MOD_ADDR_DIGI_TCM0_TYPE
	//[Description("数字温控模块0目标温度")]
	MOD_ADDR_DIGI_TCM0_DEST_TEMP
	//[Description("数字温控模块0 P参数")]
	MOD_ADDR_DIGI_TCM0_PID_P
	//[Description("数字温控模块0 I参数")]
	MOD_ADDR_DIGI_TCM0_PID_I
	//[Description("数字温控模块0 D参数")]
	MOD_ADDR_DIGI_TCM0_PID_D
	//[Description("数字温控模块0 温度上限")]
	MOD_ADDR_DIGI_TCM0_TEMP_MAX
	//[Description("数字温控模块0 温度下限")]
	MOD_ADDR_DIGI_TCM0_TEMP_MIN

	//[Description("数字温控模块1使能")]
	MOD_ADDR_DIGI_TCM1_EN
	//[Description("数字温控模块1开关")]
	MOD_ADDR_DIGI_TCM1_SW
	//[Description("数字温控模块1类型")]
	MOD_ADDR_DIGI_TCM1_TYPE
	//[Description("数字温控模块1目标温度")]
	MOD_ADDR_DIGI_TCM1_DEST_TEMP
	//[Description("数字温控模块1 P参数")]
	MOD_ADDR_DIGI_TCM1_PID_P
	//[Description("数字温控模块1 I参数")]
	MOD_ADDR_DIGI_TCM1_PID_I
	//[Description("数字温控模块1 D参数")]
	MOD_ADDR_DIGI_TCM1_PID_D
	//[Description("数字温控模块1 温度上限")]
	MOD_ADDR_DIGI_TCM1_TEMP_MAX
	//[Description("数字温控模块1 温度下限")]
	MOD_ADDR_DIGI_TCM1_TEMP_MIN

	//[Description("数字温控模块2使能")]
	MOD_ADDR_DIGI_TCM2_EN
	//[Description("数字温控模块2开关")]
	MOD_ADDR_DIGI_TCM2_SW
	//[Description("数字温控模块2类型")]
	MOD_ADDR_DIGI_TCM2_TYPE
	//[Description("数字温控模块2目标温度")]
	MOD_ADDR_DIGI_TCM2_DEST_TEMP
	//[Description("数字温控模块2 P参数")]
	MOD_ADDR_DIGI_TCM2_PID_P
	//[Description("数字温控模块2 I参数")]
	MOD_ADDR_DIGI_TCM2_PID_I
	//[Description("数字温控模块2 D参数")]
	MOD_ADDR_DIGI_TCM2_PID_D
	//[Description("数字温控模块2 温度上限")]
	MOD_ADDR_DIGI_TCM2_TEMP_MAX
	//[Description("数字温控模块2 温度下限")]
	MOD_ADDR_DIGI_TCM2_TEMP_MIN

	//[Description("数字温控模块3使能")]
	MOD_ADDR_DIGI_TCM3_EN
	//[Description("数字温控模块3开关")]
	MOD_ADDR_DIGI_TCM3_SW
	//[Description("数字温控模块3类型")]
	MOD_ADDR_DIGI_TCM3_TYPE
	//[Description("数字温控模块3目标温度")]
	MOD_ADDR_DIGI_TCM3_DEST_TEMP
	//[Description("数字温控模块3 P参数")]
	MOD_ADDR_DIGI_TCM3_PID_P
	//[Description("数字温控模块3 I参数")]
	MOD_ADDR_DIGI_TCM3_PID_I
	//[Description("数字温控模块3 D参数")]
	MOD_ADDR_DIGI_TCM3_PID_D
	//[Description("数字温控模块3 温度上限")]
	MOD_ADDR_DIGI_TCM3_TEMP_MAX
	//[Description("数字温控模块3 温度下限")]
	MOD_ADDR_DIGI_TCM3_TEMP_MIN

	//[Description("数字温控模块4使能")]
	MOD_ADDR_DIGI_TCM4_EN
	//[Description("数字温控模块4开关")]
	MOD_ADDR_DIGI_TCM4_SW
	//[Description("数字温控模块4类型")]
	MOD_ADDR_DIGI_TCM4_TYPE
	//[Description("数字温控模块4目标温度")]
	MOD_ADDR_DIGI_TCM4_DEST_TEMP
	//[Description("数字温控模块4 P参数")]
	MOD_ADDR_DIGI_TCM4_PID_P
	//[Description("数字温控模块4 I参数")]
	MOD_ADDR_DIGI_TCM4_PID_I
	//[Description("数字温控模块4 D参数")]
	MOD_ADDR_DIGI_TCM4_PID_D
	//[Description("数字温控模块4 温度上限")]
	MOD_ADDR_DIGI_TCM4_TEMP_MAX
	//[Description("数字温控模块4 温度下限")]
	MOD_ADDR_DIGI_TCM4_TEMP_MIN

	//[Description("数字温控模块5使能")]
	MOD_ADDR_DIGI_TCM5_EN
	//[Description("数字温控模块5开关")]
	MOD_ADDR_DIGI_TCM5_SW
	//[Description("数字温控模块5类型")]
	MOD_ADDR_DIGI_TCM5_TYPE
	//[Description("数字温控模块5目标温度")]
	MOD_ADDR_DIGI_TCM5_DEST_TEMP
	//[Description("数字温控模块5 P参数")]
	MOD_ADDR_DIGI_TCM5_PID_P
	//[Description("数字温控模块5 I参数")]
	MOD_ADDR_DIGI_TCM5_PID_I
	//[Description("数字温控模块5 D参数")]
	MOD_ADDR_DIGI_TCM5_PID_D
	//[Description("数字温控模块5 温度上限")]
	MOD_ADDR_DIGI_TCM5_TEMP_MAX
	//[Description("数字温控模块5 温度下限")]
	MOD_ADDR_DIGI_TCM5_TEMP_MIN

	/*DAC设置电压输出*/
	//[Description("DAC通道0设置电压(mv)")]
	MOD_ADDR_DAC_CH0_VAL
	//[Description("DAC通道1设置电压(mv)")]
	MOD_ADDR_DAC_CH1_VAL
	//[Description("DAC通道2设置电压(mv)")]
	MOD_ADDR_DAC_CH2_VAL
	//[Description("DAC通道3设置电压(mv)")]
	MOD_ADDR_DAC_CH3_VAL
	//[Description("DAC通道4设置电压(mv)")]
	MOD_ADDR_DAC_CH4_VAL
	//[Description("DAC通道5设置电压(mv)")]
	MOD_ADDR_DAC_CH5_VAL
	//[Description("DAC通道6设置电压(mv)")]
	MOD_ADDR_DAC_CH6_VAL
	//[Description("DAC通道7设置电压(mv)")]
	MOD_ADDR_DAC_CH7_VAL
	//[Description("DAC通道8设置电压(mv)")]
	MOD_ADDR_DAC_CH8_VAL
	//[Description("DAC通道9设置电压(mv)")]
	MOD_ADDR_DAC_CH9_VAL
	//[Description("DAC通道10设置电压(mv)")]
	MOD_ADDR_DAC_CH10_VAL
	//[Description("DAC通道11设置电压(mv)")]
	MOD_ADDR_DAC_CH11_VAL
	//[Description("DAC通道12设置电压(mv)")]
	MOD_ADDR_DAC_CH12_VAL
	//[Description("DAC通道13设置电压(mv)")]
	MOD_ADDR_DAC_CH13_VAL
	//[Description("DAC通道14设置电压(mv)")]
	MOD_ADDR_DAC_CH14_VAL
	//[Description("DAC通道15设置电压(mv)")]
	MOD_ADDR_DAC_CH15_VAL
	//[Description("DAC通道16设置电压(mv)")]
	MOD_ADDR_DAC_CH16_VAL
	//[Description("DAC通道17设置电压(mv)")]
	MOD_ADDR_DAC_CH17_VAL
	//[Description("DAC通道18设置电压(mv)")]
	MOD_ADDR_DAC_CH18_VAL
	//[Description("DAC通道19设置电压(mv)")]
	MOD_ADDR_DAC_CH19_VAL
	//[Description("DAC通道20设置电压(mv)")]
	MOD_ADDR_DAC_CH20_VAL
	//[Description("DAC通道21设置电压(mv)")]
	MOD_ADDR_DAC_CH21_VAL
	//[Description("DAC通道22设置电压(mv)")]
	MOD_ADDR_DAC_CH22_VAL
	//[Description("DAC通道23设置电压(mv)")]
	MOD_ADDR_DAC_CH23_VAL
	//[Description("DAC通道24设置电压(mv)")]
	MOD_ADDR_DAC_CH24_VAL
	//[Description("DAC通道25设置电压(mv)")]
	MOD_ADDR_DAC_CH25_VAL
	//[Description("DAC通道26设置电压(mv)")]
	MOD_ADDR_DAC_CH26_VAL
	//[Description("DAC通道27设置电压(mv)")]
	MOD_ADDR_DAC_CH27_VAL
	//[Description("DAC通道28设置电压(mv)")]
	MOD_ADDR_DAC_CH28_VAL
	//[Description("DAC通道29设置电压(mv)")]
	MOD_ADDR_DAC_CH29_VAL
	//[Description("DAC通道30设置电压(mv)")]
	MOD_ADDR_DAC_CH30_VAL
	//[Description("DAC通道31设置电压(mv)")]
	MOD_ADDR_DAC_CH31_VAL
	//[Description("DAC通道32设置电压(mv)")]
	MOD_ADDR_DAC_CH32_VAL
	//[Description("DAC通道33设置电压(mv)")]
	MOD_ADDR_DAC_CH33_VAL
	//[Description("DAC通道34设置电压(mv)")]
	MOD_ADDR_DAC_CH34_VAL
	//[Description("DAC通道35设置电压(mv)")]
	MOD_ADDR_DAC_CH35_VAL
	//[Description("DAC通道36设置电压(mv)")]
	MOD_ADDR_DAC_CH36_VAL
	//[Description("DAC通道37设置电压(mv)")]
	MOD_ADDR_DAC_CH37_VAL
	//[Description("DAC通道38设置电压(mv)")]
	MOD_ADDR_DAC_CH38_VAL
	//[Description("DAC通道39设置电压(mv)")]
	MOD_ADDR_DAC_CH39_VAL

	/*PD上下限制控制*/

	//[Description("PD0使能")]
	MOD_ADDR_PD0_EN
	//[Description("PD0频率上限")]
	MOD_ADDR_PD0_MAX_LIMIT
	//[Description("PD0频率下限")]
	MOD_ADDR_PD0_MIN_LIMIT

	//[Description("PD1使能")]
	MOD_ADDR_PD1_EN
	//[Description("PD1频率上限")]
	MOD_ADDR_PD1_MAX_LIMIT
	//[Description("PD1频率下限")]
	MOD_ADDR_PD1_MIN_LIMIT

	//[Description("PD2使能")]
	MOD_ADDR_PD2_EN
	//[Description("PD2频率上限")]
	MOD_ADDR_PD2_MAX_LIMIT
	//[Description("PD2频率下限")]
	MOD_ADDR_PD2_MIN_LIMIT

	//[Description("PD3使能")]
	MOD_ADDR_PD3_EN
	//[Description("PD3频率上限")]
	MOD_ADDR_PD3_MAX_LIMIT
	//[Description("PD3频率下限")]
	MOD_ADDR_PD3_MIN_LIMIT

	/*########声光设置###########*/
	//[Description("声光0使能")]
	MOD_ADDR_AOM0_EN
	//[Description("声光0寄存器值")]
	MOD_ADDR_AOM0_LEVEL
	//[Description("声光0延时参数1")]
	MOD_ADDR_AOM0_DELAY1
	//[Description("声光0延时参数2")]
	MOD_ADDR_AOM0_DELAY2
	//[Description("声光0PLL倍频值")]
	MOD_ADDR_AOM0_PLL

	//[Description("声光1使能")]
	MOD_ADDR_AOM1_EN
	//[Description("声光1寄存器值")]
	MOD_ADDR_AOM1_LEVEL
	//[Description("声光1延时参数1")]
	MOD_ADDR_AOM1_DELAY1
	//[Description("声光1延时参数2")]
	MOD_ADDR_AOM1_DELAY2
	//[Description("声光1PLL倍频值")]
	MOD_ADDR_AOM1_PLL

	//[Description("声光2使能")]
	MOD_ADDR_AOM2_EN
	//[Description("声光2寄存器值")]
	MOD_ADDR_AOM2_LEVEL
	//[Description("声光2延时参数1")]
	MOD_ADDR_AOM2_DELAY1
	//[Description("声光2延时参数2")]
	MOD_ADDR_AOM2_DELAY2
	//[Description("声光2PLL倍频值")]
	MOD_ADDR_AOM2_PLL

	//[Description("声光3使能")]
	MOD_ADDR_AOM3_EN
	//[Description("声光3寄存器值")]
	MOD_ADDR_AOM3_LEVEL
	//[Description("声光3延时参数1")]
	MOD_ADDR_AOM3_DELAY1
	//[Description("声光3延时参数2")]
	MOD_ADDR_AOM3_DELAY2
	//[Description("声光3PLL倍频值")]
	MOD_ADDR_AOM3_PLL

	//[Description("水冷机工作模式")]
	MOD_ADDR_WATER_COOLER_EN /* 0:关闭;1:使能*/
	//[Description("流量计上限")]
	MOD_ADDR_WATER_COOLER_MAX_LIMIT
	//[Description("流量计下限")]
	MOD_ADDR_WATER_COOLER_MIN_LIMIT
	//[Description("流量计补偿")]
	MOD_ADDR_WATER_COOLER_COMPENSATION

	/*########电压监控###########*/
	//[Description("电压监控0使能")]
	MOD_ADDR_VOL0_EN
	//[Description("电压监控0 ADC通道绑定")]
	MOD_ADDR_VOL0_CH_BIND
	//[Description("电压监控0 转化补偿值")]
	MOD_ADDR_VOL0_COMPENSATION_VALUE
	//[Description("电压监控0 系数值")]
	MOD_ADDR_VOL0_COEFFICIENT_VALUE
	//[Description("电压监控0 电压上限")]
	MOD_ADDR_VOL0_MAX
	//[Description("电压监控0 电压下限")]
	MOD_ADDR_VOL0_MIN
	//[Description("电压监控0 检测滤波次数")]
	MOD_ADDR_VOL0_MON_FIL_TIME

	//[Description("电压监控1使能")]
	MOD_ADDR_VOL1_EN
	//[Description("电压监控1 ADC通道绑定")]
	MOD_ADDR_VOL1_CH_BIND
	//[Description("电压监控1 转化补偿值")]
	MOD_ADDR_VOL1_COMPENSATION_VALUE
	//[Description("电压监控1 系数值")]
	MOD_ADDR_VOL1_COEFFICIENT_VALUE
	//[Description("电压监控1 电压上限")]
	MOD_ADDR_VOL1_MAX
	//[Description("电压监控1 电压下限")]
	MOD_ADDR_VOL1_MIN
	//[Description("电压监控1 检测滤波次数")]
	MOD_ADDR_VOL1_MON_FIL_TIME

	//[Description("电压监控2使能")]
	MOD_ADDR_VOL2_EN
	//[Description("电压监控2 ADC通道绑定")]
	MOD_ADDR_VOL2_CH_BIND
	//[Description("电压监控2 转化补偿值")]
	MOD_ADDR_VOL2_COMPENSATION_VALUE
	//[Description("电压监控2 系数值")]
	MOD_ADDR_VOL2_COEFFICIENT_VALUE
	//[Description("电压监控2 电压上限")]
	MOD_ADDR_VOL2_MAX
	//[Description("电压监控2 电压下限")]
	MOD_ADDR_VOL2_MIN
	//[Description("电压监控2 检测滤波次数")]
	MOD_ADDR_VOL2_MON_FIL_TIME

	//[Description("电压监控3使能")]
	MOD_ADDR_VOL3_EN
	//[Description("电压监控3 ADC通道绑定")]
	MOD_ADDR_VOL3_CH_BIND
	//[Description("电压监控3 转化补偿值")]
	MOD_ADDR_VOL3_COMPENSATION_VALUE
	//[Description("电压监控3 系数值")]
	MOD_ADDR_VOL3_COEFFICIENT_VALUE
	//[Description("电压监控3 电压上限")]
	MOD_ADDR_VOL3_MAX
	//[Description("电压监控3 电压下限")]
	MOD_ADDR_VOL3_MIN
	//[Description("电压监控3 检测滤波次数")]
	MOD_ADDR_VOL3_MON_FIL_TIME

	//[Description("电压监控4使能")]
	MOD_ADDR_VOL4_EN
	//[Description("电压监控4 ADC通道绑定")]
	MOD_ADDR_VOL4_CH_BIND
	//[Description("电压监控4 转化补偿值")]
	MOD_ADDR_VOL4_COMPENSATION_VALUE
	//[Description("电压监控4 系数值")]
	MOD_ADDR_VOL4_COEFFICIENT_VALUE
	//[Description("电压监控4 电压上限")]
	MOD_ADDR_VOL4_MAX
	//[Description("电压监控4 电压下限")]
	MOD_ADDR_VOL4_MIN
	//[Description("电压监控4 检测滤波次数")]
	MOD_ADDR_VOL4_MON_FIL_TIME

	//[Description("电压监控5使能")]
	MOD_ADDR_VOL5_EN
	//[Description("电压监控5 ADC通道绑定")]
	MOD_ADDR_VOL5_CH_BIND
	//[Description("电压监控5 转化补偿值")]
	MOD_ADDR_VOL5_COMPENSATION_VALUE
	//[Description("电压监控5 系数值")]
	MOD_ADDR_VOL5_COEFFICIENT_VALUE
	//[Description("电压监控5 电压上限")]
	MOD_ADDR_VOL5_MAX
	//[Description("电压监控5 电压下限")]
	MOD_ADDR_VOL5_MIN
	//[Description("电压监控5 检测滤波次数")]
	MOD_ADDR_VOL5_MON_FIL_TIME

	//[Description("电压监控6使能")]
	MOD_ADDR_VOL6_EN
	//[Description("电压监控6 ADC通道绑定")]
	MOD_ADDR_VOL6_CH_BIND
	//[Description("电压监控6 转化补偿值")]
	MOD_ADDR_VOL6_COMPENSATION_VALUE
	//[Description("电压监控6 系数值")]
	MOD_ADDR_VOL6_COEFFICIENT_VALUE
	//[Description("电压监控6 电压上限")]
	MOD_ADDR_VOL6_MAX
	//[Description("电压监控6 电压下限")]
	MOD_ADDR_VOL6_MIN
	//[Description("电压监控6 检测滤波次数")]
	MOD_ADDR_VOL6_MON_FIL_TIME

	//[Description("电压监控7使能")]
	MOD_ADDR_VOL7_EN
	//[Description("电压监控7 ADC通道绑定")]
	MOD_ADDR_VOL7_CH_BIND
	//[Description("电压监控7 转化补偿值")]
	MOD_ADDR_VOL7_COMPENSATION_VALUE
	//[Description("电压监控7 系数值")]
	MOD_ADDR_VOL7_COEFFICIENT_VALUE
	//[Description("电压监控7 电压上限")]
	MOD_ADDR_VOL7_MAX
	//[Description("电压监控7 电压下限")]
	MOD_ADDR_VOL7_MIN
	//[Description("电压监控7 检测滤波次数")]
	MOD_ADDR_VOL7_MON_FIL_TIME

	//[Description("电压监控8使能")]
	MOD_ADDR_VOL8_EN
	//[Description("电压监控8 ADC通道绑定")]
	MOD_ADDR_VOL8_CH_BIND
	//[Description("电压监控8 转化补偿值")]
	MOD_ADDR_VOL8_COMPENSATION_VALUE
	//[Description("电压监控8 系数值")]
	MOD_ADDR_VOL8_COEFFICIENT_VALUE
	//[Description("电压监控8 电压上限")]
	MOD_ADDR_VOL8_MAX
	//[Description("电压监控8 电压下限")]
	MOD_ADDR_VOL8_MIN
	//[Description("电压监控8 检测滤波次数")]
	MOD_ADDR_VOL8_MON_FIL_TIME

	//[Description("电压监控9使能")]
	MOD_ADDR_VOL9_EN
	//[Description("电压监控9 ADC通道绑定")]
	MOD_ADDR_VOL9_CH_BIND
	//[Description("电压监控9 转化补偿值")]
	MOD_ADDR_VOL9_COMPENSATION_VALUE
	//[Description("电压监控9 系数值")]
	MOD_ADDR_VOL9_COEFFICIENT_VALUE
	//[Description("电压监控9 电压上限")]
	MOD_ADDR_VOL9_MAX
	//[Description("电压监控9 电压下限")]
	MOD_ADDR_VOL9_MIN
	//[Description("电压监控9 检测滤波次数")]
	MOD_ADDR_VOL9_MON_FIL_TIME

	//[Description("电压监控10使能")]
	MOD_ADDR_VOL10_EN
	//[Description("电压监控10 ADC通道绑定")]
	MOD_ADDR_VOL10_CH_BIND
	//[Description("电压监控10 转化补偿值")]
	MOD_ADDR_VOL10_COMPENSATION_VALUE
	//[Description("电压监控10 系数值")]
	MOD_ADDR_VOL10_COEFFICIENT_VALUE
	//[Description("电压监控10 电压上限")]
	MOD_ADDR_VOL10_MAX
	//[Description("电压监控10 电压下限")]
	MOD_ADDR_VOL10_MIN
	//[Description("电压监控10 检测滤波次数")]
	MOD_ADDR_VOL10_MON_FIL_TIME

	//[Description("电压监控11使能")]
	MOD_ADDR_VOL11_EN
	//[Description("电压监控11 ADC通道绑定")]
	MOD_ADDR_VOL11_CH_BIND
	//[Description("电压监控11 转化补偿值")]
	MOD_ADDR_VOL11_COMPENSATION_VALUE
	//[Description("电压监控11 系数值")]
	MOD_ADDR_VOL11_COEFFICIENT_VALUE
	//[Description("电压监控11 电压上限")]
	MOD_ADDR_VOL11_MAX
	//[Description("电压监控11 电压下限")]
	MOD_ADDR_VOL11_MIN
	//[Description("电压监控11 检测滤波次数")]
	MOD_ADDR_VOL11_MON_FIL_TIME

	//[Description("电压监控12使能")]
	MOD_ADDR_VOL12_EN
	//[Description("电压监控12 ADC通道绑定")]
	MOD_ADDR_VOL12_CH_BIND
	//[Description("电压监控12 转化补偿值")]
	MOD_ADDR_VOL12_COMPENSATION_VALUE
	//[Description("电压监控12 系数值")]
	MOD_ADDR_VOL12_COEFFICIENT_VALUE
	//[Description("电压监控12 电压上限")]
	MOD_ADDR_VOL12_MAX
	//[Description("电压监控12 电压下限")]
	MOD_ADDR_VOL12_MIN
	//[Description("电压监控12 检测滤波次数")]
	MOD_ADDR_VOL12_MON_FIL_TIME

	//[Description("电压监控13使能")]
	MOD_ADDR_VOL13_EN
	//[Description("电压监控13 ADC通道绑定")]
	MOD_ADDR_VOL13_CH_BIND
	//[Description("电压监控13 转化补偿值")]
	MOD_ADDR_VOL13_COMPENSATION_VALUE
	//[Description("电压监控13 系数值")]
	MOD_ADDR_VOL13_COEFFICIENT_VALUE
	//[Description("电压监控13 电压上限")]
	MOD_ADDR_VOL13_MAX
	//[Description("电压监控13 电压下限")]
	MOD_ADDR_VOL13_MIN
	//[Description("电压监控13 检测滤波次数")]
	MOD_ADDR_VOL13_MON_FIL_TIME

	//[Description("电压监控14使能")]
	MOD_ADDR_VOL14_EN
	//[Description("电压监控14 ADC通道绑定")]
	MOD_ADDR_VOL14_CH_BIND
	//[Description("电压监控14 转化补偿值")]
	MOD_ADDR_VOL14_COMPENSATION_VALUE
	//[Description("电压监控14 系数值")]
	MOD_ADDR_VOL14_COEFFICIENT_VALUE
	//[Description("电压监控14 电压上限")]
	MOD_ADDR_VOL14_MAX
	//[Description("电压监控14 电压下限")]
	MOD_ADDR_VOL14_MIN
	//[Description("电压监控14 检测滤波次数")]
	MOD_ADDR_VOL14_MON_FIL_TIME

	//[Description("电压监控15使能")]
	MOD_ADDR_VOL15_EN
	//[Description("电压监控15 ADC通道绑定")]
	MOD_ADDR_VOL15_CH_BIND
	//[Description("电压监控15 转化补偿值")]
	MOD_ADDR_VOL15_COMPENSATION_VALUE
	//[Description("电压监控15 系数值")]
	MOD_ADDR_VOL15_COEFFICIENT_VALUE
	//[Description("电压监控15 电压上限")]
	MOD_ADDR_VOL15_MAX
	//[Description("电压监控15 电压下限")]
	MOD_ADDR_VOL15_MIN
	//[Description("电压监控15 检测滤波次数")]
	MOD_ADDR_VOL15_MON_FIL_TIME

	//[Description("电压监控16使能")]
	MOD_ADDR_VOL16_EN
	//[Description("电压监控16 ADC通道绑定")]
	MOD_ADDR_VOL16_CH_BIND
	//[Description("电压监控16 转化补偿值")]
	MOD_ADDR_VOL16_COMPENSATION_VALUE
	//[Description("电压监控16 系数值")]
	MOD_ADDR_VOL16_COEFFICIENT_VALUE
	//[Description("电压监控16 电压上限")]
	MOD_ADDR_VOL16_MAX
	//[Description("电压监控16 电压下限")]
	MOD_ADDR_VOL16_MIN
	//[Description("电压监控16 检测滤波次数")]
	MOD_ADDR_VOL16_MON_FIL_TIME

	//[Description("电压监控17使能")]
	MOD_ADDR_VOL17_EN
	//[Description("电压监控17 ADC通道绑定")]
	MOD_ADDR_VOL17_CH_BIND
	//[Description("电压监控17 转化补偿值")]
	MOD_ADDR_VOL17_COMPENSATION_VALUE
	//[Description("电压监控17 系数值")]
	MOD_ADDR_VOL17_COEFFICIENT_VALUE
	//[Description("电压监控17 电压上限")]
	MOD_ADDR_VOL17_MAX
	//[Description("电压监控17 电压下限")]
	MOD_ADDR_VOL17_MIN
	//[Description("电压监控17 检测滤波次数")]
	MOD_ADDR_VOL17_MON_FIL_TIME

	//[Description("电压监控18使能")]
	MOD_ADDR_VOL18_EN
	//[Description("电压监控18 ADC通道绑定")]
	MOD_ADDR_VOL18_CH_BIND
	//[Description("电压监控18 转化补偿值")]
	MOD_ADDR_VOL18_COMPENSATION_VALUE
	//[Description("电压监控18 系数值")]
	MOD_ADDR_VOL18_COEFFICIENT_VALUE
	//[Description("电压监控18 电压上限")]
	MOD_ADDR_VOL18_MAX
	//[Description("电压监控18 电压下限")]
	MOD_ADDR_VOL18_MIN
	//[Description("电压监控18 检测滤波次数")]
	MOD_ADDR_VOL18_MON_FIL_TIME

	//[Description("电压监控19使能")]
	MOD_ADDR_VOL19_EN
	//[Description("电压监控19 ADC通道绑定")]
	MOD_ADDR_VOL19_CH_BIND
	//[Description("电压监控19 转化补偿值")]
	MOD_ADDR_VOL19_COMPENSATION_VALUE
	//[Description("电压监控19 系数值")]
	MOD_ADDR_VOL19_COEFFICIENT_VALUE
	//[Description("电压监控19 电压上限")]
	MOD_ADDR_VOL19_MAX
	//[Description("电压监控19 电压下限")]
	MOD_ADDR_VOL19_MIN
	//[Description("电压监控19 检测滤波次数")]
	MOD_ADDR_VOL19_MON_FIL_TIME

	/*#########温度监控配置##############*/
	//[Description("温度监控0 使能")]
	MOD_ADDR_TEMP0_EN
	//[Description("温度监控0 类型")]
	MOD_ADDR_TEMP0_TYPE
	//[Description("温度监控0 通道绑定")]
	MOD_ADDR_TEMP0_CH_BIND
	//[Description("温度监控0 滤波次数")]
	MOD_ADDR_TEMP0_MON_FLT_TIME
	//[Description("温度监控0 参数1")]
	MOD_ADDR_TEMP0_PARA1
	//[Description("温度监控0 参数2")]
	MOD_ADDR_TEMP0_PARA2
	//[Description("温度监控0 参数3")]
	MOD_ADDR_TEMP0_PARA3
	//[Description("温度监控0 参数4")]
	MOD_ADDR_TEMP0_PARA4
	//[Description("温度监控0 参数5")]
	MOD_ADDR_TEMP0_PARA5
	//[Description("温度监控0 参数6")]
	MOD_ADDR_TEMP0_PARA6
	//[Description("温度监控0 检测上限")]
	MOD_ADDR_TEMP0_MAX
	//[Description("温度监控0 检测下限")]
	MOD_ADDR_TEMP0_MIN

	//[Description("温度监控1 使能")]
	MOD_ADDR_TEMP1_EN
	//[Description("温度监控1 类型")]
	MOD_ADDR_TEMP1_TYPE
	//[Description("温度监控1 通道绑定")]
	MOD_ADDR_TEMP1_CH_BIND
	//[Description("温度监控1 滤波次数")]
	MOD_ADDR_TEMP1_MON_FLT_TIME
	//[Description("温度监控1 参数1")]
	MOD_ADDR_TEMP1_PARA1
	//[Description("温度监控1 参数2")]
	MOD_ADDR_TEMP1_PARA2
	//[Description("温度监控1 参数3")]
	MOD_ADDR_TEMP1_PARA3
	//[Description("温度监控1 参数4")]
	MOD_ADDR_TEMP1_PARA4
	//[Description("温度监控1 参数5")]
	MOD_ADDR_TEMP1_PARA5
	//[Description("温度监控1 参数6")]
	MOD_ADDR_TEMP1_PARA6
	//[Description("温度监控1 检测上限")]
	MOD_ADDR_TEMP1_MAX
	//[Description("温度监控1 检测下限")]
	MOD_ADDR_TEMP1_MIN

	//[Description("温度监控2 使能")]
	MOD_ADDR_TEMP2_EN
	//[Description("温度监控2 类型")]
	MOD_ADDR_TEMP2_TYPE
	//[Description("温度监控2 通道绑定")]
	MOD_ADDR_TEMP2_CH_BIND
	//[Description("温度监控2 滤波次数")]
	MOD_ADDR_TEMP2_MON_FLT_TIME
	//[Description("温度监控2 参数1")]
	MOD_ADDR_TEMP2_PARA1
	//[Description("温度监控2 参数2")]
	MOD_ADDR_TEMP2_PARA2
	//[Description("温度监控2 参数3")]
	MOD_ADDR_TEMP2_PARA3
	//[Description("温度监控2 参数4")]
	MOD_ADDR_TEMP2_PARA4
	//[Description("温度监控2 参数5")]
	MOD_ADDR_TEMP2_PARA5
	//[Description("温度监控2 参数6")]
	MOD_ADDR_TEMP2_PARA6
	//[Description("温度监控2 检测上限")]
	MOD_ADDR_TEMP2_MAX
	//[Description("温度监控2 检测下限")]
	MOD_ADDR_TEMP2_MIN

	//[Description("温度监控3 使能")]
	MOD_ADDR_TEMP3_EN
	//[Description("温度监控3 类型")]
	MOD_ADDR_TEMP3_TYPE
	//[Description("温度监控3 通道绑定")]
	MOD_ADDR_TEMP3_CH_BIND
	//[Description("温度监控3 滤波次数")]
	MOD_ADDR_TEMP3_MON_FLT_TIME
	//[Description("温度监控3 参数1")]
	MOD_ADDR_TEMP3_PARA1
	//[Description("温度监控3 参数2")]
	MOD_ADDR_TEMP3_PARA2
	//[Description("温度监控3 参数3")]
	MOD_ADDR_TEMP3_PARA3
	//[Description("温度监控3 参数4")]
	MOD_ADDR_TEMP3_PARA4
	//[Description("温度监控3 参数5")]
	MOD_ADDR_TEMP3_PARA5
	//[Description("温度监控3 参数6")]
	MOD_ADDR_TEMP3_PARA6
	//[Description("温度监控3 检测上限")]
	MOD_ADDR_TEMP3_MAX
	//[Description("温度监控3 检测下限")]
	MOD_ADDR_TEMP3_MIN

	//[Description("温度监控4 使能")]
	MOD_ADDR_TEMP4_EN
	//[Description("温度监控4 类型")]
	MOD_ADDR_TEMP4_TYPE
	//[Description("温度监控4 通道绑定")]
	MOD_ADDR_TEMP4_CH_BIND
	//[Description("温度监控4 滤波次数")]
	MOD_ADDR_TEMP4_MON_FLT_TIME
	//[Description("温度监控4 参数1")]
	MOD_ADDR_TEMP4_PARA1
	//[Description("温度监控4 参数2")]
	MOD_ADDR_TEMP4_PARA2
	//[Description("温度监控4 参数3")]
	MOD_ADDR_TEMP4_PARA3
	//[Description("温度监控4 参数4")]
	MOD_ADDR_TEMP4_PARA4
	//[Description("温度监控4 参数5")]
	MOD_ADDR_TEMP4_PARA5
	//[Description("温度监控4 参数6")]
	MOD_ADDR_TEMP4_PARA6
	//[Description("温度监控4 检测上限")]
	MOD_ADDR_TEMP4_MAX
	//[Description("温度监控4 检测下限")]
	MOD_ADDR_TEMP4_MIN

	//[Description("温度监控5 使能")]
	MOD_ADDR_TEMP5_EN
	//[Description("温度监控5 类型")]
	MOD_ADDR_TEMP5_TYPE
	//[Description("温度监控5 通道绑定")]
	MOD_ADDR_TEMP5_CH_BIND
	//[Description("温度监控5 滤波次数")]
	MOD_ADDR_TEMP5_MON_FLT_TIME
	//[Description("温度监控5 参数1")]
	MOD_ADDR_TEMP5_PARA1
	//[Description("温度监控5 参数2")]
	MOD_ADDR_TEMP5_PARA2
	//[Description("温度监控5 参数3")]
	MOD_ADDR_TEMP5_PARA3
	//[Description("温度监控5 参数4")]
	MOD_ADDR_TEMP5_PARA4
	//[Description("温度监控5 参数5")]
	MOD_ADDR_TEMP5_PARA5
	//[Description("温度监控5 参数6")]
	MOD_ADDR_TEMP5_PARA6
	//[Description("温度监控5 检测上限")]
	MOD_ADDR_TEMP5_MAX
	//[Description("温度监控5 检测下限")]
	MOD_ADDR_TEMP5_MIN

	//[Description("温度监控6 使能")]
	MOD_ADDR_TEMP6_EN
	//[Description("温度监控6 类型")]
	MOD_ADDR_TEMP6_TYPE
	//[Description("温度监控6 通道绑定")]
	MOD_ADDR_TEMP6_CH_BIND
	//[Description("温度监控6 滤波次数")]
	MOD_ADDR_TEMP6_MON_FLT_TIME
	//[Description("温度监控6 参数1")]
	MOD_ADDR_TEMP6_PARA1
	//[Description("温度监控6 参数2")]
	MOD_ADDR_TEMP6_PARA2
	//[Description("温度监控6 参数3")]
	MOD_ADDR_TEMP6_PARA3
	//[Description("温度监控6 参数4")]
	MOD_ADDR_TEMP6_PARA4
	//[Description("温度监控6 参数5")]
	MOD_ADDR_TEMP6_PARA5
	//[Description("温度监控6 参数6")]
	MOD_ADDR_TEMP6_PARA6
	//[Description("温度监控6 检测上限")]
	MOD_ADDR_TEMP6_MAX
	//[Description("温度监控6 检测下限")]
	MOD_ADDR_TEMP6_MIN

	//[Description("温度监控7 使能")]
	MOD_ADDR_TEMP7_EN
	//[Description("温度监控7 类型")]
	MOD_ADDR_TEMP7_TYPE
	//[Description("温度监控7 通道绑定")]
	MOD_ADDR_TEMP7_CH_BIND
	//[Description("温度监控7 滤波次数")]
	MOD_ADDR_TEMP7_MON_FLT_TIME
	//[Description("温度监控7 参数1")]
	MOD_ADDR_TEMP7_PARA1
	//[Description("温度监控7 参数2")]
	MOD_ADDR_TEMP7_PARA2
	//[Description("温度监控7 参数3")]
	MOD_ADDR_TEMP7_PARA3
	//[Description("温度监控7 参数4")]
	MOD_ADDR_TEMP7_PARA4
	//[Description("温度监控7 参数5")]
	MOD_ADDR_TEMP7_PARA5
	//[Description("温度监控7 参数6")]
	MOD_ADDR_TEMP7_PARA6
	//[Description("温度监控7 检测上限")]
	MOD_ADDR_TEMP7_MAX
	//[Description("温度监控7 检测下限")]
	MOD_ADDR_TEMP7_MIN

	//[Description("温度监控8 使能")]
	MOD_ADDR_TEMP8_EN
	//[Description("温度监控8 类型")]
	MOD_ADDR_TEMP8_TYPE
	//[Description("温度监控8 通道绑定")]
	MOD_ADDR_TEMP8_CH_BIND
	//[Description("温度监控8 滤波次数")]
	MOD_ADDR_TEMP8_MON_FLT_TIME
	//[Description("温度监控8 参数1")]
	MOD_ADDR_TEMP8_PARA1
	//[Description("温度监控8 参数2")]
	MOD_ADDR_TEMP8_PARA2
	//[Description("温度监控8 参数3")]
	MOD_ADDR_TEMP8_PARA3
	//[Description("温度监控8 参数4")]
	MOD_ADDR_TEMP8_PARA4
	//[Description("温度监控8 参数5")]
	MOD_ADDR_TEMP8_PARA5
	//[Description("温度监控8 参数6")]
	MOD_ADDR_TEMP8_PARA6
	//[Description("温度监控8 检测上限")]
	MOD_ADDR_TEMP8_MAX
	//[Description("温度监控8 检测下限")]
	MOD_ADDR_TEMP8_MIN

	//[Description("温度监控9 使能")]
	MOD_ADDR_TEMP9_EN
	//[Description("温度监控9 类型")]
	MOD_ADDR_TEMP9_TYPE
	//[Description("温度监控9 通道绑定")]
	MOD_ADDR_TEMP9_CH_BIND
	//[Description("温度监控9 滤波次数")]
	MOD_ADDR_TEMP9_MON_FLT_TIME
	//[Description("温度监控9 参数1")]
	MOD_ADDR_TEMP9_PARA1
	//[Description("温度监控9 参数2")]
	MOD_ADDR_TEMP9_PARA2
	//[Description("温度监控9 参数3")]
	MOD_ADDR_TEMP9_PARA3
	//[Description("温度监控9 参数4")]
	MOD_ADDR_TEMP9_PARA4
	//[Description("温度监控9 参数5")]
	MOD_ADDR_TEMP9_PARA5
	//[Description("温度监控9 参数6")]
	MOD_ADDR_TEMP9_PARA6
	//[Description("温度监控9 检测上限")]
	MOD_ADDR_TEMP9_MAX
	//[Description("温度监控9 检测下限")]
	MOD_ADDR_TEMP9_MIN

	//[Description("温度监控10 使能")]
	MOD_ADDR_TEMP10_EN
	//[Description("温度监控10 类型")]
	MOD_ADDR_TEMP10_TYPE
	//[Description("温度监控10 通道绑定")]
	MOD_ADDR_TEMP10_CH_BIND
	//[Description("温度监控10 滤波次数")]
	MOD_ADDR_TEMP10_MON_FLT_TIME
	//[Description("温度监控10 参数1")]
	MOD_ADDR_TEMP10_PARA1
	//[Description("温度监控10 参数2")]
	MOD_ADDR_TEMP10_PARA2
	//[Description("温度监控10 参数3")]
	MOD_ADDR_TEMP10_PARA3
	//[Description("温度监控10 参数4")]
	MOD_ADDR_TEMP10_PARA4
	//[Description("温度监控10 参数5")]
	MOD_ADDR_TEMP10_PARA5
	//[Description("温度监控10 参数6")]
	MOD_ADDR_TEMP10_PARA6
	//[Description("温度监控10 检测上限")]
	MOD_ADDR_TEMP10_MAX
	//[Description("温度监控10 检测下限")]
	MOD_ADDR_TEMP10_MIN

	//[Description("温度监控11 使能")]
	MOD_ADDR_TEMP11_EN
	//[Description("温度监控11 类型")]
	MOD_ADDR_TEMP11_TYPE
	//[Description("温度监控11 通道绑定")]
	MOD_ADDR_TEMP11_CH_BIND
	//[Description("温度监控11 滤波次数")]
	MOD_ADDR_TEMP11_MON_FLT_TIME
	//[Description("温度监控11 参数1")]
	MOD_ADDR_TEMP11_PARA1
	//[Description("温度监控11 参数2")]
	MOD_ADDR_TEMP11_PARA2
	//[Description("温度监控11 参数3")]
	MOD_ADDR_TEMP11_PARA3
	//[Description("温度监控11 参数4")]
	MOD_ADDR_TEMP11_PARA4
	//[Description("温度监控11 参数5")]
	MOD_ADDR_TEMP11_PARA5
	//[Description("温度监控11 参数6")]
	MOD_ADDR_TEMP11_PARA6
	//[Description("温度监控11 检测上限")]
	MOD_ADDR_TEMP11_MAX
	//[Description("温度监控11 检测下限")]
	MOD_ADDR_TEMP11_MIN

	//[Description("温度监控12 使能")]
	MOD_ADDR_TEMP12_EN
	//[Description("温度监控12 类型")]
	MOD_ADDR_TEMP12_TYPE
	//[Description("温度监控12 通道绑定")]
	MOD_ADDR_TEMP12_CH_BIND
	//[Description("温度监控12 滤波次数")]
	MOD_ADDR_TEMP12_MON_FLT_TIME
	//[Description("温度监控12 参数1")]
	MOD_ADDR_TEMP12_PARA1
	//[Description("温度监控12 参数2")]
	MOD_ADDR_TEMP12_PARA2
	//[Description("温度监控12 参数3")]
	MOD_ADDR_TEMP12_PARA3
	//[Description("温度监控12 参数4")]
	MOD_ADDR_TEMP12_PARA4
	//[Description("温度监控12 参数5")]
	MOD_ADDR_TEMP12_PARA5
	//[Description("温度监控12 参数6")]
	MOD_ADDR_TEMP12_PARA6
	//[Description("温度监控12 检测上限")]
	MOD_ADDR_TEMP12_MAX
	//[Description("温度监控12 检测下限")]
	MOD_ADDR_TEMP12_MIN

	//[Description("温度监控13 使能")]
	MOD_ADDR_TEMP13_EN
	//[Description("温度监控13 类型")]
	MOD_ADDR_TEMP13_TYPE
	//[Description("温度监控13 通道绑定")]
	MOD_ADDR_TEMP13_CH_BIND
	//[Description("温度监控13 滤波次数")]
	MOD_ADDR_TEMP13_MON_FLT_TIME
	//[Description("温度监控13 参数1")]
	MOD_ADDR_TEMP13_PARA1
	//[Description("温度监控13 参数2")]
	MOD_ADDR_TEMP13_PARA2
	//[Description("温度监控13 参数3")]
	MOD_ADDR_TEMP13_PARA3
	//[Description("温度监控13 参数4")]
	MOD_ADDR_TEMP13_PARA4
	//[Description("温度监控13 参数5")]
	MOD_ADDR_TEMP13_PARA5
	//[Description("温度监控13 参数6")]
	MOD_ADDR_TEMP13_PARA6
	//[Description("温度监控13 检测上限")]
	MOD_ADDR_TEMP13_MAX
	//[Description("温度监控13 检测下限")]
	MOD_ADDR_TEMP13_MIN

	//[Description("温度监控14 使能")]
	MOD_ADDR_TEMP14_EN
	//[Description("温度监控14 类型")]
	MOD_ADDR_TEMP14_TYPE
	//[Description("温度监控14 通道绑定")]
	MOD_ADDR_TEMP14_CH_BIND
	//[Description("温度监控14 滤波次数")]
	MOD_ADDR_TEMP14_MON_FLT_TIME
	//[Description("温度监控14 参数1")]
	MOD_ADDR_TEMP14_PARA1
	//[Description("温度监控14 参数2")]
	MOD_ADDR_TEMP14_PARA2
	//[Description("温度监控14 参数3")]
	MOD_ADDR_TEMP14_PARA3
	//[Description("温度监控14 参数4")]
	MOD_ADDR_TEMP14_PARA4
	//[Description("温度监控14 参数5")]
	MOD_ADDR_TEMP14_PARA5
	//[Description("温度监控14 参数6")]
	MOD_ADDR_TEMP14_PARA6
	//[Description("温度监控14 检测上限")]
	MOD_ADDR_TEMP14_MAX
	//[Description("温度监控14 检测下限")]
	MOD_ADDR_TEMP14_MIN

	//[Description("温度监控15 使能")]
	MOD_ADDR_TEMP15_EN
	//[Description("温度监控15 类型")]
	MOD_ADDR_TEMP15_TYPE
	//[Description("温度监控15 通道绑定")]
	MOD_ADDR_TEMP15_CH_BIND
	//[Description("温度监控15 滤波次数")]
	MOD_ADDR_TEMP15_MON_FLT_TIME
	//[Description("温度监控15 参数1")]
	MOD_ADDR_TEMP15_PARA1
	//[Description("温度监控15 参数2")]
	MOD_ADDR_TEMP15_PARA2
	//[Description("温度监控15 参数3")]
	MOD_ADDR_TEMP15_PARA3
	//[Description("温度监控15 参数4")]
	MOD_ADDR_TEMP15_PARA4
	//[Description("温度监控15 参数5")]
	MOD_ADDR_TEMP15_PARA5
	//[Description("温度监控15 参数6")]
	MOD_ADDR_TEMP15_PARA6
	//[Description("温度监控15 检测上限")]
	MOD_ADDR_TEMP15_MAX
	//[Description("温度监控15 检测下限")]
	MOD_ADDR_TEMP15_MIN

	//[Description("温度监控16 使能")]
	MOD_ADDR_TEMP16_EN
	//[Description("温度监控16 类型")]
	MOD_ADDR_TEMP16_TYPE
	//[Description("温度监控16 通道绑定")]
	MOD_ADDR_TEMP16_CH_BIND
	//[Description("温度监控16 滤波次数")]
	MOD_ADDR_TEMP16_MON_FLT_TIME
	//[Description("温度监控16 参数1")]
	MOD_ADDR_TEMP16_PARA1
	//[Description("温度监控16 参数2")]
	MOD_ADDR_TEMP16_PARA2
	//[Description("温度监控16 参数3")]
	MOD_ADDR_TEMP16_PARA3
	//[Description("温度监控16 参数4")]
	MOD_ADDR_TEMP16_PARA4
	//[Description("温度监控16 参数5")]
	MOD_ADDR_TEMP16_PARA5
	//[Description("温度监控16 参数6")]
	MOD_ADDR_TEMP16_PARA6
	//[Description("温度监控16 检测上限")]
	MOD_ADDR_TEMP16_MAX
	//[Description("温度监控16 检测下限")]
	MOD_ADDR_TEMP16_MIN

	//[Description("温度监控17 使能")]
	MOD_ADDR_TEMP17_EN
	//[Description("温度监控17 类型")]
	MOD_ADDR_TEMP17_TYPE
	//[Description("温度监控17 通道绑定")]
	MOD_ADDR_TEMP17_CH_BIND
	//[Description("温度监控17 滤波次数")]
	MOD_ADDR_TEMP17_MON_FLT_TIME
	//[Description("温度监控17 参数1")]
	MOD_ADDR_TEMP17_PARA1
	//[Description("温度监控17 参数2")]
	MOD_ADDR_TEMP17_PARA2
	//[Description("温度监控17 参数3")]
	MOD_ADDR_TEMP17_PARA3
	//[Description("温度监控17 参数4")]
	MOD_ADDR_TEMP17_PARA4
	//[Description("温度监控17 参数5")]
	MOD_ADDR_TEMP17_PARA5
	//[Description("温度监控17 参数6")]
	MOD_ADDR_TEMP17_PARA6
	//[Description("温度监控17 检测上限")]
	MOD_ADDR_TEMP17_MAX
	//[Description("温度监控17 检测下限")]
	MOD_ADDR_TEMP17_MIN

	//[Description("温度监控18 使能")]
	MOD_ADDR_TEMP18_EN
	//[Description("温度监控18 类型")]
	MOD_ADDR_TEMP18_TYPE
	//[Description("温度监控18 通道绑定")]
	MOD_ADDR_TEMP18_CH_BIND
	//[Description("温度监控18 滤波次数")]
	MOD_ADDR_TEMP18_MON_FLT_TIME
	//[Description("温度监控18 参数1")]
	MOD_ADDR_TEMP18_PARA1
	//[Description("温度监控18 参数2")]
	MOD_ADDR_TEMP18_PARA2
	//[Description("温度监控18 参数3")]
	MOD_ADDR_TEMP18_PARA3
	//[Description("温度监控18 参数4")]
	MOD_ADDR_TEMP18_PARA4
	//[Description("温度监控18 参数5")]
	MOD_ADDR_TEMP18_PARA5
	//[Description("温度监控18 参数6")]
	MOD_ADDR_TEMP18_PARA6
	//[Description("温度监控18 检测上限")]
	MOD_ADDR_TEMP18_MAX
	//[Description("温度监控18 检测下限")]
	MOD_ADDR_TEMP18_MIN

	//[Description("温度监控19 使能")]
	MOD_ADDR_TEMP19_EN
	//[Description("温度监控19 类型")]
	MOD_ADDR_TEMP19_TYPE
	//[Description("温度监控19 通道绑定")]
	MOD_ADDR_TEMP19_CH_BIND
	//[Description("温度监控19 滤波次数")]
	MOD_ADDR_TEMP19_MON_FLT_TIME
	//[Description("温度监控19 参数1")]
	MOD_ADDR_TEMP19_PARA1
	//[Description("温度监控19 参数2")]
	MOD_ADDR_TEMP19_PARA2
	//[Description("温度监控19 参数3")]
	MOD_ADDR_TEMP19_PARA3
	//[Description("温度监控19 参数4")]
	MOD_ADDR_TEMP19_PARA4
	//[Description("温度监控19 参数5")]
	MOD_ADDR_TEMP19_PARA5
	//[Description("温度监控19 参数6")]
	MOD_ADDR_TEMP19_PARA6
	//[Description("温度监控19 检测上限")]
	MOD_ADDR_TEMP19_MAX
	//[Description("温度监控19 检测下限")]
	MOD_ADDR_TEMP19_MIN

	/*#############时间参数#############*/
	//[Description("泵浦工作时长重置")]
	MOD_ADDR_LASER_PUMP_WORK_TIME_RESET /*激光器泵浦工作时间(秒)重置*/
	//[Description("激光器出光时长重置")]
	MOD_ADDR_LASER_EMISSON_TIME_RESET /*激光器出光时间(秒)重置*/
	//[Description("激光器运行时长重置")]
	MOD_ADDR_LASER_TOTAL_UPTIME_RESET /*激光器总运行时间(秒)重置*/

	/*激光器激活时间点*/
	//[Description("激光器激活时间(年月)")]
	MOD_ADDR_LASER_ACTTIME_YEAR_MON
	//[Description("激光器激活时间(日时)")]
	MOD_ADDR_LASER_ACTTIME_DAY_HOR
	//[Description("激光器激活时间(分秒)")]
	MOD_ADDR_LASER_ACTTIME_MINU_SEC

	/*激光器出厂时间点*/
	//[Description("激光器出厂时间(年月)")]
	MOD_ADDR_LASER_FACTIME_YEAR_MON
	//[Description("激光器出厂时间(日时)")]
	MOD_ADDR_LASER_FACTIME_DAY_HOR
	//[Description("激光器出厂时间(分秒)")]
	MOD_ADDR_LASER_FACTIME_MINU_SEC

	/*激光器使用期限*/
	//[Description("激光器使用期限(年月)")]
	MOD_ADDR_LASER_SERTIME_YEAR_MON
	//[Description("激光器使用期限(日时)")]
	MOD_ADDR_LASER_SERTIME_DAY_HOR
	//[Description("激光器使用期限(分秒)")]
	MOD_ADDR_LASER_SERTIME_MINU_SEC

	/*系统时间设置*/
	//[Description("激光器使系统配置(年月)")]
	MOD_ADDR_LASER_SYS_TIME_YEAR_MON_CFG
	//[Description("激光器使系统配置(日时)")]
	MOD_ADDR_LASER_SYS_TIME_DAY_HOR_CFG
	//[Description("激光器使系统配置(分秒)")]
	MOD_ADDR_LASER_SYS_TIME_MINU_SEC_CFG
	//[Description("激光器系统时间配置执行")]
	MOD_ADDR_LASER_SYS_TIME_UPDATE

	/*解锁码输入地址激活*/
	//[Description("激光器激活密匙0")]
	MOD_ADDR_LASER_KEY_0
	//[Description("激光器激活密匙1")]
	MOD_ADDR_LASER_KEY_1
	//[Description("激光器激活密匙2")]
	MOD_ADDR_LASER_KEY_2
	//[Description("激光器激活密匙3")]
	MOD_ADDR_LASER_KEY_3
	//[Description("激光器激活密匙4")]
	MOD_ADDR_LASER_KEY_4
	//[Description("激光器激活密匙5")]
	MOD_ADDR_LASER_KEY_5
	//[Description("激光器激活密匙6")]
	MOD_ADDR_LASER_KEY_6
	//[Description("激光器激活密匙7")]
	MOD_ADDR_LASER_KEY_7
	//[Description("激光器激活密匙8")]
	MOD_ADDR_LASER_KEY_8
	//[Description("激光器激活密匙9")]
	MOD_ADDR_LASER_KEY_9
	//[Description("激光器激活执行")]
	MOD_ADDR_LASER_KEY_EN

	/*告警使能表*/
	//[Description("激光器告警使能0")]
	MOD_ADDR_LASER_ALARM_EN_0
	//[Description("激光器告警使能1")]
	MOD_ADDR_LASER_ALARM_EN_1
	//[Description("激光器告警使能2")]
	MOD_ADDR_LASER_ALARM_EN_2
	//[Description("激光器告警使能3")]
	MOD_ADDR_LASER_ALARM_EN_3
	//[Description("激光器告警使能4")]
	MOD_ADDR_LASER_ALARM_EN_4
	//[Description("激光器告警使能5")]
	MOD_ADDR_LASER_ALARM_EN_5
	//[Description("激光器告警使能6")]
	MOD_ADDR_LASER_ALARM_EN_6
	//[Description("激光器告警使能7")]
	MOD_ADDR_LASER_ALARM_EN_7
	//[Description("激光器告警使能8")]
	MOD_ADDR_LASER_ALARM_EN_8
	//[Description("激光器告警使能9")]
	MOD_ADDR_LASER_ALARM_EN_9
	//[Description("激光器告警使能10")]
	MOD_ADDR_LASER_ALARM_EN_10
	//[Description("激光器告警使能11")]
	MOD_ADDR_LASER_ALARM_EN_11

	//[Description("激光器校准参数0")]
	MOD_ADDR_POW_CALI_0
	//[Description("激光器校准参数1")]
	MOD_ADDR_POW_CALI_1
	//[Description("激光器校准参数2")]
	MOD_ADDR_POW_CALI_2
	//[Description("激光器校准参数3")]
	MOD_ADDR_POW_CALI_3
	//[Description("激光器校准参数4")]
	MOD_ADDR_POW_CALI_4
	//[Description("激光器校准参数5")]
	MOD_ADDR_POW_CALI_5
	//[Description("激光器校准参数6")]
	MOD_ADDR_POW_CALI_6
	//[Description("激光器校准参数7")]
	MOD_ADDR_POW_CALI_7
	//[Description("激光器校准参数8")]
	MOD_ADDR_POW_CALI_8
	//[Description("激光器校准参数9")]
	MOD_ADDR_POW_CALI_9
	//[Description("激光器校准参数10")]
	MOD_ADDR_POW_CALI_10
	//[Description("激光器校准参数11")]
	MOD_ADDR_POW_CALI_11
	//[Description("激光器校准参数12")]
	MOD_ADDR_POW_CALI_12
	//[Description("激光器校准参数13")]
	MOD_ADDR_POW_CALI_13
	//[Description("激光器校准参数14")]
	MOD_ADDR_POW_CALI_14
	//[Description("激光器校准参数15")]
	MOD_ADDR_POW_CALI_15
	//[Description("激光器校准参数16")]
	MOD_ADDR_POW_CALI_16
	//[Description("激光器校准参数17")]
	MOD_ADDR_POW_CALI_17
	//[Description("激光器校准参数18")]
	MOD_ADDR_POW_CALI_18
	//[Description("激光器校准参数19")]
	MOD_ADDR_POW_CALI_19

	//[Description("激光器升级")]
	MOD_ADDR_SYS_UPDATE /*系统升级设置*/
	//[Description("激光器重启")]
	MOD_ADDR_SYS_RESET /*系统重启*/

	//[Description("FPGA_调试地址配置")]
	MOD_ADDR_PPGA_DEBUG_ADDR_CFG
	//[Description("FPGA_调试写入值_0")]
	MOD_ADDR_PPGA_DEBUG_VALUE_WRITE_0
	//[Description("FPGA_调试写入值_1")]
	MOD_ADDR_PPGA_DEBUG_VALUE_WRITE_1
	//[Description("FPGA_调试读写控制")]
	MOD_ADDR_PPGA_DEBUG_WR_CMD

	//[Description("激光器模式")]
	MOD_ADDR_LASER_MODE
	//[Description("激光器脉宽上限")]
	MOD_ADDR_LASER_PULSE_WIDTH_MAX
	//[Description("激光器脉宽下限")]
	MOD_ADDR_LASER_PULSE_WIDTH_MIN
	//[Description("激光器频率上限")]
	MOD_ADDR_LASER_FREQ_MAX
	//[Description("激光器频率下限")]
	MOD_ADDR_LASER_FREQ_MIN

	//[Description("激光器参数保存")]
	MOD_ADDR_PARA_SAVE
	//[Description("激光器参数校验")]
	MOD_ADDR_PARA_CRC
	//[Description("激光器参数末尾")]
	MOD_ADDR_PARA_END

	/*###############用户使用接口###############*/
	//[Description("激光器用户使用参数头")]
	MOD_ADDR_USER_HEAD = iota - MOD_ADDR_PARA_END - 1 + USER_PARA_OFFECT
	//[Description("用户参数列表")]
	MOD_ADDR_USER_PARA_LIST_INDEX
	//[Description("用户参数保存")]
	MOD_ADDR_USER_PARA_SAVE
	//[Description("激光器选单频率")]
	MOD_ADDR_LASER_FREQ /*激光器选单频率*/
	//[Description("激光器脉宽参数")]
	MOD_ADDR_LASER_PULSE_WIDTH /*激光器脉冲宽度*/
	//[Description("激光器脉宽分频低16bit")]
	MOD_ADDR_FREQ_DIV_FACTOR_0 /*激光器分频参数低16bit*/
	//[Description("激光器脉宽分频高16bit")]
	MOD_ADDR_FREQ_DIV_FACTOR_1 /*激光器分频参数高16bit*/
	//[Description("激光器脉宽分频burst个数")]
	MOD_ADDR_LASER_PULSE_BURST /*burst 个数设定*/

	//[Description("激光器脉宽分频burst编辑0")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_0 /*burst 波形编辑*/
	//[Description("激光器脉宽分频burst编辑1")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_1
	//[Description("激光器脉宽分频burst编辑2")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_2
	//[Description("激光器脉宽分频burst编辑3")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_3
	//[Description("激光器脉宽分频burst编辑4")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_4
	//[Description("激光器脉宽分频burst编辑5")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_5
	//[Description("激光器脉宽分频burst编辑6")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_6
	//[Description("激光器脉宽分频burst编辑7")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_7
	//[Description("激光器脉宽分频burst编辑8")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_8
	//[Description("激光器脉宽分频burst编辑9")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_9
	//[Description("激光器脉宽分频burst编辑10")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_10
	//[Description("激光器脉宽分频burst编辑11")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_11
	//[Description("激光器脉宽分频burst编辑12")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_12
	//[Description("激光器脉宽分频burst编辑13")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_13
	//[Description("激光器脉宽分频burst编辑14")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_14
	//[Description("激光器脉宽分频burst编辑15")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_15
	//[Description("激光器脉宽分频burst编辑16")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_16
	//[Description("激光器脉宽分频burst编辑17")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_17
	//[Description("激光器脉宽分频burst编辑18")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_18
	//[Description("激光器脉宽分频burst编辑19")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_19
	//[Description("激光器脉宽分频burst编辑20")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_20
	//[Description("激光器脉宽分频burst编辑21")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_21
	//[Description("激光器脉宽分频burst编辑22")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_22
	//[Description("激光器脉宽分频burst编辑23")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_23
	//[Description("激光器脉宽分频burst编辑24")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_24
	//[Description("激光器脉宽分频burst编辑25")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_25
	//[Description("激光器脉宽分频burst编辑26")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_26
	//[Description("激光器脉宽分频burst编辑27")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_27
	//[Description("激光器脉宽分频burst编辑28")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_28
	//[Description("激光器脉宽分频burst编辑29")]
	MOD_ADDR_LASER_PULSE_BURST_LEVEL_29
	//[Description("激光器脉宽分频burst更新")]
	MOD_ADDR_LASER_PULSE_BURST_UPDATE /*burst 数据一次性更新*/

	//[Description("激光器告警清除")]
	MOD_ADDR_LASER_ALARM_CLEAR /*告警清除*/
	//[Description("激光器软件急停")]
	MOD_ADDR_LASER_EM_STOP /*软件急停*/

	//[Description("激光器软件急停")]
	MOD_ADDR_LASER_ERROR
	//[Description("激光器功率设定")]
	MOD_ADDR_LASER_POWER /*激光器功率*/
	//[Description("激光器软件出光使能")]
	MOD_ADDR_LASER_SW
	//[Description("激光器就绪状态")]
	MOD_ADDR_LASER_READY_STATUS /*激光器就绪状态*/
	//[Description("激光器波长输出")]
	MOD_ADDR_LASER_WAVELENGTH /*激光器波长输出*/

	//[Description("用户使能控制")]
	MOD_ADDR_USER_WATER_COOLER_EN
	//[Description("用户功率控制")]
	MOD_ADDR_USER_POW_SEL
	//[Description("用户出光控制")]
	MOD_ADDR_USER_AOM_SEL
	//[Description("一键开启进度")]
	MOD_ADDR_LASER_READY_PROGRESS
	//[Description("IO测试信号")]
	MOD_ADDR_LASER_IO_TEST_OUT

	/*保持寄存器结尾  */
	//[Description("激光器参数校验")]
	MOD_ADDR_USER_CRC
	//[Description("激光器用户参数末尾")]
	MOD_ADDR_USER_END

	MOD_ADDR_HOLDINGD_END
)

//MODBUS_HOLDING_REGISTERS;

//保持寄存器

//数据解析

package model

//################var Femto_input_reg #############//
type Femto_input_reg_db25 struct {
	Ext_data_in      uint16    `json:"Ext_data_in"`
	Ext_latch        uint16    `json:"Ext_latch"`
	Ext_alarm_cooler uint16    `json:"Ext_alarm_cooler"`
	Ext_test         uint16    `json:"Ext_test"`
	Ext_reverse      [2]uint16 `json:"Ext_reverse"`
	Ext_gate         uint16    `json:"Ext_gate"`

	Ext_pwm             uint16 `json:"Ext_pwm"`
	Ext_prr             uint16 `json:"Ext_prr"`
	Ext_trig            uint16 `json:"Ext_trig"`
	Ext_sync            uint16 `json:"Ext_sync"`
	Ext_anlog_data      uint16 `json:"Ext_anlog_data"`
	Ext_water_flow_freq uint16 `json:"Ext_water_flow_freq"`
	Ext_water_flow      uint16 `json:"Ext_water_flow"`
}

type Femto_input_reg_monitor_pump struct {
	Pump_sw    uint16 `json:"Pump_sw"`
	Actual_cur uint16 `json:"Actual_cur"`
	Fpga_cur   uint16 `json:"Fpga_cur"`
}

type Femto_input_reg_monitor_TH struct {
	Temp uint16 `json:"Temp"`
	Humi uint16 `json:"Humi"`
}

type Femto_input_reg_monitor_digi_tcm struct {
	Actual_temp uint16 `json:"Actual_temp"`
	Alarm_reg   uint16 `json:"Alarm_reg"`
}

type Femto_input_reg_monitor struct {
	Pd_freq          [4]uint16                           `json:"Pd_freq"`
	Seed_status      [6]uint16                           `json:"Seed_status"`
	Temp             [20]uint16                          `json:"Temp"`
	Pump_mon         [15]Femto_input_reg_monitor_pump    `json:"Femto_input_reg_monitor_pump"`
	Vol              [20]uint16                          `json:"Vol"`
	Motor_Actual_pos [4]uint16                           `json:"motor_Actual_pos"`
	Th               [2]Femto_input_reg_monitor_TH       `json:"Femto_input_reg_monitor_TH"`
	Digi_tcm_temp    [6]Femto_input_reg_monitor_digi_tcm `json:"Femto_input_reg_monitor_digi_tcm"`
}

type Femto_input_reg_fpga_debug struct {
	Addr_echo uint16 `json:"Addr_echo"`
	Value_0   uint16 `json:"Value_0"`
	Value_1   uint16 `json:"Value_1"`
	Cnt       uint16 `json:"Cnt"`
}

type Time_desc struct {
	Mon    uint8 `json:"Mon"`
	Year   uint8 `json:"Year"`
	Hour   uint8 `json:"Hour"`
	Day    uint8 `json:"Day"`
	Second uint8 `json:"Second"`
	Minute uint8 `json:"Minutes"`
}

type Femto_input_reg_time struct {
	Pump_work_time [2]uint16 `json:"Pump_work_time"`
	Emission_time  [2]uint16 `json:"Emission_time"`
	Uptime         [2]uint16 `json:"Uptime"`
	Total_uptime   [2]uint16 `json:"Total_uptime"`
	Sys_time       Time_desc `json:"Sys_time"`
}
type Femto_input_reg_bate struct {
	Hardware_bate uint16    `json:"Hardware_bate"`
	Mcu_boot_bate uint16    `json:"Mcu_boot_bate"`
	Mcu_app_bate  [6]uint16 `json:"Mcu_app_bate"`
	Fpga_bate     [6]uint16 `json:"Fpga_bate"`
}

type Femto_input_reg_alarm struct {
	Now     [12]uint16 `json:"Now"`
	History [12]uint16 `json:"History"`
}

type Femto_input_reg struct {
	Laser_status uint16                  `json:"laser_status"`
	Borad_online [2]uint16               `json:"board_online"`
	DB25         Femto_input_reg_db25    `json:"Femto_input_reg_db25"`
	Bate         Femto_input_reg_bate    `json:"Femto_input_reg_bate"`
	Mon          Femto_input_reg_monitor `json:"Femto_input_reg_monitor"`

	Overdue    uint16                     `json:"overdue"`
	Fpga_debug Femto_input_reg_fpga_debug `json:"Femto_input_reg_fpga_debug"`
	Time       Femto_input_reg_time       `json:"Femto_input_reg_time"`
	Alarm      Femto_input_reg_alarm      `json:"Femto_input_reg_alarm"`
}

//################var Femto_input_reg #############

//################var Femto_holding_reg#############//
type Femto_laser_para_info struct {
	Model             [20]uint8 `json:"Model"`
	SN                [20]uint8 `json:"SN"`
	PN                [20]uint8 `json:"PN"`
	Laser_serial      uint16    `json:"Laser_serial"`
	Laser_Power_level uint16    `json:"Laser_Power_level"`
}

type Femto_laser_para_esp struct {
	En       uint16     `json:"Model"`
	Mode     uint16     `json:"SN"`
	Ssid     [15]uint16 `json:"PN"`
	Wifi_ip  [4]uint16  `json:"Wifi_ip"`
	Eth_en   uint16     `json:"Eth_en"`
	Eth_mode uint16     `json:"Eth_mode"`
	Eth_ip   [4]uint16  `json:"Eth_ip"`
}

type Femto_laser_para_th struct {
	En       uint16 `json:"En"`
	Temp_max uint16 `json:"Temp_max"`
	Temp_min uint16 `json:"Temp_min"`
	Humi_max uint16 `json:"Humi_max"`
	Humi_min uint16 `json:"Humi_min"`
}

type Femto_laser_para_air_pump struct {
	En       uint16 `json:"En"`
	Open_th  uint16 `json:"Open_th"`
	Close_th uint16 `json:"Close_th"`
}

type Femto_laser_para_seed struct {
	En   uint16    `json:"En"`
	Type uint16    `json:"Type"`
	SW   uint16    `json:"SW"`
	Para [5]uint16 `json:"Para"`
}

type Femto_laser_para_pump struct {
	En               uint16 `json:"En"`
	SW               uint16 `json:"SW"`
	Priority         uint16 `json:"Priority"`
	Dest_cur         uint16 `json:"Dest_cur"`
	Compensation_val uint16 `json:"Compensation_val"`
	Coefficient_val  uint16 `json:"Coefficient_val"`
	Cur_speed        uint16 `json:"Cur_speed"`
	Cur_max_reg      uint16 `json:"Cur_max_reg"`

	//for mon
	Mon_cur_band_ch      uint16 `json:"Mon_cur_band_ch"`
	Mon_compensation_val uint16 `json:"Mon_compensation_val"`
	Mon_coefficient_val  uint16 `json:"Mon_coefficient_val"`
	Mon_cur_err_thr      uint16 `json:"Mon_cur_err_thr"`
	Mon_cur_filter_time  uint16 `json:"Mon_cur_filter_time"`
}

type Femto_laser_para_motor struct {
	En             uint16     `json:"en"`
	Freq           uint16     `json:"Freq"`
	Dest_pos       uint16     `json:"Dest_pos"`
	Power_up_reset uint16     `json:"Power_up_reset"`
	Pos_table      [20]uint16 `json:"Pos_table"`
}

type Femto_laser_para_tcm struct {
	En        uint16 `json:"en"`
	Sw        uint16 `json:"sw"`
	Type      uint16 `json:"Type"`
	Dest_temp uint16 `json:"Dest_temp"`
	Pid_p     uint16 `json:"Pid_p"`
	Pid_i     uint16 `json:"Pid_i"`
	Pid_d     uint16 `json:"Pid_d"`
	Temp_max  uint16 `json:"Temp_max"`
	Temp_min  uint16 `json:"Temp_min"`
}

type Femto_laser_para_pd struct {
	En       uint16 `json:"en"`
	Freq_max uint16 `json:"Freq_max"`
	Freq_min uint16 `json:"Freq_min"`
}

type Femto_laser_para_aom struct {
	En        uint16 `json:"en"`
	Aom_level uint16 `json:"Aom_level"`
	Delay1    uint16 `json:"Delay1"`
	Delay2    uint16 `json:"Delay2"`
	Pll       uint16 `json:"Pll"`
}

type Femto_laser_para_water_cool struct {
	En                uint16 `json:"en"`
	Flow_max          uint16 `json:"Flow_max"`
	Flow_min          uint16 `json:"Flow_min"`
	Flow_compensation uint16 `json:"Flow_compensation"`
}

type Femto_laser_para_vol struct {
	En               uint16 `json:"en"`
	Adc_band_ch      uint16 `json:"Adc_band_ch"`
	Vol_compensation uint16 `json:"Vol_compensation"`
	Vol_coefficient  uint16 `json:"vol_coefficient"`
	Vol_max          uint16 `json:"Vol_max"`
	Vol_min          uint16 `json:"Vol_min"`
	Filter_time      uint16 `json:"Filter_time"`
}

type Femto_laser_para_temp struct {
	En          uint16    `json:"en"`
	Type        uint16    `json:"Type"`
	Adc_band_ch uint16    `json:"Adc_band_ch"`
	Filter_time uint16    `json:"Filter_time"`
	Para        [6]uint16 `json:"Para"`
	Temp_max    uint16    `json:"Temp_max"`
	Temp_min    uint16    `json:"Temp_min"`
}

type Femto_laser_para_time struct {
	Pump_time_reset     uint16    `json:"Pump_time_reset"`
	Emission_time_reset uint16    `json:"Emission_time_reset"`
	Total_uptime_reset  uint16    `json:"Total_uptime_reset"`
	Acttime             Time_desc `json:"Acttime"`
	Factime             Time_desc `json:"Factime"`
	Sertime             Time_desc `json:"Sertime"`
	Rtc_time            Time_desc `json:"Rtc_time"`
}

type Femto_laser_para_key struct {
	Key_val [10]uint16 `json:"key_val"`
	Key_en  uint16     `json:"Key_en"`
}

type Femto_laser_para struct {
	Head                  uint16                       `json:"Head"`
	Laser_info            Femto_laser_para_info        `json:"Laser_info"`
	Esp_module            Femto_laser_para_esp         `json:"Esp_module"`
	Th_modelue            [2]Femto_laser_para_th       `json:"Th_modelue"`
	Air_pump_modelue      [2]Femto_laser_para_air_pump `json:"Air_pump_modelue"`
	Seed                  Femto_laser_para_seed        `json:"Seed"`
	Pump_module           [15]Femto_laser_para_pump    `json:"Pump_module"`
	Motor_module          [4]Femto_laser_para_motor    `json:"Motor_module"`
	Tcm_module            [6]Femto_laser_para_tcm      `json:"Tcm_module"`
	Dac_vol               [40]uint16                   `json:"Dac_vol"`
	Pd_modelue            [4]Femto_laser_para_pd       `json:"Pd_modelue"`
	Aom_modelue           [4]Femto_laser_para_aom      `json:"Aom_modelue"`
	Water_cool            Femto_laser_para_water_cool  `json:"Water_cool"`
	Vol_para              [20]Femto_laser_para_vol     `json:"Vol_para"`
	Temp_para             [20]Femto_laser_para_temp    `json:"Temp_para"`
	Time_para             Femto_laser_para_time        `json:"Time_para"`
	Key_para              Femto_laser_para_key         `json:"Key_para"`
	Alarm_en              [12]int16                    `json:"Alarm_en"`
	Pow_celi              [20]int16                    `json:"Pow_celi"`
	Sys_update            int16                        `json:"Sys_update"`
	Sys_reboot            int16                        `json:"Sys_reboot"`
	Fpga_addr_cfg         int16                        `json:"Fpga_addr_cfg"`
	Fpga_Value_0          uint16                       `json:"Fpga_Value_0"`
	Fpga_Value_1          uint16                       `json:"Fpga_Value_1"`
	Fpga_rw_cmd           uint16                       `json:"Fpga_rw_cmd"`
	Laser_mode            uint16                       `json:"Laser_mode"`
	Laser_pulse_width_max uint16                       `json:"Laser_pulse_width_max"`
	Laser_pulse_width_min uint16                       `json:"Laser_pulse_width_min"`
	Laser_freq_max        uint16                       `json:"Laser_freq_max"`
	Laser_freq_min        uint16                       `json:"Laser_freq_min"`
	Laser_save            uint16                       `json:"Laser_save"`
	Laser_para_crc        uint16                       `json:"Laser_para_crc"`
	Laser_para_end        uint16                       `json:"Laser_para_end"`

	//预留
	Reserved [USER_PARA_OFFECT - MOD_ADDR_PARA_END - 1]uint16 `json:"Reserved"`
}

type Femto_user_para struct {
	Head               uint16     `json:"Head"`
	List_index         uint16     `json:"List_index"`
	User_para_Save     uint16     `json:"User_para_Save"`
	Freq               uint16     `json:"Freq"`
	Puse_width         uint16     `json:"Puse_width"`
	Div_factor_0       uint16     `json:"Div_factor_0"`
	Div_factor_1       uint16     `json:"Div_factor_1"`
	Pulse_burst        uint16     `json:"Pulse_burst"`
	Pulse_level        [30]uint16 `json:"Pulse_level"`
	Pulse_burst_update uint16     `json:"Pulse_burst_update"`
	Alarm_clear        uint16     `json:"Alarm_clear"`
	Em_stop            uint16     `json:"Em_stop"`
	Laser_Err          uint16     `json:"Laser_Err"`
	Laser_power        uint16     `json:"Laser_power"`
	Laser_sw           uint16     `json:"Laser_sw"`

	Laser_ready          uint16 `json:"Laser_ready"`
	Laser_wavelength     uint16 `json:"Laser_wavelength"`
	Laser_cooler_en      uint16 `json:"Laser_cooler_en"`
	Laser_pow_sel        uint16 `json:"Laser_pow_sel"`
	Laser_aom_sel        uint16 `json:"Laser_aom_sel"`
	Laser_ready_progress uint16 `json:"Laser_ready_progress"`
	Laser_test_out       uint16 `json:"Laser_test_out"`
	Laser_user_crc       uint16 `json:"Laser_user_crc"`
	Laser_user_end       uint16 `json:"Laser_user_end"`
}

type Femto_holding_reg struct {
	Laser_para Femto_laser_para `json:"Laser_para"`
	User_para  Femto_user_para  `json:"User_para"`
}

//################var Femto_holding_reg#############//

type Femto_msg_packed struct {
	Femto_input_reg   Femto_input_reg   `json:"Femto_input_reg" gorm:"embedded"`
	Femto_holding_reg Femto_holding_reg `json:"Femto_holding_reg" gorm:"embedded"`
}

export interface Device {
  Name: string;
  Sn: string;
  Packet: Packet;
}

export interface Packet {
  Femto_input_reg: {
    laser_status: number;
    board_online:number[];
    Femto_input_reg_db25: Packet_input_DB25;
     Femto_input_reg_bate: Packet_bate;
     Femto_input_reg_monitor: Packet_mon_monitor;
     Overdue:number;
     Femto_input_reg_fpga_debug: Packet_mon_fpga;
     Femto_input_reg_time: Packet_mon_time;
     Femto_input_reg_alarm: Packet_mon_alarm;
  },

  Femto_holding_reg:{
    Laser_para : Packet_laser_para;
    User_para  : Packet_user_para;

  }

}

export interface Packet_input_DB25 {
  	Ext_data_in:number;
	Ext_latch :number;
	Ext_alarm_cooler:number;
	Ext_test     :number;
	Ext_reverse    :number;
	Ext_gate       :number;
	Ext_pwm        :number;
	Ext_prr          :number;
	Ext_trig         :number;
	Ext_sync        :number;
	Ext_anlog_data    :number;
	Ext_water_flow_freq :number;
	Ext_water_flow    :number;
}
export interface Packet_bate {
  Hardware_bate:number;
  Mcu_boot_bate:number;
  Mcu_app_bate:number[];
  Fpga_bate:number[];
}



export interface Packet_mon_pump {
	Pump_sw  :number;
	Actual_cur :number;
	Fpga_cur  :number;
}

export interface Packet_mon_th {
	Temp  :number;
	Humi  :number;
}
export interface Packet_mon_digi_tcm {
	Actual_temp  :number;
	Alarm_reg    :number;
}

export interface Packet_mon_monitor {
  Pd_freq :number[];
  Seed_status :number[];
  Temp :number[];
  Femto_input_reg_monitor_pump  :Packet_mon_pump[15];
  Vol :number[];
  motor_Actual_pos :number[];
  Femto_input_reg_monitor_TH  :Packet_mon_th[];
  Femto_input_reg_monitor_digi_tcm :Packet_mon_digi_tcm[];
}


export interface Packet_mon_fpga {
	Addr_echo  :number;
	Value_0   :number;
	Value_1    :number;
	Cnt       :number;
}

export interface Packet_mon_time {
	Pump_work_time :number[];
	Emission_time :number[];
	Uptime       :number[];
	Total_uptime :number[];
	Sys_time      :Time_desc;
}
export interface Packet_mon_alarm {
	Now    :number[];
	History:number[];
}


export interface Time_desc {
  Mon: number;
  Year: number;
  Hour: number;
  Day:  number;
  Second: number;
  Minute: number;
}








export interface Packet_laser_para{
	Head                 : number;             
	Laser_info           : Femto_laser_para_info       
	Esp_module           : Femto_laser_para_esp       
	Th_modelue           : Femto_laser_para_th [2]     
	Air_pump_modelue     :  Femto_laser_para_air_pump[2]
	Seed                 : Femto_laser_para_seed       
	Pump_module          : Femto_laser_para_pump[15] 
	Motor_module         : Femto_laser_para_motor[4]
	Tcm_module           : Femto_laser_para_tcm[6]
	Dac_vol              :  number[40];                    
	Pd_modelue           :  Femto_laser_para_pd [4]  
	Aom_modelue          :  Femto_laser_para_aom[4]
	Water_cool           :  Femto_laser_para_water_cool 
	Vol_para             :  Femto_laser_para_vol[20]
	Temp_para            :  Femto_laser_para_temp[20] 
	Time_para            :  Femto_laser_para_time       
	Key_para             :  Femto_laser_para_key       
	Alarm_en             : number[12];                
	Pow_celi             : number[20];                       
	Sys_update           : number;                           
	Sys_reboot           : number;                         
	Fpga_addr_cfg        : number;                            
	Fpga_Value_0         : number;                          
	Fpga_Value_1         : number;                             
	Fpga_rw_cmd          : number;                             
	Laser_pulse_width_max  : number;                            
	Laser_pulse_width_min  : number;                             
	Laser_freq_max         : number;                              
	Laser_freq_min         : number;                           
	Laser_save             : number;                           
	Laser_para_crc         : number;                           
	Laser_para_end         :number;                          

}


export interface Femto_laser_para_info{
	Model             : number[20];  
	SN               : number[20];  
	PN                : number[20];  
	Laser_serial      :number;       
	Laser_Power_level  :number;        
}

export interface  Femto_laser_para_esp  {
	En        :number;    
	Mode     :number;    
	Ssid     : number[15];  
	Wifi_ip   : number[4];  
	Eth_en     :number;    
	Eth_mode   :number;    
	Eth_ip  : number[4];  
}


export interface  Femto_laser_para_th  {
	En         :number;    
	Temp_max   :number;    
	Temp_min   :number;    
	Humi_max   :number;    
	Humi_min   :number;    
}

export interface  Femto_laser_para_air_pump  {
	En        :number;    
	Open_th    :number;    
	Close_th   :number;    
}

export interface Femto_laser_para_seed  {
	En     :number;    
	Type   :number;    
	SW    :number;    
	Para : number[5];  
}

export interface  Femto_laser_para_pump  {
	En                :number;    
	SW               :number;    
	Priority           :number;    
	Dest_cur           :number;    
	Compensation_val   :number;    
	Coefficient_val    :number;    
	Cur_speed         :number;    
	Cur_max_reg       :number;    

	//for mon
	Mon_cur_band_ch        :number;    
	Mon_compensation_val   :number;    
	Mon_coefficient_val    :number;    
	Mon_cur_err_thr        :number;    
	Mon_cur_filter_time    :number;    
}


export interface Femto_laser_para_motor  {
	En              :number;    
	Freq            :number;    
	Dest_pos        :number;    
	Power_up_reset   :number;    
	Pos_table       : number[20];  
}

export interface  Femto_laser_para_tcm  {
	En        :number;    
	Sw         :number;    
	Type       :number;    
	Dest_temp   :number;    
	Pid_p      :number;    
	Pid_i      :number;    
	Pid_d       :number;    
	Temp_max   :number;    
	Temp_min   :number;    
}

export interface  Femto_laser_para_pd  {
	En         :number;    
	Freq_max   :number;    
	Freq_min   :number;    
}

export interface  Femto_laser_para_aom  {
	En          :number;    
	Aom_level   :number;    
	Delay1      :number;    
	Delay2     :number;    
	Pll         :number;    
}

export interface  Femto_laser_para_water_cool  {
	En                  :number;    
	Flow_max           :number;    
	Flow_min            :number;    
	Flow_compensation   :number;    
}

export interface Femto_laser_para_vol  {
	En                 :number;    
	Adc_band_ch       :number;    
	Vol_compensation   :number;    
	Vol_coefficient   :number;    
	Vol_max           :number;    
	Vol_min           :number;    
	Filter_time       :number;    
}

export interface  Femto_laser_para_temp  {
	En           :number;    
	Type         :number;    
	Adc_band_ch   :number;    
	Filter_time   :number;    
	Para        : number[6];  
	Temp_max     :number;    
	Temp_min      :number;    
}

export interface  Femto_laser_para_time  {
	Pump_time_reset      :number;    
	Emission_time_reset   :number;    
	Total_uptime_reset    :number;    
	Acttime            : Time_desc
	Factime             :Time_desc ;  
	Sertime           :  Time_desc;  
	Rtc_time           : Time_desc ;  
}

export interface  Femto_laser_para_key  {
	Key_val  : number[10];  
	Key_en   :number;    
}

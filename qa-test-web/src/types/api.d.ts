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
  };
}



export interface PacketDB25 {
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
  Femto_input_reg_monitor_pump  :Packet_mon_pump[];
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
	Sys_time      :PacketTIme;
}
export interface Packet_mon_alarm {
	Now    :number[];
	History:number[];
}


export interface PacketTIme {
  Mon: number;
  Year: number;
  Hour: number;
  Day:  number;
  Second: number;
  Minute: number;
}



export interface PacketStatus {
  工作模式: number;
  激光器状态: number;
  种子开关: number;
  单模1开关: number;
  单模2开关: number;
  预放开关: number;
  主放1开关: number;
  主放2开关: number;
  主放3开关: number;
  主放4开关: number;
  声光开关: number;
}

export interface PacketCurrent {
  单模1电流: number;
  单模2电流: number;
  预放1电流: number;
  预放2电流: number;
  主放1电流: number;
  主放2电流: number;
  主放3电流: number;
  主放4电流: number;
  主放电流限制: number;
  种子电流: number;
  温控电流: number;
}

export interface PacketVoltage {}

export interface PacketTemp {}

export interface PacketTempBoard {
  开关: number;
  实际温度: number;
  设置温度: number;
  误差: number;
  传感器类型: number;
  "控制-P": number;
  "控制-I": number;
  "控制-D": number;
  告警: number;
  准备信号: number;
}
export interface PacketSeedModule {}

export interface PacketFPGA {
  预放电流设置: number;
  主放电流设置: number;
  aom设置: number;
  预放频率: number;
  预放占空比: number;
  主放频率: number;
  主放占空比: number;
  声光频率: number;
  声光占空比: number;
  pwm设置: number;
  声光延时: number;
  外控io1: number;
  外控io2: number;
}

export interface PacketAlarm {
  当前告警: number;
  生效告警: number;
  历史告警: number;
}

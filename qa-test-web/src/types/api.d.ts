export interface Device {
  Name: string;
  Sn: string;
  Packet: Packet;
}

export interface Packet {
  系统监控: {
    激光器时间监测: PacketTIme;
    激光器状态监测: PacketStatus;
    激光器电流监测: PacketCurrent;
    激光器电压监测: PacketVoltage;
    激光器温度监测: PacketTemp;
    激光器温控板监测: PacketTempBoard[];
    激光器种子模块监测: PacketSeedModule;
    激光器告警监测: PacketAlarm;
    激光器FPGA寄存器监测: PacketFPGA;
  };
}

export interface PacketTIme {
  年: number;
  月: number;
  日: number;
  时: number;
  分: number;
  秒: number;
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

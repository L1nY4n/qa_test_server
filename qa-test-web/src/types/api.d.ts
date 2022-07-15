
export interface Device {
    Name: string
    Sn: string
    Packet: Packet
}



export interface Packet {
    '系统监控': {
        '激光器时间监测': PacketTIme
        '激光器状态监测': PacketStatus
        '激光器电流监测': PacketCurrent
        '激光器电压监测': PacketVoltage
        '激光器温度监测': PacketTemp
        '激光器温控板监测': PacketTempBoard
        '激光器种子模块监测': PacketSeedModule
        '激光器告警监测': PacketAlarm
        '激光器FPGA寄存器监测': PacketFPGA
    }
}


export interface PacketTIme {
    '年': number
    '月': number
    '日': number
    '时': number
    '分': number
    '秒': number
}


export interface PacketStatus {
    "工作模式": number
    "激光器状态": number
    "种子开关": number
    "单模1开关": number
    "单模2开关": number
    "预放开关": number
    "主放1开关": number
    "主放2开关": number
    "主放3开关": number
    "主放4开关": number
    "声光开关": number
}

export interface PacketCurrent {
    '单模1电流': number
    '单模2电流': number
    '预放1电流': number
    '预放2电流': number
    '主放1电流': number
    '主放2电流': number
    '主放3电流': number
    '主放4电流': number
    '主放电流限制': number
    '种子电流': number
    '温控电流': number
}

export interface PacketVoltage {

}

export interface PacketTemp {

}

export interface PacketTempBoard {

}
export interface PacketSeedModule {

}

export interface PacketFPGA {

}

export interface PacketAlarm {

}
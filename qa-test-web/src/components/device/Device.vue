<template>
    <div class="device">
        <header>

          <a-row type="flex">
    <a-col flex="50px">
      <div>
       <span ref="update_ref" class='update-trigger'></span>

       </div>
    </a-col>
    <a-col flex="100px">
      <h1 class="sn"> {{ info.Sn }}</h1>

    </a-col>  
    <a-col flex="auto"></a-col>
    <a-col flex="200px">


       <h1 class="name">{{ info.Name }}</h1>

    </a-col>
    <a-col flex="50px">
      <a-button type="primary" @click="showDrawer">版本信息</a-button>
    <a-drawer
    v-model:visible="visible"
    class="custom-class"
    style="color: rgb(0, 0, 0) width: 500px;"
    title="设备信息"
    placement="right"
    @after-visible-change="afterVisibleChange"
  >

   <ul >硬件版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Hardware_bate}}</ul>
  <ul >MCU版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Mcu_app_bate}}</ul>
  <ul >FPGA版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Fpga_bate}}</ul>
                  
  </a-drawer>
  
</a-col>
<a-col flex="50px">
    
  <div class="geer" @click="setting">
                <Geer />
            </div>
</a-col>
  </a-row>




        </header>
        <div class="body">

            <div style="background-color: #f0f6f7; padding: 20px">
                    <a-row :gutter="[16,16]" style="display: flex; align-items: center;">
   
                        <a-col :span="16" style="flex: 1;" >
                            <a-card title="外部IO" hoverable :bordered="true">
                            <a-divider orientation="centor">IO电平</a-divider> 
                            <a-card-grid style="width:  25%; text-align: left">Latch:{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_latch}}</a-card-grid>
                            <a-card-grid style="width: 25%; text-align: center">PWM:{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_pwm}}</a-card-grid>
                            <a-card-grid style="width: 25%; text-align: center">GATE:{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_gate}}</a-card-grid>
                            <a-card-grid style="width: 25%; text-align: center">TRIG:{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_trig}}</a-card-grid>
                            <a-card-grid style="width: 25%; text-align: center">PRR:{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_prr}}</a-card-grid>
                            <a-card-grid style="width: 25%; text-align: center">WATERFLOW:{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_alarm_cooler}}</a-card-grid>
                            <a-divider orientation="centor">信号值</a-divider> 
                            <a-row :gutter="16">
                              <a-col :span="16" >
                              <ul >数字8BIT {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_data_in}}</ul>
                              <ul >模拟8BIT {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_anlog_data}}</ul>
                              <ul >水流量 {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_anlog_data}}</ul>
                            </a-col>
                            <a-col :span="8" >
                                <a-progress type="circle" :stroke-color="{'0%': '#108ee9','100%': '#87d068',}"  :strokeWidth = "10" :width= "100" :percent="65" />
                            </a-col>
                            </a-row>
                            </a-card>
                       </a-col>

                       <a-col :span="16"  style="flex: 1;"  >
                       
                            <!-- <a-card hoverable title="设备信息" :bordered="true">
                            <ul >硬件版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Hardware_bate}}</ul>
                            <ul >MCU版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Mcu_app_bate}}</ul>
                            <ul >FPGA版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Fpga_bate}}</ul>
                            </a-card> -->
                      
                        
                          
                             <a-card  hoverable title=" 频率检测 " :bordered="true">
                                        <a-tabs  v-model:activeKey=selectindex>
                                        <a-tab-pane v-for="(dev, index) in info.Packet['Femto_input_reg']['Femto_input_reg_monitor']['Pd_freq']" :key="index" :tab="'测频'+index">
                                          <ul >通道 {{selectindex}}: {{ dev }}</ul>
                                        </a-tab-pane>
                                    
                                        </a-tabs>
                            </a-card >
                       
                        </a-col>
                       
                    </a-row>
            </div>
    

         
            <div style="background-color: #f0f6f7; padding: 20px">
              
            
                    <a-card hoverable  title=" 泵浦监测 ">
                        <a-tabs  v-model:activeKey=selectindex1>
                       <a-tab-pane v-for="(dev , index) in info.Packet.Femto_input_reg.Femto_input_reg_monitor.Femto_input_reg_monitor_pump" :key="index" :tab="'泵浦'+index">    
                        <a-row :gutter="30">
                          <a-col :span="12" >
                        <a-card title="实时值" :bordered="true">
                        <ul > 温度:{{ info.Packet.Femto_input_reg.Femto_input_reg_monitor.Temp[selectindex1]}}</ul>
                        <ul > 电压:{{ info.Packet.Femto_input_reg.Femto_input_reg_monitor.Vol[selectindex1]}}</ul>
                        <ul > 开关状态:{{ dev.Pump_sw }}</ul>
                        <ul > 实际电流:{{ dev.Actual_cur }}</ul>
                        <ul > 设置电流:{{ dev.Fpga_cur }}</ul>
                       </a-card >
                      </a-col>
                      <a-col  :span="12" >
                       <a-card title="参数" :bordered="true">
                     
                        <ul > 使能:{{ info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].En }}</ul>
                        <ul > 开关:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].SW }}</ul>
                        <ul > 优先级:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Priority  }}</ul>
                        <ul > 目标电流:{{ info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Dest_cur  }}</ul>
                        <ul > 补偿系数:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Compensation_val  }}</ul>
                        <ul > 比例系数:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Coefficient_val  }}</ul>
                        <ul > 速度:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Cur_speed }}</ul>
                        <ul > 限流寄存器:{{ info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Cur_max_reg  }}</ul>
                    
                        <ul > 通道绑定:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Mon_cur_band_ch  }}</ul>
                        <ul > 补偿系数:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Mon_compensation_val }}</ul>
                        <ul > 比例系数:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Mon_coefficient_val  }}</ul>
                        <ul > 误差阈值:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Mon_cur_err_thr }}</ul>
                        <ul > 滤波次数:{{  info.Packet.Femto_holding_reg.Laser_para.Pump_module[selectindex1].Mon_cur_filter_time }}</ul>
                      </a-card >
                    </a-col>
                    </a-row>
                    </a-tab-pane>
                        </a-tabs>
              </a-card >




      <div class="antd-pro-pages-dashboard-analysis-twoColLayout">
          <a-row :gutter="24" type="flex" :style="{ marginTop: '24px' }">


        <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
          <a-card hoverable title="电压监控" :bordered="false">
            <a-row>
            <a-col :span="12">
              <a-slider v-model:value="indexVol" :min="0" :max="19" />
            </a-col>
            <a-col :span="4">
              <a-input-number v-model:value="indexVol" :min="0" :max="19" style="margin-left: 16px" />
            </a-col>
          </a-row> 
          <a-divider orientation="centor">监控值</a-divider>   
          <ul >通道{{indexVol}}: {{  info.Packet.Femto_input_reg.Femto_input_reg_monitor.Vol[indexVol] }}</ul>
          <a-divider orientation="centor">参数</a-divider>  
          <ul > 使能:{{ info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].En }}</ul>
          <ul > 通道绑定:{{  info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].Adc_band_ch}}</ul>
          <ul > 补偿系数:{{  info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].Vol_compensation  }}</ul>
          <ul > 比例系数:{{  info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].Vol_coefficient  }}</ul>
          <ul > 上限:{{  info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].Vol_max }}</ul>
          <ul > 下限:{{ info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].Vol_min  }}</ul>
          <ul > 滤波次数:{{ info.Packet.Femto_holding_reg.Laser_para.Vol_para[indexVol].Filter_time  }}</ul>

          </a-card>
        </a-col>

        <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
              <a-card hoverable title="温度监控" :bordered="false">
                <a-row>
              <a-col :span="12">
                <a-slider v-model:value="indexTemp" :min="0" :max="19" />
              </a-col>
              <a-col :span="4">
                <a-input-number v-model:value="indexTemp" :min="0" :max="19" style="margin-left: 16px" />
              </a-col>
            </a-row> 
            <a-divider orientation="centor">监控值</a-divider>   
            <ul >通道{{indexTemp}}: {{  info.Packet.Femto_input_reg.Femto_input_reg_monitor.Temp[indexVol] }}</ul>
            <a-divider orientation="centor">参数</a-divider>  
            <ul > 使能:{{ info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].En }}</ul>
            <ul > 通道绑定:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Adc_band_ch}}</ul>
            <ul > 参数1:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Para[0]  }}</ul>
            <ul > 参数2:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Para[1]  }}</ul>
            <ul > 参数3:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Para[2]  }}</ul>
            <ul > 参数4:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Para[3]  }}</ul>
            <ul > 参数5:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Para[4]  }}</ul>
            <ul > 参数6:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Para[5]  }}</ul>
            <ul > 上限:{{  info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Temp_max }}</ul>
            <ul > 下限:{{ info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Temp_min  }}</ul>
            <ul > 滤波次数:{{ info.Packet.Femto_holding_reg.Laser_para.Temp_para[indexTemp].Filter_time  }}</ul>

            </a-card>
            </a-col>
          </a-row>
    </div>

        </div>

            

    
            <!-- <TitleCard title="温度监测">
                <DeviceVoltage :datas="info.Packet['系统监控']['激光器温度监测']" />
            </TitleCard>
              <TitleCard title="温控板">
              <DeviceTempBoard :datas="info.Packet['系统监控']['激光器温控板监测']" />
            </TitleCard>
              <TitleCard title="种子模块">
              <DeviceSeedModule  :datas="info.Packet['系统监控']['激光器种子模块监测']"/> 
           
            </TitleCard>--> 
       
             <!-- <TitleCard title=" 激光器参数">
              <Tree :data="info.Packet['Femto_holding_reg']" />
            </TitleCard>   -->
          
            
        </div>

    </div>

</template>
<script setup lang="ts">
import moment from 'moment'
//import DataSet from '@antv/data-set'
import { Slider, Tag } from 'ant-design-vue';
import { Device ,Packet_mon_pump } from '@/types/api';
import { onUpdated, ref, computed } from 'vue';
import DeviceData from './DeviceData.vue';
import TimeVue from './Time.vue'
import StatusVue from './DeviceStatus.vue';

import Geer from '@/components/widget/svg/geer.vue'
import TitleCard from './TitleCard.vue';
import Tree from  '@/components/treelist/TreeList.vue';
import { dataTool } from 'echarts';
import { device } from '@/api';
import { InfoCircleFilled } from '@ant-design/icons-vue';


const selectindex =ref(0)
const selectindex1 =ref(0)
const indexVol=ref(0)
const indexTemp=ref(0)
const props = defineProps<{ info: Device }>()


const update_ref = ref()
onUpdated(() => {

    update_ref.value.classList.add('update')
    setTimeout(() => update_ref.value.classList.remove('update'), 150)

})


const visible =ref<boolean>(false)

const showDrawer = () => {
      visible.value = true;
    };

const afterVisibleChange = (bool: boolean) => {
console.log('visible', bool);
};






const setting= ()=>{
    alert("setting click")
}



</script>
<style lang="less" scoped>
@border-color: rgba(255, 255, 255, .1);

.device {
   position: relative;
    align-items: center;
    border: 1px solid @border-color;
    border-radius: 4px;
    background: hsl(0, 0%, 100%);
    box-shadow: 0 4px 30px rgba(218, 227, 229, 0.5);

    header {
        // display: grid;
        // grid-template-columns: 32px 160px auto 100px;
         line-height: 32px;
         text-align: center;
         border-bottom: 1px solid @border-color;

        .update-trigger {
            background-color: #d6edef;
            display: inline-block;
            width: .8rem;
            height: .8rem;
            line-height: 32px;
            margin-right: .3rem;
            border-radius: 50%;

            &.update {
                background: rgb(68, 247, 3);
                animation: updateblink .1s linear;
            }

        }

        h1 {
            font-family: Lato, sans-serif;
            font-weight: 900;
            letter-spacing: 1.57px;
            line-height: 32px;
            margin: 0;
        }
        
        .name {

            text-align: right;
            color: #c5c8d4;
        }

        .sn {
            text-align: left;
            color: #00dcfe;
        }
        .geer{
            cursor: pointer;
        }
    }

  .body{
    overflow: scroll;
    position: relative;
    .alarm{
        position: absolute;
        top: 0;
        right: 0;
    }
    .antd-pro-pages-dashboard-analysis-twoColLayout {
    position: relative;
    display: flex;
    display: block;
    flex-flow: row wrap;
  }


    
}

    &.hasAlarm{
        box-shadow: 0px 4px 15px rgba(0, 0, 0, .2);
        border-color: #c8161c;
        animation: alertblink 2s ease-in-out infinite;

        &:hover {
            background-color: rgba(252, 12, 20, 0.496);
            animation: none;
        }
    }



    @keyframes alertblink {
        0% {
            background: rgba(210, 44, 50, 0);
        }

        50% {
            background: rgba(210, 44, 50, 0.2);
        }

        100% {
            background: rgba(210, 44, 50, 0);
        }
    }

    @keyframes updateblink {
        0% {

            transform: scale(0.9);
        }

        50% {

            transform: scale(1);
            box-shadow: 0 0 0 4px rgba(118, 255, 80, 0.348);
        }

        100% {

            transform: scale(0.9);
            box-shadow: 0 0 0 0 rgba(16, 46, 27, 0.429);
        }
    }

}

.ant-card {
  border-radius: 10px;
}

</style>
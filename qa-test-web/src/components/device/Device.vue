<template>
    <div class="device">
        <header>
            <div>
                <span ref="update_ref" class='update-trigger'></span>
            </div>

            <h1 class="sn"> {{ info.Sn }}</h1>
            <h1 class="name">{{ info.Name }}</h1>
            <div class="geer" @click="setting">
                <Geer />
            </div>

        </header>
        <div class="body">

            <div style="background-color: #f0f6f7; padding: 20px">
                    <a-row :gutter="16">
                    <a-col :span="8">
                        <a-card title="设备信息" :bordered="false">
                          <a-divider orientation="left">版本</a-divider>
                        <ul >硬件版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Hardware_bate}}</ul>
                        <ul >MCU版本:   {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Mcu_app_bate}}</ul>
                        <ul >FPGA版本:  {{info.Packet['Femto_input_reg']['Femto_input_reg_bate'].Fpga_bate}}</ul>
                        </a-card>
                    </a-col>
                    <a-col :span="8">
                        <a-card title="外部IO" :bordered="false">
                          <a-divider orientation="left">IO电平</a-divider> 
                        <ul >Latch {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_latch}}</ul>
                        <ul >PWM{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_pwm}}</ul>
                        <ul >GATE{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_gate}}</ul>
                        <ul >TRIG{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_trig}}</ul>
                        <ul >PRR{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_prr}}</ul>
                        <ul >水流量计{{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_alarm_cooler}}</ul>

                        <a-divider orientation="left">信号值</a-divider> 
                        <ul >数字8BIT {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_data_in}}</ul>
                        <ul >模拟8BIT {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_anlog_data}}</ul>
                        <ul >水流量 {{info.Packet['Femto_input_reg']['Femto_input_reg_db25'].Ext_anlog_data}}</ul>
                        </a-card>
                       
                    </a-col>
                    <a-col :span="8">
                              <a-card  title=" 频率检测 ">
                                  <a-tabs  v-model:activeKey=selectindex>
                                  <a-tab-pane v-for="(dev, index) in info.Packet['Femto_input_reg']['Femto_input_reg_monitor']['Pd_freq']" :key="index" :tab="'测频'+index">
                                    <ul >通道{{selectindex}}: {{ dev }}</ul>
                                  </a-tab-pane>
                              
                                  </a-tabs>
                            </a-card >
                        </a-col>
                    </a-row>
                </div>
    

         


            <div style="background-color: #f0f6f7; padding: 20px">
              
            
                    <a-card  title=" 泵浦监测 ">
                        <a-tabs  v-model:activeKey=selectindex1>
                            <a-tab-pane v-for="(dev, index) in info.Packet['Femto_input_reg']['Femto_input_reg_monitor']['Femto_input_reg_monitor_pump']" :key="index" :tab="'泵浦'+index">
                        <ul > 温度:{{ info.Packet['Femto_input_reg']['Femto_input_reg_monitor']['Temp'][selectindex1]}}</ul>
                        <ul > 电压:{{ info.Packet['Femto_input_reg']['Femto_input_reg_monitor']['Vol'][selectindex1] }}</ul>
                        <ul > 电流:{{ info.Packet['Femto_input_reg']['Femto_input_reg_monitor']['Femto_input_reg_monitor_pump'][selectindex1]['Actual_cur'] }}</ul>
                    </a-tab-pane>
                        </a-tabs>
              </a-card >




      <div class="antd-pro-pages-dashboard-analysis-twoColLayout">
      <a-row :gutter="24" type="flex" :style="{ marginTop: '24px' }">
        <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
          <a-card title="Card title" :bordered="false">
          <p>card content</p>
        </a-card>
        </a-col>
        <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
          <a-card title="Card title" :bordered="false">
          <p>card content</p>
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
import { Device } from '@/types/api';
import { onUpdated, ref, computed } from 'vue';
import DeviceData from './DeviceData.vue';
import TimeVue from './Time.vue'
import StatusVue from './DeviceStatus.vue';
import DeviceCurrentVue from './DeviceCurrent.vue';
import Geer from '@/components/widget/svg/geer.vue'
import TitleCard from './TitleCard.vue';
import Tree from  '@/components/treelist/TreeList.vue';
import { dataTool } from 'echarts';


const selectindex =ref(0)
const selectindex1 =ref(0)
const props = defineProps<{ info: Device }>()

const update_ref = ref()
onUpdated(() => {

    update_ref.value.classList.add('update')
    setTimeout(() => update_ref.value.classList.remove('update'), 150)

})



const searchUserData :any = []
for (let i = 0; i < 7; i++) {
  searchUserData.push({
    x: moment().add(i, 'days').format('YYYY-MM-DD'),
    y: Math.ceil(Math.random() * 10)
  })
}


const sourceData = [
  { item: '家用电器', count: 32.2 },
  { item: '食用酒水', count: 21 },
  { item: '个护健康', count: 17 },
  { item: '服饰箱包', count: 13 },
  { item: '母婴产品', count: 9 },
  { item: '其他', count: 7.8 }
]


// const dv = new DataSet.View().source(sourceData)
// dv.transform({
//   type: 'percent',
//   field: 'count',
//   dimension: 'item',
//   as: 'percent'
// })
// const pieData = dv.rows
const barData  :any= []
const barData2 :any= []

const pieScale = [{
  dataKey: 'percent',
  min: 0,
  formatter: '.0%'
}]
const pieStyle=[{
        stroke: '#fff',
        lineWidth: 1
      }]

const searchUserScale = [
  {
    dataKey: 'x',
    alias: '时间'
  },
  {
    dataKey: 'y',
    alias: '用户数',
    min: 0,
    max: 10
  }]

  const searchData :any= []
for (let i = 0; i < 50; i += 1) {
  searchData.push({
    index: i + 1,
    keyword: `搜索关键词-${i}`,
    count: Math.floor(Math.random() * 1000),
    range: Math.floor(Math.random() * 100),
    status: Math.floor((Math.random() * 10) % 2)
  })
}


for (let i = 0; i < 12; i += 1) {
  barData.push({
    x: `${i + 1}月`,
    y: Math.floor(Math.random() * 1000) + 200
  })
  barData2.push({
    x: `${i + 1}月`,
    y: Math.floor(Math.random() * 1000) + 200
  })
}

const rankList :any= []
for (let i = 0; i < 7; i++) {
  rankList.push({
    name: '白鹭岛 ' + (i + 1) + ' 号店',
    total: 1234.56 - i * 100
  })
}


// const searchTableColumns=  computed(() => {
   
//         return [
//       {
//         dataIndex: 'index',
//         title: "1",
//         width: 90
//       },
//       {
//         dataIndex: 'keyword',
//         title:"2"
//       },
//       {
//         dataIndex: 'count',
//         title: "3"
//       },
//       {
//         dataIndex: 'range',
//         title:'4',
//         align: 'right',
//         sorter: (a:any, b:any) => a.range - b.range,
//         scopedSlots: { customRender: 'range' }
//       }
//       ]
//     }
// );


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
        display: grid;
        grid-template-columns: 32px 160px auto 32px;
        line-height: 32px;
        text-align: center;
        border-bottom: 1px solid @border-color;

        .update-trigger {
            background-color: #000;
            display: inline-block;
            width: .8rem;
            height: .8rem;
            line-height: 32px;
            margin-right: .3rem;
            border-radius: 50%;

            &.update {
                background: aquamarine;
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
</style>
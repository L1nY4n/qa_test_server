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
            <!-- <TimeVue class="time" :time="info.Packet['系统监控']['激光器时间监测']" />
            <Alarm  class="alarm"  :alarm=" info.Packet['系统监控']['激光器告警监测']"/> -->

            <TitleCard title="版本信息">
                <DeviceBate :datas="info.Packet['Femto_input_reg']['Femto_input_reg_bate']" />
            </TitleCard> 

            <TitleCard title="外部IO">
                <DeviceBate :datas="info.Packet['Femto_input_reg']['Femto_input_reg_db25']" />
            </TitleCard>

             <TitleCard title="FPGA测试">
                <DeviceBate :datas="info.Packet['Femto_input_reg']['Femto_input_reg_fpga_debug']" />
            </TitleCard>

            <TitleCard title="时间统计">
                <DeviceBate :datas="info.Packet['Femto_input_reg']['Femto_input_reg_time']" />
            </TitleCard>

            <TitleCard title="告警">
                <DeviceBate :datas="info.Packet['Femto_input_reg']['Femto_input_reg_alarm']" />
            </TitleCard>
            <!-- <TitleCard title="温度监测">
                <DeviceVoltage :datas="info.Packet['系统监控']['激光器温度监测']" />
            </TitleCard>
              <TitleCard title="温控板">
              <DeviceTempBoard :datas="info.Packet['系统监控']['激光器温控板监测']" />
            </TitleCard>
              <TitleCard title="种子模块">
              <DeviceSeedModule  :datas="info.Packet['系统监控']['激光器种子模块监测']"/>
           
            </TitleCard>
              <TitleCard title="FPGA寄存器">
              <DeviceFPGA :datas="info.Packet['系统监控']['激光器FPGA寄存器监测']" />
            </TitleCard> --> 
          
            
        </div>

    </div>

</template>
<script setup lang="ts">
import { Device } from '@/types/api';
import { onUpdated, ref, computed } from 'vue';
import DeviceData from './DeviceData.vue';
import TimeVue from './Time.vue'
import StatusVue from './DeviceStatus.vue';
import DeviceDB25 from './DeviceDB25.vue';
import DeviceBate from './DeviceBate.vue';
import DeviceCurrentVue from './DeviceCurrent.vue';
import Geer from '@/components/widget/svg/geer.vue'
import TitleCard from './TitleCard.vue';
import DeviceVoltage from './DeviceVoltage.vue';
import DeviceFPGA from './DeviceFPGA.vue';
import DeviceSeedModule from './DeviceSeedModule.vue';
import DeviceTempBoard from './DeviceTempBoard.vue';
import Alarm from './Alarm.vue';

const props = defineProps<{ info: Device }>()

const update_ref = ref()
onUpdated(() => {

    update_ref.value.classList.add('update')
    setTimeout(() => update_ref.value.classList.remove('update'), 150)

})

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
    background: #1d1e22;
    box-shadow: 0 4px 30px rgb(0 0 0 / 50%);

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
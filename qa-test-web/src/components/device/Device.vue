<template>
    <div class="device">
        <header>
            <div> 
            <span ref="update_ref" class='update-trigger'></span>
            </div>

            <h1 class="sn"> {{ info.Sn }}</h1>
            <h1 class="name">{{ info.Name }}</h1>
            <div class="geer"> 
             <Geer />
            </div>
          
        </header>
       <div class="time">
        <span>{{collect_time}}</span>
       </div>
        <DeviceData :info="info.Packet['系统监控']" />
    </div>

</template>
<script setup lang="ts">
import { Device } from '@/types/api';
import { onUpdated, ref,computed } from 'vue';
import DeviceData from './DeviceData.vue';
import Geer from '@/components/widget/svg/geer.vue'
const props = defineProps<{ info: Device }>()

const update_ref = ref()
onUpdated(() => {

    update_ref.value.classList.add('update')
    setTimeout(() => update_ref.value.classList.remove('update'), 150)

})

let collect_time = computed(()=>{
     const p = props.info.Packet as any
     const t = p ['系统监控']['激光器时间监测']
     const {'年': year,'月':mon,'日':day,'时': hour,'分': min,'秒': sec} = t
     return `${year}-${mon}-${day} ${hour}:${min}:${sec}`
})

</script>
<style lang="less" scoped>

@border-color: rgba(255, 255, 255, .1);
.device {

    align-items: center;
    border: 1px solid  @border-color;
    border-radius: 4px;
    background: #1d1e22;
    box-shadow: 0 4px 30px rgb(0 0 0 / 50%);

    header {
        display: grid;
        grid-template-columns: 32px 80px auto 32px;
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
            color: #d2d232;




        }
    }

    &.has-failed {
        box-shadow: 0px 4px 15px rgba(0, 0, 0, .2);
        border-color: #d22c32;
        animation: alertblink 2s ease-in-out infinite;

        &:hover {
            background-color: rgba(210, 44, 50, .2);
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
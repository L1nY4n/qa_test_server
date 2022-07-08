<template>
    <div class="device">
        <header>
            <div class="sn">{{ info.Sn }}</div>
            <div class="name">{{ info.Name }}</div>
        </header>
        <!-- <DeviceInfo :info="info.Packet['采集设备信息']" /> -->
        <DeviceData :info="info.Packet['系统监控']" />
    </div>

</template>
<script setup lang="ts">
import { Device } from '@/types/api';
//import DeviceInfo from './DeviceInfo.vue';
import DeviceData from './DeviceData.vue';
defineProps<{ info: Device }>()
</script>
<style lang="less" scoped>
.device {

    align-items: center;
    padding: 4px;
    border: 1px solid rgba(255, 255, 255, .1);
    border-radius: 4px;
    white-space: nowrap;
    background-color: rgba(255, 255, 255, .1);
    height: 500px;
    overflow: hidden auto;

    header {
        position: sticky;
        top: 0;
        display: flex;
        justify-content: space-between;
        height: 24px;
        line-height: 24px;
        padding: 0 4px;
        background-color: rgb(29, 89, 98);

        .name {

            font-weight: 600;
            color: #2eb35a;
        }

        .sn {
            color: #d2d232;

            &::before {
                content: '';
                display: inline-block;
                width: 0.6rem;
                height: 0.6rem;
                background-color: currentColor;
                border-radius: 50%;
                margin-right: 0.5rem;
                margin-top: -0.2rem;
                vertical-align: middle;
            }
        }
    }

    .devices-details {
        flex: 1;
        display: block;
        list-style: none;
        margin: 0;
        padding: 0;


        li {
            display: block;
            line-height: 1.5;
            color: #7e8794;

            .data {
                display: block;
                margin: -1.7rem 0 0 0;
                padding: 0 1rem 0 0;
                font-weight: 600;
                text-align: right;
                color: #c1c6cb;
            }

            &:last-child .data {
                font-weight: normal;
                color: rgba(230, 245, 255, .32);
            }
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
</style>
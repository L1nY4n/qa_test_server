<template>
    <div class="device">
        <header>
            <svg width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                    d="M15 6.675l-1.8-.6c-.2-.1-.3-.3-.2-.4l.9-1.7c.6-1.2-.7-2.5-1.9-1.9l-1.7.9c-.1.1-.3 0-.4-.2l-.6-1.8c-.4-1.3-2.2-1.3-2.6 0l-.6 1.8c-.1.2-.3.3-.4.2l-1.7-.9c-1.2-.6-2.5.7-1.9 1.9l.9 1.7c.1.1 0 .3-.2.4l-1.8.6c-1.3.4-1.3 2.3 0 2.7l1.8.6c.2 0 .3.2.2.3l-.9 1.7c-.6 1.2.7 2.5 1.9 1.9l1.7-.9c.2-.1.4 0 .4.2l.6 1.8c.4 1.3 2.3 1.3 2.7 0l.6-1.8c.1-.2.3-.3.4-.2l1.7.9c1.2.6 2.5-.7 1.9-1.9l-1-1.7c-.1-.2 0-.4.2-.4l1.8-.6c1.3-.4 1.3-2.2 0-2.6zm-7 3.7c-1.3 0-2.4-1.1-2.4-2.4 0-1.3 1.1-2.4 2.4-2.4 1.3 0 2.4 1.1 2.4 2.4 0 1.3-1.1 2.4-2.4 2.4z"
                    fill="#4C4F5A"></path>
            </svg>
            <h1 class="sn">{{ info.Sn }}</h1>
            <h1 class="name">{{ info.Name }}</h1>
            <svg width="16" height="8" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                    d="M8.709 7.651l6.161-5.622c.241-.22.383-.517.383-.84 0-.323-.142-.62-.383-.84A1.361 1.361 0 0 0 13.95 0c-.354 0-.68.13-.921.349l-5.27 4.808L2.492.349A1.361 1.361 0 0 0 1.57 0C1.215 0 .89.13.648.336.38.569.253.879.253 1.189c0 .297.127.595.368.84 1.615 1.486 5.807 5.325 6.09 5.596l.03.026c.509.465 1.458.465 1.968 0z"
                    fill="#4C4F5A"></path>
            </svg>
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
    background: #1d1e22;
    box-shadow: 0 4px 30px rgb(0 0 0 / 50%);

    header {

        display: grid;
        grid-template-columns: 36px  50px auto 36px;
        margin: 0;
        padding: 0;
        text-align: center;
       svg {
        margin: 10px;
       }

       h1{
            font-family: Lato,sans-serif;
    font-weight: 900;
    letter-spacing: 1.57px;
    line-height: 36px;
    margin: 0;
       }
        .name {

          text-align: right;
            color: #c5c8d4;
        }

        .sn {
            color: #d2d232;

    //       &::before {
    //           content: '';
    //           display: inline-block;
    //           width: 0.6rem;
    //           height: 0.6rem;
    //           background-color: currentColor;
    //           border-radius: 50%;
    //           margin-right: 0.5rem;
    //           margin-top: -0.2rem;
    //           vertical-align: middle;

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
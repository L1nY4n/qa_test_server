<template>
    <div class="wrapper">
        <template v-for="(v, k) in datas">

            <div v-if="k === '种子开关'||k === '告警值'" class="status-type">
                <div class="label">{{ k }}</div>
                <div>
                    <div class="status" :class="{ on: v == 1 }"></div>
                </div>
            </div>
            <div v-else-if="k === '预留参数'"></div>
            <div v-else class="value-type">
               <span class="label">{{ k }}</span> <span class="value"> {{ v }}</span>
            </div>
        </template>
    </div>
</template>
<script setup lang="ts">

import { PacketSeedModule } from '@/types/api';
defineProps<{ datas: PacketSeedModule }>()
</script>
<style lang="less" scoped>
.wrapper {
    font-size: 75%;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2px;

    .value-type {
        display: flex;
        align-items: center;
        justify-content: space-around;
        border: 1px solid #333;
     .label {
                color: #777;
                max-width: 100%;
                overflow: hidden;
                white-space: nowrap;
                text-overflow: ellipsis;
            }
        .value {
            color: rgb(22, 141, 206);
            font-weight: 500;
        }
    }

    .status-type {
        display: grid;
        grid-template-columns: 1fr 24px;
        height: 24px;
        line-height: 24px;
        border: 1px solid #333;

        &>div {
            display: flex;
            align-items: center;
            justify-content: space-around;
        }

        .label {
            color: #777;
        }

        .status {
            width: 10px;
            height: 10px;
            border-radius: 30px;
            background: black;

            &.on {
                background: #26f260;
            }
        }
    }

}
</style>
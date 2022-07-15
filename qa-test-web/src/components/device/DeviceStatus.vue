<template>
    <div class="status-wrapper">
        <template v-for="(v, k) in status">
            <div v-if="k === '工作模式'" class="mode">
                <b>工作模式: </b> <span class="mode-value"> {{ v }}</span>
            </div>
            <div v-else class="status">
                <div class="label">{{ k }}</div>
                <div>
                    <div class="circle" :class="{ on: v == 1 }"></div>
                </div>

            </div>

        </template>
    </div>
</template>
<script setup lang="ts">

import { PacketStatus } from '@/types/api';
defineProps<{ status: PacketStatus }>()
</script>
<style lang="less" scoped>
.status-wrapper {
    font-size: 75%;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap:2px;

    .mode {
        // grid-column-start: 1;
        // grid-column-end: 3;
        display: flex;
        align-items: center;
        justify-content: space-around;
        .mode-value{
            color:cornflowerblue;
            font-size: 14px;
            font-weight: 500;
        }
    }

    .status {
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

        .circle {
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
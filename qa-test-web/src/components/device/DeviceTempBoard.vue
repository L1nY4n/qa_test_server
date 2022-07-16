<template>
    <div class="wrapper">
        <div class="board" v-for="board in datas">
            <template v-for="(v, k) in board">
                <div v-if="k === '开关'" class="status-type">
                    <div class="label">{{ k }}</div>
                    <div>
                        <div class="status" :class="{ on: v == 1 }"></div>
                    </div>
                </div>
                <div v-else class="value-type">
                    <span class="label">{{ k }}</span> <span class="value"> {{ v }}</span>
                </div>
            </template>
        </div>
    </div>

</template>
<script setup lang="ts">

import { PacketTempBoard } from '@/types/api';
defineProps<{ datas: PacketTempBoard[] }>()
</script>
<style lang="less" scoped>
.wrapper {
    font-size: 75%;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 4px;

    .board {
        border: 1px solid #222;
        padding: 2px;

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
                color: rgb(22, 206, 154);
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
                max-width: 100%;
                overflow: hidden;
                white-space: nowrap;
                text-overflow: ellipsis;
            }

            .status {
                width: 10px;
                height: 10px;
                border-radius: 30px;
                background: black;

                &.on {
                    background: #eff311;
                }
            }
        }
    }


}
</style>
<template>
    <div class="wrapper">
        <template v-for="(v, k) in datas">
            <div class="current">
                <div class="label">{{ k }}</div>
                <div>
                    <div class="value" >{{ v }}</div>
                </div>
            </div>

        </template>
    </div>
</template>
<script setup lang="ts">


import { Packet_bate } from '@/types/api';
defineProps<{ datas: any }>()
</script>
<style lang="less" scoped>
.wrapper {
    font-size: 75%;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 2px;

    .current {
        display: grid;
        grid-template-columns:  repeat(2, 50%);
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

        .value {
            color: greenyellow;
        }
    }

}
</style>
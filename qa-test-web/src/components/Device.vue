<template>
    <div class="wrapper">
        <div class="header">{{ info.Name }}[{{ info.Sn }}]</div>
        <div class="packet">{{ info.Packet }}
        </div>
    </div>

</template>
<script setup lang="ts">
import { Device } from '@/types/api';
defineProps<{ info: Device }>()
</script>
<style lang="less" scoped>
.wrapper {
    width: 400px;
    background-color: aliceblue;
    .header{
        background-color: cornflowerblue;
    }
    .packet{
        height: 400px;
        overflow: auto;
    }
}
</style>
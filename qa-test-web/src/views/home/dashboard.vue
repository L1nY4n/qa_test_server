<template>
   <h1>设备列表</h1>
   <div class="list">
      <DeviceVue v-for="dev in state.list" :info="dev" />
   </div>
</template>


<script setup lang="ts">
import * as API from '@/api'
import { Device } from '@/types/api';
import DeviceVue from '@/components/Device.vue';
import { onMounted, reactive } from 'vue'

let state = reactive<{ list: Device[] }>({ list: [] })

onMounted(() => {
   setInterval(async () => {
      state.list = await API.device.list()

   }, 5000)

})


</script>
<style lang="less">
.list {
   display: flex;
   gap: 10px;
}
</style>


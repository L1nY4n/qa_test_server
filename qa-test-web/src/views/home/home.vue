<template>
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

const get_list = async () => {
   state.list = await API.device.list()
}

onMounted(() => {
   get_list()
   setInterval(get_list, 10000);
})


</script>
<style lang="less">
.list {
   display: grid;
	 grid-template-columns: repeat(auto-fill, minmax(200px, 320px));
	 grid-gap: 8px;
   height: 100%;
   overflow: auto;
}
</style>


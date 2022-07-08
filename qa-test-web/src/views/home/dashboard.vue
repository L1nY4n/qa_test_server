<template>
<h1>dashboard</h1>
<div>
   <DeviceVue v-for="dev in deviceList" :info="dev"/>
</div>
</template>


<script setup lang="ts">
import * as API from '@/api'
import { Device } from '@/types/api';
import DeviceVue from '@/components/Device.vue';
import {onMounted, reactive} from 'vue'

let deviceList = reactive<Device[]>([])

onMounted(async ()=>{
   deviceList = await API.device.list()
})


</script>


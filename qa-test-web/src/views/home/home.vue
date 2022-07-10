<template>
  <div class="wrapper">
    <div class="option">
      <span>轮训间隔：<Tag>{{ interval }}</Tag></span>
      <Slider v-model:value="interval" :min="5" :max="60" size="small" @change="resetInterval" />

    </div>
    <div class="list">
      <DeviceVue v-for="(dev, i) in state.list" :info="dev" :key="i" />
    </div>
  </div>
</template>


<script setup lang="ts">
import * as API from '@/api'
import { Slider, Tag } from 'ant-design-vue'
import { Device } from '@/types/api';
import DeviceVue from '@/components/Device.vue';
import CreateWebSocket from "@/utils/ws";
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'

let state = reactive<{ list: Device[] }>({ list: [] })
let ws: WebSocket
let interval = ref<number>(10)
let timer = -1
const get_list = async () => {
  state.list = await API.device.list()
}
const webSocketConnect = () => {
  // 订阅 ws 设备更新类型的事件
  ws = CreateWebSocket("/realtime/register/device_upload", "test");
  ws.onmessage = (event) => {
    const reader = new FileReader()
    reader.readAsText(event.data, 'utf-8')
    reader.onload = () => {
      const update_dev: Device = JSON.parse(reader.result as string)
      console.log(update_dev)
      // 根据sn 查找和替换数组设备
      const index = state.list.findIndex((dev) => dev.Sn = update_dev.Sn)
      if (index > -1) {
        state.list.splice(index, 1, update_dev)
      }

    }

  };
  ws.onopen = () => {
    console.log("ws opened")
  };
  ws.onclose = (e) => {
    console.log(
      "Socket is closed. Reconnect will be attempted in 10 second.",
      e.reason
    );
    setTimeout(() => {
      webSocketConnect();
    }, 10 * 1000);
  };
  ws.onerror = (err) => {
    console.log(
      "Socket encountered error: ",
      err);
  };
}
webSocketConnect()

onMounted(() => {
  get_list()
  timer = setInterval(get_list, interval.value * 1000);
})

onBeforeUnmount(() => {
  ws.close()
})


const resetInterval = (inteval: any) => {
  clearInterval(timer)
  timer = setInterval(get_list, inteval * 1000);
}

</script>
<style lang="less">
.wrapper {
  position: relative;
  height: 100%;

  .option {
    height: 60px;
    width: 320px;
  }

  .list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 300px));
    grid-gap: 8px;
    height: calc(100% - 60px);
    overflow: auto;
  }
}
</style>


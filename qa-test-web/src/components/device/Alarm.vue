<template>
    <div class="wrapper">
        <sub class="title">告警 </sub>
        <sub class="itme" :class="{ alarm: alarm.当前告警 > 0 }">当前:{{ alarm.当前告警 }}</sub>
        <sub class="itme" :class="{ alarm: alarm.生效告警 > 0 }">生效:{{ alarm.生效告警 }}</sub>
        <sub class="itme" :class="{ alarm: alarm.历史告警 > 0 }">历史:{{ alarm.历史告警 }}</sub>
    </div>
</template>
<script lang="ts"  setup>
import { PacketAlarm } from '@/types/api';
defineProps<{ alarm: PacketAlarm }>()
</script>
<style lang="less" scoped>
.wrapper {
    color: rgba(10, 118, 8, 0.397);
}

.itme {
    margin-right: 4px;
    &.alarm{
           color: rgb(212, 0, 0);
          border: 1px dashed rgba(176, 33, 33, 0.433);
    }
  
}
</style>
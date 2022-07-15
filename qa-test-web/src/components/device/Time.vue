<template>
    <div>{{ t.year }}-{{ t.mon }}-{{ t.day }} {{ t.hour }}:{{ t.min }}:{{ t.sec }}</div>
</template>
<script lang="ts"  setup>
import { PacketTIme } from '@/types/api';
import {computed} from 'vue'
// defineProps<{time: {
//     year: number,
//     mon: number
//     day: number
//     hour: number
//     min:number,
//     sec: number
// }}>()


const props = defineProps<{time: PacketTIme}>()

const t = computed(()=>{

    
    const { '年': year, '月': mon, '日': day, '时': hour, '分': min, '秒': sec } = props.time
    return {year,mon,day,hour,min,sec}
})

</script>
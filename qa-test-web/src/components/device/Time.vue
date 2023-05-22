<template>
    <div class="time">{{ t.year }}-{{ t.mon }}-{{ t.day }} {{ t.hour }}:{{ t.min }}:{{ t.sec }}</div>
</template>
<script lang="ts"  setup>
import { PacketTIme } from '@/types/api';
import { computed } from 'vue'
// defineProps<{time: {
//     year: number,
//     mon: number
//     day: number
//     hour: number
//     min:number,
//     sec: number
// }}>()


const props = defineProps<{ time: PacketTIme }>()

const t = computed(() => {


    const { 'Year': year, 'Mon': mon, 'Day': day, 'Hour': hour, 'Minute': min, 'Second': sec } = props.time
    return { year, mon, day, hour, min, sec }
})

</script>
<style lang="less" scoped>
.time {
    padding: 2px 4px;
    color: rgb(0, 255, 0);
    font-family: "Graduate", cursive;
    margin: auto;
    //text-shadow: 1px 1px 2px rgb(0, 255, 0);
}
</style>
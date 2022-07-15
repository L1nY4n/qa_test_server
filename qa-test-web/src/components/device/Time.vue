<template>
 <span>{{time.year}}-{{time.mon}}</span>
</template>
<script lang="ts"  setup>

defineProps<{time: {
    year: number,
    mon: number
    day: number
    hour: number
    min:number,
    sec: number
}}>()
</script>
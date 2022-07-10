<template>
    <ul class="tree">
        <Item  :title="'设备属性'" :data="data"></Item>
    </ul>
</template>
<script setup lang="ts">
import Item from './Item.vue'
    defineProps<{data: any}>()
</script>
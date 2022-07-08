<template>
    <ul class="tree">
        <Item  :title="'fdfad'" :data="data"></Item>
    </ul>
</template>
<script setup lang="ts">
import Item from './Item.vue'
    defineProps<{data: any}>()
</script>
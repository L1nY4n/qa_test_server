<template>
<div class="wrapper">
<TreeList  :data="info"/>
</div>

</template>
<script lang="ts" setup>
import TreeList from './treelist/TreeList.vue'
defineProps<{info: object|any}>()
</script>
<style>
.wrapper{
   
}
</style>
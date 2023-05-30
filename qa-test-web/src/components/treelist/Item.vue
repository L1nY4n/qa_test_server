<template>
    <li @click.stop="toggle">
       <span class="title">{{title}}</span>
        <ul v-if="Array.isArray(data) || (typeof data === 'object')" v-show="open">
            <Item v-for="(node, index) in data" :title="index" :data="node"></Item>
        </ul>
         <span v-else class="value">{{data}}</span>
    </li>
</template>
<script lang="ts">
export default {
    name: "Item"
}
</script>
<script lang="ts" setup>
import { ref } from 'vue'

 defineProps<{data: any,title: any}>()
const open = ref<boolean>(true)
const toggle = () => {
    open.value = !open.value
}
</script>
<style lang="less" scoped>

// ul {
//   display: flex;
//   flex-direction: row;
// } 

.title{
    color: #e1ca72;
    margin-right: 4px
}
.value{
    color: #74b087;
}
</style>
<template>
    <div class="wrapper">
        <div class="title">
            {{ title }}
        </div>
        <div class="content">
            <slot></slot>
        </div>
    </div>
</template>
<script setup lang="ts">
defineProps<{ title: string }>()
</script>
<style lang="less" scoped>
.wrapper {
    box-sizing: border-box;
    margin: 6px;
    
    .title {
        position: relative;
        display: inline-block;
        width: auto;
        box-shadow: inset 1px 1px 1px #ffffff ;
        text-align: center;
        padding: 2px 6px;
        font-size: 12px;
        color: rgb(7, 176, 21);
        background: hwb(0 0% 100% / 0);
        //border-radius: 5px 5px 0 0;
        border: none;
    }

    .content {
        padding: 10px;
        background: #ffffff;
        color:#000000;
    }
}
</style>
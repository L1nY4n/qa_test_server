<template>
<ul>
    <li v-for="(v,k) in info">
    {{k}}: {{v}}
    </li>
</ul>
</template>
<script lang="ts" setup>
defineProps<{info: object|any}>()
</script>
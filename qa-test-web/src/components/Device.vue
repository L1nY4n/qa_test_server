<template>
<div>{{info.name}}</div>
<div>{{info.sn}}</div>
</template>
<script setup lang="ts">
import { Device } from '@/types/api';
defineProps<{ info: Device }>()
</script>
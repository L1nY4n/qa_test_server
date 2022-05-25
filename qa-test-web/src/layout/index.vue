<template>
  <a-layout>
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible>
      <div class="logo">
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @click="onMenuClick">
        <a-sub-menu v-for="(group, i) in menuOptions" :key="i">
          <template #icon>
            <SettingOutlined />
          </template>
          <template #title>{{ group.meta?.title }}</template>

          <a-menu-item v-for="(item, k) in group.children" :key="item.name">
            <template #icon>
              <SettingOutlined />
            </template>
            <span>{{ item.meta?.title }}</span>
          </a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0">
        <menu-unfold-outlined v-if="collapsed" class="trigger" @click="() => (collapsed = !collapsed)" />
        <menu-fold-outlined v-else class="trigger" @click="() => (collapsed = !collapsed)" />
      </a-layout-header>
      <a-layout-content :style="{ margin: '8px', padding: '6px', background: '#fff', height: 'calc(100vh - 80px)' }">
        <router-view #="{ Component }">
          <component :is="Component" />
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
<script lang="ts" setup>
import { menus } from '@/router/routes'
import { defineComponent, ref, h } from 'vue';
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { MailOutlined, AppstoreOutlined, SettingOutlined,MenuFoldOutlined,MenuUnfoldOutlined } from '@ant-design/icons-vue';
const selectedKeys = ref<string[]>(['1'])
const collapsed = ref<boolean>(false)
const menuOptions = menus
const router = useRouter()
const currentRoute = useRoute();
const onMenuClick = ({ key }) => {
  if (key === currentRoute.name) return;
  if (/http(s)?:/.test(key)) {
    window.open(key);
  } else {
    console.log(key)
    router.push({ name: key}).catch(err=>{
      console.log(err)
    });
  }
}
</script>
<style>
.trigger {
  font-size: 18px;
  line-height: 64px;
  padding: 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

.logo {
  height: 32px;
  background: rgba(255, 255, 255, 0.3);
  margin: 16px;
}
</style>

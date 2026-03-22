<template>
  <a-layout class="app-layout">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      class="app-sider"
      collapsible
      :width="228"
      :collapsed-width="64"
    >
      <div class="logo">{{ APP_TITLE }}</div>
      <a-menu
        v-model:openKeys="openKeys"
        v-model:selectedKeys="selectedKeys"
        theme="dark"
        mode="inline"
        @click="onMenuClick"
      >
        <a-sub-menu v-for="group in visibleMenus" :key="String(group.name)">
          <template #icon>
            <component :is="resolveIcon(String(group.meta?.icon || group.name || ''))" />
          </template>
          <template #title>{{ group.meta?.title }}</template>

          <a-menu-item v-for="item in group.children" :key="String(item.name)">
            <template #icon>
              <component :is="resolveIcon(String(item.meta?.icon || item.name || ''))" />
            </template>
            <span>{{ item.meta?.title }}</span>
          </a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>

    <a-layout class="app-main-layout">
      <a-layout-header class="app-header">
        <menu-unfold-outlined
          v-if="collapsed"
          class="trigger"
          @click="() => (collapsed = !collapsed)"
        />
        <menu-fold-outlined
          v-else
          class="trigger"
          @click="() => (collapsed = !collapsed)"
        />

        <div class="header-title">{{ headerTitle }}</div>

        <div class="header-right">
          <a-tag color="blue">{{ roleLabel }}</a-tag>
          <a-dropdown>
            <span class="user-entry">
              <user-outlined />
              <span>{{ currentUser?.displayName || currentUser?.username || '未登录' }}</span>
            </span>
            <template #overlay>
              <a-menu>
                <a-menu-item key="user-name">{{ currentUser?.username || '-' }}</a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="logout">
                  <logout-outlined />
                  退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <a-layout-content class="main-content">
        <router-view #="{ Component }">
          <component :is="Component" />
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script lang="ts" setup>
import { menus } from '@/router/routes'
import { computed, ref, watch } from 'vue'
import { useRouter, useRoute, type RouteRecordRaw } from 'vue-router'
import {
  AppstoreOutlined,
  AreaChartOutlined,
  BarChartOutlined,
  ControlOutlined,
  DashboardOutlined,
  AlertOutlined,
  LineChartOutlined,
  KeyOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MonitorOutlined,
  RadarChartOutlined,
  SettingOutlined,
  TeamOutlined,
  ToolOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { clearSession, getCurrentUser, hasAnyRole } from '@/utils/auth'

const APP_TITLE = '\u6fc0\u5149\u5668\u7ba1\u7406\u7cfb\u7edf'

const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>([])
const collapsed = ref<boolean>(false)
const router = useRouter()
const currentRoute = useRoute()

const iconMap: Record<string, any> = {
  MonitorOutlined,
  DashboardOutlined,
  AppstoreOutlined,
  LineChartOutlined,
  AlertOutlined,
  RadarChartOutlined,
  BarChartOutlined,
  AreaChartOutlined,
  SettingOutlined,
  ControlOutlined,
  ToolOutlined,
  KeyOutlined,
  TeamOutlined,
}

const currentUser = ref(getCurrentUser())

const visibleMenus = computed<RouteRecordRaw[]>(() => {
  const role = currentUser.value?.role
  return menus
    .map((group) => {
      const children = (group.children || []).filter((item) => hasAnyRole(item.meta?.roles, role))
      return {
        ...group,
        children,
      }
    })
    .filter((group) => hasAnyRole(group.meta?.roles, role) && (group.children || []).length > 0)
})

const topMenuNameSet = computed(() => new Set(visibleMenus.value.map((group) => String(group.name || ''))))

const headerTitle = computed(() => {
  const title = currentRoute.meta?.title
  if (typeof title === 'string' && title.trim()) {
    return title
  }
  return APP_TITLE
})

const roleLabel = computed(() => {
  const role = currentUser.value?.role
  if (role === 'admin') {
    return '管理员'
  }
  if (role === 'operator') {
    return '运维'
  }
  if (role === 'viewer') {
    return '只读'
  }
  return '未登录'
})

const resolveIcon = (key: string) => {
  return iconMap[key] || AppstoreOutlined
}

const syncMenuState = () => {
  currentUser.value = getCurrentUser()

  selectedKeys.value = currentRoute.name ? [String(currentRoute.name)] : []
  const topMenuSet = topMenuNameSet.value
  const matchedTopMenu = currentRoute.matched.find((record) => {
    if (!record.name) {
      return false
    }
    return topMenuSet.has(String(record.name))
  })
  openKeys.value = matchedTopMenu?.name ? [String(matchedTopMenu.name)] : []
}

watch(
  () => currentRoute.fullPath,
  () => {
    syncMenuState()
  },
  { immediate: true }
)

const onMenuClick = (item: { key: string }) => {
  const key = String(item.key)
  if (key === String(currentRoute.name || '')) {
    return
  }

  if (/http(s)?:/.test(key)) {
    window.open(key)
    return
  }

  router.push({ name: key }).catch(() => undefined)
}

const logout = () => {
  clearSession()
  message.success('已退出登录')
  router.replace({ name: 'login' }).catch(() => undefined)
}
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
}

.app-sider {
  box-shadow: 2px 0 14px rgba(15, 23, 42, 0.22);
}

.logo {
  height: 44px;
  margin: 14px;
  border-radius: 10px;
  background: linear-gradient(120deg, #2d485f, #3c6e71);
  color: #f8fbff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.8px;
}

.app-main-layout {
  background: transparent;
}

.app-header {
  background: rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(4px);
  border-bottom: 1px solid #e8edf4;
  padding: 0 14px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-title {
  color: #324a5f;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-entry {
  color: #324a5f;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  color: #4a6075;
  transition: color 0.2s;
}

.trigger:hover {
  color: #1677ff;
}

.main-content {
  margin: 10px;
  border-radius: 12px;
  overflow: hidden;
  background: #f7fafc;
  border: 1px solid #e3ebf4;
  box-shadow: 0 8px 24px rgba(26, 45, 67, 0.08);
  height: calc(100vh - 84px);
}

@media (max-width: 900px) {
  .main-content {
    margin: 8px;
    height: calc(100vh - 80px);
  }

  .header-title {
    font-size: 12px;
  }

  .header-right {
    gap: 6px;
  }

  .user-entry span {
    max-width: 84px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>

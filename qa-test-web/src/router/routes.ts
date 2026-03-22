/* eslint-disable @typescript-eslint/promise-function-async */
import { RouteRecordRaw, RouterView } from 'vue-router'

export const menus: RouteRecordRaw[] = [
  {
    name: 'monitor',
    path: 'monitor',
    component: RouterView,
    redirect: { name: 'overview' },
    children: [
      {
        name: 'overview',
        path: 'overview',
        component: () => import('@/views/home/overview.vue'),
        meta: {
          title: '运行总览',
          icon: 'DashboardOutlined',
          roles: ['viewer'],
        },
      },
      {
        name: 'device',
        path: 'device',
        component: () => import('@/views/home/device.vue'),
        meta: {
          title: '设备监控',
          icon: 'AppstoreOutlined',
          roles: ['viewer'],
        },
      },
    ],
    meta: {
      title: '监控中心',
      icon: 'MonitorOutlined',
      roles: ['viewer'],
    },
  },
  {
    name: 'analysis',
    path: 'analysis',
    component: RouterView,
    redirect: { name: 'phm' },
    children: [
      {
        name: 'phm',
        path: 'phm',
        component: () => import('@/views/home/phm.vue'),
        meta: {
          title: 'PHM 健康管理',
          icon: 'RadarChartOutlined',
          roles: ['operator'],
        },
      },
      {
        name: 'alerts',
        path: 'alerts',
        component: () => import('@/views/home/alerts.vue'),
        meta: {
          title: '告警中心',
          icon: 'AlertOutlined',
          roles: ['operator'],
        },
      },
      {
        name: 'trends',
        path: 'trends',
        component: () => import('@/views/home/trends.vue'),
        meta: {
          title: '趋势分析',
          icon: 'AreaChartOutlined',
          roles: ['operator'],
        },
      },
    ],
    meta: {
      title: '运维分析',
      icon: 'LineChartOutlined',
      roles: ['operator'],
    },
  },
  {
    name: 'governance',
    path: 'governance',
    component: RouterView,
    redirect: { name: 'rules' },
    children: [
      {
        name: 'rules',
        path: 'rules',
        component: () => import('@/views/home/rules.vue'),
        meta: {
          title: '规则策略',
          icon: 'ControlOutlined',
          roles: ['operator'],
        },
      },
      {
        name: 'system',
        path: 'system',
        component: () => import('@/views/home/system.vue'),
        meta: {
          title: '系统状态',
          icon: 'ToolOutlined',
          roles: ['operator'],
        },
      },
      {
        name: 'decrypt',
        path: 'decrypt',
        component: () => import('@/views/home/decrypt.vue'),
        meta: {
          title: '时间密钥解密',
          icon: 'KeyOutlined',
          roles: ['operator'],
        },
      },
      {
        name: 'users',
        path: 'users',
        component: () => import('@/views/home/users.vue'),
        meta: {
          title: '用户管理',
          icon: 'TeamOutlined',
          roles: ['admin'],
        },
      },
    ],
    meta: {
      title: '策略与系统',
      icon: 'SettingOutlined',
      roles: ['operator'],
    },
  },
]

const routes: RouteRecordRaw[] = [
  {
    name: 'login',
    path: '/login',
    component: () => import('@/views/login.vue'),
    meta: {
      title: '系统登录',
      public: true,
    },
  },
  {
    name: 'layout',
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: { name: 'overview' },
    children: menus,
  },
  {
    name: 'not-found',
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/error.vue'),
    meta: {
      title: '页面不存在',
      public: true,
    },
  },
]

export default routes

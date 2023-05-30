/* eslint-disable @typescript-eslint/promise-function-async */
import { RouteRecordRaw, RouterView } from 'vue-router'


export const menus: RouteRecordRaw[] = [

  {
    name: 'home',
    path: 'home',
    component: RouterView,
    redirect: 'home',
    children: [
      {
        name: 'device',
        path: 'device',
        component: () => import('@/views/home/device.vue'),
        meta: {
          title: 'device',
        }
      },
      {
        name: 'dashboard',
        path: 'dashboard',
        component: () => import('@/views/home/dashboard.vue'),
        meta: {
          title: 'dashboard',
        }
      },

      {
        name: 'dashboard2',
        path: 'dashboard2',
        component: () => import('@/views/home/dashboard.vue'),
        meta: {
          title: 'dashboard2',
        }
      }

    ],
    meta: {
      title: 'Home'
    }
  }
]


//静态路由
const routes: RouteRecordRaw[] = [
  {
    name: 'login',
    path: '/login',
    component: () => import('@/views/login.vue'),
    meta: {
      title: '系统登陆'
    }
  },
  {
    name: 'layout',
    path: '/',
    component: () => import('@/layout/index.vue'),
    children: menus
  },
  // ## not found page
  {
    name: 'not-found',
    path: '/:path*',
    component: () => import('@/views/error.vue'),
    meta: {
      title: 'Oh no!'
    }
  }
]

export default routes

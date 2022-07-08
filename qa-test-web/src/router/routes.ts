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
        name: 'home',
        path: 'home',
        component: () => import('@/views/home/home.vue'),
        meta: {
          title: 'home',
        }
      },
      {
        name: 'dashboard',
        path: 'dashboard',
        component: () => import('@/views/home/dashboard.vue'),
        meta: {
          title: 'dashboard',
        }
      }

    ],
    meta: {
      title: 'Home'
    }
  }
]

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

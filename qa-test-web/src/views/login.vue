<template>
  <div class="auth-page">
    <div class="auth-bg"></div>
    <a-card class="auth-card" :bordered="false">
      <template #title>
        <div class="card-title">{{ APP_TITLE }}</div>
      </template>

      <a-tabs v-model:activeKey="tab">
        <a-tab-pane key="login" tab="登录">
          <a-form :model="loginForm" layout="vertical" @finish="onLogin">
            <a-form-item
              label="用户名"
              name="username"
              :rules="[{ required: true, message: '请输入用户名' }]"
            >
              <a-input v-model:value="loginForm.username" placeholder="请输入用户名" />
            </a-form-item>
            <a-form-item
              label="密码"
              name="password"
              :rules="[{ required: true, message: '请输入密码' }]"
            >
              <a-input-password v-model:value="loginForm.password" placeholder="请输入密码" />
            </a-form-item>

            <a-button type="primary" html-type="submit" :loading="submitting" block>
              登录
            </a-button>
          </a-form>
        </a-tab-pane>

        <a-tab-pane key="register" tab="注册">
          <a-form :model="registerForm" layout="vertical" @finish="onRegister">
            <a-form-item
              label="用户名"
              name="username"
              :rules="[
                { required: true, message: '请输入用户名' },
                { min: 3, max: 32, message: '用户名长度需在 3-32 之间' },
              ]"
            >
              <a-input v-model:value="registerForm.username" placeholder="支持字母数字 ._-" />
            </a-form-item>
            <a-form-item label="显示名" name="displayName">
              <a-input v-model:value="registerForm.displayName" placeholder="不填则默认使用用户名" />
            </a-form-item>
            <a-form-item
              label="密码"
              name="password"
              :rules="[
                { required: true, message: '请输入密码' },
                { min: 6, message: '密码至少 6 位' },
              ]"
            >
              <a-input-password v-model:value="registerForm.password" placeholder="至少 6 位" />
            </a-form-item>
            <a-form-item
              label="确认密码"
              name="confirmPassword"
              :rules="[
                { required: true, message: '请再次输入密码' },
                { validator: validateConfirmPassword },
              ]"
            >
              <a-input-password
                v-model:value="registerForm.confirmPassword"
                placeholder="再次输入密码"
              />
            </a-form-item>

            <a-button type="primary" html-type="submit" :loading="submitting" block>
              注册并登录
            </a-button>
          </a-form>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import * as API from '@/api'
import { saveSession } from '@/utils/auth'

const APP_TITLE = '\u6fc0\u5149\u5668\u7ba1\u7406\u7cfb\u7edf'

interface LoginForm {
  username: string
  password: string
}

interface RegisterForm {
  username: string
  displayName: string
  password: string
  confirmPassword: string
}

const router = useRouter()
const route = useRoute()

const tab = ref<'login' | 'register'>('login')
const submitting = ref(false)

const loginForm = reactive<LoginForm>({
  username: '',
  password: '',
})

const registerForm = reactive<RegisterForm>({
  username: '',
  displayName: '',
  password: '',
  confirmPassword: '',
})

const redirectPath = computed(() => {
  const redirect = route.query.redirect
  return typeof redirect === 'string' ? redirect : ''
})

const redirectAfterAuth = async () => {
  if (redirectPath.value) {
    await router.replace(redirectPath.value)
    return
  }
  await router.replace({ name: 'overview' })
}

const onLogin = async () => {
  submitting.value = true
  try {
    const payload = await API.auth.login({
      username: loginForm.username.trim(),
      password: loginForm.password,
    })
    saveSession(payload.token, payload.user, payload.expiresAt)
    message.success('登录成功')
    await redirectAfterAuth()
  } catch (error: any) {
    message.error(error?.message || '登录失败')
  } finally {
    submitting.value = false
  }
}

const validateConfirmPassword = async (): Promise<void> => {
  if (registerForm.confirmPassword !== registerForm.password) {
    return Promise.reject(new Error('两次输入的密码不一致'))
  }
  return Promise.resolve()
}

const onRegister = async () => {
  submitting.value = true
  try {
    const payload = await API.auth.register({
      username: registerForm.username.trim(),
      displayName: registerForm.displayName.trim() || undefined,
      password: registerForm.password,
    })
    saveSession(payload.token, payload.user, payload.expiresAt)
    message.success('注册成功，已自动登录')
    await redirectAfterAuth()
  } catch (error: any) {
    message.error(error?.message || '注册失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.auth-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 20% 20%, rgba(22, 119, 255, 0.18), transparent 40%),
    radial-gradient(circle at 80% 80%, rgba(60, 110, 113, 0.2), transparent 42%),
    linear-gradient(135deg, #ecf4ff, #f7fbff 55%, #f0f7ff);
}

.auth-card {
  width: min(460px, 100%);
  border-radius: 12px;
  box-shadow: 0 18px 48px rgba(25, 59, 96, 0.18);
  position: relative;
}

.card-title {
  text-align: center;
  font-size: 20px;
  font-weight: 700;
  color: #254463;
}
</style>

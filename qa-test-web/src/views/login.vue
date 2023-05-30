<template>
    <a-form :form="loginForm" @submit="handleSubmit">
      <a-form-item label="用户名" name="username" :rules="[ { required: true, message: '请输入用户名' } ]">
        <a-input v-model:value="loginForm.username" placeholder="请输入用户名" />
      </a-form-item>
      <a-form-item label="密码" name="password" :rules="[ { required: true, message: '请输入密码' } ]">
        <a-input-password v-model:value="loginForm.password" placeholder="请输入密码" />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit">登录</a-button>
      </a-form-item>
    </a-form>
  </template>
  
  <script lang="ts">
  import { defineComponent } from 'vue';
  import { Form, Input, Button } from 'ant-design-vue';
  
  interface LoginForm {
    username: string;
    password: string;
  }
  
  export default defineComponent({
    name: 'Login',
    components: {
      'a-form': Form,
      'a-form-item': Form.Item,
      'a-input': Input,
      'a-input-password': Input.Password,
      'a-button': Button,
    },
    data() {
      return {
        loginForm: {
          username: '',
          password: '',
        } as LoginForm,
      };
    },
    methods: {
    handleSubmit() {
    //   this.$refs.loginForm.validate((valid) => {
    //     if (valid) {
    //       // 发送登录请求
    //       console.log(this.loginForm);
    //     }
    //   });
    },
  },
});
</script>
  
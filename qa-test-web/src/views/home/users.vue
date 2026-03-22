<template>
  <div class="users-page">
    <div class="toolbar">
      <a-input-search
        v-model:value="keyword"
        allow-clear
        class="search"
        placeholder="按用户名或显示名搜索"
        @search="onSearch"
      />
      <a-space>
        <a-button @click="loadUsers" :loading="state.loading">刷新</a-button>
        <a-button type="primary" @click="openCreate">新增用户</a-button>
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="state.items"
      :loading="state.loading"
      :pagination="false"
      row-key="id"
      size="middle"
      :scroll="{ x: 980 }"
      class="table"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'role'">
          <a-tag :color="roleColor(record.role)">{{ roleLabel(record.role) }}</a-tag>
        </template>

        <template v-else-if="column.dataIndex === 'enabled'">
          <a-switch
            :checked="record.enabled"
            :disabled="isCurrentUser(record)"
            :loading="switchingMap[record.id]"
            @change="(checked: boolean) => onToggleEnabled(record, checked)"
          />
        </template>

        <template v-else-if="column.dataIndex === 'createdAt'">
          {{ formatTime(record.createdAt) }}
        </template>

        <template v-else-if="column.dataIndex === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button type="link" size="small" @click="openResetPassword(record)">重置密码</a-button>
            <a-button
              danger
              type="link"
              size="small"
              :disabled="isCurrentUser(record)"
              @click="onDelete(record)"
            >
              删除
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <div class="pager" v-if="state.total > pageSize">
      <a-pagination
        v-model:current="page"
        :total="state.total"
        :page-size="pageSize"
        size="small"
        @change="onPageChange"
      />
    </div>

    <a-modal
      v-model:visible="createVisible"
      title="新增用户"
      :confirm-loading="modalSubmitting"
      @ok="submitCreate"
      @cancel="closeCreate"
    >
      <a-form layout="vertical" :model="createForm">
        <a-form-item label="用户名" required>
          <a-input v-model:value="createForm.username" placeholder="3-32 位，支持 ._-" />
        </a-form-item>
        <a-form-item label="显示名">
          <a-input v-model:value="createForm.displayName" placeholder="可选" />
        </a-form-item>
        <a-form-item label="角色" required>
          <a-select v-model:value="createForm.role" :options="roleOptions" />
        </a-form-item>
        <a-form-item label="初始密码" required>
          <a-input-password v-model:value="createForm.password" placeholder="至少 6 位" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="editVisible"
      title="编辑用户"
      :confirm-loading="modalSubmitting"
      @ok="submitEdit"
      @cancel="closeEdit"
    >
      <a-form layout="vertical" :model="editForm">
        <a-form-item label="用户名">
          <a-input :value="editForm.username" disabled />
        </a-form-item>
        <a-form-item label="显示名">
          <a-input v-model:value="editForm.displayName" />
        </a-form-item>
        <a-form-item label="角色" required>
          <a-select v-model:value="editForm.role" :options="roleOptions" :disabled="editForm.isSelf" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="resetVisible"
      title="重置密码"
      :confirm-loading="modalSubmitting"
      @ok="submitResetPassword"
      @cancel="closeResetPassword"
    >
      <a-form layout="vertical" :model="resetForm">
        <a-form-item label="目标用户">
          <a-input :value="resetForm.username" disabled />
        </a-form-item>
        <a-form-item label="新密码" required>
          <a-input-password v-model:value="resetForm.password" placeholder="至少 6 位" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import * as API from '@/api'
import { getCurrentUser, type CurrentUser, type UserRole } from '@/utils/auth'

const keyword = ref('')
const page = ref(1)
const pageSize = 12
const modalSubmitting = ref(false)

const state = reactive<{
  items: CurrentUser[]
  total: number
  loading: boolean
}>({
  items: [],
  total: 0,
  loading: false,
})

const switchingMap = reactive<Record<number, boolean>>({})
const currentUser = computed(() => getCurrentUser())

const roleOptions: Array<{ label: string; value: UserRole }> = [
  { label: '只读', value: 'viewer' },
  { label: '运维', value: 'operator' },
  { label: '管理员', value: 'admin' },
]

const columns = [
  { title: '用户名', dataIndex: 'username', width: 140 },
  { title: '显示名', dataIndex: 'displayName', width: 160 },
  { title: '角色', dataIndex: 'role', width: 110 },
  { title: '启用', dataIndex: 'enabled', width: 90 },
  { title: '创建时间', dataIndex: 'createdAt', width: 190 },
  { title: '操作', dataIndex: 'actions', width: 220 },
]

const createVisible = ref(false)
const editVisible = ref(false)
const resetVisible = ref(false)

const createForm = reactive<{
  username: string
  displayName: string
  role: UserRole
  password: string
}>({
  username: '',
  displayName: '',
  role: 'viewer',
  password: '',
})

const editForm = reactive<{
  id: number
  username: string
  displayName: string
  role: UserRole
  isSelf: boolean
}>({
  id: 0,
  username: '',
  displayName: '',
  role: 'viewer',
  isSelf: false,
})

const resetForm = reactive<{
  id: number
  username: string
  password: string
}>({
  id: 0,
  username: '',
  password: '',
})

const formatTime = (value: string): string => {
  if (!value) {
    return '-'
  }
  const time = new Date(value).getTime()
  if (!Number.isFinite(time)) {
    return value
  }
  return new Date(time).toLocaleString()
}

const roleLabel = (role: string): string => {
  if (role === 'admin') {
    return '管理员'
  }
  if (role === 'operator') {
    return '运维'
  }
  return '只读'
}

const roleColor = (role: string): string => {
  if (role === 'admin') {
    return 'red'
  }
  if (role === 'operator') {
    return 'blue'
  }
  return 'default'
}

const isCurrentUser = (record: CurrentUser): boolean => {
  return record.id === currentUser.value?.id
}

const loadUsers = async () => {
  state.loading = true
  try {
    const res = await API.user.list({
      keyword: keyword.value.trim() || undefined,
      offset: (page.value - 1) * pageSize,
      limit: pageSize,
    })
    state.items = res.items || []
    state.total = Number(res.total || 0)
  } catch (error: any) {
    message.error(error?.message || '加载用户列表失败')
  } finally {
    state.loading = false
  }
}

const onSearch = () => {
  page.value = 1
  void loadUsers()
}

const onPageChange = () => {
  void loadUsers()
}

const openCreate = () => {
  createForm.username = ''
  createForm.displayName = ''
  createForm.password = ''
  createForm.role = 'viewer'
  createVisible.value = true
}

const closeCreate = () => {
  createVisible.value = false
}

const submitCreate = async () => {
  const username = createForm.username.trim()
  if (username.length < 3) {
    message.warning('用户名至少 3 位')
    return
  }
  if (createForm.password.length < 6) {
    message.warning('密码至少 6 位')
    return
  }

  modalSubmitting.value = true
  try {
    await API.user.create({
      username,
      displayName: createForm.displayName.trim() || undefined,
      password: createForm.password,
      role: createForm.role,
    })
    message.success('用户创建成功')
    createVisible.value = false
    page.value = 1
    await loadUsers()
  } catch (error: any) {
    message.error(error?.message || '创建用户失败')
  } finally {
    modalSubmitting.value = false
  }
}

const openEdit = (record: CurrentUser) => {
  editForm.id = record.id
  editForm.username = record.username
  editForm.displayName = record.displayName
  editForm.role = record.role
  editForm.isSelf = isCurrentUser(record)
  editVisible.value = true
}

const closeEdit = () => {
  editVisible.value = false
}

const submitEdit = async () => {
  if (!editForm.id) {
    return
  }

  modalSubmitting.value = true
  try {
    await API.user.update(editForm.id, {
      displayName: editForm.displayName.trim(),
      role: editForm.role,
    })
    message.success('用户信息已更新')
    editVisible.value = false
    await loadUsers()
  } catch (error: any) {
    message.error(error?.message || '更新用户失败')
  } finally {
    modalSubmitting.value = false
  }
}

const openResetPassword = (record: CurrentUser) => {
  resetForm.id = record.id
  resetForm.username = record.username
  resetForm.password = ''
  resetVisible.value = true
}

const closeResetPassword = () => {
  resetVisible.value = false
}

const submitResetPassword = async () => {
  if (!resetForm.id) {
    return
  }
  if (resetForm.password.length < 6) {
    message.warning('密码至少 6 位')
    return
  }

  modalSubmitting.value = true
  try {
    await API.user.resetPassword(resetForm.id, resetForm.password)
    message.success('密码已重置')
    resetVisible.value = false
  } catch (error: any) {
    message.error(error?.message || '重置密码失败')
  } finally {
    modalSubmitting.value = false
  }
}

const onDelete = (record: CurrentUser) => {
  Modal.confirm({
    title: `确认删除用户 ${record.username} 吗？`,
    content: '删除后该用户将无法登录系统。',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      try {
        await API.user.remove(record.id)
        message.success('用户已删除')
        if (state.items.length === 1 && page.value > 1) {
          page.value -= 1
        }
        await loadUsers()
      } catch (error: any) {
        message.error(error?.message || '删除用户失败')
        return Promise.reject(error)
      }
    },
  })
}

const onToggleEnabled = async (record: CurrentUser, checked: boolean) => {
  if (isCurrentUser(record) && !checked) {
    message.warning('不能禁用当前登录账号')
    return
  }

  switchingMap[record.id] = true
  try {
    await API.user.update(record.id, { enabled: checked })
    record.enabled = checked
    message.success(`已${checked ? '启用' : '禁用'} ${record.username}`)
  } catch (error: any) {
    message.error(error?.message || '更新启用状态失败')
  } finally {
    switchingMap[record.id] = false
  }
}

onMounted(() => {
  void loadUsers()
})
</script>

<style scoped>
.users-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 14px;
  background: #f8fbff;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.search {
  max-width: 360px;
}

.table {
  flex: 1;
}

.pager {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

@media (min-width: 1600px) {
  .users-page {
    padding: 18px;
  }

  .toolbar {
    margin-bottom: 14px;
    gap: 14px;
  }

  .search {
    max-width: 460px;
  }

  .pager {
    margin-top: 16px;
  }
}

@media (min-width: 2200px) {
  .users-page {
    padding: 22px;
  }

  .search {
    max-width: 560px;
  }
}

@media (max-width: 900px) {
  .users-page {
    padding: 10px;
  }

  .toolbar {
    flex-wrap: wrap;
  }

  .search {
    max-width: 100%;
    width: 100%;
  }

  .pager {
    justify-content: center;
  }
}

@media (max-width: 680px) {
  .toolbar :deep(.ant-space) {
    width: 100%;
  }

  .toolbar :deep(.ant-space .ant-space-item) {
    width: 50%;
  }

  .toolbar :deep(.ant-space .ant-btn) {
    width: 100%;
  }
}
</style>

<script setup lang="ts">
import { ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import {
  createSecret,
  deleteSecret,
  getSecretStorePath,
  listSecrets,
  revealSecret,
  revealSecretExtra,
  updateSecret,
  type SecretInput,
  type SecretItem,
} from '../../api/secretManager'

interface SecretForm {
  id: string
  title: string
  type: string
  account: string
  secret: string
  secretExtra: string
  note: string
}

const items = ref<SecretItem[]>([])
const keyword = ref('')
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const statusMessage = ref('')
const storePath = ref('')
const visibleSecrets = ref<Record<string, string>>({})
const formDialogVisible = ref(false)
let statusTimer: number | undefined
const form = reactive<SecretForm>({
  id: '',
  title: 'GitHub',
  type: 'password',
  account: '',
  secret: '',
  secretExtra: '',
  note: '',
})
const secretTypeOptions = [
  { label: '用户名密码', value: 'password' },
  { label: 'AK / SK', value: 'ak-sk' },
]

const isEditing = computed(() => Boolean(form.id))
const filteredItems = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLocaleLowerCase()

  return items.value.filter((item) => {
    return normalizedKeyword
      ? `${item.title} ${item.account} ${item.note}`
        .toLocaleLowerCase()
        .includes(normalizedKeyword)
      : true
  })
})
const isAkSk = computed(() => form.type === 'ak-sk')

onMounted(async () => {
  await refreshItems()
  storePath.value = await getSecretStorePath()
})

async function refreshItems() {
  loading.value = true
  errorMessage.value = ''

  try {
    items.value = await listSecrets()
  } catch (error) {
    errorMessage.value = readError(error, '读取账号密钥失败')
  } finally {
    loading.value = false
  }
}

async function saveItem() {
  saving.value = true
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    const input = buildInput()
    if (form.id) {
      await updateSecret(form.id, input)
      showStatus('账号密钥已更新')
    } else {
      await createSecret(input)
      showStatus('账号密钥已添加')
    }
    resetForm()
    formDialogVisible.value = false
    visibleSecrets.value = {}
    await refreshItems()
  } catch (error) {
    errorMessage.value = readError(error, '保存账号密钥失败')
  } finally {
    saving.value = false
  }
}

async function removeItem(item: SecretItem) {
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    await ElMessageBox.confirm(`确定要删除「${item.title}」吗？`, '删除账号密钥', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteSecret(item.id)
    if (form.id === item.id) {
      resetForm()
    }
    const { [item.id]: _removed, ...nextVisibleSecrets } = visibleSecrets.value
    visibleSecrets.value = nextVisibleSecrets
    await refreshItems()
    showStatus('账号密钥已删除')
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    errorMessage.value = readError(error, '删除账号密钥失败')
  }
}

async function copySecret(item: SecretItem) {
  await copySecretValue(item, 'secret')
}

async function copySecretExtra(item: SecretItem) {
  await copySecretValue(item, 'extra')
}

async function copySecretValue(item: SecretItem, field: 'secret' | 'extra') {
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    const secret = field === 'secret' ? await revealSecret(item.id) : await revealSecretExtra(item.id)
    await navigator.clipboard.writeText(secret)
    showStatus(field === 'secret' ? '已复制密钥' : '已复制 Secret Key')
  } catch (error) {
    errorMessage.value = readError(error, '复制密钥失败')
  }
}

function showStatus(message: string) {
  statusMessage.value = message
  if (statusTimer) {
    window.clearTimeout(statusTimer)
  }
  statusTimer = window.setTimeout(() => {
    statusMessage.value = ''
  }, 1800)
}

async function toggleSecret(item: SecretItem) {
  await toggleSecretValue(item, 'secret')
}

async function toggleSecretExtra(item: SecretItem) {
  await toggleSecretValue(item, 'extra')
}

async function toggleSecretValue(item: SecretItem, field: 'secret' | 'extra') {
  const key = `${item.id}:${field}`
  if (visibleSecrets.value[key]) {
    const { [key]: _removed, ...nextVisibleSecrets } = visibleSecrets.value
    visibleSecrets.value = nextVisibleSecrets
    return
  }

  try {
    const secret = field === 'secret' ? await revealSecret(item.id) : await revealSecretExtra(item.id)
    visibleSecrets.value = { ...visibleSecrets.value, [key]: secret }
  } catch (error) {
    errorMessage.value = readError(error, '读取密钥失败')
  }
}

function editItem(item: SecretItem) {
  form.id = item.id
  form.title = item.title
  form.type = item.type
  form.account = item.account
  form.secret = ''
  form.secretExtra = ''
  form.note = item.note
  showStatus('编辑时需要重新输入密钥')
  errorMessage.value = ''
  formDialogVisible.value = true
}

function resetForm() {
  form.id = ''
  form.title = ''
  form.type = 'password'
  form.account = ''
  form.secret = ''
  form.secretExtra = ''
  form.note = ''
}

function openCreateDialog() {
  resetForm()
  formDialogVisible.value = true
}

function buildInput(): SecretInput {
  return {
    title: form.title,
    type: form.type,
    account: form.account,
    secret: form.secret,
    secretExtra: form.secretExtra,
    url: '',
    tags: [],
    note: form.note,
  }
}

function secretTypeLabel(value: string) {
  return secretTypeOptions.find((option) => option.value === value)?.label ?? '密码'
}

function formatTime(value: string) {
  if (!value) {
    return '-'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleString()
}

function readError(error: unknown, fallback: string) {
  if (error instanceof Error && error.message) {
    return error.message
  }
  return fallback
}
</script>

<template>
  <section class="tool-page secret-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Security</p>
        <h1>账号密钥管理</h1>
        <p class="summary">本地保存用户名密码和 AK/SK，敏感字段加密后写入数据文件。</p>
      </div>
      <div class="tool-actions">
        <el-button type="primary" @click="openCreateDialog">新建密钥</el-button>
        <el-button @click="refreshItems">刷新</el-button>
      </div>
    </header>

    <section class="secret-layout">
      <section class="secret-list-panel">
        <div class="secret-toolbar">
          <el-input v-model="keyword" class="secret-search" clearable placeholder="搜索名称、账号或备注" />
        </div>

        <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
        <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />

        <div v-if="filteredItems.length" class="secret-list">
          <article v-for="item in filteredItems" :key="item.id" class="secret-item">
            <div class="secret-item__main">
              <h2>{{ item.title }}</h2>
              <el-tag size="small" effect="plain">{{ secretTypeLabel(item.type) }}</el-tag>
              <p v-if="item.type !== 'ak-sk' && item.account" class="secret-account">{{ item.account }}</p>
              <p class="secret-value">
                <span>{{ item.type === 'ak-sk' ? 'AK' : '密钥' }}</span>
                {{ visibleSecrets[`${item.id}:secret`] || '••••••••••••' }}
              </p>
              <p v-if="item.type === 'ak-sk'" class="secret-value">
                <span>SK</span>
                {{ visibleSecrets[`${item.id}:extra`] || '••••••••••••' }}
              </p>
              <p v-if="item.note" class="secret-note">{{ item.note }}</p>
              <div class="secret-item__meta">
                <span>更新于 {{ formatTime(item.updatedAt) }}</span>
              </div>
            </div>

            <div class="secret-item__actions">
              <el-button
                class="action-icon"
                :class="visibleSecrets[`${item.id}:secret`] ? 'action-icon--hide' : 'action-icon--view'"
                :title="visibleSecrets[`${item.id}:secret`] ? '隐藏' : '查看'"
                :aria-label="visibleSecrets[`${item.id}:secret`] ? '隐藏' : '查看'"
                @click.stop="toggleSecret(item)"
              >
                {{ visibleSecrets[`${item.id}:secret`] ? '●' : '○' }}
              </el-button>
              <el-button
                class="action-icon action-icon--copy"
                type="primary"
                title="复制"
                aria-label="复制"
                @click.stop="copySecret(item)"
              >
                ⧉
              </el-button>
              <el-button
                v-if="item.type === 'ak-sk'"
                class="action-icon"
                :class="visibleSecrets[`${item.id}:extra`] ? 'action-icon--hide' : 'action-icon--view'"
                :title="visibleSecrets[`${item.id}:extra`] ? '隐藏 SK' : '查看 SK'"
                :aria-label="visibleSecrets[`${item.id}:extra`] ? '隐藏 SK' : '查看 SK'"
                @click.stop="toggleSecretExtra(item)"
              >
                {{ visibleSecrets[`${item.id}:extra`] ? '●' : '○' }}
              </el-button>
              <el-button
                v-if="item.type === 'ak-sk'"
                class="action-icon action-icon--copy"
                title="复制 SK"
                aria-label="复制 SK"
                @click.stop="copySecretExtra(item)"
              >
                ⧉
              </el-button>
              <el-button
                class="action-icon action-icon--edit"
                title="编辑"
                aria-label="编辑"
                @click.stop="editItem(item)"
              >
                ✎
              </el-button>
              <el-button
                class="action-icon action-icon--delete"
                type="danger"
                plain
                title="删除"
                aria-label="删除"
                @click.stop="removeItem(item)"
              >
                ×
              </el-button>
            </div>
          </article>
        </div>

        <el-empty v-else :description="loading ? '读取中' : '暂无账号密钥'" />

        <p v-if="storePath" class="secret-store-path">数据文件：{{ storePath }}</p>
      </section>
    </section>

    <el-dialog v-model="formDialogVisible" :title="isEditing ? '编辑账号密钥' : '添加账号密钥'" width="520px">
      <form class="secret-form" @submit.prevent="saveItem">
        <label class="secret-field">
          <span>名称</span>
          <el-input v-model="form.title" clearable placeholder="GitHub" />
        </label>

        <label class="secret-field">
          <span>类型</span>
          <el-select v-model="form.type">
            <el-option
              v-for="option in secretTypeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </label>

        <label v-if="!isAkSk" class="secret-field">
          <span>账号</span>
          <el-input v-model="form.account" clearable placeholder="用户名、邮箱或账号标识" />
        </label>

        <label class="secret-field">
          <span>{{ isAkSk ? 'Access Key' : '密码' }}</span>
          <el-input
            v-model="form.secret"
            clearable
            show-password
            :placeholder="isAkSk ? 'Access Key，不会明文保存' : '密码不会明文保存'"
          />
        </label>

        <label v-if="isAkSk" class="secret-field">
          <span>Secret Key</span>
          <el-input v-model="form.secretExtra" clearable show-password placeholder="Secret Key，不会明文保存" />
        </label>

        <label class="secret-field">
          <span>备注</span>
          <el-input v-model="form.note" type="textarea" :rows="4" clearable placeholder="可选" />
        </label>

        <div class="secret-form__actions">
          <el-button @click="formDialogVisible = false">取消</el-button>
          <el-button native-type="submit" type="primary" :loading="saving">
            {{ isEditing ? '保存修改' : '添加密钥' }}
          </el-button>
        </div>
      </form>
    </el-dialog>
  </section>
</template>

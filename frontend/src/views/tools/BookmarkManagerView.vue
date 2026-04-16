<script setup lang="ts">
import { ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import {
  createBookmark,
  deleteBookmark,
  getBookmarkStorePath,
  listBookmarks,
  openBookmark,
  updateBookmark,
  type Bookmark,
  type BookmarkInput,
} from '../../api/bookmarks'

interface BookmarkForm {
  id: string
  title: string
  url: string
  tagsText: string
  note: string
}

const bookmarks = ref<Bookmark[]>([])
const keyword = ref('')
const activeTag = ref('')
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const statusMessage = ref('')
const storePath = ref('')
const formDialogVisible = ref(false)
let statusTimer: number | undefined
const form = reactive<BookmarkForm>({
  id: '',
  title: 'GitHub',
  url: 'github.com',
  tagsText: '开发,代码',
  note: '',
})

const isEditing = computed(() => Boolean(form.id))
const allTags = computed(() => {
  const tags = new Set<string>()
  bookmarks.value.forEach((bookmark) => bookmark.tags.forEach((tag) => tags.add(tag)))
  return Array.from(tags).sort((first, second) => first.localeCompare(second, 'zh-Hans-CN'))
})
const filteredBookmarks = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLocaleLowerCase()

  return bookmarks.value.filter((bookmark) => {
    const matchesKeyword = normalizedKeyword
      ? `${bookmark.title} ${bookmark.url} ${bookmark.tags.join(' ')} ${bookmark.note}`
        .toLocaleLowerCase()
        .includes(normalizedKeyword)
      : true
    const matchesTag = activeTag.value ? bookmark.tags.includes(activeTag.value) : true

    return matchesKeyword && matchesTag
  })
})

onMounted(async () => {
  await refreshBookmarks()
  storePath.value = await getBookmarkStorePath()
})

async function refreshBookmarks() {
  loading.value = true
  errorMessage.value = ''

  try {
    bookmarks.value = await listBookmarks()
  } catch (error) {
    errorMessage.value = readError(error, '读取书签失败')
  } finally {
    loading.value = false
  }
}

async function saveBookmark() {
  saving.value = true
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    const input = buildInput()
    if (form.id) {
      await updateBookmark(form.id, input)
      showStatus('书签已更新')
    } else {
      await createBookmark(input)
      showStatus('书签已添加')
    }
    resetForm()
    formDialogVisible.value = false
    await refreshBookmarks()
  } catch (error) {
    errorMessage.value = readError(error, '保存书签失败')
  } finally {
    saving.value = false
  }
}

async function removeBookmark(bookmark: Bookmark) {
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    await ElMessageBox.confirm(`确定要删除「${bookmark.title}」吗？`, '删除书签', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteBookmark(bookmark.id)
    if (form.id === bookmark.id) {
      resetForm()
    }
    await refreshBookmarks()
    showStatus('书签已删除')
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return
    }
    errorMessage.value = readError(error, '删除书签失败')
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

async function jumpToBrowser(bookmark: Bookmark) {
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    await openBookmark(bookmark.id)
  } catch (error) {
    errorMessage.value = readError(error, '打开浏览器失败')
  }
}

function editBookmark(bookmark: Bookmark) {
  form.id = bookmark.id
  form.title = bookmark.title
  form.url = bookmark.url
  form.tagsText = bookmark.tags.join(', ')
  form.note = bookmark.note
  formDialogVisible.value = true
}

function resetForm() {
  form.id = ''
  form.title = ''
  form.url = 'github.com'
  form.tagsText = ''
  form.note = ''
}

function openCreateDialog() {
  resetForm()
  formDialogVisible.value = true
}

function buildInput(): BookmarkInput {
  return {
    title: form.title,
    url: form.url,
    tags: form.tagsText.split(/[,，;；]/).map((tag) => tag.trim()).filter(Boolean),
    note: form.note,
  }
}

function selectTag(tag: string) {
  activeTag.value = activeTag.value === tag ? '' : tag
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
  <section class="tool-page bookmark-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Browser</p>
        <h1>浏览器书签</h1>
        <p class="summary">保存常用网站，按关键词和标签查找，并跳转到系统默认浏览器。</p>
      </div>
      <div class="tool-actions">
        <el-button type="primary" @click="openCreateDialog">新建书签</el-button>
        <el-button @click="refreshBookmarks">刷新</el-button>
      </div>
    </header>

    <section class="bookmark-layout">
      <section class="bookmark-list-panel">
        <div class="bookmark-toolbar">
          <el-input v-model="keyword" class="bookmark-search" clearable placeholder="搜索标题、网址、标签或备注" />
          <div v-if="allTags.length" class="bookmark-tags">
            <el-tag
              v-for="tag in allTags"
              :key="tag"
              :effect="activeTag === tag ? 'dark' : 'plain'"
              class="bookmark-tag"
              @click="selectTag(tag)"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>

        <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
        <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />

        <div v-if="filteredBookmarks.length" class="bookmark-list">
          <article v-for="bookmark in filteredBookmarks" :key="bookmark.id" class="bookmark-item">
            <div class="bookmark-item__main">
              <h2>{{ bookmark.title }}</h2>
              <p class="bookmark-url">{{ bookmark.url }}</p>
              <p v-if="bookmark.note" class="bookmark-note">{{ bookmark.note }}</p>
              <div class="bookmark-item__meta">
                <el-tag v-for="tag in bookmark.tags" :key="tag" size="small" effect="plain">{{ tag }}</el-tag>
                <span>更新于 {{ formatTime(bookmark.updatedAt) }}</span>
              </div>
            </div>

            <div class="bookmark-item__actions">
              <el-button
                class="action-icon action-icon--open"
                type="primary"
                title="打开"
                aria-label="打开"
                @click.stop="jumpToBrowser(bookmark)"
              >
                ↗
              </el-button>
              <el-button
                class="action-icon action-icon--edit"
                title="编辑"
                aria-label="编辑"
                @click.stop="editBookmark(bookmark)"
              >
                ✎
              </el-button>
              <el-button
                class="action-icon action-icon--delete"
                type="danger"
                plain
                title="删除"
                aria-label="删除"
                @click.stop="removeBookmark(bookmark)"
              >
                ×
              </el-button>
            </div>
          </article>
        </div>

        <el-empty v-else :description="loading ? '读取中' : '暂无书签'" />

        <p v-if="storePath" class="bookmark-store-path">数据文件：{{ storePath }}</p>
      </section>
    </section>

    <el-dialog v-model="formDialogVisible" :title="isEditing ? '编辑书签' : '添加书签'" width="520px">
      <form class="bookmark-form" @submit.prevent="saveBookmark">
        <label class="bookmark-field">
          <span>标题</span>
          <el-input v-model="form.title" clearable placeholder="GitHub" />
        </label>

        <label class="bookmark-field">
          <span>网址</span>
          <el-input v-model="form.url" clearable placeholder="github.com" />
        </label>

        <label class="bookmark-field">
          <span>标签</span>
          <el-input v-model="form.tagsText" clearable placeholder="开发,代码" />
        </label>

        <label class="bookmark-field">
          <span>备注</span>
          <el-input v-model="form.note" type="textarea" :rows="4" clearable placeholder="可选" />
        </label>

        <div class="bookmark-form__actions">
          <el-button @click="formDialogVisible = false">取消</el-button>
          <el-button native-type="submit" type="primary" :loading="saving">
            {{ isEditing ? '保存修改' : '添加书签' }}
          </el-button>
        </div>
      </form>
    </el-dialog>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const now = new Date()
const timestampInput = ref(String(now.getTime()))
const datetimeInput = ref(toLocalDatetimeValue(now))
const statusMessage = ref('')
const errorMessage = ref('')

const timestampResult = computed(() => {
  const value = timestampInput.value.trim()
  if (!value) {
    return null
  }

  const numericValue = Number(value)
  if (!Number.isFinite(numericValue)) {
    return null
  }

  const milliseconds = value.length <= 10 ? numericValue * 1000 : numericValue
  const date = new Date(milliseconds)

  if (Number.isNaN(date.getTime())) {
    return null
  }

  return {
    seconds: Math.floor(date.getTime() / 1000),
    milliseconds: date.getTime(),
    local: date.toLocaleString(),
    iso: date.toISOString(),
  }
})

const datetimeResult = computed(() => {
  if (!datetimeInput.value) {
    return null
  }

  const date = new Date(datetimeInput.value)
  if (Number.isNaN(date.getTime())) {
    return null
  }

  return {
    seconds: Math.floor(date.getTime() / 1000),
    milliseconds: date.getTime(),
    iso: date.toISOString(),
  }
})

function toLocalDatetimeValue(date: Date) {
  const offset = date.getTimezoneOffset()
  const local = new Date(date.getTime() - offset * 60 * 1000)
  return local.toISOString().slice(0, 19)
}

function useCurrentTime() {
  const current = new Date()
  timestampInput.value = String(current.getTime())
  datetimeInput.value = toLocalDatetimeValue(current)
  statusMessage.value = '已使用当前时间'
  errorMessage.value = ''
}

function syncDatetimeFromTimestamp() {
  if (!timestampResult.value) {
    errorMessage.value = '请输入有效时间戳'
    return
  }

  datetimeInput.value = toLocalDatetimeValue(new Date(timestampResult.value.milliseconds))
  statusMessage.value = '已同步到日期时间'
  errorMessage.value = ''
}

function syncTimestampFromDatetime() {
  if (!datetimeResult.value) {
    errorMessage.value = '请输入有效日期时间'
    return
  }

  timestampInput.value = String(datetimeResult.value.milliseconds)
  statusMessage.value = '已同步到时间戳'
  errorMessage.value = ''
}

async function copyText(value: string | number | undefined, label: string) {
  if (value === undefined || value === '') {
    statusMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(String(value))
    statusMessage.value = `已复制${label}`
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Developer</p>
        <h1>时间戳转换</h1>
        <p class="summary">秒级/毫秒级时间戳和本地日期时间互相转换。</p>
      </div>
      <div class="tool-actions">
        <el-button type="primary" @click="useCurrentTime">当前时间</el-button>
      </div>
    </header>

    <section class="timestamp-layout">
      <el-card class="timestamp-card" shadow="never">
        <template #header>时间戳转日期</template>
        <label class="timestamp-field">
          <span>时间戳</span>
          <el-input
            v-model="timestampInput"
            clearable
            placeholder="支持秒或毫秒"
            @clear="statusMessage = ''; errorMessage = ''"
          />
        </label>

        <div class="timestamp-actions">
          <el-button @click="syncDatetimeFromTimestamp">同步到日期时间</el-button>
          <el-button @click="copyText(timestampResult?.milliseconds, '毫秒时间戳')">复制毫秒</el-button>
          <el-button @click="copyText(timestampResult?.seconds, '秒时间戳')">复制秒</el-button>
        </div>

        <dl class="timestamp-result">
          <div>
            <dt>本地时间</dt>
            <dd>{{ timestampResult?.local ?? '-' }}</dd>
          </div>
          <div>
            <dt>ISO 时间</dt>
            <dd>{{ timestampResult?.iso ?? '-' }}</dd>
          </div>
          <div>
            <dt>秒</dt>
            <dd>{{ timestampResult?.seconds ?? '-' }}</dd>
          </div>
          <div>
            <dt>毫秒</dt>
            <dd>{{ timestampResult?.milliseconds ?? '-' }}</dd>
          </div>
        </dl>
      </el-card>

      <el-card class="timestamp-card" shadow="never">
        <template #header>日期转时间戳</template>
        <label class="timestamp-field">
          <span>本地日期时间</span>
          <input v-model="datetimeInput" class="datetime-input" type="datetime-local" step="1" />
        </label>

        <div class="timestamp-actions">
          <el-button @click="syncTimestampFromDatetime">同步到时间戳</el-button>
          <el-button @click="copyText(datetimeResult?.milliseconds, '毫秒时间戳')">复制毫秒</el-button>
          <el-button @click="copyText(datetimeResult?.seconds, '秒时间戳')">复制秒</el-button>
        </div>

        <dl class="timestamp-result">
          <div>
            <dt>秒</dt>
            <dd>{{ datetimeResult?.seconds ?? '-' }}</dd>
          </div>
          <div>
            <dt>毫秒</dt>
            <dd>{{ datetimeResult?.milliseconds ?? '-' }}</dd>
          </div>
          <div>
            <dt>ISO 时间</dt>
            <dd>{{ datetimeResult?.iso ?? '-' }}</dd>
          </div>
        </dl>
      </el-card>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type Preset = {
  name: string
  description: string
  expression: string
}

const presets: Preset[] = [
  { name: '每分钟', description: '每分钟执行', expression: '* * * * *' },
  { name: '每小时', description: '每小时第0分钟执行', expression: '0 * * * *' },
  { name: '每天', description: '每天0点执行', expression: '0 0 * * *' },
  { name: '每周', description: '每周一0点执行', expression: '0 0 * * 1' },
  { name: '每月', description: '每月1号0点执行', expression: '0 0 1 * *' },
  { name: '工作日', description: '工作日每小时执行', expression: '0 * * * 1-5' },
  { name: '周末', description: '周末每小时执行', expression: '0 * * * 0,6' },
  { name: '早晚高峰', description: '早8晚6执行', expression: '0 8,18 * * *' },
  { name: '每季度', description: '每季度第1天执行', expression: '0 0 1 1,4,7,10 *' },
  { name: '每半年', description: '每年1月和7月执行', expression: '0 0 1 1,7 *' },
]

const minute = ref('*')
const hour = ref('*')
const day = ref('*')
const month = ref('*')
const weekday = ref('*')

const statusMessage = ref('')
const errorMessage = ref('')

const expression = computed(() => 
  `${minute.value} ${hour.value} ${day.value} ${month.value} ${weekday.value}`
)

function usePreset(preset: Preset) {
  const parts = preset.expression.split(' ')
  minute.value = parts[0]
  hour.value = parts[1]
  day.value = parts[2]
  month.value = parts[3]
  weekday.value = parts[4]
  statusMessage.value = '已应用预设'
  errorMessage.value = ''
}

function clearAll() {
  minute.value = '*'
  hour.value = '*'
  day.value = '*'
  month.value = '*'
  weekday.value = '*'
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyExpression() {
  try {
    await navigator.clipboard.writeText(expression.value)
    statusMessage.value = '已复制表达式'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

function validate(): boolean {
  const fields = [minute.value, hour.value, day.value, month.value, weekday.value]
  const ranges = [[0, 59], [0, 23], [1, 31], [1, 12], [0, 6]]
  
  for (let i = 0; i < fields.length; i++) {
    const field = fields[i]
    const [min, max] = ranges[i]
    
    if (!isValidField(field, min, max)) {
      return false
    }
  }
  
  return true
}

function isValidField(field: string, min: number, max: number): boolean {
  if (field === '*') return true
  
  const values = field.split(',')
  
  for (const value of values) {
    if (value.includes('/')) {
      const [step, period] = value.split('/')
      if (!isValidRangeOrNumber(step, min, max) || !isValidNumber(period, 1, max)) {
        return false
      }
    } else if (value.includes('-')) {
      const [start, end] = value.split('-')
      if (!isValidNumber(start, min, max) || !isValidNumber(end, min, max)) {
        return false
      }
      const s = parseInt(start)
      const e = parseInt(end)
      if (s >= e) return false
    } else {
      if (!isValidNumber(value, min, max)) return false
    }
  }
  
  return true
}

function isValidRangeOrNumber(value: string, min: number, max: number): boolean {
  if (value === '*') return true
  return isValidNumber(value, min, max)
}

function isValidNumber(value: string, min: number, max: number): boolean {
  const num = parseInt(value)
  if (isNaN(num)) return false
  if (num < min || num > max) return false
  return true
}

const isValid = computed(() => validate())

const nextRuns = computed(() => {
  if (!isValid.value) return []
  try {
    const runs: Date[] = []
    let date = new Date()
    
    for (let i = 0; i < 5; i++) {
      const next = calculateNextRun(date)
      if (!next) break
      runs.push(next)
      date = new Date(next.getTime() + 60000)
    }
    
    return runs
  } catch {
    return []
  }
})

function calculateNextRun(start: Date): Date | null {
  let date = new Date(start.getTime())
  
  for (let i = 0; i < 100000; i++) {
    date = new Date(date.getTime() + 60000)
    
    if (matches(date)) {
      return date
    }
  }
  
  return null
}

function matches(date: Date): boolean {
  const m = date.getMinutes()
  const h = date.getHours()
  const d = date.getDate()
  const mo = date.getMonth() + 1
  const w = date.getDay()
  
  return (
    matchesField(minute.value, m) &&
    matchesField(hour.value, h) &&
    matchesField(day.value, d) &&
    matchesField(month.value, mo) &&
    matchesField(weekday.value, w)
  )
}

function matchesField(field: string, value: number): boolean {
  if (field === '*') return true
  
  const values = field.split(',')
  
  for (const v of values) {
    if (v.includes('/')) {
      const [step, period] = v.split('/')
      const range = step === '*' ? value : parseInt(step)
      
      if (step === '*') {
        if (value % parseInt(period) === 0) return true
      } else {
        const start = parseInt(step)
        const p = parseInt(period)
        if (value >= start && (value - start) % p === 0) return true
      }
    } else if (v.includes('-')) {
      const [start, end] = v.split('-')
      const s = parseInt(start)
      const e = parseInt(end)
      if (value >= s && value <= e) return true
    } else {
      if (parseInt(v) === value) return true
    }
  }
  
  return false
}

function formatDate(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  const year = date.getFullYear()
  const month = pad(date.getMonth() + 1)
  const day = pad(date.getDate())
  const hour = pad(date.getHours())
  const minute = pad(date.getMinutes())
  return `${year}-${month}-${day} ${hour}:${minute}`
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Cron</p>
        <h1>Cron 表达式生成器</h1>
        <p class="summary">可视化配置定时任务，生成标准 Cron 表达式</p>
      </div>
    </header>

    <section class="cron-layout">
      <aside class="cron-sidebar">
        <div class="cron-sidebar-section">
          <h3>快速预设</h3>
          <div class="cron-presets">
            <button 
              v-for="preset in presets" 
              :key="preset.expression"
              class="preset-btn"
              @click="usePreset(preset)"
            >
              <strong>{{ preset.name }}</strong>
              <span>{{ preset.description }}</span>
            </button>
          </div>
        </div>
      </aside>

      <section class="cron-content">
        <div class="cron-editor">
          <div class="cron-fields">
            <div class="cron-field">
              <label>分钟</label>
              <el-input v-model="minute" placeholder="0-59" />
              <span class="field-hint">0-59</span>
            </div>
            <div class="cron-field">
              <label>小时</label>
              <el-input v-model="hour" placeholder="0-23" />
              <span class="field-hint">0-23</span>
            </div>
            <div class="cron-field">
              <label>日期</label>
              <el-input v-model="day" placeholder="1-31" />
              <span class="field-hint">1-31</span>
            </div>
            <div class="cron-field">
              <label>月份</label>
              <el-input v-model="month" placeholder="1-12" />
              <span class="field-hint">1-12</span>
            </div>
            <div class="cron-field">
              <label>星期</label>
              <el-input v-model="weekday" placeholder="0-6" />
              <span class="field-hint">0=周日, 6=周六</span>
            </div>
          </div>

          <div class="cron-expression">
            <div class="expression-label">生成表达式</div>
            <div class="expression-value">
              <code>{{ expression }}</code>
              <el-button 
                type="primary" 
                @click="copyExpression"
                :disabled="!isValid"
              >
                复制
              </el-button>
            </div>
            <el-alert 
              v-if="!isValid"
              type="warning" 
              title="表达式可能无效，请检查各字段格式"
              :closable="false"
              show-icon
            />
          </div>

          <div class="cron-next-runs">
            <h3>下次执行</h3>
            <div class="next-runs-list">
              <div v-for="(run, index) in nextRuns" :key="index" class="run-item">
                <span class="run-index">#{{ index + 1 }}</span>
                <span class="run-time">{{ formatDate(run) }}</span>
              </div>
              <div v-if="nextRuns.length === 0" class="empty-runs">
                <span>无法计算下次执行时间</span>
              </div>
            </div>
          </div>
        </div>

        <div class="cron-footer">
          <el-button @click="clearAll">重置</el-button>
        </div>
      </section>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
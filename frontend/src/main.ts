import { createApp } from 'vue'
import {
  ElAlert,
  ElAside,
  ElButton,
  ElCard,
  ElCheckbox,
  ElCol,
  ElContainer,
  ElDialog,
  ElEmpty,
  ElInput,
  ElInputNumber,
  ElMain,
  ElMenu,
  ElMenuItem,
  ElOption,
  ElRow,
  ElSelect,
  ElStatistic,
  ElTag,
} from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import './styles/main.css'

const app = createApp(App)

app.use(ElAlert)
app.use(ElAside)
app.use(ElButton)
app.use(ElCard)
app.use(ElCheckbox)
app.use(ElCol)
app.use(ElContainer)
app.use(ElDialog)
app.use(ElEmpty)
app.use(ElInput)
app.use(ElInputNumber)
app.use(ElMain)
app.use(ElMenu)
app.use(ElMenuItem)
app.use(ElOption)
app.use(ElRow)
app.use(ElSelect)
app.use(ElStatistic)
app.use(ElTag)

app.mount('#app')

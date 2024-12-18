import { createApp } from 'vue'
import * as pdfjsLib from 'pdfjs-dist'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import i18n from './i18n'
// 配置 PDF.js worker
const pdfjsWorkerSrc = new URL(
  'pdfjs-dist/build/pdf.worker.min.js',
  import.meta.url
).toString()
pdfjsLib.GlobalWorkerOptions.workerSrc = pdfjsWorkerSrc

import App from './App.vue'

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(ElementPlus)
app.use(i18n)
app.mount('#app')

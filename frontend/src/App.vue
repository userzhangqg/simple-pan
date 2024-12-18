<template>
  <div class="app-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1>{{ $t('message.welcome') }}</h1>
          <el-button 
            class="lang-button"
            type="primary" 
            plain 
            @click="toggleLanguage"
            :icon="Management">
            {{ currentLanguage === 'zh-CN' ? 'English' : '中文' }}
          </el-button>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <div class="upload-section">
          <el-upload
            class="upload-area"
            drag
            :action="uploadUrl"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            :on-progress="handleUploadStart"
            accept=".jpg,.jpeg,.png,.gif,.pdf,.doc,.docx,.xls,.xlsx,.txt"
            multiple
            :disabled="uploadLoading">
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              {{ $t('message.file.upload.dragText') }}
              <div class="upload-tip">{{ $t('message.file.upload.supportedTypes') }}</div>
            </div>
          </el-upload>
        </div>

        <div class="file-list">
          <!-- 搜索框 -->
          <div class="search-box">
            <el-input
              v-model="searchQuery"
              :placeholder="$t('message.file.search.placeholder')"
              clearable
              @update:model-value="handleSearch">
              <template #prefix>
                <el-icon><search /></el-icon>
              </template>
            </el-input>
          </div>

          <el-table 
            :data="fileList" 
            style="width: 100%"
            :row-class-name="tableRowClassName"
            @row-click="handleRowClick"
            v-loading="tableLoading"
            element-loading-text="Loading..."
            element-loading-background="rgba(255, 255, 255, 0.9)">
            <el-table-column prop="name" :label="$t('message.file.name')" min-width="200" header-align="center">
              <template #default="scope">
                <div class="file-name">
                  <el-icon class="file-icon">
                    <component :is="getFileIcon(scope.row)"></component>
                  </el-icon>
                  <span class="file-name-text">{{ scope.row.name }}</span>
                </div>
                <div v-if="isMobileDevice" class="file-info">
                  <span>{{ scope.row.type }}</span>
                  <span>{{ scope.row.size }}</span>
                  <span>{{ scope.row.modified }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column v-if="!isMobileDevice" prop="type" :label="$t('message.file.type')" width="100" align="center" header-align="center" />
            <el-table-column v-if="!isMobileDevice" prop="size" :label="$t('message.file.size')" width="100" align="center" header-align="center" />
            <el-table-column v-if="!isMobileDevice" prop="modified" :label="$t('message.file.lastModified')" width="180" align="center" header-align="center" />
            <el-table-column :label="$t('message.file.actions')" width="200" fixed="right" header-align="center">
              <template #default="scope">
                <div class="action-buttons">
                  <div class="button-wrapper">
                    <el-tooltip
                      :content="$t('message.download')"
                      placement="top"
                      :disabled="!isMobileDevice">
                      <el-button
                        size="small"
                        type="primary"
                        plain
                        @click.stop="handleDownload(scope.row)"
                        :icon="Download"
                        :loading="downloadLoading.has(scope.row.name)">
                        <span class="button-text">{{ downloadLoading.has(scope.row.name) ? $t('message.downloading') : $t('message.download') }}</span>
                      </el-button>
                    </el-tooltip>
        <el-tooltip
          :content="$t('message.preview')"
          placement="top"
          :disabled="!isMobileDevice || (!isImage(scope.row.name) && !isPdf(scope.row.name))">
          <el-button
            size="small"
            type="primary"
            @click.stop="handlePreview(scope.row)"
            :icon="ZoomIn"
            :disabled="!isImage(scope.row.name) && !isPdf(scope.row.name)"
            :title="!isImage(scope.row.name) && !isPdf(scope.row.name) ? $t('message.file.preview.notSupported') : ''">
            <span class="button-text">{{ $t('message.preview') }}</span>
          </el-button>
        </el-tooltip>
                    <el-tooltip
                      :content="$t('message.delete')"
                      placement="top"
                      :disabled="!isMobileDevice">
                      <el-button
                        size="small"
                        type="danger"
                        plain
                        @click.stop="handleDelete(scope.row)"
                        :icon="Delete"
                        :loading="deleteLoading.has(scope.row.name)">
                        <span class="button-text">{{ deleteLoading.has(scope.row.name) ? $t('message.deleting') : $t('message.delete') }}</span>
                      </el-button>
                    </el-tooltip>
                  </div>
                </div>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="totalFiles"
              @current-change="handlePageChange"
              @size-change="handleSizeChange"
              :page-sizes="[5, 10, 20, 50]"
              layout="sizes, total, prev, pager, next"
              :total-template="$t('message.file.pagination.total', { total: totalFiles })"
              background />
          </div>
        </div>

        <!-- 预览对话框 -->
        <el-dialog 
          v-model="previewVisible" 
          :title="$t('message.preview')" 
          width="90%"
          class="preview-dialog"
          :close-on-click-modal="true"
          :show-close="true">
          <template v-if="!isPdfPreview">
            <div class="preview-controls">
              <el-button-group>
                <el-button :icon="ZoomOut" @click="handleZoomOut" />
                <el-button>{{ Math.round(zoomLevel * 100) }}%</el-button>
                <el-button :icon="ZoomIn" @click="handleZoomIn" />
                <el-button :icon="RefreshRight" @click="resetZoom" />
              </el-button-group>
            </div>
            <div class="preview-container" ref="previewContainer">
              <el-image
                ref="previewImage"
                class="preview-image"
                :src="previewUrl"
                :preview-src-list="[previewUrl]"
                fit="contain"
                :initial-index="0"
                :style="{ transform: `scale(${zoomLevel})` }"
                @load="handleImageLoad"
              />
            </div>
          </template>
          <template v-else>
            <div class="pdf-preview-container">
              <vue-pdf-embed
                :source="previewUrl"
                :width="800"
                :height="1000"
                style="width: 100%; height: 100%;"
              />
            </div>
          </template>
        </el-dialog>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, shallowRef, nextTick } from 'vue'
import VuePdfEmbed from 'vue-pdf-embed'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  UploadFilled, 
  Download, 
  Delete, 
  ZoomIn,
  ZoomOut,
  RefreshRight,
  Document,
  Picture,
  VideoCamera,
  Files,
  Management,
  Search
} from '@element-plus/icons-vue'
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api'
const uploadUrl = `${API_BASE_URL}/files/upload`
const fileList = ref([])
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(5)  // 默认每页显示5条
const totalFiles = ref(0)
const tableLoading = ref(false)
const uploadLoading = ref(false)
const deleteLoading = ref(new Set())
const downloadLoading = ref(new Set())
const isMobileDevice = ref(window.innerWidth <= 768)
const imageLoading = ref(false)

// 防抖函数
const debounce = (fn, delay) => {
  let timer = null
  return (...args) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

// 搜索处理
const handleSearch = debounce(() => {
  currentPage.value = 1
  fetchFileList()
}, 300)

// 分页处理
const handlePageChange = (page) => {
  currentPage.value = page
  fetchFileList()
}

// 处理每页显示数量变化
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1 // 重置到第一页
  fetchFileList()
}

// 获取文件列表
const fetchFileList = async () => {
  console.log('Fetching file list with:', {
    page: currentPage.value,
    pageSize: pageSize.value,
    search: searchQuery.value
  })
  
  tableLoading.value = true
  try {
    const params = new URLSearchParams()
    params.append('page', currentPage.value.toString())
    params.append('pageSize', pageSize.value.toString())
    
    if (searchQuery.value.trim()) {
      params.append('search', searchQuery.value.trim())
    }
    
    const url = `${API_BASE_URL}/files/list?${params.toString()}`
    console.log('Request URL:', url)
    
    const response = await axios.get(url)
    console.log('Response:', response.data)
    
    if (response.data && typeof response.data === 'object') {
      // 处理文件列表数据
      const files = response.data.files || []
      fileList.value = files.map(file => ({
        name: file.name || '',
        type: file.type || '',
        size: file.size || '',
        modified: file.modified ? new Date(file.modified).toLocaleString() : ''
      }))
      
      // 处理分页数据
      totalFiles.value = parseInt(response.data.total) || 0
      currentPage.value = parseInt(response.data.page) || 1
      
      if (fileList.value.length === 0) {
        if (searchQuery.value.trim()) {
          ElMessage.info(t('message.file.search.noMatch'))
        } else {
          ElMessage.info(t('message.file.list.empty'))
        }
      }
    } else {
      throw new Error('Invalid response format')
    }
  } catch (error) {
    console.error('Error fetching file list:', error)
    ElMessage.error(t('message.file.list.error'))
    fileList.value = []
    totalFiles.value = 0
  } finally {
    tableLoading.value = false
  }
}

const { t, locale } = useI18n()
const currentLanguage = ref(locale.value)

// 监听窗口大小变化
window.addEventListener('resize', () => {
  isMobileDevice.value = window.innerWidth <= 768
})

// 获取文件图标
const getFileIcon = (file) => {
  const type = file.type.toLowerCase()
  if (type.includes('image')) return Picture
  if (type.includes('video')) return VideoCamera
  if (type.includes('document')) return Document
  return Files
}

// 表格行样式
const tableRowClassName = ({ row, rowIndex }) => {
  return 'file-row'
}

// 行点击事件
const handleRowClick = (row) => {
  // 如果是图片，则预览
  if (isImage(row.name)) {
    handlePreview(row)
  }
}

// 图片加载处理
const handleImageLoad = () => {
  imageLoading.value = false
}

const toggleLanguage = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  currentLanguage.value = locale.value
}

// 上传相关处理
const handleUploadSuccess = (response) => {
  uploadLoading.value = false
  ElMessage.success(t('message.file.upload.success'))
  fetchFileList()
}

const handleUploadError = () => {
  uploadLoading.value = false
  ElMessage.error(t('message.file.upload.error'))
}

const handleUploadStart = () => {
  uploadLoading.value = true
}

const SUPPORTED_FILE_TYPES = [
  'image/jpeg',
  'image/png',
  'image/gif',
  'application/pdf',
  'application/msword',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  'application/vnd.ms-excel',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'text/plain'
]

const beforeUpload = (file) => {
  // 检查文件大小
  const maxSize = 100 * 1024 * 1024 // 100MB
  if (file.size > maxSize) {
    ElMessage.error(t('message.file.upload.sizeError'))
    return false
  }

  // 检查文件类型
  if (!SUPPORTED_FILE_TYPES.includes(file.type)) {
    ElMessage({
      message: `${t('message.file.upload.typeError')}\n${t('message.file.upload.supportedTypes')}`,
      type: 'error',
      duration: 5000
    })
    return false
  }

  return true
}

// 下载文件
const handleDownload = async (file) => {
  if (downloadLoading.value.has(file.name)) return
  downloadLoading.value.add(file.name)
  
  try {
    const response = await axios({
      url: `${API_BASE_URL}/files/download/${file.name}`,
      method: 'GET',
      responseType: 'blob'
    })
    
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', file.name)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success(t('message.file.download.success'))
  } catch (error) {
    ElMessage.error(t('message.file.download.error'))
    console.error('Error downloading file:', error)
  } finally {
    downloadLoading.value.delete(file.name)
  }
}

// 判断文件类型
const isImage = (filename) => {
  const ext = filename.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)
}

const isPdf = (filename) => {
  const ext = filename.toLowerCase().split('.').pop()
  return ext === 'pdf'
}

// 预览相关
const previewVisible = ref(false)
const previewUrl = ref('')
const zoomLevel = ref(1)
const isPdfPreview = ref(false)

const MIN_ZOOM = 0.1
const MAX_ZOOM = 3
const ZOOM_STEP = 0.1

const handleZoomIn = () => {
  if (zoomLevel.value < MAX_ZOOM) {
    zoomLevel.value = Math.min(zoomLevel.value + ZOOM_STEP, MAX_ZOOM)
  }
}

const handleZoomOut = () => {
  if (zoomLevel.value > MIN_ZOOM) {
    zoomLevel.value = Math.max(zoomLevel.value - ZOOM_STEP, MIN_ZOOM)
  }
}

const resetZoom = () => {
  zoomLevel.value = 1
}


const handlePreview = async (file) => {
  try {
    previewUrl.value = `${API_BASE_URL}/files/preview/${file.name}`
    previewVisible.value = true
    isPdfPreview.value = isPdf(file.name)
    if (!isPdfPreview.value) {
      resetZoom()
    }
  } catch (error) {
    console.error('Error previewing file:', error)
    ElMessage.error(t('message.file.preview.error'))
  }
}

const handleDelete = async (file) => {
  if (deleteLoading.value.has(file.name)) return
  
  try {
    deleteLoading.value.add(file.name)
    await ElMessageBox.confirm(
      t('message.file.delete.confirm'),
      t('message.warning'),
      {
        confirmButtonText: t('message.confirm'),
        cancelButtonText: t('message.cancel'),
        type: 'warning'
      }
    )
    
    await axios.delete(`${API_BASE_URL}/files/delete/${file.name}`)
    ElMessage.success(t('message.file.delete.success'))
    fetchFileList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('message.file.delete.error'))
      console.error('Error deleting file:', error)
    }
  } finally {
    deleteLoading.value.delete(file.name)
  }
}

onMounted(() => {
  fetchFileList()
})
</script>

<style scoped>
.app-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
  background-color: #f8f9fa;
  width: 100%;
  box-sizing: border-box;
}

@media screen and (max-width: 768px) {
  .app-container {
    padding: 10px;
  }

  .header-content h1 {
    font-size: 20px !important;
  }

  .el-upload-dragger {
    height: 150px !important;
  }

  .el-table {
    font-size: 14px;
  }

  .action-buttons .el-button {
    padding: 6px 12px;
    font-size: 12px;
  }

  .action-buttons .el-button + .el-button {
    margin-left: 4px;
  }

  .preview-dialog :deep(.el-dialog) {
    width: 95% !important;
    margin: 10px auto !important;
  }
}

@media screen and (max-width: 480px) {
  .app-container {
    padding: 5px;
  }

  .header-content {
    flex-direction: column;
    gap: 10px;
    padding: 10px !important;
  }

  .header-content h1 {
    font-size: 18px !important;
  }

  .el-upload-dragger {
    height: 120px !important;
  }

  .el-icon--upload {
    font-size: 36px !important;
  }

  .el-upload__text {
    font-size: 14px !important;
  }

  .file-name-text {
    max-width: 120px;
  }

  .action-buttons {
    flex-wrap: wrap;
  }

  .action-buttons .el-button {
    padding: 4px 8px;
    font-size: 12px;
  }

  .preview-dialog :deep(.el-dialog) {
    width: 98% !important;
    margin: 5px auto !important;
  }

  .preview-dialog :deep(.el-dialog__body) {
    padding: 12px;
  }
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.header:hover {
  box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 100%;
}

.header-content h1 {
  margin: 0;
  font-size: 24px;
  color: #409EFF;
  font-weight: 600;
}

.lang-button {
  transition: all 0.3s ease;
}

.main-content {
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.upload-section {
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(145deg, #ffffff, #f8f9fa);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.upload-area {
  width: 100%;
  border-radius: 12px;
  transition: all 0.3s ease;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
}

.upload-area:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(64, 158, 255, 0.1);
}

.el-upload {
  width: 100%;
}

.el-upload-dragger {
  width: 100%;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 2px dashed #e4e7ed;
  transition: all 0.3s ease;
}

.el-upload-dragger:hover {
  border-color: #409EFF;
  background-color: rgba(64, 158, 255, 0.05);
}

.el-icon--upload {
  font-size: 48px;
  color: #409EFF;
  margin-bottom: 16px;
  transition: all 0.3s ease;
}

.el-upload-dragger:hover .el-icon--upload {
  transform: translateY(-5px);
  color: #66b1ff;
}

.el-upload__text {
  color: #606266;
  font-size: 16px;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.el-upload-dragger:hover .el-upload__text {
  color: #409EFF;
}

.el-upload-dragger:hover .upload-tip {
  color: #79bbff;
}

@media screen and (max-width: 480px) {
  .upload-tip {
    font-size: 11px;
    padding: 0 12px;
    text-align: center;
  }
}

.file-list {
  margin-top: 20px;
  border-radius: 8px;
  overflow: hidden;
}

.search-box {
  margin-bottom: 20px;
}

.search-box .el-input {
  max-width: 300px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

@media screen and (max-width: 768px) {
  .search-box {
    margin-bottom: 16px;
  }

  .search-box .el-input {
    max-width: 100%;
  }

  .pagination-container {
    margin-top: 16px;
  }

  .pagination-container :deep(.el-pagination) {
    font-size: 12px;
  }
}

.file-list .el-table {
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.file-list .el-table::before {
  display: none;
}

.file-list .el-table th.el-table__cell {
  background-color: #f8f9fa;
  color: #606266;
  font-weight: 600;
  border-bottom: 2px solid #ebeef5;
  padding: 12px;
  text-align: center !important;
}

.file-list .el-table td.el-table__cell {
  padding: 12px;
}

@media screen and (max-width: 768px) {
  .hide-on-mobile {
    display: none;
  }

  .file-list .el-table td.el-table__cell {
    padding: 8px;
  }

  .file-list :deep(.el-table__header-wrapper) {
    display: none;
  }

  .file-list :deep(.el-table__body) {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
    padding: 16px;
  }

  .file-list :deep(.el-table__row) {
    display: flex;
    flex-direction: column;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    padding: 16px;
    margin: 0;
  }

  .file-list :deep(.el-table__row:hover) {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }

  .file-list :deep(.el-table__cell) {
    display: block;
    padding: 4px 0 !important;
    border: none !important;
  }

  .file-info {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin: 8px 0;
    color: #909399;
    font-size: 12px;
  }

  .file-info span {
    background: #f5f7fa;
    padding: 2px 8px;
    border-radius: 4px;
  }

  .action-buttons {
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid #ebeef5;
  }
}

@media screen and (max-width: 480px) {
  .file-list :deep(.el-table__body) {
    grid-template-columns: 1fr;
    padding: 12px;
    gap: 12px;
  }

  .file-list :deep(.el-table__row) {
    padding: 12px;
  }
}

.file-list .el-table--enable-row-hover .el-table__body tr:hover > td.el-table__cell {
  background-color: #f5f7fa;
}

.file-row {
  transition: all 0.3s ease;
}

.file-row:hover {
  background-color: #f5f7fa;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-name-text {
  font-size: 14px;
  color: #606266;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-icon {
  font-size: 20px;
  color: #909399;
  flex-shrink: 0;
}

.action-buttons {
  display: flex;
  justify-content: center;
}

.button-wrapper {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.action-buttons .el-button {
  padding: 8px 16px;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.action-buttons .el-button:hover {
  transform: translateY(-2px);
}

.action-buttons .el-button[disabled] {
  cursor: not-allowed;
  opacity: 0.6;
}

.action-buttons .el-button[disabled]:hover {
  transform: none;
}

@media screen and (max-width: 768px) {
  .action-buttons .el-button {
    padding: 8px;
  }
  
  .button-text {
    display: none;
  }
  
  .action-buttons .el-button .el-icon {
    margin: 0;
  }
}

@media screen and (max-width: 480px) {
  .button-wrapper {
    gap: 4px;
  }
  
  .action-buttons .el-button {
    padding: 6px;
    min-width: 32px;
  }
}

.preview-dialog :deep(.el-dialog) {
  display: flex;
  flex-direction: column;
  max-height: 90vh;
  margin: 5vh auto !important;
  border-radius: 8px;
  overflow: hidden;
}

.preview-dialog :deep(.el-dialog__header) {
  margin: 0;
  padding: 16px 20px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e4e7ed;
}

.preview-dialog :deep(.el-dialog__body) {
  flex: 1;
  overflow: hidden;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.preview-controls {
  padding: 10px;
  display: flex;
  justify-content: center;
  background: #f8f9fa;
  border-bottom: 1px solid #e4e7ed;
}

.preview-controls .el-button-group {
  display: flex;
  align-items: center;
}

.preview-controls .el-button {
  padding: 8px 15px;
}

.preview-container {
  flex: 1;
  overflow: auto;
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f8f9fa;
}

.preview-image {
  max-width: 100%;
  max-height: calc(90vh - 140px);
  object-fit: contain;
  transition: transform 0.3s ease;
  transform-origin: center center;
}

.pdf-preview-container {
  flex: 1;
  overflow: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: #f8f9fa;
}

/* 添加按钮动画效果 */
.el-button {
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 添加表格动画效果 */
.el-table {
  transition: all 0.3s ease;
}

.el-table:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #909399;
}
</style>

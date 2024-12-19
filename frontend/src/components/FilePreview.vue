<template>
  <el-dialog 
    v-model="visible" 
    :title="$t('message.preview')" 
    width="90%"
    class="preview-dialog"
    :close-on-click-modal="true"
    :show-close="true"
    @closed="handleClose">
    <template v-if="isImage">
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
    <template v-else-if="isPdf">
      <PDFViewer :url="previewUrl" class="pdf-preview-container" />
    </template>
    <template v-else-if="isEpub">
      <div class="epub-preview-container">
        <EPUBViewer :url="previewUrl" />
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import PDFViewer from './PDFViewer.vue'
import EPUBViewer from './EPUBViewer.vue'
import { ZoomIn, ZoomOut, RefreshRight } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  file: {
    type: Object,
    default: () => ({})
  },
  apiBaseUrl: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const zoomLevel = ref(1)
const MIN_ZOOM = 0.1
const MAX_ZOOM = 3
const ZOOM_STEP = 0.1

const previewUrl = computed(() => {
  if (!props.file?.name) return ''
  return `${props.apiBaseUrl}/files/preview/${props.file.name}`
})

const isImage = computed(() => {
  const ext = props.file?.name?.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)
})

const isPdf = computed(() => {
  const ext = props.file?.name?.toLowerCase().split('.').pop()
  return ext === 'pdf'
})

const isEpub = computed(() => {
  const ext = props.file?.name?.toLowerCase().split('.').pop()
  return ext === 'epub'
})

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

const handleImageLoad = () => {
  // 图片加载完成后的处理
}

const handleClose = () => {
  resetZoom()
}
</script>

<style scoped>
.preview-dialog :deep(.el-dialog) {
  display: flex;
  flex-direction: column;
  max-height: 95vh;
  margin: 2.5vh auto !important;
  border-radius: 12px;
  overflow: hidden;
  background: #f8f9fa;
}

.preview-dialog :deep(.el-dialog__header) {
  margin: 0;
  padding: 12px 20px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.preview-dialog :deep(.el-dialog__headerbtn) {
  top: 14px;
}

.preview-dialog :deep(.el-dialog__body) {
  flex: 1;
  padding: 0;
  overflow: hidden;
  min-height: 0;
  background: #f8f9fa;
}

.preview-dialog :deep(.el-dialog__title) {
  font-weight: 500;
  color: #303133;
}

.preview-controls {
  padding: 10px;
  display: flex;
  justify-content: center;
  background: #fff;
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
  overflow: hidden;
  height: calc(95vh - 80px);
  display: flex;
  flex-direction: column;
  background: #f8f9fa;
}

.preview-dialog :deep(.el-dialog__body > .pdf-preview-container) {
  min-height: 0;
  height: calc(95vh - 80px);
}

.epub-preview-container {
  flex: 1;
  overflow: hidden;
  padding: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  background: #fff;
  height: calc(95vh - 100px);
}

.epub-preview-container :deep(.epub-viewer) {
  flex: 1;
  height: 100%;
  margin: 0;
  border-radius: 0;
}

@media screen and (max-width: 768px) {
  .preview-dialog :deep(.el-dialog) {
    width: 95% !important;
    margin: 10px auto !important;
  }

  .preview-controls .el-button {
    padding: 6px 12px;
  }
}

@media screen and (max-width: 480px) {
  .preview-dialog :deep(.el-dialog) {
    width: 98% !important;
    margin: 5px auto !important;
  }

  .preview-dialog :deep(.el-dialog__body) {
    padding: 0;
  }

  .preview-controls .el-button {
    padding: 4px 8px;
  }
}
</style>

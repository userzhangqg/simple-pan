<template>
  <div class="pdf-viewer">
    <div class="pdf-container">
      <div class="pdf-content">
        <vue-pdf-embed
          :source="url"
          class="pdf-document"
          @rendered="handleRendered"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import VuePdfEmbed from 'vue-pdf-embed'

defineProps({
  url: {
    type: String,
    required: true
  }
})

const handleRendered = () => {
  // 添加渲染完成的类名，用于淡入动画
  const pdfEmbed = document.querySelector('.vue-pdf-embed')
  if (pdfEmbed) {
    pdfEmbed.classList.add('rendered')
  }
}
</script>

<style scoped>
.pdf-viewer {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f8f9fa;
  position: relative;
}

.pdf-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  scrollbar-width: thin;
  scrollbar-color: #909399 #f4f4f5;
  -webkit-overflow-scrolling: touch;
}

.pdf-container::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.pdf-container::-webkit-scrollbar-track {
  background: #f4f4f5;
  border-radius: 3px;
}

.pdf-container::-webkit-scrollbar-thumb {
  background: #909399;
  border-radius: 3px;
}

.pdf-container::-webkit-scrollbar-thumb:hover {
  background: #606266;
}

.pdf-content {
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.pdf-content :deep(.vue-pdf-embed) {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.pdf-content :deep(.vue-pdf-embed.rendered) {
  opacity: 1;
}

.pdf-content :deep(.vue-pdf-embed > div) {
  width: 100%;
  max-width: 100%;
  background: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s;
}

.pdf-content :deep(.vue-pdf-embed > div:hover) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.pdf-content :deep(canvas) {
  width: 100% !important;
  height: auto !important;
  display: block;
}

@media screen and (max-width: 768px) {
  .pdf-container {
    padding: 10px;
  }
  
  .pdf-content {
    gap: 12px;
  }
  
  .pdf-content :deep(.vue-pdf-embed) {
    gap: 12px;
  }
}

@supports (-webkit-touch-callout: none) {
  .pdf-container {
    padding: 20px 16px;
  }
}
</style>

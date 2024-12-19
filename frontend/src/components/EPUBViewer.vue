<template>
  <div class="epub-viewer">
    <div class="epub-container">
      <div class="viewer-container" ref="viewerRef"></div>
      <div class="epub-controls">
        <el-button-group>
          <el-button @click="prevPage" :disabled="!book" :title="$t('message.file.preview.prevPage')">
            <el-icon><arrow-left /></el-icon>
          </el-button>
          <el-button @click="nextPage" :disabled="!book" :title="$t('message.file.preview.nextPage')">
            <el-icon><arrow-right /></el-icon>
          </el-button>
        </el-button-group>
        <span class="page-info" v-if="currentPage && totalPages">
          {{ $t('message.file.preview.pageOf', { current: currentPage, total: totalPages }) }}
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import ePub from 'epubjs'
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'EPUBViewer',
  components: {
    ArrowLeft,
    ArrowRight
  },
  props: {
    url: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const { t } = useI18n()
    const viewerRef = ref(null)
    const book = ref(null)
    const rendition = ref(null)
    const currentPage = ref(0)
    const totalPages = ref(0)

    const initBook = async () => {
      if (!props.url) return

      try {
        book.value = ePub(props.url)
        
        rendition.value = book.value.renderTo(viewerRef.value, {
          width: '100%',
          height: '100%',
          spread: 'none',
          flow: 'paginated'
        })

        await rendition.value.display()
        currentPage.value = 1

        // 等待电子书加载完成
        await book.value.ready

        // 生成页码定位信息
        await book.value.locations.generate(1600)
        
        // 设置总页数
        totalPages.value = book.value.locations.total

        // 监听位置变化更新页码
        rendition.value.on('relocated', (location) => {
          if (location && location.start) {
            currentPage.value = book.value.locations.locationFromCfi(location.start.cfi)
          }
        })

        // Add keyboard event listeners
        window.addEventListener('keyup', handleKeyPress)
      } catch (error) {
        console.error('Error initializing EPUB:', error)
      }
    }

    const handleKeyPress = (e) => {
      if (!rendition.value) return
      
      switch(e.key) {
        case 'ArrowLeft':
          prevPage()
          break
        case 'ArrowRight':
          nextPage()
          break
      }
    }

    const prevPage = () => {
      if (rendition.value) {
        rendition.value.prev()
      }
    }

    const nextPage = () => {
      if (rendition.value) {
        rendition.value.next()
      }
    }

    // 监听 url 变化
    watch(() => props.url, (newUrl) => {
      if (book.value) {
        book.value.destroy()
        book.value = null
        rendition.value = null
        currentPage.value = 0
        totalPages.value = 0
      }
      if (newUrl) {
        initBook()
      }
    })

    onMounted(() => {
      initBook()
    })

    onBeforeUnmount(() => {
      window.removeEventListener('keyup', handleKeyPress)
      if (book.value) {
        book.value.destroy()
      }
    })

    return {
      viewerRef,
      book,
      currentPage,
      totalPages,
      prevPage,
      nextPage
    }
  }
}
</script>

<style scoped>
.epub-viewer {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  position: relative;
}

.epub-container {
  flex: 1;
  display: flex;
  position: relative;
  overflow: hidden;
}

.viewer-container {
  flex: 1;
  height: 100%;
  overflow: hidden;
}

.epub-controls {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 15px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid #eee;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 1000;
  min-width: 180px;
  backdrop-filter: blur(8px);
}

.epub-controls :deep(.el-button-group) {
  display: flex;
  width: 100%;
}

.epub-controls :deep(.el-button) {
  flex: 1;
  height: 36px;
  padding: 8px 16px;
  font-size: 16px;
}

.epub-controls :deep(.el-icon) {
  font-size: 18px;
  width: 18px;
  height: 18px;
}

.page-info {
  text-align: center;
  font-size: 14px;
  color: #606266;
  padding: 4px 8px;
  white-space: nowrap;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 4px;
  min-width: 120px;
}

@media screen and (max-width: 768px) {
  .epub-controls {
    bottom: 15px;
    min-width: 160px;
    padding: 12px;
  }

  .epub-controls :deep(.el-button) {
    height: 32px;
    padding: 6px 12px;
    font-size: 14px;
  }

  .epub-controls :deep(.el-icon) {
    font-size: 16px;
    width: 16px;
    height: 16px;
  }
  
  .page-info {
    font-size: 13px;
    min-width: 100px;
  }
}
</style>

export default {
  message: {
    hello: '你好',
    welcome: '网盘系统',
    upload: '上传',
    download: '下载',
    delete: '删除',
    preview: '预览',
    cancel: '取消',
    confirm: '确认',
    success: '成功',
    error: '错误',
    warning: '警告',
    downloading: '下载中...',
    deleting: '删除中...',
    file: {
      name: '文件名',
      size: '大小',
      type: '类型',
      lastModified: '最后修改时间',
      actions: '操作',
      upload: {
        success: '文件上传成功',
        error: '文件上传失败',
        dragText: '拖拽文件到此处或点击上传',
        sizeError: '文件大小不能超过100MB',
        typeError: '不支持的文件类型，仅支持以下格式：',
        supportedTypes: '支持的文件类型：.jpg, .jpeg, .png, .gif, .pdf, .doc, .docx, .xls, .xlsx, .txt, .epub'
      },
      download: {
        success: '文件下载成功',
        error: '文件下载失败'
      },
      delete: {
        confirm: '确定要删除这个文件吗？',
        success: '文件删除成功',
        error: '文件删除失败'
      },
      list: {
        error: '获取文件列表失败'
      },
      preview: {
        notSupported: '该文件类型不支持预览',
        outline: '目录',
        error: '文件加载失败',
        loading: '正在加载文件...',
        prevPage: '上一页',
        nextPage: '下一页',
        pageOf: '第 {current} 页，共 {total} 页'
      },
      search: {
        placeholder: '搜索文件名...',
        noMatch: '没有找到匹配的文件'
      },
      pagination: {
        total: '共 {total} 个文件'
      }
    }
  }
}

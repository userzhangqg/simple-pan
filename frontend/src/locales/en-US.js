export default {
  message: {
    hello: 'Hello',
    welcome: 'File System',
    upload: 'Upload',
    download: 'Download',
    delete: 'Delete',
    preview: 'Preview',
    cancel: 'Cancel',
    confirm: 'Confirm',
    success: 'Success',
    error: 'Error',
    warning: 'Warning',
    downloading: 'Downloading...',
    deleting: 'Deleting...',
    file: {
      name: 'File Name',
      size: 'Size',
      type: 'Type',
      lastModified: 'Last Modified',
      actions: 'Actions',
      upload: {
        success: 'File uploaded successfully',
        error: 'Failed to upload file',
        dragText: 'Drop files here or click to upload',
        sizeError: 'File size cannot exceed 100MB',
        typeError: 'Unsupported file type, only the following formats are supported:',
        supportedTypes: 'Supported file types: .jpg, .jpeg, .png, .gif, .pdf, .doc, .docx, .xls, .xlsx, .txt, .epub'
      },
      download: {
        success: 'File downloaded successfully',
        error: 'Failed to download file'
      },
      delete: {
        confirm: 'Are you sure you want to delete this file?',
        success: 'File deleted successfully',
        error: 'Failed to delete file'
      },
      list: {
        error: 'Failed to get file list'
      },
      preview: {
        notSupported: 'Preview not supported for this file type',
        outline: 'Outline',
        error: 'Failed to load file',
        loading: 'Loading file...',
        prevPage: 'Previous Page',
        nextPage: 'Next Page',
        pageOf: 'Page {current} of {total}'
      },
      search: {
        placeholder: 'Search files...',
        noMatch: 'No matching files found'
      },
      pagination: {
        total: 'Total {total} files'
      }
    }
  }
}

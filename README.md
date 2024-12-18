# 简单网盘系统

一个基于 Vue 3 + Go 的现代网盘系统，支持文件上传、下载、预览等功能。

## 功能特性

- 文件上传：支持拖拽上传和点击上传
- 文件预览：支持图片和 PDF 文件预览
- 文件下载：支持单文件下载
- 文件管理：支持文件删除和搜索
- 国际化：支持中文和英文界面切换
- 响应式设计：支持移动端和桌面端

## Roadmap

### 1.0.x (当前版本)
- [x] 基础文件管理功能
- [x] 图片预览
- [x] PDF 预览
- [x] 国际化支持
- [x] 响应式设计

### 1.1.x (下一版本)
- [ ] 用户系统
  - [ ] 用户注册
  - [ ] 用户登录
  - [ ] 权限管理
- [ ] 文件分享功能
  - [ ] 生成分享链接
  - [ ] 设置分享密码
  - [ ] 设置分享有效期
- [ ] 文件管理增强
  - [ ] 批量操作
  - [ ] 文件夹支持
  - [ ] 文件移动

### 1.2.x
- [ ] 在线预览增强
  - [ ] Office 文档预览
  - [ ] 音视频预览
  - [ ] 代码文件预览
- [ ] 存储管理
  - [ ] 存储容量限制
  - [ ] 存储空间统计
  - [ ] 回收站功能

### 2.0.x
- [ ] 多端同步
  - [ ] 桌面客户端
  - [ ] 移动端 App
  - [ ] 文件自动同步
- [ ] 协作功能
  - [ ] 团队空间
  - [ ] 文件协作编辑
  - [ ] 版本控制

### 未来计划
- [ ] 云存储集成
  - [ ] 对象存储支持
  - [ ] 多云存储
- [ ] AI 功能
  - [ ] 智能文件分类
  - [ ] 图片识别
  - [ ] 内容搜索
- [ ] 安全增强
  - [ ] 文件加密
  - [ ] 病毒扫描
  - [ ] 访问审计

## 技术栈

### 前端
- Vue 3
- Vite
- Element Plus
- Vue I18n
- Vue PDF Embed
- Axios

### 后端
- Go
- Gin
- SQLite

## 快速开始

### 使用 Docker Compose

1. 克隆仓库：
\`\`\`bash
git clone git@github.com:userzhangqg/simple-pan.git
cd simple-pan
\`\`\`

2. 启动服务：
\`\`\`bash
docker-compose up -d
\`\`\`

3. 访问系统：
打开浏览器访问 http://localhost:5173

### 手动部署

#### 前端

1. 安装依赖：
\`\`\`bash
cd frontend
npm install
\`\`\`

2. 开发模式：
\`\`\`bash
npm run dev
\`\`\`

3. 构建生产版本：
\`\`\`bash
npm run build
\`\`\`

#### 后端

1. 安装依赖：
\`\`\`bash
cd backend
go mod download
\`\`\`

2. 运行服务：
\`\`\`bash
go run main.go
\`\`\`

## 环境变量

### 前端
- `VITE_API_BASE_URL`: API 基础 URL，默认为 http://localhost:8081/api

### 后端
- `PORT`: 服务端口，默认为 8081

## 目录结构

\`\`\`
.
├── frontend/               # 前端项目
│   ├── src/               # 源代码
│   ├── public/            # 静态资源
│   └── package.json       # 项目配置
├── backend/               # 后端项目
│   ├── handlers/          # 请求处理器
│   ├── utils/             # 工具函数
│   └── main.go           # 入口文件
├── docker-compose.yml    # Docker Compose 配置
├── Dockerfile.frontend   # 前端 Dockerfile
└── Dockerfile.backend    # 后端 Dockerfile
\`\`\`

## 贡献指南

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/AmazingFeature`
3. 提交更改：`git commit -m 'Add some AmazingFeature'`
4. 推送分支：`git push origin feature/AmazingFeature`
5. 提交 Pull Request

## 许可证

MIT License

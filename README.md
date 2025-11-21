# 🎲 彩票数据管理系统

一个功能强大的彩票数据管理、分析和智能推荐系统，支持**双色球**和**大乐透**两种彩票类型。

**🎯 双模式支持**: Web应用 + 跨平台桌面应用（Windows/macOS）

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![Wails](https://img.shields.io/badge/Wails-v2.11.0-DF0032?style=flat)](https://wails.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📋 目录

- [功能特性](#-功能特性)
- [技术栈](#-技术栈)
- [运行模式](#-运行模式)
- [快速开始](#-快速开始)
- [详细使用](#-详细使用)
- [项目结构](#-项目结构)
- [开发指南](#-开发指南)
- [部署](#-部署)
- [文档](#-文档)
- [贡献](#-贡献)
- [许可证](#-许可证)

---

## ✨ 功能特性

### 🎰 彩票数据管理
- ✅ **数据获取**: 从官方API自动获取最新开奖数据
- ✅ **批量导入**: 支持批量获取历史数据（1-2000期）
- ✅ **智能去重**: 自动跳过已存在数据，避免重复
- ✅ **本地存储**: SQLite数据库，轻量高效

### 📊 数据分析
- ✅ **走势分析**: 可视化展示号码出现趋势
- ✅ **统计分析**: 热号、冷号、遗漏值统计
- ✅ **图表展示**: 直观的图表展示分析结果
- ✅ **自定义周期**: 支持自定义分析期数

### 🎯 智能推荐
- ✅ **多策略推荐**: 基于热号、冷号、冷热结合三种策略
- ✅ **批量生成**: 一次生成多注推荐号码
- ✅ **置信度标注**: 每注号码附带推荐依据
- ✅ **实时更新**: 基于最新100期数据动态推荐

### 💎 用户体验
- ✅ **响应式设计**: 完美适配桌面、平板、手机
- ✅ **暗黑模式**: 护眼模式，可自由切换
- ✅ **PWA支持**: Web版可安装到桌面（离线使用）
- ✅ **双模式**: Web浏览器 + 原生桌面应用

---

## 🛠️ 技术栈

### 后端
- **Go 1.25+**: 高性能后端服务
- **Gin**: Web框架
- **GORM**: ORM数据库操作
- **SQLite**: 轻量级数据库
- **Redis**: 缓存服务（Web模式可选）
- **Wails v2**: 桌面应用框架

### 前端
- **Vue 3**: 渐进式前端框架
- **Vite**: 快速构建工具
- **Element Plus**: UI组件库
- **ECharts**: 数据可视化
- **Axios**: HTTP客户端
- **Pinia**: 状态管理

### 桌面应用
- **Wails v2.11.0**: Go + Web技术构建原生应用
- **API适配层**: 自动适配Web/桌面环境
- **统一代码库**: 一套代码，多种运行方式

---

## 🚀 运行模式

本项目支持**两种运行模式**，使用同一套代码库：

### 1️⃣ Web应用模式

**特点**:
- 浏览器访问
- 支持PWA（可安装）
- 支持Redis缓存
- 支持异步批量任务
- 支持Docker部署

**适用场景**: 
- 需要多用户访问
- 云端部署
- 移动端访问

### 2️⃣ 桌面应用模式

**特点**:
- 原生应用体验
- 无需浏览器
- 数据本地存储
- 跨平台支持
- 性能更优

**适用场景**:
- 个人使用
- 离线使用
- 隐私要求高

---

## 🏃 快速开始

### 方式一: 桌面应用（推荐个人用户）

**直接下载运行**（无需安装Go/Node）：

1. **下载应用**
   - 从 [Releases](https://github.com/your-repo/releases) 下载对应平台的安装包
   - macOS (ARM): `彩票助手-macOS-arm64.app`
   - macOS (Intel): `彩票助手-macOS-amd64.app`
   - Windows: `彩票助手-Windows-amd64.exe`

2. **安装运行**
   ```bash
   # macOS - 直接双击.app运行
   # 或命令行运行
   open 彩票助手-macOS-arm64.app
   
   # Windows - 直接双击.exe运行
   ```

3. **开始使用** 🎉
   - 应用会自动创建数据库
   - 数据存储位置:
     - macOS: `~/Library/Application Support/caipiao/`
     - Windows: `%APPDATA%/caipiao/`

### 方式二: Web应用（推荐团队/服务器部署）


#### 前置要求
- Go 1.25+
- Node.js 18+
- Redis（可选）

#### 安装步骤

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/caipiao.git
cd caipiao

# 2. 启动后端
cd backend
go mod download
go run .

# 3. 启动前端（新终端）
cd frontend
npm install
npm run dev

# 4. 访问
# 打开浏览器访问: http://localhost:5173
```

#### 使用Docker（最简单）

```bash
# 1. 使用docker-compose启动
docker-compose up -d

# 2. 访问
# 打开浏览器访问: http://localhost
```

---

## 📖 详细使用

### 桌面应用

#### 获取数据

1. **获取最新数据**
   - 点击"获取最新数据"按钮
   - 自动从官方API获取最新一期

2. **批量获取历史数据**
   - 点击"批量获取历史数据"
   - 输入期数（如100）
   - 点击"开始获取"
   - 耐心等待（每期约1秒）

3. **指定期数获取**
   - 点击"指定期数获取"
   - 输入期号（如2024001）
   - 点击"确定"

#### 数据分析

1. **查看走势**
   - 点击顶部导航"双色球走势"或"大乐透走势"
   - 可视化图表展示号码趋势
   - 支持自定义分析期数

2. **统计分析**
   - 在数据列表页点击"数据统计"
   - 查看热号、冷号频率
   - 红球/蓝球分别统计

#### 智能推荐

- 点击顶部导航"双色球推荐"或"大乐透推荐"
- 选择生成注数（1-10注）
- 点击"生成推荐"
- 查看推荐号码及依据

### Web应用

使用方式与桌面应用基本相同，额外功能：

- **PWA安装**: 浏览器地址栏点击安装图标
- **多标签页**: 可同时查看多个页面
- **异步任务**: 批量获取大量数据时支持后台任务

---

## 📁 项目结构

```
caipiao/
├─ backend/              # Go后端
│  ├─ api/              # API路由和处理器
│  ├─ config/           # 配置管理
│  ├─ database/         # 数据库初始化
│  ├─ middleware/       # 中间件
│  ├─ models/           # 数据模型
│  ├─ services/         # 业务逻辑
│  └─ main.go          # Web服务入口
│
├─ frontend/            # Vue前端
│  ├─ public/          # 静态资源
│  ├─ src/             
│  │  ├─ api/         # API调用（含适配层）
│  │  ├─ assets/      # 图片、样式
│  │  ├─ components/  # 组件
│  │  ├─ router/      # 路由
│  │  ├─ stores/      # 状态管理
│  │  ├─ utils/       # 工具函数
│  │  ├─ views/       # 页面
│  │  ├─ App.vue      # 根组件
│  │  └─ main.js      # 入口文件
│  ├─ vite.config.js  # Vite配置
│  └─ package.json
│
├─ build/               # Wails构建配置
├─ .github/workflows/   # GitHub Actions
├─ app.go              # Wails应用入口
├─ main.go             # Wails主函数
├─ utils.go            # Wails工具函数
├─ wails.json          # Wails配置
├─ go.mod              # Go依赖
├─ docker-compose.yml  # Docker配置
└─ README.md           # 本文件
```

---

## 👨‍💻 开发指南

### 桌面应用开发

#### 环境要求
- Go 1.25+
- Node.js 18+
- Wails CLI v2.11.0

```bash
# 安装Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 验证安装
wails doctor
```

#### 开发调试

```bash
# 开发模式（热重载）
wails dev
```

#### 构建应用

```bash
# 构建当前平台
wails build

# 构建所有平台（需要相应工具链）
wails build -platform windows/amd64
wails build -platform darwin/universal
```

#### 添加新功能

1. **后端**: 在 `app.go` 添加新方法
   ```go
   func (a *App) YourNewMethod(param string) (interface{}, error) {
       // 实现逻辑
       return result, nil
   }
   ```

2. **重新生成绑定**
   ```bash
   wails generate module
   ```

3. **前端**: 更新API适配层 (`frontend/src/api/adapter.js`)
   ```javascript
   async yourNewMethod(param) {
       const wails = await getWailsAPI()
       if (wails) {
           return await wails.YourNewMethod(param)
       }
       return axios.get('/api/your-endpoint')
   }
   ```

### Web应用开发

#### 本地开发

```bash
# 启动后端（终端1）
cd backend
go run .

# 启动前端（终端2）
cd frontend
npm run dev
```

#### 添加API

1. **后端**: `backend/api/handlers.go` 添加处理器
2. **前端**: `frontend/src/api/` 添加API调用

### 代码规范

- **Go**: 使用 `go fmt` 格式化
- **Vue**: 使用 ESLint + Prettier
- **提交**: 使用语义化提交消息

---

## 🚢 部署

### Docker部署（推荐）

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 手动部署

#### 后端

```bash
cd backend
go build -o caipiao
./caipiao
```

#### 前端

```bash
cd frontend
npm run build
# 将 dist/ 目录部署到Nginx/Apache
```

### 桌面应用分发

#### GitHub Actions自动构建

推送tag自动构建多平台版本：

```bash
git tag v1.0.0
git push origin v1.0.0
```

自动生成：
- Windows安装包
- macOS应用（Intel + ARM）
- 自动发布到GitHub Releases

#### 手动构建

```bash
# macOS
wails build -platform darwin/universal

# Windows (需要在Windows上或配置交叉编译)
wails build -platform windows/amd64

# Linux
wails build -platform linux/amd64
```

---

## 📚 文档

- [Wails桌面应用设置](WAILS_SETUP.md) - 技术实现详情
- [Wails使用指南](WAILS_USAGE.md) - API适配层使用
- [批量获取功能](BATCH_FETCH_FIXED.md) - 批量获取说明
- [调试指南](DEBUG_BATCH_FETCH.md) - 问题诊断
- [成功总结](WAILS_SUCCESS.md) - 项目完成情况

---

## 🎯 路线图

- [x] 基础数据管理功能
- [x] 走势分析和统计
- [x] 智能推荐算法
- [x] Web应用完整功能
- [x] PWA支持
- [x] Docker部署
- [x] Wails桌面应用
- [x] GitHub Actions自动构建
- [ ] 更多推荐策略
- [ ] 数据导出功能
- [ ] 用户系统（Web版）
- [ ] 移动端优化

---

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

## 🙏 致谢

- [Wails](https://wails.io) - 优秀的Go桌面应用框架
- [Vue.js](https://vuejs.org) - 渐进式JavaScript框架
- [Element Plus](https://element-plus.org) - Vue 3 UI组件库
- [ECharts](https://echarts.apache.org) - 强大的可视化库

---

## 📧 联系方式

- 项目主页: [GitHub](https://github.com/your-repo/caipiao)
- 问题反馈: [Issues](https://github.com/your-repo/caipiao/issues)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个Star支持一下！**

Made with ❤️ by [Your Name]

</div>

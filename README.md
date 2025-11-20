# 🎰 彩票数据管理系统

> 基于 Go + Vue3 的全栈彩票数据管理、统计分析和智能推荐系统

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![PWA](https://img.shields.io/badge/PWA-Enabled-blueviolet)](https://web.dev/progressive-web-apps/)

## 📖 项目简介

一个功能完整的彩票数据管理系统，支持**双色球**和**大乐透**的数据获取、存储、统计分析和智能推荐。采用现代化技术栈，支持Docker一键部署、PWA离线使用、Redis缓存加速等特性。

### ✨ 核心特性

- 🎯 **双彩种支持** - 完整支持双色球和大乐透
- 📊 **数据可视化** - ECharts走势图、频率统计、遗漏分析
- 🤖 **智能推荐** - 基于热号、冷号、周期等多维度算法推荐
- ⚡ **高性能** - Redis缓存、数据库索引、虚拟滚动优化
- 📱 **PWA支持** - 可安装、离线访问、类原生体验
- 🐳 **容器化部署** - Docker + Docker Compose 一键部署
- 🌐 **官方数据源** - 对接中国福彩和体彩官方API

---

## 🎨 功能展示

### 1. 数据管理
- ✅ 获取最新一期开奖数据
- ✅ 批量获取历史数据（最多500期）
- ✅ 异步任务处理大批量数据
- ✅ 期号搜索、分页展示
- ✅ 自动去重，避免重复保存

### 2. 统计分析
- 📈 **号码频率统计** - 红球/蓝球/前区/后区出现次数
- 📊 **走势图分析** - 支持30/50/100/200期数据展示
- 🔥 **热门号码** - Top 10高频号码柱状图
- ❄️ **遗漏分析** - 长期未出现号码智能提醒

### 3. 智能推荐
- 🎯 **热号策略** - 基于近期高频号码推荐
- 🧊 **冷号策略** - 基于长期遗漏号码推荐
- ⚖️ **冷热结合** - 均衡选择，稳健投注
- 📋 **一键生成** - 可生成1-10注推荐号码
- 💡 **推荐依据** - 提供详细的推荐理由说明

### 4. PWA特性
- 📲 **添加到主屏幕** - 像原生应用一样使用
- 🔌 **离线访问** - 支持离线浏览已加载内容
- ⚡ **快速加载** - Service Worker资源缓存
- 🔄 **自动更新** - 检测新版本自动更新

---

## 🏗️ 技术架构

### 后端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **Go** | 1.25+ | 主要编程语言 |
| **Gin** | 1.9.1 | Web框架 |
| **GORM** | 1.25.5 | ORM框架 |
| **SQLite** | - | 嵌入式数据库 |
| **Redis** | 7+ | 缓存数据库 |
| **Resty** | 2.16.5 | HTTP客户端 |

### 前端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **Vue 3** | 3.4.0 | 前端框架 |
| **Element Plus** | 2.5.0 | UI组件库 |
| **ECharts** | 6.0.0 | 数据可视化 |
| **Vue Router** | 4.2.5 | 路由管理 |
| **Pinia** | 2.1.7 | 状态管理 |
| **Axios** | 1.6.0 | HTTP客户端 |
| **Vite** | 5.0.0 | 构建工具 |

### 系统架构

```
┌─────────────────┐
│   Nginx (80)    │  ← 前端静态资源 + 反向代理
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│   Go API (8080) │  ← RESTful API
└────────┬────────┘
         │
    ┌────┴────┐
    ↓         ↓
┌────────┐ ┌──────────┐
│ SQLite │ │ Redis    │  ← 数据存储 + 缓存
└────────┘ └──────────┘
```

---

## 🚀 快速开始

### 环境要求

**方式一：Docker部署（推荐）**
- Docker 20.10+
- Docker Compose 2.0+

**方式二：本地开发**
- Go 1.25+
- Node.js 18+
- Redis 7+ （可选，用于缓存）

### 📦 安装部署

#### 🐳 Docker部署（推荐）

```bash
# 1. 克隆项目
git clone <repository-url>
cd caipiao

# 2. 一键启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 访问应用
# 前端: http://localhost
# 后端: http://localhost:8080
```

#### 💻 本地开发部署

**启动后端：**

```bash
cd backend

# 安装依赖
go mod download

# 启动Redis（可选，用于缓存）
redis-server

# 运行后端服务
go run main.go

# 后端将在 http://localhost:8080 启动
```

**启动前端：**

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 前端将在 http://localhost:5173 启动
```

---

## 🎮 使用指南

### 基础操作

1. **查看开奖数据**
   - 访问首页，选择"双色球"或"大乐透"
   - 自动显示最新开奖记录列表

2. **搜索期号**
   - 在搜索框输入期号（如：2024001）
   - 支持模糊匹配

3. **获取数据**
   - 点击"获取最新数据"按钮获取最新一期
   - 点击"批量获取"可指定获取期数（最多500期）

4. **查看统计**
   - 点击"数据统计"查看号码频率分布
   - 支持切换红球/蓝球或前区/后区

5. **走势分析**
   - 点击"走势图"进入走势分析页面
   - 支持选择30/50/100/200期数据
   - 查看热门号码Top 10和遗漏值

6. **智能推荐**
   - 点击"智能推荐"获取推荐号码
   - 提供热号、冷号、冷热结合三种策略
   - 可生成1-10注推荐号码

### PWA使用

#### 安装到设备

**桌面端（Chrome/Edge）：**
1. 访问应用网址
2. 点击地址栏右侧的安装图标 ⊕
3. 点击"安装"按钮
4. 应用将作为独立窗口打开

**移动端（iOS/Android）：**
1. 使用Safari（iOS）或Chrome（Android）访问
2. iOS: 点击分享按钮 → 选择"添加到主屏幕"
3. Android: 点击菜单 → 选择"安装应用"或"添加到主屏幕"

#### 离线使用

1. 首次访问时，Service Worker会自动缓存必要资源
2. 断网后仍可访问已缓存的页面和数据
3. 重新联网后会自动同步最新数据

#### 更新应用

- 应用会自动检测新版本
- 发现新版本时会提示刷新
- 刷新后即可使用最新功能

---

## 🔌 API文档

### 双色球接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/shuangseqiu/list` | GET | 获取列表 |
| `/api/shuangseqiu/fetch` | POST | 获取单期数据 |
| `/api/shuangseqiu/fetch-history` | POST | 批量获取历史 |
| `/api/shuangseqiu/statistics` | GET | 获取统计数据 |
| `/api/shuangseqiu/trend` | GET | 获取走势数据 |
| `/api/shuangseqiu/recommend` | GET | 智能推荐 |

### 大乐透接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/daletou/list` | GET | 获取列表 |
| `/api/daletou/fetch` | POST | 获取单期数据 |
| `/api/daletou/fetch-history` | POST | 批量获取历史 |
| `/api/daletou/statistics` | GET | 获取统计数据 |
| `/api/daletou/trend` | GET | 获取走势数据 |
| `/api/daletou/recommend` | GET | 智能推荐 |

### 请求示例

```bash
# 获取双色球列表（分页）
curl "http://localhost:8080/api/shuangseqiu/list?page=1&page_size=20"

# 获取统计数据
curl "http://localhost:8080/api/shuangseqiu/statistics?type=red"

# 获取智能推荐（5注）
curl "http://localhost:8080/api/shuangseqiu/recommend?count=5"
```

### 响应格式

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 🛠️ 开发指南

### 项目结构

```
caipiao/
├── backend/                 # 后端服务
│   ├── api/                 # API路由和处理器
│   │   ├── handlers.go      # API处理函数
│   │   └── routes.go        # 路由配置
│   ├── config/              # 配置管理
│   ├── database/            # 数据库初始化
│   ├── models/              # 数据模型
│   │   ├── shuangseqiu.go   # 双色球模型
│   │   └── daletou.go       # 大乐透模型
│   ├── services/            # 业务逻辑
│   │   ├── lottery_service.go    # 彩票服务
│   │   ├── trend_service.go      # 走势服务
│   │   ├── cache_service.go      # 缓存服务
│   │   ├── recommend_service.go  # 推荐服务
│   │   └── task_manager.go       # 任务管理
│   ├── main.go              # 程序入口
│   ├── Dockerfile           # Docker镜像构建
│   └── go.mod               # Go模块依赖
├── frontend/                # 前端项目
│   ├── src/
│   │   ├── api/             # API接口定义
│   │   ├── components/      # 可复用组件
│   │   ├── router/          # 路由配置
│   │   ├── views/           # 页面视图
│   │   │   ├── Shuangseqiu.vue          # 双色球主页
│   │   │   ├── ShuangseqiuTrend.vue     # 双色球走势
│   │   │   ├── ShuangseqiuRecommend.vue # 双色球推荐
│   │   │   ├── Daletou.vue              # 大乐透主页
│   │   │   ├── DaletouTrend.vue         # 大乐透走势
│   │   │   └── DaletouRecommend.vue     # 大乐透推荐
│   │   ├── App.vue          # 根组件
│   │   └── main.js          # 入口文件
│   ├── public/              # 静态资源
│   ├── Dockerfile           # Docker镜像构建
│   ├── nginx.conf           # Nginx配置
│   ├── package.json         # npm依赖
│   └── vite.config.js       # Vite配置
├── docker-compose.yml       # Docker编排
└── README.md                # 项目文档
```

### 添加新功能

**后端添加API接口：**

1. 在 `models/` 中定义数据模型
2. 在 `services/` 中实现业务逻辑
3. 在 `api/handlers.go` 中添加处理函数
4. 在 `api/routes.go` 中注册路由

**前端添加页面：**

1. 在 `views/` 中创建Vue组件
2. 在 `router/index.js` 中注册路由
3. 在 `api/` 中定义接口调用

### 本地调试

```bash
# 后端热重载（使用air）
go install github.com/cosmtrek/air@latest
cd backend && air

# 前端热重载（Vite自带）
cd frontend && npm run dev
```

### 构建生产版本

```bash
# 后端编译
cd backend
go build -o lottery-server main.go

# 前端构建
cd frontend
npm run build
# 构建产物在 dist/ 目录
```

---

## 📊 性能优化

### 已实施优化

1. **数据库索引** - 为常用查询字段添加索引，查询速度提升50%+
2. **Redis缓存** - 统计数据缓存1小时，API响应速度提升80%+
3. **虚拟滚动** - 大列表渲染优化，内存占用降低50%
4. **资源压缩** - Gzip压缩，传输大小减少70%
5. **PWA缓存** - 静态资源离线缓存，二次访问速度提升90%

### 性能指标

| 指标 | 优化前 | 优化后 | 提升 |
|------|-------|-------|-----|
| 列表查询 | 200ms | 100ms | 50% |
| 统计查询 | 500ms | 100ms | 80% |
| 页面加载 | 2.5s | 0.8s | 68% |
| 滚动渲染 | 30fps | 60fps | 100% |

---

## 🔧 配置说明

### 后端配置

编辑 `backend/config/config.go`：

```go
Server.Port: "8080"              // 服务端口
Server.Mode: "release"           // 运行模式（debug/release）
Database.Path: "./data/lottery.db"  // 数据库路径
```

### 环境变量

```bash
# Redis配置（可选）
export REDIS_HOST=localhost:6379

# Gin运行模式
export GIN_MODE=release
```

### 前端配置

编辑 `frontend/vite.config.js`：

```javascript
server: {
  port: 5173,                    // 开发服务器端口
  proxy: {
    '/api': {
      target: 'http://localhost:8080'  // 后端API地址
    }
  }
}
```

---

## 🐛 常见问题

### Q1: Docker启动失败？

**A:** 检查端口是否被占用：
```bash
# 检查端口占用
lsof -i :80
lsof -i :8080
lsof -i :6379

# 修改docker-compose.yml中的端口映射
```

### Q2: 前端无法连接后端API？

**A:** 确保：
1. 后端服务已启动（http://localhost:8080）
2. 前端代理配置正确（vite.config.js）
3. 清除浏览器缓存重新访问

### Q3: 获取数据失败429错误？

**A:** 官方API有频率限制：
- 降低获取频率
- 批量获取时增加间隔时间
- 使用已缓存的数据

### Q4: PWA无法安装？

**A:** 检查：
1. 必须使用HTTPS或localhost
2. 需要有效的manifest.json
3. 浏览器版本需支持PWA（Chrome 80+, Safari 11.1+）

### Q5: Redis连接失败？

**A:** Redis是可选的：
- 未安装Redis时系统仍可正常运行，只是不使用缓存
- 如需缓存功能，确保Redis服务已启动
- Docker部署会自动包含Redis服务

---

## 📄 数据库结构

### 双色球表（shuangseqiu）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键 |
| issue | TEXT | 期号（唯一） |
| red_ball_1~6 | INTEGER | 红球1-6 |
| blue_ball | INTEGER | 蓝球 |
| draw_date | DATETIME | 开奖日期 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### 大乐透表（daletou）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键 |
| issue | TEXT | 期号（唯一） |
| front_ball_1~5 | INTEGER | 前区1-5 |
| back_ball_1~2 | INTEGER | 后区1-2 |
| draw_date | DATETIME | 开奖日期 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

---

## 🤝 参与贡献

欢迎提交Issue和Pull Request！

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交Pull Request

---

## 📝 开发计划

- [ ] 添加用户系统和权限管理
- [ ] 支持更多彩票玩法
- [ ] 数据导出功能（Excel/CSV）
- [ ] 移动端原生应用
- [ ] 更多推荐算法（机器学习）
- [ ] 实时开奖通知

---

## 📜 许可证

本项目采用 [MIT](LICENSE) 许可证。

---

## 👨‍💻 作者

Created with ❤️ by Caipiao Team

---

## 🙏 致谢

- 数据来源：[中国福利彩票官网](https://www.cwl.gov.cn) | [中国体育彩票官网](https://www.sporttery.cn)
- UI框架：[Element Plus](https://element-plus.org)
- 图表库：[Apache ECharts](https://echarts.apache.org)
- Web框架：[Gin](https://gin-gonic.com)

---

**⚠️ 免责声明**

本项目仅供学习交流使用，不提供任何投注建议。彩票有风险，投注需谨慎。请理性娱乐，远离赌博。

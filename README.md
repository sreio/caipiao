# 彩票数据管理系统

## 项目简介

这是一个基于 Go + Gin + SQLite + Vue3 + Element Plus 开发的彩票数据管理系统，支持双色球和大乐透数据的展示、查询、统计等功能。

## 技术栈

### 后端
- **Go 1.25** - 编程语言
- **Gin** - Web框架
- **Resty** - HTTP客户端库
- **SQLite** - 本地数据库
- **GORM** - ORM框架

### 前端
- **Vue 3** - 前端框架
- **Element Plus** - UI组件库
- **ECharts** - 数据可视化图表库
- **Vue Router** - 路由管理
- **Pinia** - 状态管理
- **Axios** - HTTP客户端
- **Vite** - 构建工具

## 项目结构

```
caipiao/
├── backend/                    # 后端服务
│   ├── api/                   # API路由和控制器
│   │   ├── handlers.go        # API处理器
│   │   └── routes.go          # 路由配置
│   ├── config/                # 配置文件
│   │   └── config.go          # 应用配置
│   ├── database/              # 数据库
│   │   └── database.go        # 数据库初始化
│   ├── models/                # 数据模型
│   │   ├── shuangseqiu.go    # 双色球模型
│   │   └── daletou.go         # 大乐透模型
│   ├── services/              # 业务服务
│   │   └── lottery_service.go # 彩票服务
│   ├── main.go                # 主程序入口
│   └── go.mod                 # Go模块配置
├── frontend/                   # 前端项目
│   ├── src/
│   │   ├── api/               # API接口
│   │   ├── components/        # 组件
│   │   ├── router/            # 路由
│   │   ├── views/             # 页面
│   │   ├── App.vue            # 根组件
│   │   └── main.js            # 入口文件
│   ├── index.html             # HTML模板
│   ├── package.json           # 依赖配置
│   └── vite.config.js         # Vite配置
├── start-backend.sh           # 后端启动脚本
├── start-frontend.sh          # 前端启动脚本
└── README.md                  # 项目说明
```

## 功能特性

### 🎯 核心功能
- ✅ **双色球管理** - 完整的双色球数据管理
- ✅ **大乐透管理** - 完整的大乐透数据管理
- ✅ **数据查询** - 支持按期号搜索
- ✅ **数据统计** - 红球/蓝球/前区/后区出现频率统计
- ✅ **数据获取** - 手动或指定期数获取最新数据
- ✅ **批量获取** - 支持批量获取历史数据（最多500期）
- ✅ **异步任务** - 大批量数据自动异步处理，实时显示进度
- ✅ **分页展示** - 支持分页和每页条数设置
- ✅ **数据去重** - 自动检测重复数据，避免重复保存

### 📊 走势分析（NEW!）
- 📈 **号码走势图** - 基于ECharts的号码走势可视化
- 📊 **频率分析** - 热门号码统计（Top 10柱状图）
- ⚠️ **遗漏值提醒** - 长期未出现号码智能提醒
- 🎯 **多期切换** - 支持30/50/100/200期数据展示
- 🔄 **实时更新** - 数据实时刷新
- 📱 **响应式图表** - 自适应窗口大小

### 🎨 界面特点
- 🎨 现代化UI设计
- 📱 响应式布局
- 🌈 彩色球号展示
- 📊 可视化统计图表
- ⚡ 流畅的交互体验
- 🎯 智能提醒和状态标识

## 快速开始

### 环境要求

**后端:**
- Go 1.25 或更高版本

**前端:**
- Node.js 16+ 
- npm 或 yarn

### 测试API连接（可选）

在启动项目前，可以先测试API是否可用：

```bash
./test-api.sh
```

### 一键启动

#### 方式一：使用启动脚本（推荐）

**启动后端:**
```bash
./start-backend.sh
```

**启动前端:**
```bash
./start-frontend.sh
```

#### 方式二：手动启动

**后端启动:**
```bash
cd backend
go mod tidy           # 下载依赖
go run main.go        # 启动服务
```
后端服务将在 `http://localhost:8080` 启动

**前端启动:**
```bash
cd frontend
npm install           # 安装依赖（首次运行）
npm install echarts   # 安装图表库（首次运行）
npm run dev          # 启动开发服务器
```
前端服务将在 `http://localhost:5173` 启动

### 访问应用

启动成功后，在浏览器中打开 `http://localhost:5173` 即可访问应用。

## API接口文档

### 双色球接口

#### 1. 获取双色球列表
```
GET /api/shuangseqiu/list
参数:
  - page: 页码（默认1）
  - page_size: 每页条数（默认20）
  - issue: 期号（可选，模糊搜索）
```

#### 2. 获取双色球数据
```
POST /api/shuangseqiu/fetch
参数:
  - issue: 期号（可选，不传则获取最新）
```

#### 3. 双色球统计
```
GET /api/shuangseqiu/statistics
参数:
  - type: 类型（red=红球, blue=蓝球）
```

### 大乐透接口

#### 1. 获取大乐透列表
```
GET /api/daletou/list
参数:
  - page: 页码（默认1）
  - page_size: 每页条数（默认20）
  - issue: 期号（可选，模糊搜索）
```

#### 2. 获取大乐透数据
```
POST /api/daletou/fetch
参数:
  - issue: 期号（可选，不传则获取最新）
```

#### 3. 大乐透统计
```
GET /api/daletou/statistics
参数:
  - type: 类型（front=前区, back=后区）
```

## 数据库结构

### 双色球表 (shuangseqiu)
```sql
id          INTEGER     主键
issue       TEXT        期号（唯一）
red_ball_1  INTEGER     红球1
red_ball_2  INTEGER     红球2
red_ball_3  INTEGER     红球3
red_ball_4  INTEGER     红球4
red_ball_5  INTEGER     红球5
red_ball_6  INTEGER     红球6
blue_ball   INTEGER     蓝球
draw_date   DATETIME    开奖日期
created_at  DATETIME    创建时间
updated_at  DATETIME    更新时间
```

### 大乐透表 (daletou)
```sql
id           INTEGER     主键
issue        TEXT        期号（唯一）
front_ball_1 INTEGER     前区1
front_ball_2 INTEGER     前区2
front_ball_3 INTEGER     前区3
front_ball_4 INTEGER     前区4
front_ball_5 INTEGER     前区5
back_ball_1  INTEGER     后区1
back_ball_2  INTEGER     后区2
draw_date    DATETIME    开奖日期
created_at   DATETIME    创建时间
updated_at   DATETIME    更新时间
```

## 配置说明

### 后端配置

编辑 `backend/config/config.go` 修改配置：

```go
Server.Port: "8080"              // 服务端口
Server.Mode: "debug"             // 运行模式（debug/release）
Database.Path: "./data/lottery.db"  // 数据库路径

// API配置（已对接官方API，无需修改）
API.ShuangseqiuURL: "https://www.cwl.gov.cn/..."  // 双色球官方API
API.DaletouURL: "https://webapi.sporttery.cn/..." // 大乐透官方API
```

### 前端配置

编辑 `frontend/vite.config.js` 修改配置：

```javascript
server.port: 5173                // 开发服务器端口
server.proxy: {                  // API代理配置
  '/api': {
    target: 'http://localhost:8080'
  }
}
```

## 开发指南

### 添加新功能

1. **后端添加接口:**
   - 在 `models/` 中定义数据模型
   - 在 `services/` 中实现业务逻辑
   - 在 `api/handlers.go` 中添加处理器
   - 在 `api/routes.go` 中注册路由

2. **前端添加页面:**
   - 在 `views/` 中创建页面组件
   - 在 `router/index.js` 中注册路由
   - 在 `api/` 中定义接口调用

### 构建部署

**后端构建:**
```bash
cd backend
go build -o lottery-server main.go
./lottery-server
```

**前端构建:**
```bash
cd frontend
npm run build
# 构建产物在 dist/ 目录
```

## 📚 详细文档

- **[前端新功能使用指南](./frontend/前端新功能使用指南.md)** - 异步任务、走势图、频率分析、遗漏值提醒详细说明
- **[前端功能完成总结](./前端功能完成总结.md)** - 所有前端功能实现总结
- **[前端功能快速测试指南](./前端功能快速测试指南.md)** - 快速测试所有功能
- **[部署说明](./部署说明.md)** - 生产环境部署指南

## 常见问题

### Q1: 后端启动失败？
- 检查Go版本是否为1.25+
- 确保8080端口未被占用
- 运行 `go mod tidy` 更新依赖

### Q2: 前端访问API失败？
- 确保后端服务已启动
- 检查 `vite.config.js` 中的代理配置
- 查看浏览器控制台错误信息

### Q3: 获取数据失败？
- 检查外部API配置是否正确
- 确认网络连接正常
- 查看后端日志错误信息

### Q4: 走势图不显示？
- 确保已安装ECharts：`npm install echarts`
- 确保数据库有数据（至少30期）
- 查看浏览器控制台错误信息

### Q5: 异步任务一直pending？
- 后端正在处理中，耐心等待（最多5分钟）
- 查看后端日志确认进度
- 如长时间无响应，重启后端服务

## 许可证

MIT License

## 联系方式

如有问题或建议，欢迎提Issue。


# Wails桌面应用使用说明

## ✅ 成功打包！

恭喜！您已成功将彩票系统打包为macOS桌面应用。

## 📍 应用位置

构建的应用位于：
```
build/bin/彩票助手.app
```

## 🔧 前端API适配

已创建 `frontend/src/api/adapter.js` API适配层，自动检测运行环境：
- **桌面应用模式**: 使用Wails Go绑定直接调用后端方法
- **Web模式**: 使用axios发送HTTP请求

## 📝 使用API适配层

### 在Vue组件中使用

将现有的axios调用替换为adapter:

**之前（仅Web模式）**:
```javascript
import axios from 'axios'

// 获取列表
const res = await axios.get('/api/shuangseqiu/list', {
  params: { page: 1, page_size: 20 }
})
```

**之后（Web + 桌面模式）**:
```javascript
import api from '@/api/adapter'

// 获取列表 - 自动适配环境
const res = await api.getShuangseqiuList(1, 20, '')
```

### API方法列表

**双色球**:
- `api.getShuangseqiuList(page, pageSize, issue)` - 获取列表
- `api.fetchShuangseqiu(issue)` - 获取指定期数据
- `api.getShuangseqiuStatistics(ballType)` - 获取统计
- `api.getShuangseqiuTrend(limit)` - 获取走势
- `api.getShuangseqiuRecommendation(count)` - 获取推荐

**大乐透**:
- `api.getDaletouList(page, pageSize, issue)` - 获取列表
- `api.fetchDaletou(issue)` - 获取指定期数据
- `api.getDaletouStatistics(ballType)` - 获取统计
- `api.getDaletouTrend(limit)` - 获取走势
- `api.getDaletouRecommendation(count)` - 获取推荐

## 🔄 更新现有组件

需要更新以下组件以使用API适配层：
1. `frontend/src/views/Shuangseqiu.vue`
2. `frontend/src/views/Daletou.vue`
3. `frontend/src/views/ShuangseqiuTrend.vue`
4. `frontend/src/views/DaletouTrend.vue`
5. `frontend/src/views/ShuangseqiuRecommend.vue`
6. `frontend/src/views/DaletouRecommend.vue`

## 🚀 重新构建步骤

```bash
# 1. 构建前端
cd frontend
npm run build

# 2. 返回根目录构建桌面应用
cd ..
wails build -s -skipbindings

# 3. 运行应用
open build/bin/彩票助手.app
```

## 🐛 常见问题

### Q: 为什么会出现405错误？
A: 桌面应用模式下，前端直接调用Go方法，不走HTTP API。必须使用API适配层。

### Q: 如何同时支持Web和桌面模式？
A: 使用 `frontend/src/api/adapter.js`，它会自动检测环境并选择正确的调用方式。

### Q: 数据库文件在哪里？
A: 桌面应用的数据库位置：
- macOS: `~/Library/Application Support/caipiao/lottery.db`
- Windows: `%APPDATA%\caipiao\lottery.db`

### Q: 如何调试桌面应用？
A: 使用开发模式:
```bash
wails dev
```

## 📊 性能对比

| 模式 | API调用方式 | 性能 |
|------|------------|------|
| Web模式 | HTTP (axios) | 正常 |
| 桌面模式 | Go绑定 (直接调用) | 更快 |

## 🎯 下一步

1. 更新所有Vue组件使用API适配层
2. 重新构建前端和桌面应用
3. 测试所有功能
4. 准备发布

## 🔗 相关文件

- API适配层: `frontend/src/api/adapter.js`
- Wails配置: `wails.json`
- Go应用入口: `app.go`
- 主函数: `main.go`
- GitHub Actions: `.github/workflows/build-desktop.yml`

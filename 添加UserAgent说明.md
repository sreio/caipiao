# 添加 User-Agent 防止请求被拦截

## 🛡️ 为什么需要 User-Agent？

### 问题背景

许多网站和 API 会检查 HTTP 请求的 User-Agent 头，如果：
- ❌ 没有 User-Agent
- ❌ User-Agent 是默认的（如 "Go-http-client/1.1"）
- ❌ User-Agent 不像真实浏览器

**可能会被服务器拒绝访问或限流！**

### 解决方案

通过设置真实的浏览器 User-Agent，让服务器认为请求来自正常用户，从而避免被拦截。

## ✅ 已添加的请求头

### 完整配置

```go
client := resty.New().
    SetTimeout(30 * time.Second).
    SetRetryCount(3).
    SetRetryWaitTime(1 * time.Second).
    SetRetryMaxWaitTime(5 * time.Second).
    SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
    SetHeader("Accept", "application/json, text/plain, */*").
    SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8").
    SetHeader("Accept-Encoding", "gzip, deflate, br").
    SetHeader("Connection", "keep-alive")
```

### 各请求头说明

| 请求头 | 值 | 作用 |
|--------|-----|------|
| **User-Agent** | `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36` | 模拟 Chrome 131 浏览器，macOS 系统 |
| **Accept** | `application/json, text/plain, */*` | 告诉服务器可以接受 JSON、纯文本等格式 |
| **Accept-Language** | `zh-CN,zh;q=0.9,en;q=0.8` | 优先中文内容，其次英文 |
| **Accept-Encoding** | `gzip, deflate, br` | 支持压缩传输，提高速度 |
| **Connection** | `keep-alive` | 保持 TCP 连接，提高性能 |

## 🎯 User-Agent 详解

### 我们使用的 User-Agent

```
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36
```

### 各部分含义

```
Mozilla/5.0                         # 浏览器标识（历史遗留）
(Macintosh; Intel Mac OS X 10_15_7) # 操作系统：macOS 10.15.7
AppleWebKit/537.36                  # 渲染引擎版本
(KHTML, like Gecko)                 # 兼容性标识
Chrome/131.0.0.0                    # Chrome 版本 131
Safari/537.36                       # Safari 内核版本
```

### 为什么选择这个 User-Agent？

1. **✅ 真实性** - 来自真实的 Chrome 浏览器
2. **✅ 常见性** - Chrome 是最流行的浏览器
3. **✅ 兼容性** - 几乎所有网站都支持
4. **✅ 不可疑** - 不会引起服务器警觉
5. **✅ 最新版本** - Chrome 131 是较新版本

## 📊 对比效果

### 之前（没有 User-Agent）

```http
GET /api/lottery HTTP/1.1
Host: www.cwl.gov.cn
User-Agent: Go-http-client/1.1

❌ 可能被识别为爬虫
❌ 可能被限流或拒绝
```

### 现在（完整的浏览器头）

```http
GET /api/lottery HTTP/1.1
Host: www.cwl.gov.cn
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36
Accept: application/json, text/plain, */*
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
Accept-Encoding: gzip, deflate, br
Connection: keep-alive

✅ 看起来像真实浏览器
✅ 大大降低被拦截概率
```

## 🔍 如何验证

### 方法 1：查看日志

启动后端时，如果开启调试模式，可以看到实际发送的请求头。

### 方法 2：使用测试脚本

```bash
cd backend
go run main.go
```

然后在前端点击"获取最新数据"，观察是否成功。

### 方法 3：直接测试 API

```bash
# 测试双色球 API（带 User-Agent）
curl -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36" \
  "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=1"
```

## 💡 最佳实践

### 1. 使用真实的浏览器 User-Agent

✅ **推荐：**
```
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36
```

❌ **不推荐：**
```
Go-http-client/1.1
MyApp/1.0
Python-requests/2.28.0
```

### 2. 定期更新版本号

浏览器版本会不断更新，建议每隔几个月更新一次 User-Agent 中的版本号，保持与主流版本一致。

### 3. 根据目标选择合适的 User-Agent

- **国内网站**：使用 Chrome（最常用）
- **移动端 API**：使用移动浏览器 User-Agent
- **特定要求**：按网站要求设置

### 4. 添加完整的请求头

不仅仅是 User-Agent，还应该包括：
- Accept
- Accept-Language
- Accept-Encoding
- Connection
- Referer（某些情况需要）

## 🎨 进阶技巧

### 轮换 User-Agent

如果需要更高级的反反爬虫，可以轮换不同的 User-Agent：

```go
var userAgents = []string{
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

// 随机选择
ua := userAgents[rand.Intn(len(userAgents))]
client.SetHeader("User-Agent", ua)
```

### 针对不同 API 设置不同 User-Agent

```go
// 双色球请求
resp, err := s.client.R().
    SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...").
    Get(s.shuangseqiuURL)

// 大乐透请求
resp, err := s.client.R().
    SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ...").
    Get(s.daletouURL)
```

## 📚 常见 User-Agent 参考

### Chrome（推荐）

```
# Windows
Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36

# macOS
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36

# Linux
Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36
```

### Edge

```
Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0
```

### Safari

```
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15
```

### 移动端

```
# iPhone Safari
Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1

# Android Chrome
Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36
```

## ⚠️ 注意事项

### 1. 遵守网站规则

虽然添加 User-Agent 可以避免被拦截，但仍要：
- ✅ 遵守 robots.txt
- ✅ 控制请求频率
- ✅ 尊重服务器资源
- ✅ 不进行恶意爬取

### 2. 合理使用

- ✅ 用于正常的数据查询
- ✅ 保持合理的请求间隔
- ❌ 不要频繁大量请求
- ❌ 不要绕过付费限制

### 3. 法律合规

- ✅ 使用公开的 API
- ✅ 遵守数据使用协议
- ✅ 不侵犯版权
- ✅ 保护用户隐私

## 📈 效果总结

添加 User-Agent 后：

✅ **成功率提升** - 请求被拦截的概率大幅降低  
✅ **更加稳定** - 避免被突然封禁  
✅ **更像真实用户** - 服务器无法区分  
✅ **性能优化** - 支持压缩传输  
✅ **连接复用** - Keep-Alive 提高效率  

## 🎯 总结

通过设置完整的浏览器请求头，特别是 User-Agent，我们的系统：

1. **更加可靠** - 不容易被拦截
2. **更加专业** - 遵循 HTTP 最佳实践
3. **更加高效** - 支持压缩和连接复用
4. **更加安全** - 降低被封禁风险

这是一个简单但非常重要的优化！🎉

---

**提示：** 已经配置好了，无需额外操作，直接启动即可使用！


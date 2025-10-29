# API接口说明文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: JSON

## 统一响应格式

### 成功响应
```json
{
  "code": 0,
  "msg": "success",
  "data": {...}
}
```

### 错误响应
```json
{
  "code": -1,
  "msg": "错误信息"
}
```

---

## 双色球接口

### 1. 获取双色球列表

**接口地址:** `GET /api/shuangseqiu/list`

**请求参数:**

| 参数名 | 类型 | 必填 | 说明 | 默认值 |
|--------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页条数 | 20 |
| issue | string | 否 | 期号（模糊搜索） | - |

**请求示例:**
```bash
curl "http://localhost:8080/api/shuangseqiu/list?page=1&page_size=20"
curl "http://localhost:8080/api/shuangseqiu/list?issue=2024"
```

**响应示例:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "issue": "2024001",
        "red_ball_1": 5,
        "red_ball_2": 12,
        "red_ball_3": 18,
        "red_ball_4": 23,
        "red_ball_5": 28,
        "red_ball_6": 31,
        "blue_ball": 8,
        "draw_date": "2024-01-02T00:00:00Z",
        "created_at": "2024-01-02T10:00:00Z",
        "updated_at": "2024-01-02T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 2. 获取双色球数据（从外部API）

**接口地址:** `POST /api/shuangseqiu/fetch`

**请求参数:**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| issue | string | 否 | 期号，不传则获取最新一期 |

**请求示例:**
```bash
# 获取最新一期
curl -X POST "http://localhost:8080/api/shuangseqiu/fetch"

# 获取指定期号
curl -X POST "http://localhost:8080/api/shuangseqiu/fetch?issue=2024001"
```

**响应示例:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "issue": "2024001",
    "red_ball_1": 5,
    "red_ball_2": 12,
    "red_ball_3": 18,
    "red_ball_4": 23,
    "red_ball_5": 28,
    "red_ball_6": 31,
    "blue_ball": 8,
    "draw_date": "2024-01-02T00:00:00Z",
    "created_at": "2024-01-02T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

### 3. 双色球统计

**接口地址:** `GET /api/shuangseqiu/statistics`

**请求参数:**

| 参数名 | 类型 | 必填 | 说明 | 默认值 |
|--------|------|------|------|--------|
| type | string | 否 | 统计类型：red=红球, blue=蓝球 | red |

**请求示例:**
```bash
# 红球统计
curl "http://localhost:8080/api/shuangseqiu/statistics?type=red"

# 蓝球统计
curl "http://localhost:8080/api/shuangseqiu/statistics?type=blue"
```

**响应示例:**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "number": 1,
      "count": 25
    },
    {
      "number": 2,
      "count": 18
    }
  ]
}
```

---

## 大乐透接口

### 1. 获取大乐透列表

**接口地址:** `GET /api/daletou/list`

**请求参数:**

| 参数名 | 类型 | 必填 | 说明 | 默认值 |
|--------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页条数 | 20 |
| issue | string | 否 | 期号（模糊搜索） | - |

**请求示例:**
```bash
curl "http://localhost:8080/api/daletou/list?page=1&page_size=20"
curl "http://localhost:8080/api/daletou/list?issue=2024"
```

**响应示例:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "issue": "2024001",
        "front_ball_1": 5,
        "front_ball_2": 12,
        "front_ball_3": 18,
        "front_ball_4": 23,
        "front_ball_5": 28,
        "back_ball_1": 3,
        "back_ball_2": 8,
        "draw_date": "2024-01-02T00:00:00Z",
        "created_at": "2024-01-02T10:00:00Z",
        "updated_at": "2024-01-02T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 2. 获取大乐透数据（从外部API）

**接口地址:** `POST /api/daletou/fetch`

**请求参数:**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| issue | string | 否 | 期号，不传则获取最新一期 |

**请求示例:**
```bash
# 获取最新一期
curl -X POST "http://localhost:8080/api/daletou/fetch"

# 获取指定期号
curl -X POST "http://localhost:8080/api/daletou/fetch?issue=2024001"
```

**响应示例:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "issue": "2024001",
    "front_ball_1": 5,
    "front_ball_2": 12,
    "front_ball_3": 18,
    "front_ball_4": 23,
    "front_ball_5": 28,
    "back_ball_1": 3,
    "back_ball_2": 8,
    "draw_date": "2024-01-02T00:00:00Z",
    "created_at": "2024-01-02T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

### 3. 大乐透统计

**接口地址:** `GET /api/daletou/statistics`

**请求参数:**

| 参数名 | 类型 | 必填 | 说明 | 默认值 |
|--------|------|------|------|--------|
| type | string | 否 | 统计类型：front=前区, back=后区 | front |

**请求示例:**
```bash
# 前区统计
curl "http://localhost:8080/api/daletou/statistics?type=front"

# 后区统计
curl "http://localhost:8080/api/daletou/statistics?type=back"
```

**响应示例:**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "number": 1,
      "count": 25
    },
    {
      "number": 2,
      "count": 18
    }
  ]
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| -1 | 通用错误（具体错误信息见msg字段） |

---

## 外部API对接说明

本系统已对接官方彩票数据API，配置在 `backend/config/config.go` 中：

```go
API: APIConfig{
    // 双色球官方API
    ShuangseqiuURL: "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice",
    // 大乐透官方API
    DaletouURL: "https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry",
}
```

### 双色球官方API

**接口地址：**
```
https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice
```

**请求参数：**
- `name`: 游戏名称（固定值：ssq）
- `issueCount`: 返回期数（默认：1）
- `issueStart`: 起始期号（可选）
- `issueEnd`: 结束期号（可选）

**请求示例：**
```
https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=1
```

**响应格式：**
```json
{
  "state": 0,
  "message": "成功",
  "total": 1,
  "result": [
    {
      "name": "双色球",
      "code": "2024001",
      "date": "2024-01-02",
      "week": "星期二",
      "red": "01,05,12,18,23,31",
      "blue": "08",
      "sales": "362043434",
      "poolmoney": "868123456",
      "prizegrades": [...]
    }
  ]
}
```

### 大乐透官方API

**接口地址：**
```
https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry
```

**请求参数：**
- `gameNo`: 游戏编号（固定值：85）
- `provinceId`: 省份ID（默认：0）
- `pageSize`: 每页条数（默认：1）
- `pageNo`: 页码（默认：1）
- `isVerify`: 是否验证（固定值：1）
- `issueStart`: 起始期号（可选）
- `issueEnd`: 结束期号（可选）

**请求示例：**
```
https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=1&isVerify=1&pageNo=1
```

**响应格式：**
```json
{
  "success": true,
  "errorCode": "0",
  "message": "成功",
  "value": {
    "total": 1,
    "num": 1,
    "pages": 1,
    "list": [
      {
        "lotteryDrawNum": "24001",
        "lotteryDrawTime": "2024-01-01 20:30:00",
        "lotteryDrawResult": "01 05 12 18 23 # 03 08",
        "lotterySaleAmount": "362043434",
        "lotteryPoolAmount": "868123456",
        "lotteryDrawStatus": 2
      }
    ]
  }
}
```

---

## 测试建议

### 使用Postman测试

1. 导入以下环境变量：
   - `base_url`: `http://localhost:8080`

2. 创建请求集合，按顺序测试：
   - 先测试 `fetch` 接口获取数据
   - 再测试 `list` 接口查看数据
   - 最后测试 `statistics` 接口查看统计

### 使用curl测试

```bash
# 1. 获取最新双色球数据
curl -X POST http://localhost:8080/api/shuangseqiu/fetch

# 2. 查看双色球列表
curl http://localhost:8080/api/shuangseqiu/list

# 3. 查看统计数据
curl http://localhost:8080/api/shuangseqiu/statistics?type=red
```

---

## 注意事项

1. **CORS配置**: 后端已配置CORS，允许前端（localhost:5173）跨域访问
2. **数据库**: 数据库文件位于 `backend/data/lottery.db`
3. **日志**: 开发模式下会输出详细日志
4. **性能**: 统计接口在数据量大时可能较慢，建议添加缓存
5. **数据唯一性**: 期号作为唯一索引，重复插入会失败


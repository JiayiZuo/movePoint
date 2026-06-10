# MovePoint

MovePoint 是一个攀岩训练记录与分析项目。当前仓库包含 Go + Gin 后端服务，以及根据现有接口补齐的微信小程序前端。

## 功能说明

### 用户与认证

- 邮箱密码注册、登录。
- 登录成功后保存 JWT，小程序请求会自动携带 `Authorization: Bearer <token>`。
- Token 失效或未授权时，小程序会清理登录态并回到登录页。

### 攀岩记录

- 新增、查看、编辑、删除攀岩记录。
- 支持抱石 `bouldering` 和难度攀爬 `sport_climbing`。
- 记录字段包含开始时间、结束时间、难度、线路颜色、尝试次数、是否完成、评分、地点、备注和媒体 URL。
- 后端根据开始时间与结束时间自动计算训练时长。
- 后端根据攀爬类型和时长估算热量消耗：抱石约 `12 kcal/min`，难度攀爬约 `8 kcal/min`。
- 记录列表支持分页和按日期范围过滤。

### 数据分析

- 默认展示最近 3 个月的训练分析。
- 支持自定义开始日期和结束日期。
- 分析内容包括训练次数、总时长、总热量、最高难度、难度分布、各难度成功率和月度趋势。

### 个人中心

- 查看并编辑个人资料：头像 URL、体重、身高、生日、个人简介。
- 查看用户统计：总训练次数、总时长、最高难度、本周训练次数、最近活动时间。
- 查看与刷新成就进度。

### 微信小程序页面

- `pages/login`：登录。
- `pages/register`：注册。
- `pages/dashboard`：首页概览和最近记录。
- `pages/records`：记录列表、日期过滤和分页加载。
- `pages/record-form`：新增、编辑和删除记录。
- `pages/analysis`：训练数据分析。
- `pages/profile`：个人资料、成就和退出登录。

## 项目结构

```text
.
├── controller/api/main.go          # 后端入口和路由注册
├── internal
│   ├── database/connection.go      # MySQL 连接和 GORM 自动迁移
│   ├── handlers                    # HTTP Handler
│   ├── models                      # 数据模型
│   └── services                    # 业务逻辑
├── pkg
│   ├── middleware/auth.go          # JWT 鉴权中间件
│   └── utils/jwt.go                # JWT 生成与校验
└── miniprogram                     # 微信小程序前端
    ├── app.js / app.json / app.wxss
    ├── services/api.js             # 后端 API 封装
    ├── utils/request.js            # 请求封装与鉴权处理
    └── pages                       # 小程序页面
```

## 后端接口

默认服务地址为 `http://127.0.0.1:8080`，接口前缀为 `/api`。

### 公开接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/api/register` | 用户注册 |
| `POST` | `/api/login` | 用户登录 |

### 需要登录的接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `POST` | `/api/records` | 创建攀岩记录 |
| `GET` | `/api/records?page=1&limit=20&from=2026-01-01&to=2026-03-31` | 获取记录列表 |
| `GET` | `/api/records/:id` | 获取单条记录 |
| `PUT` | `/api/records/:id` | 更新记录 |
| `DELETE` | `/api/records/:id` | 删除记录 |
| `GET` | `/api/analysis/climbing?from=2026-01-01&to=2026-03-31` | 获取攀岩分析 |
| `GET` | `/api/profile` | 获取个人资料 |
| `PUT` | `/api/profile` | 更新个人资料 |
| `GET` | `/api/profile/stats` | 获取用户统计 |
| `GET` | `/api/profile/achievements` | 获取成就 |
| `POST` | `/api/profile/check-achievements` | 检查并刷新成就 |

## 数据字段

### 注册请求

```json
{
  "username": "alex",
  "email": "alex@example.com",
  "password": "123456",
  "birth_date": "2000-01-01",
  "weight": 65,
  "height": 172,
  "avatar_url": "",
  "bio": ""
}
```

### 登录响应

```json
{
  "user_id": 1,
  "username": "alex",
  "email": "alex@example.com",
  "token": "jwt-token"
}
```

### 攀岩记录

```json
{
  "type": "bouldering",
  "start_time": "2026-06-10T18:00:00+08:00",
  "end_time": "2026-06-10T19:30:00+08:00",
  "grade": "V4",
  "color": "blue",
  "attempts": "2-3",
  "success": true,
  "rating": 4,
  "location": "本地岩馆",
  "notes": "今天核心发力不错",
  "media_urls": "[]"
}
```

字段说明：

- `type`：`bouldering` 表示抱石，`sport_climbing` 表示难度攀爬。
- `attempts`：可选 `flash`、`2-3`、`4-6`、`7+`、`failed`。
- `rating`：1 到 5。
- `duration` 和 `calories` 由后端计算，前端无需提交。

## 本地运行后端

### 1. 准备 MySQL

先确认本机 MySQL 已启动，并且 `.env` 里的账号有创建数据库权限。项目根目录已经提供 `.env.example`，本地运行可以复制为 `.env`：

```env
DB_DSN=root:123456@tcp(127.0.0.1:3306)/movepoint?charset=utf8mb4&parseTime=True&loc=Local
PORT=8080
JWT_SECRET=move-point-secret
```

如果你的 MySQL 密码不是 `123456`，请修改 `.env` 中的 `DB_DSN`。

### 2. 启动服务

```bash
go mod download
go run ./controller/api
```

启动时会先执行：

1. 读取 `.env` 中的 `DB_DSN`。
2. 连接 MySQL 服务。
3. 如果 `movepoint` 数据库不存在，则自动创建。
4. 连接 `movepoint` 数据库。
5. 执行 GORM 自动迁移。

会自动迁移以下表：

- `users`
- `climbing_records`
- `climbing_analyses`

也可以查看手动 SQL 迁移文件：

```text
migrations/001_init_schema.sql
```

## 运行微信小程序

1. 打开微信开发者工具。
2. 选择“导入项目”。
3. 项目目录选择仓库中的 `miniprogram` 目录。
4. AppID 可以先使用测试号，或保留 `project.config.json` 中的 `touristappid`。
5. 确认后端已启动在 `http://127.0.0.1:8080`。
6. 在开发者工具中开启“不校验合法域名、web-view、TLS 版本以及 HTTPS 证书”。
7. 编译运行。

如果后端地址不是本机 `8080`，修改：

```js
// miniprogram/app.js
globalData: {
  baseURL: 'http://你的服务地址/api'
}
```

## 小程序前端实现说明

- `utils/request.js` 统一处理请求、JSON Header、JWT Header、401 登录态失效。
- `services/api.js` 按后端路由封装认证、记录、分析、个人资料和成就接口。
- 记录表单会把微信 `date` / `time` picker 结果拼成 RFC3339 时间，例如 `2026-06-10T18:00:00+08:00`，便于 Go 的 `time.Time` 解析。
- 首页、记录列表、分析页和个人中心都会在未登录时自动跳转到登录页。

## 常见问题

### 小程序提示无法连接服务器

- 确认后端服务已经启动。
- 确认 `miniprogram/app.js` 中的 `baseURL` 正确。
- 微信开发者工具本地调试时需要开启“不校验合法域名”。
- 真机调试或上线时，需要使用 HTTPS 域名，并在微信公众平台配置 request 合法域名。

### 登录后仍然跳回登录页

- 检查后端是否返回了 `token` 字段。
- 检查请求头是否包含 `Authorization: Bearer <token>`。
- 检查 `JWT_SECRET` 是否在服务重启前后发生变化。

### 新增记录失败

- 确认 `start_time` 和 `end_time` 是合法时间，且结束时间晚于开始时间。
- 确认 `type` 为 `bouldering` 或 `sport_climbing`。
- 确认 `attempts` 为 `flash`、`2-3`、`4-6`、`7+` 或 `failed`。

## 后续可扩展方向

- 接入微信登录，减少邮箱密码登录步骤。
- 增加图片上传，将 `media_urls` 改为真实媒体列表。
- 完善难度比较算法，区分 V 级、YDS、法式难度等体系。
- 增加训练计划、目标管理和社区分享。
- 为后端补充单元测试和接口测试。

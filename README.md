# qqbot

轻量级 QQ 通知机器人，支持 HTTP 与 Go SDK 接入，提供同步发送、异步队列、状态查询与目标自动采集能力，内置失败重试，部署简单、使用方便。

## 1. 快速开始

- 已开通 QQ 机器人并拿到 `app_id`、`client_secret`
- 开通地址：https://q.qq.com/qqbot/openclaw/login.html

使用docker-compose部署：

```yaml
services:
  qqbot:
    image: ghcr.io/sky22333/qqbot
    container_name: qqbot
    restart: always
    ports:
      - "8080:8080"
    volumes:
      - ./configs/config.toml:/root/config.toml
      - ./data:/root/data
```

`config.toml`最小配置示例：
```toml
[qqbot]
app_id = "你的AppID"
client_secret = "你的ClientSecret"
markdown = false

[server]
listen_addr = ":8080"
api_token = "强密码token"
```

启动服务后，需要用自己的 QQ 给机器人发一次消息，让系统自动采集目标ID后才能正常推送信息。

## 2. 推送测试
```
# 异步推送
curl -X POST "http://127.0.0.1:8080/api/v1/messages" \
  -H "Authorization: Bearer 接口鉴权token" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d '{"content":"这是一条异步推送测试"}'

# 同步推送
curl -X POST "http://127.0.0.1:8080/api/v1/messages/send" \
  -H "Authorization: Bearer 接口鉴权token" \
  -H "Content-Type: application/json; charset=utf-8" \
  -d '{"content":"这是一条同步推送测试"}'
```


## 3. 常用命令

```bash
# 运行测试
go test ./...

# 静态检查
go vet ./...

# 直接运行
go run ./cmd/qqbotd -config configs/config.toml

# 构建
go build ./cmd/qqbotd

# 生产环境构建
go build -trimpath -ldflags "-s -w -buildid=" -o qqbotd ./cmd/qqbotd
```

## 4. HTTP 调用方式

除健康检查外，业务接口都需要鉴权：

```http
Authorization: Bearer 你的api_token
Content-Type: application/json
```

#### 4.1 同步发送 `POST /api/v1/messages/send`

```json
{
  "target_type": "c2c",
  "content": "这是一条通知"
}
```

`target_id` 可选；不传时自动使用最近采集目标。  
`target_type` 也可选；不传时自动使用最近采集目标的类型与 ID。

#### 4.2 异步入队 `POST /api/v1/messages`

```json
{
  "target_type": "group",
  "content": "这是一条群通知"
}
```

`target_id` 可选；不传时自动使用最近采集目标。

#### 4.3 查询状态 `GET /api/v1/messages/{request_id}`

#### 4.4 查询目标 `GET /api/v1/targets`

可选参数：`target_type=c2c|group|channel`

#### 4.5 健康检查（无需鉴权）

- `GET /healthz`
- `GET /readyz`

## 5. 目标采集

启动服务后，用自己的 QQ 给机器人发消息，系统会自动采集目标并写入 `targets.file_path` 对应的文件（默认 `data/targets.json`）。  
可通过 `GET /api/v1/targets` 查看。

## 6. 时效与可达性说明

- `target_id`（如 `user_openid`）作为发送目标标识可长期保存，不按 TTL 失效
- 服务端调用凭证 `access_token` 会过期，项目已自动刷新
- 当前项目仅发送文本或 markdown，不包含富媒体 `file_info` 链路
- 是否可送达受平台规则影响：用户关闭主动消息、频控超限都会导致发送失败
- 建议用异步接口发送，并通过 `GET /api/v1/messages/{request_id}` 跟踪状态

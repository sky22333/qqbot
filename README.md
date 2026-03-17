# qqbot

轻量级 QQ 通知机器人，支持 HTTP 服务与 Go SDK 两种使用方式。

## 1. 快速开始

- Go 1.26+
- 已开通 QQ 机器人并拿到 `app_id`、`client_secret`
- 开通地址：https://q.qq.com/qqbot/openclaw/login.html

最小配置示例：
```toml
[qqbot]
app_id = "你的AppID"
client_secret = "你的ClientSecret"
markdown = false

[server]
listen_addr = ":8080"
api_token = "强密码token"
```

## 2. HTTP 调用方式

除健康检查外，业务接口都需要鉴权：

```http
Authorization: Bearer 你的api_token
Content-Type: application/json
```

### 2.1 同步发送 `POST /api/v1/messages/send`

```json
{
  "target_type": "c2c",
  "content": "这是一条通知"
}
```

`target_id` 可选；不传时自动使用最近采集目标。  
`target_type` 也可选；不传时自动使用最近采集目标的类型与 ID。

### 2.2 异步入队 `POST /api/v1/messages`

```json
{
  "target_type": "group",
  "content": "这是一条群通知"
}
```

`target_id` 可选；不传时自动使用最近采集目标。

### 2.3 查询状态 `GET /api/v1/messages/{request_id}`

### 2.4 查询目标 `GET /api/v1/targets`

可选参数：`target_type=c2c|group|channel`

### 2.5 健康检查（无需鉴权）

- `GET /healthz`
- `GET /readyz`

## 3. SDK 调用方式

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/sky22333/qqbot"
)

func main() {
	client, err := qqbot.NewFromConfigFile("configs/config.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err = client.Send(ctx, qqbot.PushRequest{
		TargetType: qqbot.TargetC2C,
		Content:    "SDK通知消息",
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

`target_id` 可选；不传时 SDK 会自动按 `target_type` 回填最近采集目标。  
`target_type` 与 `target_id` 都不传时，SDK 会使用最近一次采集到的目标类型与 ID。

如需在短生命周期任务中关闭采集器，可使用：

```go
client, err := qqbot.NewWithOptions(cfg, qqbot.ClientOptions{
	StartCollector: false,
})
```

## 4. 目标采集

启动服务后，用自己的 QQ 给机器人发消息，系统会自动采集目标并写入 `targets.file_path` 对应的文件（默认 `data/targets.json`）。  
可通过 `GET /api/v1/targets` 查看。

## 5. 常用命令

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

## 推送测试
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

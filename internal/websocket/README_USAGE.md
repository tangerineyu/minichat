# websocket 使用说明（配合 internal/websocket/hub.go）

这份文档是对 `internal/websocket/hub.go` 的“使用侧”补充说明：

- `Hub` 负责维护 *在线连接表*（`map[userID]*Client`）
- `Client` 对应一个 websocket 连接
- 业务侧只需要调用：`websocket.GlobalHub.PushMessage(userID, payload)`

> 约定：客户端收到的 payload 是服务端写入的 `[]byte`。本项目 `message_service` 目前是 `json.Marshal(msg)` 后推送。

---

## 1. 在 main 启动 Hub

因为 `Hub` 通过 `Register/Unregister` channel 串行消费事件，所以必须先启动事件循环：

- 在 `cmd/main.go`（或你的启动入口）里：
  - `go websocket.GlobalHub.Run()`

常见坑：

- 如果忘了启动 `Run()`，你往 `GlobalHub.Register <- client` 写入会阻塞（没人消费）。

---

## 2. 在 HTTP handler 里升级 websocket，并注册 Client

你需要一个“建连接口”（通常是 `GET /ws`），整体流程像这样：

1) **鉴权**：从 token / session 中拿到当前用户 `userID`
2) **升级协议**：`Upgrader.Upgrade(w, r, nil)` 得到 `*websocket.Conn`
3) **创建 Client**：设置 `Send` 为 *带缓冲* 的 channel
4) **注册到 Hub**：`websocket.GlobalHub.Register <- client`
5) **启动两个 goroutine**：
   - `go client.WritePump()`：专门负责写（包含 ping）
   - `go client.ReadPump()`：专门负责读（包含 pong 续期 + 感知断开）

伪代码示意：

- `Send` 缓冲建议先用 256 或 1024（取决于你的推送频率和消息大小）
- `ReadPump()` 退出时会自动 `Unregister`，并关闭 `Send`，从而让 `WritePump()` 也退出

---

## 3. 业务侧推送消息（SendMessage）

当前项目里：`internal/service/message/message_service.go` 的 `SendMessage` 里已经做了：

- 落库成功后：`json.Marshal(msg)`
- 然后 `go websocket.GlobalHub.PushMessage(receiverId, payload)`

注意点：

- `PushMessage` 是“尽力而为”：
  - 用户不在线：忽略
  - 用户在线但 `Send` 队列满：触发 `Unregister` 主动断开（避免队列无限堆积）

---

## 4. 客户端该怎么配合（心跳）

服务端 `WritePump` 会定期发送 `Ping`。

- 大多数 websocket 客户端库会自动响应 `Pong`
- 如果你的客户端库不自动回 `Pong`，需要你显式处理 `Ping` 并回复 `Pong`

服务端 `ReadPump` 依赖 `Pong` 来续期 `ReadDeadline`，否则会在 `pongWait` 超时后断开。

---

## 5. 常见扩展点

### 5.1 支持多端同时在线（同一个 userID 多连接）

当前结构：`map[int64]*Client` 只会保留一个连接。

要支持多端：

- 改成：`map[int64]map[string]*Client` 或 `map[int64][]*Client`
- `PushMessage` 需要给该 userID 的所有 client.Send 投递

### 5.2 服务端接收客户端消息

现在 `ReadPump` 只是“读掉但不处理”。

要支持客户端发消息到服务端：

- 在 `ReadMessage()` 后拿到 payload
- 解析 JSON/protobuf
- 做鉴权/校验/限流
- 调用 service 落库/转发

---

## 6. 与 Gin 集成时的几个注意点

- `Upgrade` 需要 `http.ResponseWriter` 和 `*http.Request`：
  - Gin 里分别是 `c.Writer` 和 `c.Request`
- 需要放开跨域 / Origin 校验：
  - `websocket.Upgrader.CheckOrigin` 默认可能拒绝跨域
  - 生产环境建议只允许你的前端域名


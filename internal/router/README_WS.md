# WebSocket 聊天（HTTP 查历史 + WS 实时收发）

## 1. 路由

- WebSocket 连接：`GET /ws/chat`
  - 鉴权：
    - Header：`Authorization: Bearer <access_token>`（推荐）
    - 或 Query：`/ws/chat?token=<access_token>`（浏览器方便）

- HTTP 历史消息：`GET /api/message/history`
  - 需要登录（JWT 中间件）
  - 参数：
    - `target_id`：私聊对方 userId 或群 id
    - `session_type`：1 私聊，2 群聊
    - `cursor`：滚动分页游标（上一页返回的 next_cursor）
    - `limit`：返回条数（默认 20，最大 100）

> 现状：HTTP 的 `/api/message/send` 仍保留（便于调试），但推荐用 WS 发送以获得实时性。

## 2. WS 协议

### 2.1 客户端 -> 服务端：发送消息

```json
{
  "type": "send_message",
  "data": {
    "receiver_id": 2,
    "session_type": 1,
    "msg_type": 1,
    "content": "hello"
  }
}
```

### 2.2 服务端 -> 客户端：ack

```json
{
  "type": "ack",
  "data": {
    "message": {"id": 123, "sender_id": 1, "receiver_id": 2, "content": "hello"}
  }
}
```

### 2.3 服务端 -> 客户端：error

```json
{
  "type": "error",
  "data": {"code": 500, "msg": "你们不是好友关系，无法发送消息"}
}
```

### 2.4 服务端推送新消息

`MessageService.SendMessage` 落库后会通过 `websocket.GlobalHub.PushMessage(receiverId, payload)` 推送给在线的接收方。

推送 payload 当前就是 `model.Message` 的 JSON。

## 3. 业务规则

- 私聊（`session_type = 1`）：必须是好友（friend.status == 1），否则发送失败。
- 群聊（`session_type = 2`）：目前不做成员校验（后续可以补：必须是群成员才能发）。

## 4. 最快测试方式

1) 用 HTTP 登录接口拿到 access token。

2) 打开两个客户端分别用两个 token 建立 WS：
- A：`ws://<host>/ws/chat?token=<tokenA>`
- B：`ws://<host>/ws/chat?token=<tokenB>`

3) A 发送 `send_message`，B 若在线会收到推送。

4) 用 HTTP `/api/message/history` 验证落库与滚动分页。


package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"minichat/internal/handler/response"
	"minichat/internal/req"
	messageService "minichat/internal/service/message"
	appws "minichat/internal/websocket"
	jwt_util "minichat/util/jwt"

	"github.com/gin-gonic/gin"
	gws "github.com/gorilla/websocket"
)

// WSHandler 负责 websocket 入口（HTTP -> WS upgrade）以及 WS 消息的协议适配。
//
// 分层约定：
//   - handler：鉴权、解析客户端帧(JSON)、调用 service、回写 ack/错误
//   - service：业务规则（比如“不是好友不能发消息”）+ 落库 + 推送
//   - websocket(Hub/Client)：连接生命周期与推送通道
//
// 注意：Conn 不能并发写。这个项目通过 Client.WritePump 串行写入。WSHandler 不直接写 Conn。
// （除非是 handshake 失败/鉴权失败时返回 HTTP JSON）
type WSHandler struct {
	messageService messageService.MessageServiceInterface
}

var _ WSHandlerInterface = (*WSHandler)(nil)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域（开发方便）。生产建议按 Origin 白名单严格校验。
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Connect 建立 websocket 连接。
//
// 鉴权策略：
//  1. 优先读取 Header: Authorization: Bearer <token>
//  2. 其次读取 query: ?token=<access_token>
//
// 成功后：
//   - 创建 websocket.Client（含 Send channel）
//   - 注册到 GlobalHub
//   - 启动 client.WritePump / ReadPump
//   - 同时启动一个业务 ReadLoop：用于读取客户端业务帧（send_message）并调用 service
func (h *WSHandler) Connect(c *gin.Context) {
	userID, err := authUserID(c)
	if err != nil {
		// WS 握手前还在 HTTP 阶段，可以正常回 JSON。
		response.Fail(c, 401, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &appws.Client{
		Hub:    appws.GlobalHub,
		Conn:   conn,
		UserID: userID,
		// 给一个适度缓冲，避免短时间 burst push 阻塞业务协程。
		Send: make(chan []byte, 256),
	}

	// 注册到 hub。
	appws.GlobalHub.Register <- client

	// 所有写都走 Client.WritePump。
	go client.WritePump()

	// 注意：不能有两个 goroutine 同时 Read 同一个 Conn。
	// 所以我们在这里用 runReadLoop 统一承担：
	//   - 业务帧解析（send_message）
	//   - 心跳（pong handler + read deadline）
	//   - 断开时注销
	go h.runReadLoop(client)
}

// authUserID 从请求中解析 access token 并返回 userID。
func authUserID(c *gin.Context) (int64, error) {
	// header: Authorization: Bearer xxx
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && strings.TrimSpace(parts[1]) != "" {
			claims, err := jwt_util.ValidateAccessToken(strings.TrimSpace(parts[1]))
			if err == nil {
				return claims.Id, nil
			}
		}
	}

	// query: /ws/chat?token=xxx
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		return 0, errors.New("missing token")
	}
	claims, err := jwt_util.ValidateAccessToken(token)
	if err != nil {
		return 0, errors.New("invalid or expired token")
	}
	return claims.Id, nil
}

// --- WS 业务协议 ---

// 入站帧统一格式：
// {"type":"send_message","data":{...SendMessageReq...}}
//
// 你也可以扩展：ping、typing、read_receipt 等。

type inboundFrame struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// 出站 ack 格式：
// {"type":"ack","data":{"message":{...}}}
// {"type":"error","data":{"code":400,"msg":"..."}}

type outboundFrame struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// writeJSON 只负责把 payload 写入 client.Send，真正写 socket 由 Client.WritePump 完成。
func writeJSON(client *appws.Client, v any) {
	b, _ := json.Marshal(v)
	select {
	case client.Send <- b:
	default:
		// 队列满了，交给 hub 清理（避免阻塞）。
		go func() { client.Hub.Unregister <- client }()
	}
}

// runReadLoop 负责：
//   - 读取客户端帧
//   - 维护 read deadline + pong handler（心跳）
//   - 处理 send_message：调用 service.SendMessage，service 里会做好友校验+落库+推送
func (h *WSHandler) runReadLoop(client *appws.Client) {
	defer func() {
		client.Hub.Unregister <- client
		_ = client.Conn.Close()
	}()

	client.Conn.SetReadLimit(8 * 1024)
	_ = client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		_ = client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, payload, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}

		var f inboundFrame
		if err := json.Unmarshal(payload, &f); err != nil {
			writeJSON(client, outboundFrame{Type: "error", Data: gin.H{"code": 400, "msg": "invalid frame"}})
			continue
		}

		switch f.Type {
		case "send_message":
			var in req.SendMessageReq
			if err := json.Unmarshal(f.Data, &in); err != nil {
				writeJSON(client, outboundFrame{Type: "error", Data: gin.H{"code": 400, "msg": "invalid data"}})
				continue
			}
			msg, err := h.messageService.SendMessage(
				contextBackground(),
				client.UserID,
				in.ReceiverId,
				in.SessionType,
				in.MsgType,
				in.Content,
			)
			if err != nil {
				writeJSON(client, outboundFrame{Type: "error", Data: gin.H{"code": 500, "msg": err.Error()}})
				continue
			}
			writeJSON(client, outboundFrame{Type: "ack", Data: gin.H{"message": msg}})

		default:
			writeJSON(client, outboundFrame{Type: "error", Data: gin.H{"code": 400, "msg": "unknown type"}})
		}
	}
}

// contextBackground 抽出来避免在 WS 层引入 gin.Context。
// 这里用 Background，未来你可以改成带 trace/span 的 ctx。
func contextBackground() context.Context {
	return context.Background()
}

func NewWSHandler(messageService messageService.MessageServiceInterface) WSHandlerInterface {
	return &WSHandler{messageService: messageService}
}

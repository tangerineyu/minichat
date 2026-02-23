package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait: 每次写 websocket 帧允许的最大耗时。
	// 如果网络堵塞/对端不读导致写阻塞，超过该时间就会失败并触发连接关闭。
	writeWait = 10 * time.Second

	// pongWait: 允许客户端多久不回复 Pong。
	// 我们会周期性发送 Ping，如果在 pongWait 内收不到 Pong，则认为连接不可用。
	pongWait = 60 * time.Second

	// pingPeriod: 服务端发送 Ping 的周期。
	// 通常要小于 pongWait，这样才能在超时之前发出 ping 并等待 pong。
	pingPeriod = (pongWait * 9) / 10
)

// Client 表示一个在线连接。
//
// 关键点：
//   - Conn:   底层 websocket 连接。
//   - UserID: 该连接所属用户（用来定位推送目标）。
//   - Send:   “服务端 -> 客户端”的消息队列。业务侧 push 时只往这里塞 []byte。
//     WritePump 会从该 channel 读取并真正写入 Conn。
//
// 注意：Send 建议使用有缓冲 channel（在创建 client 时设置），否则高并发推送时更容易阻塞。
// 这里的文件只定义结构体，不负责具体创建逻辑。
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID int64
	Send   chan []byte
}

// WritePump 是“写协程”。
//
// 它只做两件事：
//  1. 从 c.Send 收到应用层消息后，写入 websocket
//  2. 定时发送 Ping，维持心跳；若写失败，则退出并关闭连接
//
// 为什么要单独写协程？
//   - gorilla/websocket 明确要求：同一个 Conn 不能并发写。
//   - 所以所有写操作都集中在这里，业务侧只能向 c.Send 投递。
func (c *Client) WritePump() {
	// ticker 驱动定时 ping。
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// defer 保证退出时资源正确释放。
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// 每次写之前设置截止时间，避免对端不读导致写永久阻塞。
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// Send channel 被关闭意味着 Hub 已经把该 client 清理掉了。
				// 发送一个 close frame 给对端，随后退出。
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 这里默认把 payload 当作 text message。
			// 如果你传的是二进制（protobuf/图片），改成 websocket.BinaryMessage。
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				// 写失败通常意味着连接已断开或网络不可用，直接退出触发清理。
				return
			}

		case <-ticker.C:
			// 定时发 ping，让中间层（NAT/LB）不把长连接当作“空闲”断开，
			// 同时也能更快发现死连接。
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump 是“读协程”。
//
// 这个项目里 ReadPump 目前不处理客户端发来的业务消息（只读掉即可），
// 主要用途是：
//   - 读取控制帧（尤其是 Pong），配合 SetPongHandler 维护 read deadline
//   - 一旦 ReadMessage 返回错误（对端断开/超时），触发注销和资源清理
//
// 如果你未来想支持“客户端发消息到服务端”，就在 ReadMessage 后解析 payload，
// 然后调用业务 service（注意鉴权、参数校验、限流、落库）。
func (c *Client) ReadPump() {
	defer func() {
		// 任何退出路径都会走到这里：通知 Hub 注销 + 关闭连接。
		// 注销会关闭 Send channel，从而让 WritePump 也退出。
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	// 控制客户端单条消息大小上限，防止恶意大包占用内存。
	c.Conn.SetReadLimit(512)

	// 设置首个读超时截止时间。
	// 只要后续不断收到 Pong，就会在 PongHandler 里续期。
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// 当收到客户端的 Pong 时，延长 read deadline。
	// 这就是“心跳续命”的关键：只要对端还活着并回 Pong，连接就不会因为超时被我们断开。
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// ReadMessage 会阻塞直到：
		//   - 收到下一帧
		//   - 发生错误（包括超时、对端 close 等）
		// 在这里我们不使用消息内容，所以直接丢弃。
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

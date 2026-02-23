package websocket

import (
	"sync"
)

// 这份实现采用了 gorilla/websocket 推荐的经典模型：
//
//   - 每个连接（Client）各自跑两个 goroutine：
//       1) ReadPump:  负责“读”客户端发来的控制帧/数据帧（这里我们主要用来维护心跳、感知断开）
//       2) WritePump: 负责“写”服务端要推送给客户端的消息（从 c.Send channel 取数据写入 socket）
//
//   - Hub 是一个“连接注册中心”，维护 userID -> *Client 映射。
//     业务侧（比如 message_service）只需要调用 Hub.PushMessage(userID, payload)
//     就能把消息推给在线用户；不在线则忽略（你也可以在这里扩展为落库/离线消息）。
//
// 并发与线程安全约定：
//   - 一个 websocket.Conn 不允许并发写。我们通过“只有 WritePump 会写 conn”来保证这一点。
//   - 业务侧不会直接写 conn，而是写入 client.Send（一个 channel）。
//   - Hub.Clients map 由 Hub 自己维护，读写都受互斥锁保护。

// Hub 负责管理所有在线连接。
//
// Clients:
//   - key: userID
//   - val: 对应在线连接（一个 userID 同时只保留一个连接；如果你允许多端登录，
//     可改为 map[int64]map[connID]*Client 或 map[int64][]*Client）
//
// Register/Unregister:
//   - 用 channel 做串行化“注册/注销”事件，避免业务侧直接操作 Clients map。
//   - 但注意：本实现中 PushMessage 仍会直接读/写 Clients map，因此仍使用 mu 保护。
//
// Broadcast:
//   - 用于未来扩展“广播给所有在线用户”的能力。
//   - 当前 Run() 没处理该 channel（如果要用，需要在 Run 里加一个 case）。
type Hub struct {
	Clients    map[int64]*Client
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
	Broadcast  chan []byte
}

// GlobalHub 是一个全局单例 Hub。
//
// 使用方式：
//   - main 启动时 go GlobalHub.Run()
//   - websocket 建连时：GlobalHub.Register <- client
//   - 断开时：GlobalHub.Unregister <- client（ReadPump defer 会做）
//   - 推送时：GlobalHub.PushMessage(userID, payload)
var GlobalHub = &Hub{
	Clients:    make(map[int64]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	// Broadcast 当前没用到，但初始化出来能避免后续扩展时忘记 make。
	Broadcast: make(chan []byte),
}

// Run 是 Hub 的“事件循环”。
//
// 典型用法：在 main 初始化时启动：
//
//	go websocket.GlobalHub.Run()
//
// 它会持续处理注册/注销事件并维护 Clients 映射。
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// 注册：把连接放入在线表。
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			// 注销：把连接从在线表移除，并关闭 Send channel。
			// 关闭 Send 的作用：让 WritePump 读到 ok=false 并退出。
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

// PushMessage 向指定 userId 推送消息。
//
// 行为约定：
//   - 用户在线：尝试写入 client.Send
//   - 用户不在线：什么都不做（可扩展为离线消息）
//   - 用户在线但 Send 队列满了：认为该连接“消费不过来”，触发注销（防止内存堆积）
//
// 为什么 default 分支要 Unregister？
//   - 如果 Send 是有缓冲 channel，写满说明对端写不出去（网络慢/客户端卡死）。
//   - 持续堆积会占用内存，所以选择主动断开，让客户端重连。
func (h *Hub) PushMessage(userId int64, message []byte) {
	// 这里用写锁是因为我们要读取 Clients map。
	// 其实只读的话可以用 RLock，但考虑到后续可能会在这里做更多维护逻辑，
	// 目前保持简单。若要优化可改为 RLock。
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.Clients[userId]; ok {
		select {
		case client.Send <- message:
			// 投递成功。
		default:
			// 投递失败（队列满/无人接收）。异步触发注销，避免在锁内阻塞。
			go func() {
				h.Unregister <- client
			}()
		}
	}
}

package websocket

import (
	"encoding/json"
	"time"
)

type WsPushMsg struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
	Ts    int64  `json:"ts"`
}

func BuildPushBytes(event string, data any) []byte {
	msg := WsPushMsg{
		Event: event,
		Data:  data,
		Ts:    time.Now().UnixMilli(),
	}
	bytes, _ := json.Marshal(msg)
	return bytes
}

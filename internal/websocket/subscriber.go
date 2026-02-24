package websocket

import (
	"context"
	"fmt"

	"minichat/internal/event"
)

func InitSubscriber() {
	wsHandler := func(ctx context.Context, payload any) error {
		env, ok := payload.(event.WsEnvelope)
		if !ok {
			return fmt.Errorf("websocket subscriber: invalid payload type %T", payload)
		}

		wsBytes := BuildPushBytes(env.Event, env.Data)

		// 1) 批量推送优先（群聊通常走这里）
		if len(env.TargetUserIDs) > 0 {
			GlobalHub.PushMessages(env.TargetUserIDs, wsBytes)
			return nil
		}
		if len(env.GroupMemberIDs) > 0 {
			GlobalHub.PushMessages(env.GroupMemberIDs, wsBytes)
			return nil
		}

		// 2) 兼容旧逻辑：只填 TargetID -> 默认当作 user 推送
		GlobalHub.PushMessage(env.TargetID, wsBytes)
		return nil
	}

	event.GlobalBus.Subscribe(event.EventMessageNew, wsHandler)
	event.GlobalBus.Subscribe(event.EventMessageWithdraw, wsHandler)
	event.GlobalBus.Subscribe(event.EventFriendApplyNew, wsHandler)
	// 你后续如果要把“入群申请/审核/踢人/解散/系统踢下线”等也统一走 ws 推送，
	// 可以继续在这里 Subscribe 对应事件并复用同一个 wsHandler。
}

package event

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type EventType string

type EventHandler func(ctx context.Context, payload any) error

type EventBus struct {
	mu   sync.RWMutex
	subs map[EventType][]EventHandler
}

var GlobalBus = &EventBus{
	subs: make(map[EventType][]EventHandler),
}

func (b *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[eventType] = append(b.subs[eventType], handler)
}

func (b *EventBus) Publish(ctx context.Context, eventType EventType, payload any) {
	b.mu.RLock()
	handlers, ok := b.subs[eventType]
	b.mu.RUnlock()

	if !ok || len(handlers) == 0 {
		return // 没有订阅者，不算错误
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error("event Panic", zap.String("eventType", string(eventType)), zap.Any("payload", payload))
				}
			}()

			if err := h(ctx, payload); err != nil {
				zap.L().Error("事件处理失败", zap.Error(err), zap.String("eventType", string(eventType)))
			}
		}(handler)
	}
}

package xcontext

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// NewDetachedContext создаёт новый контекст, отвязанный от переданного, но имеющий доступ к его значениям.
func NewDetachedContext(ctx context.Context) context.Context {
	return WithContextValues(opentracing.ContextWithSpan(context.Background(), opentracing.SpanFromContext(ctx)), ctx)
}

// WithContextValues создаёт потомка контекста ctx, который будет иметь доступ к Values обоих контекстов.
// Сначала значение ищется в ctx, затем, если не найдено, в ctxWithValues. При этом не используется таймаут контекста ctxWithValues.
// Функция полезна, когда нужно запустить в фоне процесс, которому нужна информация из оригинального контекста (например, метаданные запроса).
func WithContextValues(ctx, ctxWithValues context.Context) context.Context {
	return &composed{Context: ctx, ctxValues: ctxWithValues}
}

type composed struct {
	context.Context
	ctxValues context.Context
}

func (c *composed) Value(key interface{}) interface{} {
	if v := c.Context.Value(key); v != nil {
		return v
	}
	return c.ctxValues.Value(key)
}

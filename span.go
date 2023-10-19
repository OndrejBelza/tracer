package tracer

import "context"

type Span interface {
	SetName(string)
	Finish()
	// Tag sets Tag with given key and value to the Span. If key already exists in
	// the Span the value will be overridden except for error tags where the first
	// value is persisted.
	Tag(key, value string)
}

type ctxKey struct{}

var spanKey = ctxKey{}

// NewContext stores a Span into Go's context propagation mechanism.
func NewContext(ctx context.Context, s Span) context.Context {
	return context.WithValue(ctx, spanKey, s)
}

// SpanFromContext retrieves a Span from Go's context propagation
// mechanism if found. If not found, returns nil.
func SpanFromContext(ctx context.Context) Span {
	if s, ok := ctx.Value(spanKey).(Span); ok {
		return s
	}
	return nil
}

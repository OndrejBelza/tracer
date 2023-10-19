package tracer

import "context"

type TrafficTracer interface {
	StartSpan(name string) Span
	StartSpanFromContext(ctx context.Context, name string) (Span, context.Context)
}

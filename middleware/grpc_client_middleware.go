package middleware

import (
	"context"

	"github.com/OndrejBelza/tracer"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

type GrpcClientZipkinHandler struct {
	tracer            tracer.TrafficTracer
	remoteServiceName string
}

func NewGrpcClientZipkinHandler(tracer tracer.TrafficTracer) stats.Handler {
	return &GrpcClientZipkinHandler{
		tracer: tracer,
	}
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (c *GrpcClientZipkinHandler) HandleConn(_ context.Context, _ stats.ConnStats) {
	// no-op
}

// TagConn exists to satisfy gRPC stats.Handler.
func (c *GrpcClientZipkinHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c *GrpcClientZipkinHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	handleRPC(ctx, rs)
}

// TagRPC implements per-RPC context management.
func (c *GrpcClientZipkinHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	switch t := c.tracer.(type) {
	case *tracer.ZipkinTracer:
		var span *tracer.ZipkinSpan

		ep := remoteEndpointFromContext(ctx, c.remoteServiceName)

		name := spanName(rti)
		span, ctx = t.StartSpanFromContextOptions(ctx, name, zipkin.Kind(model.Client), zipkin.RemoteEndpoint(ep))

		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.New(nil)
		}
		_ = b3.InjectGRPC(&md)(span.Context())

		// inject baggage fields from span context into the outgoing gRPC request metadata
		if span.Context().Baggage != nil {
			span.Context().Baggage.Iterate(func(key string, values []string) {
				md.Set(key, values...)
			})
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return ctx
}

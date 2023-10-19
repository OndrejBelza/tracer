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

type GrpcServerZipkinHandler struct {
	tracer tracer.TrafficTracer
}

func NewGrpcServerZipkinHandler(tracer tracer.TrafficTracer) stats.Handler {
	return &GrpcServerZipkinHandler{
		tracer: tracer,
	}
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (s *GrpcServerZipkinHandler) HandleConn(_ context.Context, _ stats.ConnStats) {
	// no-op
}

// TagConn exists to satisfy gRPC stats.Handler.
func (s *GrpcServerZipkinHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (s *GrpcServerZipkinHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	handleRPC(ctx, rs)
}

// TagRPC implements per-RPC context management.
func (s *GrpcServerZipkinHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	// In practice, ok never seems to be false but add a defensive check.
	if !ok {
		md = metadata.New(nil)
	}

	switch t := s.tracer.(type) {
	case *tracer.ZipkinTracer:
		name := spanName(rti)

		spanContext := t.Extract(b3.ExtractGRPC(&md))

		span := t.StartSpanOptions(
			name,
			zipkin.Kind(model.Server),
			zipkin.Parent(spanContext),
			zipkin.RemoteEndpoint(remoteEndpointFromContext(ctx, "")),
		)

		return tracer.NewContext(ctx, span)
	}

	return ctx

}

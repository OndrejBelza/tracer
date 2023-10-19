package middleware

import (
	"context"
	"strings"

	"github.com/OndrejBelza/tracer"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

func handleRPC(ctx context.Context, rs stats.RPCStats) {
	span := tracer.SpanFromContext(ctx)

	switch rs := rs.(type) {
	case *stats.End:
		s, ok := status.FromError(rs.Error)
		// rs.Error should always be convertable to a status, this is just a defensive check.
		if ok {
			if s.Code() != codes.OK {
				// Uppercase for consistency with Brave
				c := strings.ToUpper(s.Code().String())
				span.Tag(string(tracer.TagGRPCStatusCode), c)
				span.Tag(string(tracer.TagError), c)
			}
		} else {
			span.Tag(string(tracer.TagError), rs.Error.Error())
		}
		span.Finish()
	}
}

func spanName(rti *stats.RPCTagInfo) string {
	name := strings.TrimPrefix(rti.FullMethodName, "/")
	name = strings.Replace(name, "/", ".", -1)
	nameParts := strings.Split(name, ".")
	name = strings.Join(nameParts[len(nameParts)-1:], ".")
	return name
}

func remoteEndpointFromContext(ctx context.Context, name string) *model.Endpoint {
	remoteAddr := ""

	p, ok := peer.FromContext(ctx)
	if ok {
		remoteAddr = p.Addr.String()
	}

	ep, _ := zipkin.NewEndpoint(name, remoteAddr)
	return ep
}

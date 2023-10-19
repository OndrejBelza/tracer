package tracer

import (
	"context"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

type ZipkinTracer struct {
	endpointURL string
	serviceName string
	portHost    string
	tracer      *zipkin.Tracer
}

func NewZipkinTracer(endpointURL, serviceName, portHost string) (TrafficTracer, error) {

	reporter := reporterhttp.NewReporter(endpointURL)
	endpoint, err := zipkin.NewEndpoint(serviceName, portHost)
	if err != nil {
		return nil, err
	}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	tracer, err := zipkin.NewTracer(reporter, zipkin.WithSampler(sampler), zipkin.WithLocalEndpoint(endpoint))

	if err != nil {
		return nil, err
	}

	return &ZipkinTracer{
		endpointURL: endpointURL,
		serviceName: serviceName,
		portHost:    portHost,
		tracer:      tracer,
	}, nil

}

func (t *ZipkinTracer) StartSpan(name string) Span {
	return t.StartSpanOptions(name)
}

func (t *ZipkinTracer) StartSpanOptions(name string, options ...zipkin.SpanOption) Span {
	return newZipkinSpan(t, name, options...)
}

func (t *ZipkinTracer) StartSpanFromContext(ctx context.Context, name string) (Span, context.Context) {
	return t.StartSpanFromContextOptions(ctx, name)
}

func (t *ZipkinTracer) StartSpanFromContextOptions(ctx context.Context, name string, options ...zipkin.SpanOption) (*ZipkinSpan, context.Context) {
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		if s, ok := parentSpan.(*ZipkinSpan); ok {
			options = append(options, zipkin.Parent(s.zipkinSpan.Context()))
		}
	}
	span := newZipkinSpan(t, name, options...)
	return span, NewContext(ctx, span)
}

func (t *ZipkinTracer) Extract(extractor propagation.Extractor) (sc model.SpanContext) {
	psc, err := extractor()
	if psc != nil {
		sc = *psc
	}
	sc.Err = err
	return
}

var _ TrafficTracer = (*ZipkinTracer)(nil)

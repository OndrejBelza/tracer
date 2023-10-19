package tracer

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
)

type ZipkinSpan struct {
	zipkinSpan zipkin.Span
}

func newZipkinSpan(t *ZipkinTracer, name string, options ...zipkin.SpanOption) *ZipkinSpan {
	return &ZipkinSpan{
		zipkinSpan: t.tracer.StartSpan(name, options...),
	}
}

func (s *ZipkinSpan) Tag(key, value string) {
	s.zipkinSpan.Tag(key, value)
}

func (s *ZipkinSpan) SetName(name string) {
	s.zipkinSpan.SetName(name)
}

func (s *ZipkinSpan) Finish() {
	s.zipkinSpan.Finish()
}

func (s *ZipkinSpan) Context() model.SpanContext {
	return s.zipkinSpan.Context()
}

func (s *ZipkinSpan) TraceID() string {
	return s.zipkinSpan.Context().TraceID.String()
}

var _ Span = (*ZipkinSpan)(nil)

package util

import (
	"encoding/json"

	"github.com/ShahoBashoki/kucoin/object"
	"go.opentelemetry.io/otel/trace"
)

type spanContext struct {
	trace.SpanContext
}

var (
	_ json.Marshaler = (*spanContext)(nil)
	_ object.GetMap  = (*spanContext)(nil)
)

// NewSpanContext is a function.
func NewSpanContext(
	traceSpan trace.Span,
) *spanContext {
	if traceSpan == nil {
		return &spanContext{
			SpanContext: trace.SpanContext{},
		}
	}

	return &spanContext{
		SpanContext: traceSpan.SpanContext(),
	}
}

// GetMap is a function.
func (spanContext *spanContext) GetMap() map[string]any {
	return map[string]any{
		"trace_id":    spanContext.TraceID(),
		"span_id":     spanContext.SpanID(),
		"trace_flags": spanContext.TraceFlags(),
		"trace_state": spanContext.TraceState(),
		"remote":      spanContext.IsRemote(),
	}
}

// MarshalJSON is a function.
// read more https://pkg.go.dev/encoding/json#Marshaler
func (spanContext *spanContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(spanContext.GetMap())
}

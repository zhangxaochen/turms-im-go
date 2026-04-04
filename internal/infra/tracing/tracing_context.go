package tracing

import (
	"context"
)

type contextKey struct{}

var tracingContextKey = contextKey{}

// TracingContext represents a tracing scope.
type TracingContext struct {
	TraceId string
}

// NewTracingContext creates a new TracingContext.
func NewTracingContext(traceId string) *TracingContext {
	return &TracingContext{TraceId: traceId}
}

// FromContext extracts a TracingContext from a Context if present.
func FromContext(ctx context.Context) *TracingContext {
	tc, ok := ctx.Value(tracingContextKey).(*TracingContext)
	if ok {
		return tc
	}
	return nil
}

// WithTracingContext returns a copy of parent context with the TracingContext.
func WithTracingContext(ctx context.Context, tc *TracingContext) context.Context {
	return context.WithValue(ctx, tracingContextKey, tc)
}

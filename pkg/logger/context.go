package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// key is an unexported type to prevent collisions in context
type key int

const (
	loggerKey key = iota
	traceIDKey
)

// WithContext returns a new context with the logger embedded
func WithContext(ctx context.Context, l *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext retrieves the logger from context.
// If missing, it returns the Global Logger (safe fallback).
func From(ctx context.Context) *zerolog.Logger {
	if l, ok := ctx.Value(loggerKey).(*zerolog.Logger); ok {
		return l
	}
	return &log.Logger
}

// WithTraceID injects just the ID (useful if you need it separately)
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// TraceIDFrom retrieves the ID
func TraceIDFrom(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return ""
}

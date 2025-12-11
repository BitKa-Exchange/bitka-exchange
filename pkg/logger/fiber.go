package logger

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Try to get existing ID (from Traefik), otherwise generate new one
		traceID := c.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		// Set header for client/downstream services
		c.Set("X-Request-ID", traceID)

		// Create Context-Aware Logger
		reqLogger := log.With().Str("trace_id", traceID).Logger()

		// INJECT INTO GO CONTEXT
		// Create a new context derived from the request's context
		ctx := c.UserContext()
		ctx = WithContext(ctx, &reqLogger) // Embed Logger
		ctx = WithTraceID(ctx, traceID)    // Embed ID string
		c.SetUserContext(ctx)              // Save back to Fiber

		// Process Request
		err := c.Next()

		// Metrics & Status Calculation
		latencyMs := float64(time.Since(start).Nanoseconds()) / 1e6

		// Default HTTP status
		httpCode := c.Response().StatusCode()
		log.Info().Err(err).Int("http_code", httpCode).Msg("Request completed")
		if err != nil {
			if httpCode == 200 { // Force 500 if error occurred but status wasn't set
				httpCode = 500
			}
		}

		// Determine Logic Status (success/error) and Log Level
		statusStr := "success"
		event := log.Info()

		if httpCode >= 500 {
			event = log.Error()
			statusStr = "error"
		} else if httpCode >= 400 {
			event = log.Warn()
			statusStr = "error"
		}

		// Write the Log
		event.
			Str("trace_id", traceID).
			Str("action", "http_request").
			Str("status", statusStr). // "success" or "error"

			// HTTP Details
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Int("http_code", httpCode).
			Dur("latency", time.Duration(latencyMs)).

			// Final Message
			Msg("Incoming Request")

		return err
	}
}

package logger

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Handle the request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code (default to 500 if error exists but status not set)
		status := c.Response().StatusCode()
		if err != nil {
			if status == 200 { // If error but status is still 200, force 500
				status = 500
			}
		}

		// Determine log level based on status
		event := log.Info()
		if status >= 500 {
			event = log.Error()
		} else if status >= 400 {
			event = log.Warn()
		}

		// Write the log
		event.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.IP()).
			Str("req_id", c.GetRespHeader("X-Request-ID")). // Trace ID
			Msg("Incoming Request")

		return err
	}
}

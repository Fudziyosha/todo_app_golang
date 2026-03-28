package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/sirupsen/logrus"
)

// default logger for fiber
func logrusLoggerInstance(c fiber.Ctx, data *logger.Data, cfg *logger.Config) error {
	// Check if Skip is defined and call it.
	// Now, if Skip(c) == true, we SKIP logging:
	if cfg.Skip != nil && cfg.Skip(c) {
		return nil // Skip logging if Skip returns true
	}
	logrus.WithFields(logrus.Fields{
		"status_code": c.Response().StatusCode(),
		"latency":     data.Stop.Sub(data.Start),
		"ip":          c.IP(),
		"method":      c.Method(),
	}).Info(c.Path())

	return nil
}

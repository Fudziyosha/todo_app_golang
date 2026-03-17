package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/spf13/viper"
)

func InitMiddleware(app *fiber.App) {
	redisStorage := redis.New(redis.Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Username: "",
		Password: "",
		Database: 0,
	})

	app.Use(logger.New())
	app.Use(compress.New())
	app.Use(session.New(session.Config{
		Storage:         redisStorage,
		CookieSecure:    true,
		CookieHTTPOnly:  true,
		AbsoluteTimeout: viper.GetDuration("middleware.AbsoluteCookieTimeout"),
	}))
}

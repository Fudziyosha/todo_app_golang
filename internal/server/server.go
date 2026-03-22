package server

import (
	"fmt"
	"web_todos/internal/config"
	"web_todos/internal/handler"
	"web_todos/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/template/html/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	app     *fiber.App
	handler *handler.Handler
	config  *config.Config
}

func NewServer(config *config.Config) *Server {
	engine := html.New("./internal/html/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return &Server{
		app:    app,
		config: config,
	}
}

func (s *Server) Server(repo *repository.Repository) error {
	statics := viper.GetBool("server.statics")
	if statics == true {
		s.app.Use("/static", static.New("./static"))
		s.app.Use("/uploads", static.New("./uploads"))
	}

	s.InitMiddleware()
	s.handler = handler.NewHandler(repo)

	err := s.RegisterRoutes()
	if err != nil {
		logrus.Error("server: failed register routes ", err)
		return err
	}

	handleAddr := fmt.Sprintf("%v:%v", s.config.Host, s.config.Port)

	log.Fatal(s.app.Listen(handleAddr), fiber.ListenConfig{
		EnablePrefork: viper.GetBool("server.Prefork"),
	})
	return nil
}

func (s *Server) RegisterRoutes() error {
	s.app.Get("/", s.handler.CheckCookieAuthenticated, s.handler.TodoHandler.GetHome)
	s.app.Post("/", s.handler.CheckCookieAuthenticated, s.handler.TodoHandler.CreateListInHomePage)

	// Todos route
	lists := s.app.Group("/list")

	lists.Get("/:id/:filters", s.handler.CheckCookieAuthenticated, s.handler.TodoHandler.GetTasksByUser)
	lists.Post("/:id/:filters", s.handler.CheckCookieAuthenticated, s.handler.TodoHandler.TaskHandler)

	// User route
	user := s.app.Group("/user")

	user.Get("/register", s.handler.UserHandler.GetRegistrationPage)
	user.Post("/register", s.handler.UserHandler.UserRegistration)

	user.Get("/login", s.handler.UserHandler.GetUserLogin)
	user.Post("/login", s.handler.UserHandler.UserLogin)

	user.Post("/settings", s.handler.CheckCookieAuthenticated, s.handler.UserHandler.UpdateUserNameAndAvatar)

	user.Get("/change-password", s.handler.CheckCookieAuthenticated, s.handler.UserHandler.ChangePassword)
	user.Post("/change-password", s.handler.CheckCookieAuthenticated, s.handler.UserHandler.UpdateUserPass)

	user.Post("/logout", s.handler.CheckCookieAuthenticated, s.handler.UserHandler.Logout)

	return nil
}

func (s *Server) InitMiddleware() {
	redisStorage := redis.New(redis.Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Username: "",
		Password: "",
		Database: 0,
	})

	s.app.Use(logger.New())
	s.app.Use(compress.New())
	s.app.Use(session.New(session.Config{
		Storage:         redisStorage,
		CookieSecure:    true,
		CookieHTTPOnly:  true,
		AbsoluteTimeout: viper.GetDuration("server.AbsoluteCookieTimeout"),
	}))
}

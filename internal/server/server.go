package server

import (
	"web_todos/internal/handler"
	"web_todos/internal/middleware"
	"web_todos/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	app     *fiber.App
	handler *Handler
}

func NewServer() *Server {
	engine := html.New("./internal/html/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return &Server{
		app: app,
	}
}

type Handler struct {
	todoHandler *handler.TodoHandler
	userHandler *handler.UserHandler
	validate    *handler.StructValidator
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		validate:    &handler.StructValidator{Validator: validator.New()},
		todoHandler: handler.NewTodoHandler(repo),
		userHandler: handler.NewUserHandler(repo),
	}
}

func (s *Server) Server(repo *repository.Repository) {
	statics := viper.GetBool("server.statics")
	if statics == true {
		s.app.Use("/static", static.New("./static"))
		s.app.Use("/uploads", static.New("./uploads"))
	}

	middleware.InitMiddleware(s.app)
	s.handler = NewHandler(repo)

	err := s.RegisterRoutes()
	if err != nil {
		logrus.Error("server: failed register routes ", err)
	}

	log.Fatal(s.app.Listen(":3000"), fiber.ListenConfig{
		EnablePrefork: viper.GetBool("server.Prefork"),
	})
}

func (s *Server) RegisterRoutes() error {
	s.app.Get("/", s.handler.todoHandler.GetHome)
	s.app.Post("/", s.handler.todoHandler.CreateListInHomePage)

	// Todos route
	lists := s.app.Group("/list")

	lists.Get("/:id/:filters", s.handler.todoHandler.GetTasksByUser)
	lists.Post("/:id/:filters", s.handler.todoHandler.TaskHandler)

	// User route
	user := s.app.Group("/user")

	user.Get("/register", s.handler.userHandler.GetRegistrationPage)
	user.Post("/register", s.handler.userHandler.UserRegistration)

	user.Get("/login", s.handler.userHandler.GetUserLogin)
	user.Post("/login", s.handler.userHandler.UserLogin)

	user.Post("/settings", s.handler.userHandler.UpdateUserNameAndAvatar)

	user.Get("/change-password", s.handler.userHandler.ChangePassword)
	user.Post("/change-password", s.handler.userHandler.UpdateUserPass)

	user.Post("/logout", s.handler.userHandler.Logout)

	return nil
}

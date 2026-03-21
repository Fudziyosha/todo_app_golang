package server

import (
	"web_todos/internal/handler"
	"web_todos/internal/middleware"
	"web_todos/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	app *fiber.App
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

func (s *Server) Server(repo *repository.Repository) {
	statics := viper.GetBool("server.statics")
	if statics == true {
		s.app.Use("/static", static.New("./static"))
		s.app.Use("/uploads", static.New("./uploads"))
	}

	middleware.InitMiddleware(s.app)
	newHandler := handler.NewHandler(repo)

	err := s.RegisterRoutes(newHandler, s.app)
	if err != nil {
		logrus.Error("server: failed register routes ", err)
	}

	log.Fatal(s.app.Listen(":3000"), fiber.ListenConfig{
		EnablePrefork: viper.GetBool("server.Prefork"),
	})
}

func (s *Server) RegisterRoutes(h *handler.Handler, app *fiber.App) error {
	app.Get("/", h.GetHome)
	app.Post("/", h.CreateListInHomePage)

	// Todos route
	lists := app.Group("/list")

	lists.Get("/:id/:filters", h.GetTasksByUser)
	lists.Post("/:id/:filters", h.TaskHandler)

	// User route
	user := app.Group("/user")

	user.Get("/register", h.GetRegistrationPage)
	user.Post("/register", h.UserRegistration)

	user.Get("/login", h.GetUserLogin)
	user.Post("/login", h.UserLogin)

	user.Post("/settings", h.UpdateUserNameAndAvatar)

	user.Get("/change-password", h.ChangePassword)
	user.Post("/change-password", h.UpdateUserPass)

	user.Post("/logout", h.Logout)

	return nil
}

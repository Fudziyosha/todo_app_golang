package routes

import (
	"web_todos/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(h *handler.Handler, app *fiber.App) error {
	app.Get("/", h.GetHome)
	app.Post("/", h.CreateListInHomePage)

	//Todos route
	lists := app.Group("/list")

	lists.Get("/:id/:filters", h.GetTasksByUser)
	lists.Post("/:id/:filters", h.TaskHandler)

	//User route
	user := app.Group("/user")

	user.Get("/register", h.UserRegistration)
	user.Post("/register", h.UserRegistration)

	user.Get("/login", h.GetUserLogin)
	user.Post("/login", h.UserLogin)

	user.Post("/settings", h.UpdateUserSettings)

	user.Get("/change-password", h.ChangePassword)
	user.Post("/change-password", h.UpdateUserPass)

	user.Post("/logout", h.Logout)

	return nil
}

package serve

import (
	"web_todos/internal/handler"
	"web_todos/internal/middleware"
	"web_todos/internal/repository"
	"web_todos/internal/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Serve(todoModel *repository.Repository) {
	engine := html.New("./internal/html/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	statics := viper.GetBool("serve.statics")
	if statics == true {
		app.Use("/static", static.New("./static"))
		app.Use("/uploads", static.New("./uploads"))
	}

	middleware.InitMiddleware(app)
	newHandler := handler.NewHandler(todoModel)

	err := routes.RegisterRoutes(newHandler, app)
	if err != nil {
		logrus.Error("serve: failed register routes ", err)
	}

	log.Fatal(app.Listen(":3000"), fiber.ListenConfig{
		EnablePrefork: viper.GetBool("serve.Prefork"),
	})
}

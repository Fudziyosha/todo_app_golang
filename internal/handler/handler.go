package handler

import (
	"web_todos/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	TodoHandler *TodoHandler
	UserHandler *UserHandler
	Validate    *StructValidator
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Validate:    &StructValidator{Validator: validator.New()},
		TodoHandler: NewTodoHandler(repo),
		UserHandler: NewUserHandler(repo),
	}
}

func (h *Handler) CheckCookieAuthenticated(c fiber.Ctx) error {
	sess := session.FromContext(c)
	authenticated := sess.Get(sessionAuthenticated)

	if authenticated == nil {
		return c.Redirect().To("user/login")
	}
	return c.Next()
}

func GetUserIDInSession(c fiber.Ctx) (uuid.UUID, error) {
	sess := session.FromContext(c)
	userID := sess.Get(sessionUserIDKey)

	stringUserID := userID.(string)
	parseUUID, err := uuid.Parse(stringUserID)
	if err != nil {
		logrus.Error("handler: failed parse uuid ", err)
		return uuid.UUID{}, err
	}
	return parseUUID, nil
}

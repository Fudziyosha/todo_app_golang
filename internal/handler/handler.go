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
	repo     *repository.Repository
	validate *StructValidator
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo:     repo,
		validate: &StructValidator{validate: validator.New()},
	}
}

func (h *Handler) GetUserIdInSession(c fiber.Ctx) (uuid.UUID, error) {
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

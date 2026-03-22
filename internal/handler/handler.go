package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func GetUserIdInSession(c fiber.Ctx) (uuid.UUID, error) {
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

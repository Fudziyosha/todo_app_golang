package entities

import uuid "github.com/jackc/pgtype/ext/gofrs-uuid"

type User struct {
	ID        uuid.UUID
	Name      string `form:"user-name"`
	Surname   string `form:"user-surname"`
	Email     string `form:"user-email" validate:"required,email"`
	Password  string `form:"user-password" validate:"required,min=8"`
	ImageName *string
	PathImage *string
}

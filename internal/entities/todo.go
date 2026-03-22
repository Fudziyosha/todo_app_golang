package entities

import (
	"time"

	"github.com/google/uuid"
)

type List struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Todos     []Todo
}

type Todo struct {
	ID            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Status        bool      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedInList uuid.UUID `json:"created_in_list"`
}

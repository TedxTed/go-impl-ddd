package entity

import uuid "github.com/google/uuid"

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
}

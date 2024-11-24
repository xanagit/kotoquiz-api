package dto

import "github.com/google/uuid"

type WordIdsList struct {
	Ids []uuid.UUID `json:"ids"`
}

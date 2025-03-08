package domain

import (
	"time"
)

type Film struct {
	ID            uint `gorm:"primaryKey"`
	CreatorUserID int
	Title         string    `json:"title"`
	Director      string    `json:"director" validate:"required"`
	Release       time.Time `json:"release" validate:"required"`
}

type FilmFilter struct {
	Title    string    `schema:"title"`
	Director string    `schema:"director"`
	Release  time.Time `schema:"release"`
}

package models

import "time"

// AstronomicalResearcher представляет исследователя экзопланет
type AstronomicalResearcher struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	FullName     string    `json:"full_name" db:"full_name"`
	Institution  string    `json:"institution" db:"institution"`
	Email        string    `json:"email" db:"email"`
	IsModerator  bool      `json:"is_moderator" db:"is_moderator"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}



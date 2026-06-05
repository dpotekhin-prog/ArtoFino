package models

import "time"

type AuthorValidationRequest struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Status    string    `gorm:"type:text;not null" json:"status"`
	Comment   *string   `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

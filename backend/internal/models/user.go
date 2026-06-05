package models

import "time"

type User struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email      string    `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Name       string    `gorm:"type:text;not null" json:"name"`
	KeycloakID string    `gorm:"type:text;uniqueIndex;not null" json:"keycloak_id"`
	CreatedAt  time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

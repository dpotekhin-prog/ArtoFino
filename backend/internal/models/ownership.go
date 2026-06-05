package models

import "time"

type Ownership struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	ObjectID  string    `gorm:"type:text;not null;index" json:"object_id"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
}

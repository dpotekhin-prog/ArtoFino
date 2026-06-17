package models

import "time"

type Transfer struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ObjectID    string    `gorm:"type:varchar(24);not null;index" json:"objectId"`           // Ссылка на MongoDB ID картины
	RequesterID string    `gorm:"type:uuid;not null;index" json:"requesterId"`               // Кто запрашивает картину
	Destination string    `gorm:"type:text;not null" json:"destination"`                     // Адрес назначения (галерея, выставка)
	Status      string    `gorm:"type:varchar(50);not null;default:'pending'" json:"status"` // pending, approved, rejected
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}

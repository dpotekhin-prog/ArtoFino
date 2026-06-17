package models

import "time"

type Transaction struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ObjectID    string    `gorm:"type:varchar(24);not null;index" json:"objectId"`
	FromUserID  string    `gorm:"type:uuid;not null;index" json:"from_user_id"`
	ToUserID    string    `gorm:"type:uuid;not null;index" json:"to_user_id"`
	AmountCents int64     `gorm:"type:bigint;not null" json:"amount_cents"`
	SharePct    float64   `gorm:"type:numeric(5,2);not null" json:"sharePct"`
	Currency    string    `gorm:"type:text;not null;default:'EUR'" json:"currency"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

type ArtShareBalance struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ObjectID  string    `gorm:"type:varchar(24);not null;uniqueIndex:idx_user_object" json:"objectId"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_user_object" json:"userId"`
	SharePct  float64   `gorm:"type:numeric(5,2);not null;default:0.00" json:"sharePct"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

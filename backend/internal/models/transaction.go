package models

import "time"

type Transaction struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FromUserID  string    `gorm:"type:uuid;not null;index" json:"from_user_id"`
	ToUserID    string    `gorm:"type:uuid;not null;index" json:"to_user_id"`
	AmountCents int64     `gorm:"type:bigint;not null" json:"amount_cents"`
	Currency    string    `gorm:"type:text;not null;default:'EUR'" json:"currency"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

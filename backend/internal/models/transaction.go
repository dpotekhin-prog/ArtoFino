package models

import "time"

type Transaction struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ObjectID    string    `gorm:"type:varchar(24);not null;index" json:"objectId"` // Ссылка на Hex ID из MongoDB
	FromUserID  string    `gorm:"type:uuid;not null;index" json:"from_user_id"`    // Продавец (обычно ArtistID)
	ToUserID    string    `gorm:"type:uuid;not null;index" json:"to_user_id"`      // Покупатель (кто берет долю)
	AmountCents int64     `gorm:"type:bigint;not null" json:"amount_cents"`        // Сумма в центах (копейках)
	SharePct    float64   `gorm:"type:numeric(5,2);not null" json:"sharePct"`      // Процент купленной доли
	Currency    string    `gorm:"type:text;not null;default:'EUR'" json:"currency"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

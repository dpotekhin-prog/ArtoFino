package models

import "time"

type AuthorApplication struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       string    `gorm:"type:uuid;not null;index" json:"userId"`
	Bio          string    `gorm:"type:text;not null" json:"bio"`
	PortfolioURL string    `gorm:"type:text;not null" json:"portfolioUrl"`
	SocialLinks  string    `gorm:"type:text;not null" json:"socialLinks"`
	Status       string    `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}

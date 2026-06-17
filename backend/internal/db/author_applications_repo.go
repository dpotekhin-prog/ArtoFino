package db

import (
	"context"

	"ArtoFino/backend/internal/models"

	"gorm.io/gorm"
)

type AuthorApplicationsRepository struct {
	db *gorm.DB
}

func NewAuthorApplicationsRepository(db *gorm.DB) *AuthorApplicationsRepository {
	return &AuthorApplicationsRepository{db: db}
}

func (r *AuthorApplicationsRepository) Create(ctx context.Context, app *models.AuthorApplication) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *AuthorApplicationsRepository) FindByID(ctx context.Context, id string) (*models.AuthorApplication, error) {
	var app models.AuthorApplication
	err := r.db.WithContext(ctx).First(&app, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *AuthorApplicationsRepository) Update(ctx context.Context, app *models.AuthorApplication) error {
	return r.db.WithContext(ctx).Save(app).Error
}

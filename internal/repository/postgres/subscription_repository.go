package postgres

import (
	"github.com/H1vee/WeatherAPI/internal/models"
	"github.com/H1vee/WeatherAPI/internal/repository"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) repository.SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (r *subscriptionRepository) Create(subscription models.Subscription) error {
	return r.db.Create(&subscription).Error
}

func (r *subscriptionRepository) FindByToken(token string) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := r.db.Where("token =?", token).First(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) UpdateConfirmation(token string, confirmed bool) error {
	return r.db.Model(&models.Subscription{}).Where("token=?", token).Update("confirmed", confirmed).Error
}

func (r *subscriptionRepository) FindAllConfirmed() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("confirmed = ?", true).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *subscriptionRepository) Delete(token string) error {
	return r.db.Where("token =?", token).Delete(&models.Subscription{}).Error
}

package repository

import (
	"github.com/H1vee/WeatherAPI/internal/models"
)

type SubscriptionRepository interface {
	Create(subscription models.Subscription) error
	FindByToken(token string) (*models.Subscription, error)
	UpdateConfirmation(token string, confirmed bool) error
	Delete(token string) error
}

package services

import (
	"github.com/H1vee/WeatherAPI/internal/models"
)

type SubscriptionService interface {
	Subscribe(Subscription models.Subscription) error
	ConfirmSubscription(token string) error
	UnSubscribe(token string) error
}

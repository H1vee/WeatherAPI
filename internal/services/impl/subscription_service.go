package impl

import (
	"github.com/H1vee/WeatherAPI/internal/repository"
)

type subscriptionService struct {
	repo        repository.SubscriptionRepository
	emailSender EmailSender
}

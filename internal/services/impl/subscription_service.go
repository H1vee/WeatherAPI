package impl

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/H1vee/WeatherAPI/internal/models"
	"github.com/H1vee/WeatherAPI/internal/repository"
	"github.com/H1vee/WeatherAPI/internal/services"
)

type subscriptionService struct {
	repo        repository.SubscriptionRepository
	emailSender services.EmailSender
}

func NewSubscriptionService(repo repository.SubscriptionRepository, emailSender services.EmailSender) *subscriptionService {
	return &subscriptionService{
		repo:        repo,
		emailSender: emailSender,
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil
	}
	return hex.EncodeToString(bytes), nil
}

func (s *subscriptionService) Subscribe(subscription models.Subscription) error {
	token, err := generateToken()
	if err != nil {
		fmt.Errorf("failed to generate token: %w", err)
	}

	subscription.Token = token
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	subscription.Confirmed = false

	if err := s.repo.Create(subscription); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	if err := s.emailSender.SendConfirmationEmail(subscription.Email, subscription.City, token); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	return nil
}

func (s *subscriptionService) ConfirmSubscription(token string) error {
	subscription, err := s.repo.FindByToken(token)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}
	if subscription.Confirmed {
		return errors.New("subscription is already confirmed")
	}

	if err := s.repo.UpdateConfirmation(token, true); err != nil {
		return fmt.Errorf("failed to confirm subscription: %w", err)
	}
	return nil
}

func (s *subscriptionService) UnSubscribe(token string) error {
	_, err := s.repo.FindByToken(token)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}
	if err := s.repo.Delete(token); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

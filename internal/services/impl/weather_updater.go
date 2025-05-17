package impl

import (
	"fmt"
	"log"
	"time"

	"github.com/H1vee/WeatherAPI/internal/repository"
	"github.com/H1vee/WeatherAPI/internal/services"
)

type WeatherUpdater struct {
	subscriptionRepo repository.SubscriptionRepository
	weatherService   services.WeatherService
	emailSender      services.EmailSender
	hourlyTicker     *time.Ticker
	dailyTicker      *time.Ticker
	stopChan         chan struct{}
}

func NewWeatherUpdater(subscriptionRepo repository.SubscriptionRepository, weatherService services.WeatherService, emailSender services.EmailSender) *WeatherUpdater {
	return &WeatherUpdater{
		subscriptionRepo: subscriptionRepo,
		weatherService:   weatherService,
		emailSender:      emailSender,
		stopChan:         make(chan struct{}),
	}
}

func (u *WeatherUpdater) sendUpdates(frequency string) error {
	subscriptions, err := u.subscriptionRepo.FindAllConfirmed()
	if err != nil {
		return fmt.Errorf("failed to get confirmed subscription: %w", err)
	}
	for _, subscription := range subscriptions {
		if subscription.Frequency != frequency {
			continue
		}

		weatherData, err := u.weatherService.GetCurrentWeather(subscription.City)
		if err != nil {
			log.Printf("Failed to get weather for %s: %v", subscription.City, err)
			continue
		}
		if err := u.emailSender.SendWeatherUpdate(subscription.Email, subscription.City, subscription.Token, weatherData); err != nil {
			log.Printf("Failed to send weather update to %s: %v", subscription.Email, err)

		}
	}
	return nil
}

func (u *WeatherUpdater) Start() {
	u.hourlyTicker = time.NewTicker(1 * time.Hour)
	u.dailyTicker = time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-u.hourlyTicker.C:
				if err := u.sendUpdates("hourly"); err != nil {
					log.Printf("Error sending hourly updates: %w", err)
				}
			case <-u.dailyTicker.C:
				if err := u.sendUpdates("daily"); err != nil {
					log.Printf("Error sending daily updates: %w", err)
				}
			case <-u.stopChan:
				u.hourlyTicker.Stop()
				u.dailyTicker.Stop()
				return
			}
		}
	}()
}

func (u *WeatherUpdater) Stop() {
	close(u.stopChan)
}

package main

import (
	"fmt"
	"log"

	"github.com/H1vee/WeatherAPI/internal/config"
	"github.com/H1vee/WeatherAPI/internal/db"
	"github.com/H1vee/WeatherAPI/internal/email"
	"github.com/H1vee/WeatherAPI/internal/http/controllers"
	"github.com/H1vee/WeatherAPI/internal/repository/postgres"
	"github.com/H1vee/WeatherAPI/internal/services/impl"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load configuration
	cfg := config.Load("cmd/config/config.yaml")

	// Database connection
	database, err := db.ConnectDB(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.RunMigrations(cfg.Database.URL, cfg.Database.MigrationsDir); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	subscriptionRepo := postgres.NewSubscriptionRepository(database)

	// Initialize services
	weatherService := impl.NewWeatherService(cfg.Weather.APIKey)

	emailConfig := email.Config{
		Host:       cfg.Email.Host,
		Port:       int(cfg.Email.Port),
		Username:   cfg.Email.Username,
		Password:   cfg.Email.Password,
		FromEmail:  cfg.Email.FromEmail,
		WebsiteURL: cfg.Email.WebsiteURL,
	}
	emailSender := email.NewEmailSender(emailConfig)

	subscriptionService := impl.NewSubscriptionService(subscriptionRepo, emailSender)

	// Initialize weather updater
	weatherUpdater := impl.NewWeatherUpdater(subscriptionRepo, weatherService, emailSender)
	weatherUpdater.Start()
	defer weatherUpdater.Stop()

	// Initialize controllers
	weatherController := controllers.NewWeatherController(weatherService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	// Setup Echo
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	api := e.Group("/api")
	api.GET("/weather", weatherController.GetWeather)
	api.POST("/subscribe", subscriptionController.Subscribe)
	api.GET("/confirm/:token", subscriptionController.ConfirmSubscription)
	api.GET("/unsubscribe/:token", subscriptionController.UnSubscribe)

	// Start server
	log.Printf("Server starting on port %d", cfg.Server.Port)
	if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package controllers

import (
	"net/http"

	"github.com/H1vee/WeatherAPI/internal/models"
	"github.com/H1vee/WeatherAPI/internal/services"
	"github.com/labstack/echo/v4"
)

type SubscriptionController struct {
	subscriptionService services.SubscriptionService
}

func NewSubscriptionController(SubscriptionService services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		subscriptionService: SubscriptionService,
	}
}

type SubscriptionRequest struct {
	Email     string `json:"email" validate:"required,email"`
	City      string `json:"city" validate:"required"`
	Frequency string `json:"frequency" validate:"required,oneof=daily hourly"`
}

func (c *SubscriptionController) Subscribe(ctx echo.Context) error {
	var req SubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	subscription := models.Subscription{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
	}
	if err := c.subscriptionService.Subscribe(subscription); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, map[string]string{"message": "Subscription created. Please check your email to confirm."})
}

func (c *SubscriptionController) ConfirmSubscription(ctx echo.Context) error {
	token := ctx.Param("token")
	if token == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}
	if err := c.subscriptionService.ConfirmSubscription(token); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Subscription confirmed successfully"})
}

func (c *SubscriptionController) UnSubscribe(ctx echo.Context) error {
	token := ctx.Param("token")
	if token == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}
	if err := c.subscriptionService.UnSubscribe(token); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Unsubscribed successfully"})
}

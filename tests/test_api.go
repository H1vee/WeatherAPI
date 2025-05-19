package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/H1vee/WeatherAPI/internal/db"
	"github.com/H1vee/WeatherAPI/internal/http/controllers"
	"github.com/H1vee/WeatherAPI/internal/repository"
	"github.com/H1vee/WeatherAPI/internal/repository/postgres"
	"github.com/H1vee/WeatherAPI/internal/services"
	"github.com/H1vee/WeatherAPI/internal/services/impl"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendConfirmationEmail(email, city, token string) error {
	args := m.Called(email, city, token)
	return args.Error(0)
}

func (m *MockEmailSender) SendWeatherUpdate(email, city, token string, weatherData *services.WeatherData) error {
	args := m.Called(email, city, token, weatherData)
	return args.Error(0)
}

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetCurrentWeather(city string) (*services.WeatherData, error) {
	args := m.Called(city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.WeatherData), args.Error(1)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type APITestSuite struct {
	suite.Suite
	DB               *gorm.DB
	Echo             *echo.Echo
	SubscriptionRepo repository.SubscriptionRepository
	WeatherService   *MockWeatherService
	EmailSender      *MockEmailSender
	Tokens           map[string]string
}

func (suite *APITestSuite) SetupSuite() {
	testDBUrl := os.Getenv("TEST_DB_URL")
	if testDBUrl == "" {
		testDBUrl = "postgres://postgres:postgres@localhost:5432/weatherapi_test?sslmode=disable"
	}

	db, err := db.ConnectDB(testDBUrl)
	if err != nil {
		suite.T().Fatalf("Failed to connect to test database: %v", err)
	}
	suite.DB = db

	suite.WeatherService = &MockWeatherService{}
	suite.EmailSender = &MockEmailSender{}

	suite.SubscriptionRepo = postgres.NewSubscriptionRepository(suite.DB)

	suite.Echo = echo.New()
	suite.Echo.Validator = &CustomValidator{validator: validator.New()}

	suite.Tokens = make(map[string]string)
}

func (suite *APITestSuite) SetupTest() {
	suite.DB.Exec("TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
}

func (suite *APITestSuite) TearDownSuite() {
	sqlDB, _ := suite.DB.DB()
	sqlDB.Close()
}

func (suite *APITestSuite) makeRequest(method, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	rec := httptest.NewRecorder()
	suite.Echo.ServeHTTP(rec, req)
	return rec, err
}

func (suite *APITestSuite) TestGetWeatherSuccess() {
	weatherController := controllers.NewWeatherController(suite.WeatherService)
	suite.Echo.GET("/api/weather", weatherController.GetWeather)

	mockWeatherData := &services.WeatherData{
		Temperature: 25.5,
		Humidity:    60,
		Description: "Partly cloudy",
	}
	suite.WeatherService.On("GetCurrentWeather", "London").Return(mockWeatherData, nil)

	rec, err := suite.makeRequest(http.MethodGet, "/api/weather?city=London", nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var response services.WeatherData
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockWeatherData.Temperature, response.Temperature)
	assert.Equal(suite.T(), mockWeatherData.Humidity, response.Humidity)
	assert.Equal(suite.T(), mockWeatherData.Description, response.Description)

	suite.WeatherService.AssertCalled(suite.T(), "GetCurrentWeather", "London")
}

func (suite *APITestSuite) TestGetWeatherMissingCity() {
	weatherController := controllers.NewWeatherController(suite.WeatherService)
	suite.Echo.GET("/api/weather", weatherController.GetWeather)

	rec, err := suite.makeRequest(http.MethodGet, "/api/weather", nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "error")
	assert.Equal(suite.T(), "city parameter is required", response["error"])
}

func (suite *APITestSuite) TestGetWeatherServiceError() {
	weatherController := controllers.NewWeatherController(suite.WeatherService)
	suite.Echo.GET("/api/weather", weatherController.GetWeather)

	suite.WeatherService.On("GetCurrentWeather", "InvalidCity").Return(nil, fmt.Errorf("city not found"))

	rec, err := suite.makeRequest(http.MethodGet, "/api/weather?city=InvalidCity", nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "error")
	assert.Equal(suite.T(), "city not found", response["error"])
}

func (suite *APITestSuite) TestSubscriptionWorkflow() {
	subscriptionService := impl.NewSubscriptionService(suite.SubscriptionRepo, suite.EmailSender)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	suite.Echo.POST("/api/subscribe", subscriptionController.Subscribe)
	suite.Echo.GET("/api/confirm/:token", subscriptionController.ConfirmSubscription)
	suite.Echo.GET("/api/unsubscribe/:token", subscriptionController.UnSubscribe)

	suite.EmailSender.On("SendConfirmationEmail", "test@example.com", "Berlin", mock.AnythingOfType("string")).
		Run(func(args mock.Arguments) {
			token := args.Get(2).(string)
			suite.Tokens["confirmToken"] = token
		}).Return(nil)

	subscriptionData := map[string]string{
		"email":     "test@example.com",
		"city":      "Berlin",
		"frequency": "daily",
	}

	rec, err := suite.makeRequest(http.MethodPost, "/api/subscribe", subscriptionData)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	suite.EmailSender.AssertCalled(suite.T(), "SendConfirmationEmail", "test@example.com", "Berlin", mock.AnythingOfType("string"))

	confirmToken := suite.Tokens["confirmToken"]
	assert.NotEmpty(suite.T(), confirmToken)

	rec, err = suite.makeRequest(http.MethodGet, "/api/confirm/"+confirmToken, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	subscription, err := suite.SubscriptionRepo.FindByToken(confirmToken)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), subscription.Confirmed)

	rec, err = suite.makeRequest(http.MethodGet, "/api/unsubscribe/"+confirmToken, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	_, err = suite.SubscriptionRepo.FindByToken(confirmToken)
	assert.Error(suite.T(), err)
}

func (suite *APITestSuite) TestSubscribeInvalidData() {
	subscriptionService := impl.NewSubscriptionService(suite.SubscriptionRepo, suite.EmailSender)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	suite.Echo.POST("/api/subscribe", subscriptionController.Subscribe)

	testCases := []struct {
		name string
		data map[string]string
	}{
		{
			name: "Missing email",
			data: map[string]string{
				"city":      "Berlin",
				"frequency": "daily",
			},
		},
		{
			name: "Invalid email",
			data: map[string]string{
				"email":     "not-an-email",
				"city":      "Berlin",
				"frequency": "daily",
			},
		},
		{
			name: "Missing city",
			data: map[string]string{
				"email":     "test@example.com",
				"frequency": "daily",
			},
		},
		{
			name: "Invalid frequency",
			data: map[string]string{
				"email":     "test@example.com",
				"city":      "Berlin",
				"frequency": "weekly",
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			rec, err := suite.makeRequest(http.MethodPost, "/api/subscribe", tc.data)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

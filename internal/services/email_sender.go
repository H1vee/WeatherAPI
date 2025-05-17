package services

type EmailSender interface {
	SendConfirmationEmail(email, city, token string) error
	SendWeatherUpdate(email, city, token string, weatherData *WeatherData) error
}

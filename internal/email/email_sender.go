package email

import (
	"fmt"
	"net/smtp"

	"github.com/H1vee/WeatherAPI/internal/services"
)

type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	FromEmail  string
	WebsiteURL string
}

type EmailSender struct {
	config Config
}

func NewEmailSender(config Config) *EmailSender {
	return &EmailSender{
		config: config,
	}
}

func (s *EmailSender) SendConfirmationEmail(email, city, token string) error {
	subject := "Confirm Your Weather Update Subscription"
	confirmURL := fmt.Sprintf("%s/api/confirm/%s", s.config.WebsiteURL, token)

	body := fmt.Sprintf(`Hello,

Thank you for subscribing to weather updates for %s.

Please confirm your subscription by clicking the link below:
%s

If you did not request this subscription, please ignore this email.

Best regards`, city, confirmURL)

	return s.sendEmail(email, subject, body)
}

func (s *EmailSender) SendWeatherUpdate(email, city, token string, weatherData *services.WeatherData) error {
	subject := fmt.Sprintf("Weather Update for %s", city)
	unsubscribeURL := fmt.Sprintf("%s/api/unsubscribe/%s", s.config.WebsiteURL, token)

	body := fmt.Sprintf(`Hello,

Here is your weather update for %s:

Temperature: %.1f°C
Humidity: %d%%
Conditions: %s

To unsubscribe from these updates, click the link below:
%s

Best regards`,
		city,
		weatherData.Temperature,
		weatherData.Humidity,
		weatherData.Description,
		unsubscribeURL)

	return s.sendEmail(email, subject, body)
}

func (s *EmailSender) sendEmail(to, subject, body string) error {
	headers := make(map[string]string)
	headers["From"] = s.config.FromEmail
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	smtpAddr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	err := smtp.SendMail(smtpAddr, auth, s.config.FromEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

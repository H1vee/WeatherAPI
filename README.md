# WeatherAPI

A robust weather subscription service built with Go that allows users to subscribe to weather updates for their favorite cities. Users receive email notifications with current weather information at their chosen frequency (hourly or daily).

## Features

-  **Real-time Weather Data**: Get current weather information for any city
-  **Email Subscriptions**: Subscribe to weather updates via email
-  **Flexible Frequency**: Choose between hourly or daily notifications  
-  **Email Confirmation**: Double opt-in subscription process
-  **Secure Unsubscribe**: Easy one-click unsubscribe functionality
-  **Docker Support**: Full containerization with Docker Compose
-  **PostgreSQL Database**: Reliable data persistence
-  **Comprehensive Testing**: Unit and integration tests included

## Tech Stack

- **Backend**: Go (Golang) with Echo framework
- **Database**: PostgreSQL with GORM ORM
- **Email**: SMTP integration (Gmail configured)
- **Weather API**: WeatherAPI.com integration
- **Containerization**: Docker & Docker Compose
- **Migration**: golang-migrate
- **Validation**: go-playground/validator

## Project Structure

```
WeatherAPI/
├── cmd/
│   ├── server/
│   │   └── main.go              # Application entry point
│   └── config/
│       └── config.yaml          # Configuration file
├── internal/
│   ├── config/                  # Configuration management
│   ├── db/                      # Database connection and migrations
│   ├── email/                   # Email service implementation
│   ├── http/
│   │   └── controllers/         # HTTP request handlers
│   ├── models/                  # Data models
│   ├── repository/              # Data access layer
│   └── services/               # Business logic layer
├── migrations/                  # Database migrations
├── tests/                      # Test files
├── Dockerfile                  # Docker build configuration
├── docker-compose.yml          # Docker Compose configuration
└── README.md                   # Project documentation
```

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Weather API key from [WeatherAPI.com](https://www.weatherapi.com/)
- Gmail account with app password for SMTP

### 1. Clone the Repository

```bash
git clone https://github.com/H1vee/WeatherAPI.git
cd WeatherAPI
```

### 2. Configure the Application

Update `cmd/config/config.yaml` with your credentials:

```yaml
server:
  port: 8080

database:
  url: "postgres://postgres:postgres@postgres:5432/weatherapi?sslmode=disable"
  migrations_dir: "./migrations"

weather:
  api_key: "your_weather_api_key_here"

email:
  host: "smtp.gmail.com"
  port: 587
  username: "your_email@gmail.com"
  password: "your_app_password"
  from_email: "your_email@gmail.com"
  website_url: "http://localhost:8080"
```

### 3. Start the Application

```bash
# Start all services
docker compose up

# Or run in detached mode
docker compose up -d

# View logs
docker compose logs -f
```

The application will be available at `http://localhost:8080`

## API Endpoints

### Weather

- **GET** `/api/weather?city={city_name}` - Get current weather for a city

### Subscriptions

- **POST** `/api/subscribe` - Subscribe to weather updates
- **GET** `/api/confirm/{token}` - Confirm email subscription
- **GET** `/api/unsubscribe/{token}` - Unsubscribe from updates

### Example Usage

#### Get Weather Information

```bash
curl "http://localhost:8080/api/weather?city=London"
```

Response:
```json
{
  "temperature": 15.5,
  "humidity": 65,
  "description": "Partly cloudy"
}
```

#### Subscribe to Weather Updates

```bash
curl -X POST "http://localhost:8080/api/subscribe" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "city": "London",
    "frequency": "daily"
  }'
```

Response:
```json
{
  "message": "Subscription successful. Confirmation email sent."
}
```

## Development

### Local Development Setup

1. **Install Dependencies**
   ```bash
   go mod download
   ```

2. **Setup Local Database**
   ```bash
   # Start only PostgreSQL
   docker compose up postgres -d
   ```

3. **Run the Application**
   ```bash
   go run cmd/server/main.go
   ```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test suite
go test ./tests/...
```

### Database Management

#### Manual Migrations

```bash
# Apply migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/weatherapi?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/weatherapi?sslmode=disable" down
```

#### Create New Migration

```bash
migrate create -ext sql -dir migrations -seq create_new_table
```

## Configuration

### Environment Variables

You can override configuration values using environment variables:

- `DB_URL` - Database connection string
- `WEATHER_API_KEY` - Weather API key
- `EMAIL_USERNAME` - SMTP username
- `EMAIL_PASSWORD` - SMTP password

### Email Setup (Gmail)

1. Enable 2-Factor Authentication on your Gmail account
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a password for "Mail"
3. Use the generated password in your configuration

## Monitoring and Logging

The application includes:
- Structured logging with request/response details
- Health check endpoints for containers
- Graceful shutdown handling
- Error recovery middleware

## Security Considerations

- Email confirmation required for subscriptions
- Secure token generation for subscriptions
- CORS middleware enabled
- SQL injection protection via GORM
- Input validation on all endpoints

## Docker Commands

```bash
# Build and start services
docker compose up --build

# Stop services
docker compose down

# View service logs
docker compose logs app
docker compose logs postgres

# Execute commands in running container
docker compose exec app sh
docker compose exec postgres psql -U postgres -d weatherapi

# Remove all containers and volumes
docker compose down -v
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure PostgreSQL container is healthy
   - Check database URL in configuration
   - Verify port 5432 is not in use by another service

2. **Email Not Sending**
   - Verify Gmail app password is correct
   - Check SMTP configuration
   - Ensure less secure app access is enabled (if not using app password)

3. **Weather API Errors**
   - Verify your WeatherAPI.com API key
   - Check API rate limits
   - Ensure internet connectivity from container

### Debug Mode

To run with debug logging:

```bash
# Add to docker-compose.yml environment section
environment:
  - LOG_LEVEL=debug
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [WeatherAPI.com](https://www.weatherapi.com/) for weather data
- [Echo Framework](https://echo.labstack.com/) for HTTP handling
- [GORM](https://gorm.io/) for database operations
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

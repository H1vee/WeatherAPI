# WeatherAPI

A RESTful API application that provides weather forecasts and allows users to subscribe to updates for specific cities.

## Features

- Get current weather for any city
- Subscribe to weather updates via email (hourly or daily)
- Email confirmation of subscriptions
- Unsubscribe from updates

## API Endpoints

- `GET /api/weather?city={city}` - Get current weather for a city
- `POST /api/subscribe` - Subscribe to weather updates
- `GET /api/confirm/{token}` - Confirm subscription
- `GET /api/unsubscribe/{token}` - Unsubscribe from updates

## Tech Stack

- Go
- Echo (web framework)
- GORM (ORM)
- PostgreSQL
- Golang-migrate (database migrations)
- WeatherAPI.com (weather data provider)

## Setup and Installation

### Prerequisites

- Go 1.20 or newer
- PostgreSQL
- [WeatherAPI.com](https://www.weatherapi.com/) API key

### Installation

1. Clone the repository
```
git clone https://github.com/H1vee/WeatherAPI.git
cd WeatherAPI
```

2. Install dependencies
```
go mod download
```

3. Configure the application
   - Copy `config/config.example.yaml` to `config/config.yaml`
   - Edit `config/config.yaml` with your database and WeatherAPI.com credentials

4. Create the database
```
createdb weatherapi
```

5. Run the server
```
go run cmd/server/main.go
```

### Running with Docker

```
# Build the Docker image
docker build -t weatherapi .

# Run the container
docker run -p 8080:8080 --env-file .env weatherapi
```

## Configuration

Configuration is handled via a YAML file in `config/config.yaml`:

```yaml
server:
  port: 8080

database:
  url: "postgres://postgres:postgres@localhost:5432/weatherapi?sslmode=disable"
  migrations_dir: "./migrations"

weather:
  api_key: "your_weatherapi_com_key_here"

email:
  host: "smtp.example.com"
  port: 587
  username: "your_email@example.com"
  password: "your_email_password"
  from_email: "weather@example.com"
  website_url: "http://localhost:8080"
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── config/
│   └── config.yaml           # Configuration file
├── internal/
│   ├── db/                   # Database connection and migrations
│   ├── email/                # Email service
│   ├── http/
│   │   └── controllers/      # HTTP controllers
│   ├── models/               # Data models
│   ├── repository/           # Data access layer
│   │   └── postgres/         # PostgreSQL implementation
│   └── services/             # Business logic
│       └── impl/             # Service implementations
├── migrations/               # Database migrations
└── Dockerfile                # Docker configuration
```

## License

MIT

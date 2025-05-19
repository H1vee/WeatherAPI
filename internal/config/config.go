package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		URL           string `yaml:"url"`
		MigrationsDir string `yaml:"migrations_dir"`
	}
	Weather struct {
		APIKey string `yaml:"api_key"`
	}
	Email struct {
		Host       string `yaml:"host"`
		Port       uint   `yaml:"port"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		FromEmail  string `yaml:"from_email"`
		WebsiteURL string `yaml:"website_url"`
	}
}

func Load(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening configuration file: %v", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Error decoding yaml: %v", err)
	}
	return &cfg
}

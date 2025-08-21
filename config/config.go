// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App       AppConfig
	DB        DBConfig
	Auth      GRPCConfig
	User      GRPCConfig
	Store     GRPCConfig
	Product   GRPCConfig
	Order     GRPCConfig
	Company   GRPCConfig
	Promotion GRPCConfig
	Redis     RedisConfig
	Rabbit    RabbitMQConfig
}

type AppConfig struct {
	Env       string
	Port      string
	JWTSecret string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type GRPCConfig struct {
	Port string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Get memuat konfigurasi berdasarkan lingkungan aplikasi.
// Urutan prioritas pemuatan file: .env.local -> .env.staging -> .env
// Untuk 'production', disarankan untuk tidak menggunakan file .env sama sekali.
func Get() *Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default ke development jika tidak diset
	}

	// Untuk production, kita tidak memuat file .env sama sekali.
	// Konfigurasi harus disediakan melalui environment variables di server.
	if env != "production" {
		filenames := []string{}
		if env == "development" {
			// Untuk development, utamakan .env.local
			filenames = append(filenames, ".env.local")
		}
		if env == "staging" {
			// Untuk staging, gunakan .env.staging
			filenames = append(filenames, ".env.staging")
		}
		// Selalu tambahkan .env sebagai fallback
		filenames = append(filenames, ".env")

		// Muat file .env berdasarkan urutan prioritas
		// godotenv akan memuat yang pertama kali ditemukannya.
		err := godotenv.Load(filenames...)
		if err != nil {
			log.Printf("Warning: Could not load any of the following .env files: %v. Relying on system environment variables.", filenames)
		}
	}

	return &Config{
		App: AppConfig{
			Env:       os.Getenv("APP_ENV"),
			Port:      os.Getenv("APP_PORT"),
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Auth: GRPCConfig{
			Port: os.Getenv("AUTH_SERVICE_PORT"),
		},
		User: GRPCConfig{
			Port: os.Getenv("USER_SERVICE_PORT"),
		},
		Store: GRPCConfig{
			Port: os.Getenv("STORE_SERVICE_PORT"),
		},
		Product: GRPCConfig{
			Port: os.Getenv("PRODUCT_SERVICE_PORT"),
		},
		Order: GRPCConfig{
			Port: os.Getenv("ORDER_SERVICE_PORT"),
		},
		Company: GRPCConfig{
			Port: os.Getenv("COMPANY_SERVICE_PORT"),
		},
		Promotion: GRPCConfig{
			Port: os.Getenv("PROMOTION_SERVICE_PORT"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Rabbit: RabbitMQConfig{
			Host:     os.Getenv("RABBITMQ_HOST"),
			Port:     os.Getenv("RABBITMQ_PORT"),
			User:     os.Getenv("RABBITMQ_USER"),
			Password: os.Getenv("RABBITMQ_PASSWORD"),
		},
	}
}

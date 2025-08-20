package database

import (
	"fmt"
	"log"

	// Kita akan menggunakan objek config yang sudah kita buat
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Mengubah nama fungsi dan menambahkan parameter config
func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	// Variabel diambil dari objek cfg, bukan os.Getenv langsung
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

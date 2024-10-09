package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a package-level variable that holds the database connection
var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	log.Printf("Attempting to connect to database with DSN: host=%s port=%s user=%s dbname=%s", host, port, user, dbname)

	var err error
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/5): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after 5 attempts:", err)
	}

	log.Println("Successfully connected to the database")

	err = MigrateSchema()
	if err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func MigrateSchema() error {
	// List of all models that should be migrated
	models := []interface{}{
		&models.Character{},
		&models.Zkill{},
		&models.Kill{},
		&models.Region{},
		&models.System{},
		&models.Constellation{},
		&models.ESIItem{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}

	fmt.Println("Schema migration completed successfully")
	return nil
}

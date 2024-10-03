package db

import (
	"fmt"
	"log"
	"os"

	"github.com/tadeasf/eve-ran/src/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a package-level variable that holds the database connection
var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to the database")

	err = MigrateSchema()
	if err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func MigrateSchema() error {
	// List of all models that should be migrated
	models := []interface{}{
		&models.Character{},
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

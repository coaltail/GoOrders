package database

import (
	"fmt"
	"log"
	"os"

	"github.com/coaltail/GoOrders/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Zagreb",
		goDotEnvVariable("DB_USER"),
		goDotEnvVariable("DB_PASSWORD"),
		goDotEnvVariable("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Connected!")

	db.Logger = logger.Default.LogMode(logger.Info)
	models.AutoMigrate(db)

	DB = Dbinstance{
		Db: db,
	}
}

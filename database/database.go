package database

import (
	"flag"
	"github.com/MihaPecnik/restaurant/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() (*Database, error) {
	migrate := flag.Bool("migrate", false, "database migration")
	database := flag.String("database", "", "database connection")
	flag.Parse()
	conn, err := gorm.Open(postgres.Open(*database), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	log.Println("successfully created database connection")

	if *migrate {
		err := conn.AutoMigrate(
			&models.Product{},
			&models.Order{},
			&models.OrderMapping{},
			&models.Cost{})
		if err != nil {
			log.Println(err.Error())
		}
	}

	return &Database{db: conn}, nil
}

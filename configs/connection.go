package configs

import (
	"github.com/hodukihugi/winglets-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connection() *gorm.DB {
	databaseURI := os.Getenv("DATABASE_URI")
	db, err := gorm.Open(mysql.Open(databaseURI), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(
		&models.UserModel{},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

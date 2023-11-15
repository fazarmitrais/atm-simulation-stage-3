package postgre

import (
	"fmt"
	"log"
	"time"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/envLib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	username := envLib.GetEnv("DB_USERNAME")
	password := envLib.GetEnv("DB_PASSWORD")
	host := envLib.GetEnv("DB_HOST")
	dbName := envLib.GetEnv("DB_NAME")
	db, err := gorm.Open(postgres.Open("postgres://"+username+":"+password+"@"+host+"/"+dbName+"?sslmode=disable"),
		&gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&entity.Account{}, &trxEntity.Transaction{})
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqlDb.SetMaxIdleConns(2)
	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetConnMaxLifetime(time.Hour * 1)

	fmt.Println("Database successfully connected!")

	return db
}

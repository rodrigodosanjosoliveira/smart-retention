package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"smart-retention/internal/model"
)

var dsn string

func init() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" || appEnv == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Println("⚠️ Ignorando .env.development: não encontrado.")
		}
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("variaveis de ambiente do banco não configuradas")
	}

	dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s  port=%s sslmode=require",
		host,
		user,
		password,
		dbname,
		port,
	)
}

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("erro ao conectar no banco: ", err)
	}
	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Cliente{},
		&model.Item{},
		&model.Compra{},
		&model.CompraItem{},
		&model.DiaCompraCliente{},
	)
	if err != nil {
		log.Fatal("erro ao migrar o banco: ", err)
	} else {
		log.Println("Banco migrado com sucesso")
	}
}

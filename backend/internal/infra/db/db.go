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

func Connect() *gorm.DB {
	appEnv := os.Getenv("APP_ENV")

	fmt.Println("⚙️ APP_ENV =", appEnv)

	if appEnv == "" || appEnv == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Println("erro ao carregar o arquivo .env.development: ", err)
		}
	}

	dsn := getDSN()

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

func getDSN() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" || sslmode == "" {
		log.Fatal("variaveis de ambiente do banco não configuradas")
	}

	return "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
}

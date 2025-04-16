package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"smart-retention/internal/model"
)

func Connect() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=clientes port=5432 sslmode=disable"
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

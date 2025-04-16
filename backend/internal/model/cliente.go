package model

import "time"

type (
	Cliente struct {
		ID         string             `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
		CNPJ       string             `gorm:"unique;not null" json:"cnpj"`
		Nome       string             `json:"nome"`
		Telefone   string             `json:"telefone"`
		Email      string             `json:"email"`
		Endereco   string             `json:"endereco"`
		Itens      []Item             `gorm:"many2many:cliente_itens" json:"itens"`
		DiasCompra []DiaCompraCliente `gorm:"foreignKey:ClienteID;constraint:OnDelete:CASCADE" json:"dias_compra"`
	}

	Item struct {
		ID   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
		Nome string `json:"nome"`
	}

	Compra struct {
		ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		ClienteID  string
		Cliente    Cliente
		DataCompra time.Time
		Itens      []CompraItem
	}

	CompraItem struct {
		ID       string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
		CompraID string
		ItemID   string
		Item     Item
		Preco    float64
	}

	DiaCompraCliente struct {
		ClienteID string `gorm:"primaryKey" json:"-"`
		DiaSemana int    `gorm:"primaryKey" json:"dia_semana"`
	}
)

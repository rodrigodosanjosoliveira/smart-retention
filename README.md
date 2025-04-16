# 🧠 Smart Retention

Sistema inteligente de retenção de clientes para estabelecimentos comerciais. Permite o cadastro de clientes, registro de compras, alertas automáticos e dashboards com métricas.

---

## 📁 Estrutura do Projeto

```
smart-retention/
├── backend/     # Projeto em Go (Golang) com GORM e PostgreSQL
└── frontend/    # Interface em React + Vite + TailwindCSS
```

---

## 🚀 Como Rodar

### 1. Clonar o projeto

```bash
git clone https://github.com/SEU_USUARIO/smart-retention.git
cd smart-retention
```

### 2. Backend (Go)

#### Pré-requisitos

- Go 1.21+
- PostgreSQL rodando
- [Goose](https://github.com/pressly/goose) (para migrações)

#### Setup

```bash
cd backend
cp .env.example .env        # configure o banco
go mod tidy
go run main.go              # inicia a API
```

> Por padrão roda em: `http://localhost:8080`

### 3. Frontend (React)

#### Pré-requisitos

- Node.js 18+
- Yarn ou npm

#### Setup

```bash
cd frontend
cp .env.example .env        # configure a URL do backend
npm install
npm run dev
```

> Por padrão roda em: `http://localhost:5173`

---

## 📦 Funcionalidades

- Cadastro de cliente com dias de compra e itens recorrentes
- Registro de compras com múltiplos itens
- Visualização do histórico de compras
- Alertas inteligentes:
    - Cliente inativo
    - Ausente no dia previsto
    - Itens deixados de comprar
- Dashboard com métricas
- Notificações em tempo real com WebSocket

---

## 🛠️ Tecnologias

- **Backend:** Go, GORM, PostgreSQL, Gin, WebSocket
- **Frontend:** React, Vite, TailwindCSS, Axios

---

## 📄 Licença

Este projeto está sob a licença MIT.

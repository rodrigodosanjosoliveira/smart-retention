#!/bin/bash

set -e

echo ""
echo "=== ETAPA 1: Build do frontend React (Vite) ==="
cd frontend
npm install
npm run build
cd ..

echo ""
echo "=== ETAPA 2: Copiar dist para backend/web ==="
rm -rf backend/web
mkdir -p backend/web
cp -r frontend/dist/* backend/web/

echo ""
echo "=== ETAPA 3: Build do backend Go com embed ==="
cd backend
go build -o ../app
cd ..

echo ""
echo "=== ETAPA 4: Gerar pacote .zip final ==="
rm -rf pacote-final
mkdir pacote-final
cp app pacote-final/
cd pacote-final
zip -r ../pacote-app.zip ./*
cd ..

echo ""
echo "✅ Build finalizado com sucesso: pacote-app.zip"

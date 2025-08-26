#!/bin/bash

# Script para configurar variables de entorno secretas en Fly.io
# Ejecutar después de crear la app en Fly.io

echo "Configurando variables de entorno secretas en Fly.io..."

# Verificar que flyctl esté instalado
if ! command -v flyctl &> /dev/null; then
    echo "Error: flyctl no está instalado. Instálalo primero con:"
    echo "curl -L https://fly.io/install.sh | sh"
    exit 1
fi

# Verificar que estemos en el directorio correcto
if [ ! -f "fly.toml" ]; then
    echo "Error: No se encontró fly.toml. Ejecuta este script desde el directorio del backend."
    exit 1
fi

echo "Configurando variables de entorno..."

# Configurar DATABASE_URL (reemplazar con la URL real de PostgreSQL)
echo "Configurando DATABASE_URL..."
read -p "Ingresa la DATABASE_URL de PostgreSQL: " DATABASE_URL
flyctl secrets set DATABASE_URL="$DATABASE_URL"

# Configurar JWT_SECRET
echo "Configurando JWT_SECRET..."
read -s -p "Ingresa el JWT_SECRET: " JWT_SECRET
echo
flyctl secrets set JWT_SECRET="$JWT_SECRET"

# Configurar ALLOWED_ORIGINS
echo "Configurando ALLOWED_ORIGINS..."
read -p "Ingresa el dominio del frontend (ej: https://myapp.vercel.app): " ALLOWED_ORIGINS
flyctl secrets set ALLOWED_ORIGINS="$ALLOWED_ORIGINS"

# Configurar otras variables opcionales
flyctl secrets set APP_ENV="production"
flyctl secrets set LOG_LEVEL="info"
flyctl secrets set DB_MAX_OPEN_CONNS="25"
flyctl secrets set DB_MAX_IDLE_CONNS="5"
flyctl secrets set DB_CONN_MAX_LIFETIME="300s"

echo "✅ Variables de entorno configuradas exitosamente!"
echo "Puedes verificar las variables con: flyctl secrets list"
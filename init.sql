-- Script de inicialización para PostgreSQL
-- Este script se ejecuta automáticamente cuando se crea el contenedor

-- Crear la base de datos si no existe
SELECT 'CREATE DATABASE raborimet_crm'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'raborimet_crm')\gexec

-- Conectar a la base de datos
\c raborimet_crm;

-- Crear extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Configurar zona horaria
SET timezone = 'America/Mexico_City';

-- Crear esquema si no existe
CREATE SCHEMA IF NOT EXISTS public;

-- Otorgar permisos al usuario postgres
GRANT ALL PRIVILEGES ON DATABASE raborimet_crm TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO postgres;

-- Configurar permisos por defecto para objetos futuros
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO postgres;

-- Mensaje de confirmación
SELECT 'Base de datos raborimet_crm inicializada correctamente' AS mensaje;
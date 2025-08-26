# 🚀 Deployment Rápido en Fly.io

## Resumen Ejecutivo

Este backend está configurado para deployment en **Fly.io** con **PostgreSQL**. Fly.io es más rápido y económico que Heroku.

## 🏃‍♂️ Deployment en 3 pasos

### 1. Instalar flyctl
```bash
curl -L https://fly.io/install.sh | sh
flyctl auth login
```

### 2. Configurar PostgreSQL
```bash
# Crear base de datos
flyctl postgres create --name raborimet-crm-db --region mia

# Guardar la DATABASE_URL que se muestra
```

### 3. Deployment automático
```bash
cd raborimet-crm/backend
./deploy.sh
```

## 📁 Archivos de configuración incluidos

- ✅ `fly.toml` - Configuración de la app
- ✅ `Dockerfile` - Imagen optimizada para producción
- ✅ `.env.example` - Variables de entorno necesarias
- ✅ `fly-secrets.sh` - Script para configurar secretos
- ✅ `deploy.sh` - Deployment automatizado
- ✅ `DEPLOYMENT.md` - Documentación completa

## 🔧 Configuración manual de secretos

Si prefieres configurar manualmente:

```bash
# Variables críticas
flyctl secrets set DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=require"
flyctl secrets set JWT_SECRET="tu-clave-secreta"
flyctl secrets set ALLOWED_ORIGINS="https://tu-frontend.vercel.app"
```

## 💰 Costos estimados

- **PostgreSQL**: ~$1.94/mes
- **Backend**: ~$1.94/mes
- **Total**: ~$4/mes

## 🆘 Comandos útiles

```bash
# Ver logs
flyctl logs -f

# Estado de la app
flyctl status

# Conectar a la DB
flyctl postgres connect -a raborimet-crm-db

# Redeploy
flyctl deploy
```

## 🔗 Después del deployment

1. **Actualizar frontend**: Cambiar la URL del API a `https://raborimet-crm-backend.fly.dev`
2. **Probar API**: `curl https://raborimet-crm-backend.fly.dev/health`
3. **Configurar dominio personalizado** (opcional)

---

**¿Problemas?** Consulta `DEPLOYMENT.md` para la guía completa.
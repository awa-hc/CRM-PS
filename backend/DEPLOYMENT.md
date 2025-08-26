# Deployment en Fly.io - Raborimet CRM Backend

Esta guía te ayudará a deployar el backend del CRM en Fly.io con PostgreSQL.

## Prerrequisitos

- Cuenta en [Fly.io](https://fly.io)
- Go 1.21+ instalado localmente
- Git instalado

## 1. Instalación de flyctl

### macOS/Linux
```bash
curl -L https://fly.io/install.sh | sh
```

### Windows (PowerShell)
```powershell
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"
```

### Verificar instalación
```bash
flyctl version
```

## 2. Autenticación en Fly.io

```bash
flyctl auth login
```

Esto abrirá tu navegador para autenticarte.

## 3. Configuración de PostgreSQL

### Crear una base de datos PostgreSQL
```bash
# Crear una app de PostgreSQL
flyctl postgres create --name raborimet-crm-db --region mia
```

**Importante:** Guarda la información de conexión que se muestra, especialmente:
- Username
- Password
- Hostname
- Database name

### Obtener la URL de conexión
```bash
flyctl postgres connect -a raborimet-crm-db
```

La URL tendrá el formato:
```
postgres://username:password@hostname:5432/database_name?sslmode=require
```

## 4. Preparación del Backend

### Navegar al directorio del backend
```bash
cd raborimet-crm/backend
```

### Verificar archivos de configuración
Asegúrate de que existan estos archivos:
- `fly.toml` ✅
- `Dockerfile` ✅
- `.env.example` ✅
- `fly-secrets.sh` ✅

## 5. Crear la aplicación en Fly.io

```bash
# Crear la app (esto usará la configuración de fly.toml)
flyctl apps create raborimet-crm-backend
```

## 6. Configurar variables de entorno

### Opción A: Usar el script automático
```bash
./fly-secrets.sh
```

### Opción B: Configurar manualmente
```bash
# Configurar DATABASE_URL (usar la URL obtenida en el paso 3)
flyctl secrets set DATABASE_URL="postgres://username:password@hostname:5432/database_name?sslmode=require"

# Configurar JWT_SECRET (generar una clave segura)
flyctl secrets set JWT_SECRET="tu-clave-jwt-super-secreta-aqui"

# Configurar CORS para el frontend
flyctl secrets set ALLOWED_ORIGINS="https://tu-frontend-domain.vercel.app"

# Variables adicionales
flyctl secrets set APP_ENV="production"
flyctl secrets set LOG_LEVEL="info"
```

### Verificar variables configuradas
```bash
flyctl secrets list
```

## 7. Deployment inicial

```bash
# Hacer el primer deploy
flyctl deploy
```

Este comando:
1. Construirá la imagen Docker
2. La subirá a Fly.io
3. Creará y iniciará las máquinas virtuales
4. Ejecutará las migraciones (si existen)

## 8. Verificar el deployment

### Verificar estado de la app
```bash
flyctl status
```

### Ver logs
```bash
flyctl logs
```

### Probar la API
```bash
# Obtener la URL de la app
flyctl info

# Probar endpoint de salud
curl https://raborimet-crm-backend.fly.dev/health
```

## 9. Comandos útiles

### Escalar la aplicación
```bash
# Escalar a 2 instancias
flyctl scale count 2

# Cambiar el tamaño de VM
flyctl scale vm shared-cpu-1x
```

### Conectarse a la base de datos
```bash
# Conectar directamente a PostgreSQL
flyctl postgres connect -a raborimet-crm-db
```

### Ejecutar migraciones manualmente
```bash
# Si necesitas ejecutar migraciones
flyctl ssh console
# Dentro del contenedor:
./migrate
```

### Monitoreo
```bash
# Ver métricas
flyctl metrics

# Ver logs en tiempo real
flyctl logs -f
```

## 10. Configuración del Frontend

Una vez que el backend esté funcionando, actualiza la configuración del frontend para usar la nueva URL:

```typescript
// En el frontend, actualizar la URL del API
const API_URL = 'https://raborimet-crm-backend.fly.dev/api/v1';
```

## Troubleshooting

### Error de conexión a la base de datos
1. Verificar que la DATABASE_URL esté correctamente configurada
2. Asegurar que la base de datos PostgreSQL esté funcionando:
   ```bash
   flyctl status -a raborimet-crm-db
   ```

### Error de build
1. Verificar que el Dockerfile esté correcto
2. Revisar los logs de build:
   ```bash
   flyctl logs --app raborimet-crm-backend
   ```

### App no responde
1. Verificar que el puerto 8080 esté configurado correctamente
2. Revisar los health checks en fly.toml
3. Ver logs de la aplicación:
   ```bash
   flyctl logs -f
   ```

## Costos estimados

- **PostgreSQL**: ~$1.94/mes (shared-cpu-1x, 1GB RAM)
- **Backend App**: ~$1.94/mes (shared-cpu-1x, 256MB RAM)
- **Total aproximado**: ~$4/mes

*Los costos pueden variar según el uso y la región.*

## Recursos adicionales

- [Documentación oficial de Fly.io](https://fly.io/docs/)
- [Fly.io Go Guide](https://fly.io/docs/languages-and-frameworks/golang/)
- [Fly Postgres](https://fly.io/docs/postgres/)
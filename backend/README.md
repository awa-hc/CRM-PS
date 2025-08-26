# Raborimet CRM - Backend

Backend del sistema CRM para empresas de construcción Raborimet, desarrollado en Go con Gin framework.

## 🚀 Características

- **Framework**: Gin (Go)
- **Base de datos**: PostgreSQL con GORM
- **Autenticación**: JWT (JSON Web Tokens)
- **Arquitectura**: REST API
- **Documentación**: Swagger/OpenAPI
- **Middleware**: CORS, Rate Limiting, Logging

## 📋 Requisitos previos

- Go 1.21 o superior
- PostgreSQL 12 o superior
- Git

## 🛠️ Instalación

### 1. Clonar el repositorio
```bash
git clone <repository-url>
cd raborimet-crm/backend
```

### 2. Instalar dependencias
```bash
go mod download
```

### 3. Configurar base de datos

#### Crear base de datos PostgreSQL:
```sql
CREATE DATABASE raborimet_crm;
CREATE USER raborimet_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE raborimet_crm TO raborimet_user;
```

### 4. Configurar variables de entorno

Copiar el archivo `.env` y ajustar las configuraciones:
```bash
cp .env.example .env
```

Editar `.env` con tus configuraciones:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=raborimet_user
DB_PASSWORD=your_password
DB_NAME=raborimet_crm
JWT_SECRET=your_jwt_secret_key
```

### 5. Ejecutar migraciones

Las migraciones se ejecutan automáticamente al iniciar la aplicación.

### 6. Ejecutar la aplicación

#### Modo desarrollo:
```bash
go run cmd/main.go
```

#### Compilar y ejecutar:
```bash
go build -o bin/raborimet-crm cmd/main.go
./bin/raborimet-crm
```

La API estará disponible en: `http://localhost:8080`

## 📚 Estructura del proyecto

```
backend/
├── cmd/
│   └── main.go              # Punto de entrada de la aplicación
├── config/
│   └── database.go          # Configuración de base de datos
├── controllers/
│   ├── auth_controller.go   # Controlador de autenticación
│   ├── client_controller.go # Controlador de clientes
│   └── project_controller.go# Controlador de proyectos
├── middleware/
│   └── auth.go              # Middleware de autenticación
├── models/
│   ├── user.go              # Modelo de usuario
│   ├── client.go            # Modelo de cliente
│   ├── project.go           # Modelo de proyecto
│   ├── quote.go             # Modelo de cotización
│   ├── invoice.go           # Modelo de factura
│   └── material.go          # Modelo de material
├── routes/
│   └── routes.go            # Definición de rutas
├── services/
│   └── auth_service.go      # Servicios de autenticación
├── utils/
│   └── (utilidades)         # Funciones de utilidad
├── .env                     # Variables de entorno
├── go.mod                   # Dependencias de Go
├── go.sum                   # Checksums de dependencias
└── README.md                # Este archivo
```

## 🔗 Endpoints principales

### Autenticación
- `POST /api/v1/auth/register` - Registrar usuario
- `POST /api/v1/auth/login` - Iniciar sesión
- `GET /api/v1/auth/profile` - Obtener perfil
- `PUT /api/v1/auth/profile` - Actualizar perfil
- `POST /api/v1/auth/change-password` - Cambiar contraseña

### Clientes
- `GET /api/v1/clients` - Listar clientes
- `GET /api/v1/clients/:id` - Obtener cliente
- `POST /api/v1/clients` - Crear cliente
- `PUT /api/v1/clients/:id` - Actualizar cliente
- `DELETE /api/v1/clients/:id` - Eliminar cliente
- `GET /api/v1/clients/stats` - Estadísticas de clientes

### Proyectos
- `GET /api/v1/projects` - Listar proyectos
- `GET /api/v1/projects/:id` - Obtener proyecto
- `POST /api/v1/projects` - Crear proyecto
- `PUT /api/v1/projects/:id` - Actualizar proyecto
- `DELETE /api/v1/projects/:id` - Eliminar proyecto
- `GET /api/v1/projects/stats` - Estadísticas de proyectos
- `GET /api/v1/projects/:id/materials` - Materiales del proyecto

### Utilidades
- `GET /health` - Health check
- `GET /api/info` - Información de la API

## 🔐 Autenticación

La API utiliza JWT (JSON Web Tokens) para la autenticación. Para acceder a endpoints protegidos:

1. Registrarse o iniciar sesión para obtener un token
2. Incluir el token en el header `Authorization`:
   ```
   Authorization: Bearer <your-jwt-token>
   ```

## 🗄️ Modelos de datos

### Usuario
- ID, Email, Password (hash)
- FirstName, LastName, Role
- IsActive, CreatedAt, UpdatedAt

### Cliente
- ID, Name, Email, Phone
- Address, City, State, ZipCode
- Company, TaxID, ContactType
- Notes, IsActive, CreatedAt, UpdatedAt

### Proyecto
- ID, Name, Description, ClientID
- Status, Priority, Type
- Address, City, State, ZipCode
- StartDate, EndDate, Budget
- EstimatedCost, ActualCost, Progress
- Notes, Code, CreatedAt, UpdatedAt

## 🧪 Testing

```bash
# Ejecutar tests
go test ./...

# Ejecutar tests con coverage
go test -cover ./...

# Ejecutar tests específicos
go test ./controllers -v
```

## 📦 Deployment

### Docker
```bash
# Construir imagen
docker build -t raborimet-crm-backend .

# Ejecutar contenedor
docker run -p 8080:8080 --env-file .env raborimet-crm-backend
```

### Producción
1. Configurar variables de entorno de producción
2. Compilar la aplicación: `go build -o bin/raborimet-crm cmd/main.go`
3. Ejecutar: `./bin/raborimet-crm`

## 🔧 Configuración avanzada

### Variables de entorno importantes

- `PORT`: Puerto del servidor (default: 8080)
- `GIN_MODE`: Modo de Gin (debug/release)
- `DB_*`: Configuración de base de datos
- `JWT_SECRET`: Clave secreta para JWT
- `CORS_ALLOWED_ORIGINS`: Orígenes permitidos para CORS

### Logs

Los logs se guardan en `logs/app.log` y se rotan automáticamente.

## 🤝 Contribución

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 📞 Soporte

Para soporte técnico, contactar a: [support@raborimet.com](mailto:support@raborimet.com)

---

**Raborimet CRM Backend v1.0.0**
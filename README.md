# Raborimet CRM - Sistema de Gestión Integral para Constructora

## Descripción
Sistema CRM profesional diseñado específicamente para la constructora Raborimet, que incluye gestión de clientes, proyectos de construcción, cotizaciones, seguimiento de obras y facturación.

## Arquitectura Técnica

### Frontend
- **Framework**: Angular 17+
- **UI Library**: Angular Material Design
- **Lenguaje**: TypeScript
- **Estilos**: SCSS
- **Estado**: NgRx (para gestión de estado compleja)

### Backend
- **Lenguaje**: Go (Golang)
- **Framework**: Gin
- **Base de Datos**: PostgreSQL
- **Autenticación**: JWT
- **ORM**: GORM

### Estructura del Proyecto
```
raborimet-crm/
├── frontend/          # Aplicación Angular
├── backend/           # API en Golang
├── docs/             # Documentación
└── README.md         # Este archivo
```

## Funcionalidades Principales

### Gestión de Clientes
- Registro y administración de clientes
- Historial de interacciones
- Segmentación de clientes
- Información de contacto completa

### Gestión de Proyectos
- Creación y seguimiento de proyectos de construcción
- Cronogramas y etapas de construcción
- Asignación de recursos
- Control de avance de obras

### Sistema de Cotizaciones
- Generación de cotizaciones detalladas
- Gestión de materiales y costos
- Aprobación de presupuestos
- Conversión a proyectos

### Dashboard Ejecutivo
- Métricas de ventas y proyectos
- Indicadores de rendimiento (KPIs)
- Reportes financieros
- Gráficos y analytics

### Facturación
- Generación de facturas
- Control de pagos
- Estados de cuenta
- Reportes financieros

## Instalación y Configuración

### Prerrequisitos
- Node.js 18+
- Angular CLI 17+
- Go 1.21+
- PostgreSQL 14+
- Docker (opcional)

### Frontend (Angular)
```bash
cd frontend
npm install
ng serve
```

### Backend (Golang)
```bash
cd backend
go mod tidy
go run main.go
```

## Configuración de Base de Datos

### PostgreSQL
1. Crear base de datos: `raborimet_crm`
2. Configurar variables de entorno
3. Ejecutar migraciones

## Variables de Entorno

### Backend (.env)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=raborimet_crm
JWT_SECRET=your-secret-key
PORT=8080
```

### Frontend (environment.ts)
```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api'
};
```

## Desarrollo

### Comandos Útiles

#### Frontend
- `ng serve` - Servidor de desarrollo
- `ng build` - Build de producción
- `ng test` - Ejecutar tests
- `ng lint` - Linter

#### Backend
- `go run main.go` - Ejecutar servidor
- `go test ./...` - Ejecutar tests
- `go build` - Compilar aplicación

## Contribución

1. Fork del proyecto
2. Crear rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## Licencia

Proyecto propietario de Raborimet Constructora.

## Contacto

Para soporte técnico o consultas sobre el sistema, contactar al equipo de desarrollo.
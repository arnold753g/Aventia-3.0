# Andaria - Guía de Docker

Esta guía te ayudará a ejecutar el proyecto Andaria usando Docker.

## Estructura de Archivos Docker

```
.
├── backend/
│   ├── Dockerfile              # Imagen Docker del backend (Go)
│   └── .dockerignore           # Archivos excluidos del build
├── frontend/
│   ├── Dockerfile              # Imagen Docker del frontend (Nuxt)
│   └── .dockerignore           # Archivos excluidos del build
├── nginx/
│   ├── conf.d/
│   │   └── andaria.conf        # Configuración del reverse proxy
│   └── ssl/                    # Certificados SSL (vacío inicialmente)
├── scripts/
│   ├── backup-db.sh            # Script de backup de base de datos
│   ├── backup-uploads.sh       # Script de backup de archivos
│   └── restore-db.sh           # Script de restauración
├── docker-compose.yml          # Configuración para desarrollo
├── docker-compose.prod.yml     # Configuración para producción
└── .env.production             # Plantilla de variables de entorno
```

## Requisitos Previos

- Docker Desktop instalado
- Docker Compose v3.8+
- 4GB RAM mínimo
- 10GB espacio en disco

## Inicio Rápido (Desarrollo)

### 1. Configurar Variables de Entorno

```bash
# Copiar plantilla
cp .env.production .env

# Generar secretos (usar Git Bash o WSL en Windows)
openssl rand -hex 32  # Para JWT_SECRET
openssl rand -hex 32  # Para OTP_SECRET
```

Edita `.env` y configura:
- `DB_PASSWORD`: Contraseña segura para PostgreSQL
- `JWT_SECRET`: Secreto generado arriba
- `OTP_SECRET`: Secreto generado arriba
- `SMTP_USER`: Tu email de Gmail
- `SMTP_PASS`: App Password de Gmail (no tu contraseña normal)
- `SMTP_FROM`: Nombre y email para envíos

### 2. Iniciar Servicios

```bash
# Construir imágenes
docker-compose build

# Iniciar todos los servicios
docker-compose up -d

# Ver logs en tiempo real
docker-compose logs -f
```

### 3. Verificar Servicios

```bash
# Ver estado
docker-compose ps

# Probar backend
curl http://localhost:5750/health

# Probar frontend
# Abre http://localhost:3000 en tu navegador
```

## Comandos Útiles

### Gestión de Servicios

```bash
# Iniciar servicios
docker-compose up -d

# Detener servicios
docker-compose down

# Reiniciar un servicio específico
docker-compose restart backend

# Ver logs
docker-compose logs -f backend      # Backend
docker-compose logs -f frontend     # Frontend
docker-compose logs -f postgres     # Base de datos

# Reconstruir después de cambios en código
docker-compose build backend
docker-compose up -d backend
```

### Base de Datos

```bash
# Conectar a PostgreSQL
docker-compose exec postgres psql -U postgres -d Andaria_01

# Ver tablas
docker-compose exec postgres psql -U postgres -d Andaria_01 -c "\dt"

# Ejecutar consulta
docker-compose exec postgres psql -U postgres -d Andaria_01 -c "SELECT COUNT(*) FROM usuarios;"
```

### Acceso a Contenedores

```bash
# Shell en backend
docker-compose exec backend sh

# Shell en frontend
docker-compose exec frontend sh

# Shell en base de datos
docker-compose exec postgres sh
```

### Limpieza

```bash
# Detener y eliminar contenedores
docker-compose down

# Detener y eliminar contenedores + volúmenes (CUIDADO: Borra datos)
docker-compose down -v

# Limpiar imágenes no usadas
docker system prune -a

# Ver uso de espacio
docker system df
```

## Backups

### Backup Automático

Los scripts están en la carpeta `scripts/`:

```bash
# Backup de base de datos
bash scripts/backup-db.sh

# Backup de archivos subidos
bash scripts/backup-uploads.sh
```

Los backups se guardan en `./backups/` y se rotan automáticamente (se mantienen 7 días).

### Restaurar Backup

```bash
# Restaurar base de datos
bash scripts/restore-db.sh backups/andaria_backup_20250131_020000.sql.gz
```

### Configurar Backups Automáticos (Opcional)

En Linux/WSL, usa cron:

```bash
# Editar crontab
crontab -e

# Agregar backups diarios a las 2 AM
0 2 * * * cd /ruta/al/proyecto && bash scripts/backup-db.sh
0 3 * * * cd /ruta/al/proyecto && bash scripts/backup-uploads.sh
```

## Producción

### 1. Configurar Variables de Entorno

Edita `.env` y actualiza:

```env
# Cambiar a tu dominio
NUXT_PUBLIC_API_BASE=https://tudominio.com/api/v1
NUXT_PUBLIC_WS_BASE=wss://tudominio.com

# URL interna del frontend
FRONTEND_BASE_URL=http://frontend:3000

# Asegurar que APP_ENV es production
APP_ENV=production
NODE_ENV=production
```

### 2. Configurar SSL (Opcional pero Recomendado)

```bash
# Opción 1: Let's Encrypt (Recomendado)
# Necesitas tener el dominio apuntando a tu servidor

# Opción 2: Certificados propios
# Coloca tus certificados en nginx/ssl/
# - nginx/ssl/cert.pem
# - nginx/ssl/key.pem
```

Luego, descomenta la configuración SSL en `nginx/conf.d/andaria.conf`.

### 3. Iniciar en Producción

```bash
# Construir imágenes
docker-compose -f docker-compose.prod.yml build

# Iniciar servicios
docker-compose -f docker-compose.prod.yml up -d

# Ver logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 4. Actualizar Aplicación

```bash
# Obtener últimos cambios
git pull

# Reconstruir y reiniciar
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d
```

## Servicios y Puertos

### Desarrollo (docker-compose.yml)

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| Frontend | 3000 | Aplicación Nuxt (SSR) |
| Backend | 5750 | API Go |
| PostgreSQL | 5432 | Base de datos |

Accesos:
- Frontend: http://localhost:3000
- Backend API: http://localhost:5750/api/v1
- Health Check: http://localhost:5750/health

### Producción (docker-compose.prod.yml)

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| Nginx | 80, 443 | Reverse proxy |
| Frontend | - | Solo acceso interno |
| Backend | - | Solo acceso interno |
| PostgreSQL | - | Solo acceso interno |

Accesos:
- Todo: http://localhost (o tu dominio)
- Health Check: http://localhost/health

## Volúmenes Docker

| Volumen | Propósito |
|---------|-----------|
| `andaria-postgres-data` | Datos de PostgreSQL |
| `andaria-backend-uploads` | Archivos subidos (fotos, comprobantes) |

### Inspeccionar Volúmenes

```bash
# Ver información del volumen
docker volume inspect andaria-postgres-data
docker volume inspect andaria-backend-uploads

# Ver contenido (usando contenedor temporal)
docker run --rm -v andaria-backend-uploads:/data alpine ls -lah /data
```

## Solución de Problemas

### El backend no inicia

```bash
# Ver logs detallados
docker-compose logs backend

# Posibles causas:
# 1. Base de datos no está lista → Esperar a que health check pase
# 2. Variables de entorno faltantes → Verificar .env
# 3. Puerto en uso → Cambiar BACKEND_PORT en .env
```

### El frontend no puede conectarse al backend

```bash
# Verificar que backend está corriendo
curl http://localhost:5750/health

# Verificar variables de entorno
docker-compose exec frontend env | grep NUXT_PUBLIC

# Verificar red Docker
docker network inspect andaria-network
```

### Problemas con la base de datos

```bash
# Ver logs de PostgreSQL
docker-compose logs postgres

# Verificar conexión
docker-compose exec postgres psql -U postgres -d Andaria_01

# Recrear base de datos (CUIDADO: Borra datos)
docker-compose down -v
docker-compose up -d
```

### Los archivos subidos no persisten

```bash
# Verificar volumen
docker volume inspect andaria-backend-uploads

# Verificar permisos
docker-compose exec backend ls -la /app/uploads

# Recrear volumen si es necesario
docker-compose down
docker volume rm andaria-backend-uploads
docker-compose up -d
```

### WebSocket no funciona

```bash
# En desarrollo, verificar directamente
# ws://localhost:5750/api/v1/ws

# En producción, verificar configuración de Nginx
docker-compose -f docker-compose.prod.yml exec nginx nginx -t
```

## Monitoreo

### Ver Recursos

```bash
# Uso de CPU/Memoria en tiempo real
docker stats

# Uso de disco
docker system df

# Logs de Nginx (producción)
docker-compose -f docker-compose.prod.yml logs -f nginx
```

### Health Checks

Todos los servicios tienen health checks configurados:

```bash
# Ver estado de salud
docker-compose ps

# Healthy = ✓
# Unhealthy = ✗
```

## Mejores Prácticas

1. **Nunca** commitear `.env` a git
2. Usar secretos fuertes (32+ caracteres)
3. Hacer backups regulares antes de actualizaciones
4. En producción, usar HTTPS (SSL/TLS)
5. Rotar secretos periódicamente
6. Monitorear logs regularmente
7. Mantener Docker actualizado

## Arquitectura

### Desarrollo

```
[Navegador] → [Frontend:3000] → [Backend:5750] → [PostgreSQL:5432]
```

### Producción

```
[Internet] → [Nginx:80/443] → [Frontend:3000]
                            → [Backend:5750] → [PostgreSQL:5432]
```

## Recursos Adicionales

- Plan completo: `C:\Users\arnol\.claude\plans\crystalline-snacking-salamander.md`
- Documentación Docker: https://docs.docker.com
- Documentación Nginx: https://nginx.org/en/docs/
- PostgreSQL en Docker: https://hub.docker.com/_/postgres

## Soporte

Si encuentras problemas:

1. Revisa los logs: `docker-compose logs -f`
2. Verifica el estado: `docker-compose ps`
3. Consulta esta guía
4. Revisa el plan de implementación completo

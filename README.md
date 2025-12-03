# Ionix

Ionix es un servicio backend estructurado siguiendo Clean Architecture para separar la lógica de negocio, casos de uso y la infraestructura.

## Versión de Go
Requiere Go 1.25 (ver `go.mod` para confirmar la versión exacta utilizada).

## Estructura de carpetas (resumen)
- `cmd/` — puntos de entrada (binaries / servidores).
- `internal/controller/` — handlers HTTP (Gin).
- `internal/usecase/` — casos de uso / servicios (interfaz y lógica de negocio).
- `internal/infra/` — implementaciones de repositorios / acceso a datos.
- `internal/domain/` — modelos y entidades.
- `internal/util/` — utilidades y errores comunes.
- `internal/infra/mock/` — mocks generados con mockery (para tests).
- `.github/workflows/` — pipelines de CI / CD (GitHub Actions).
- `Makefile` — tareas comunes (build, run, test, fmt, mock, tidy).

## Clean Architecture
La app sigue Clean Architecture:
- Entidades en `internal/domain`.
- Casos de uso en `internal/usecase`.
- Controladores/HTTP en `internal/controller`.
- Infraestructura (DB, repos) en `internal/infra`.
Esto facilita pruebas unitarias y mocks en cada capa.

## GitHub Actions (workflows)
El proyecto usa GitHub Actions. Los workflows típicos incluyen:
1. CI (`.github/workflows/ci.yml`):
   - Checkout del repositorio.
   - Setup de Go (versión definida en `go.mod` / workflow).
   - Cache de módulos.
   - `go mod download`.
   - Revision estática: `go fmt`, `go vet`, linter y SAST para estructura y deteccion de vulnerabilidades.
   - Ejecutar tests: `go test ./...`.
   - Compilar artefactos (binaries) para distintos SO/arch si procede.
   - Publicar resultados (informes, badges).
2. Release / CD (`.github/workflows/release.yml`):
   - Se dispara al crear un tag (p.ej. `vX.Y.Z`).
   - Compila binarios para plataformas objetivo.
   - Publica un GitHub Release con los artefactos (binaries).
   - Construye y publica imagen Docker en AWS ECR.
   - Publica la imagen Docker en un cluster AWS ECS 

## Releases
Se crea un release cuando se empuja un tag de versión (semver). El workflow de release compila artefactos y los adjunta al Release en GitHub. Sigue el flujo de tagging para publicar nuevas versiones.

## Requisitos para correr la app
- Go 1.25+
- make (opcional, para targets del Makefile)
- mockery (opcional, para regenerar mocks) — `go install github.com/vektra/mockery/v2/...@latest`
- (Opcional) Delve para debugging en VS Code: `go install github.com/go-delve/delve/cmd/dlv@latest`

Variables de entorno típicas (ajustar según `internal/infra` y configuración real):
- `CONNECTION_STRING_DB` — conexión a la base de datos.

## Cómo ejecutar localmente
1. Instalar dependencias:
   - `go mod tidy`
2. Ejecutar:
   - Con Make: `make run` (si existe target)
   - Directo: `go run ./cmd/...`
3. Ejecutar tests:
   - `make test` o `go test ./... -v`
4. Ver reporte de cobertura
   - `make cover`
5. Ver analisis SAST Scan (deteccion de vulnerabilidades, pii, etc.)
   - `make sast`

## VS Code (arrancar y debug)
- Abrir la carpeta del proyecto en VS Code.
- Instalar la extensión "Go" (golang.go).
- Asegurarse que Go 1.25 está disponible en PATH.
- Instalar Delve (`dlv`) para depuración.
- Añadir/usar configuración en `.vscode/launch.json` apuntando al paquete `main` (p.ej. `${workspaceFolder}/cmd/<app>`), luego ejecutar "Run and Debug" (F5).

Ejemplo mínimo de `launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch ionix (Go)",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/...",
      "env": {
                "CONNECTION_STRING_DB": "USER:PASSWORD@tcp(database-1.ASDGQWEQWE.us-east-2.rds.amazonaws.com:3306)/ionix"
            }
    }
  ]
}
```

## Endpoints implementados (resumen)
Controlador de tareas expone:
- POST /login — login de usuarios (payload JSON -> domain.User)
- POST /api/tasks — crear una tarea (payload JSON -> domain.Task).
- GET /api/tasks/:id — obtener tarea por ID.
- GET /api/tasks — listar todas las tareas.
- PUT /api/tasks/:id — actualizar tarea (payload JSON -> domain.Task) - Not implemented.
- DELETE /api/tasks/:id — eliminar tarea por ID - Not implemented.

Cada endpoint devuelve códigos HTTP estándar:
- 200 OK, 
- 201 Created
- 400 Bad Request
- 401 Not Authorized
- 404 Not Found
- 500 Internal Server Error según el caso.

## Tests y mocks
- Los tests unitarios usan mocks generados por mockery (`internal/infra/mock`, `internal/usecase/mock`).
- Ejecutar tests: `go test ./...`, `make test` o `make cover` para ver reporte de cobertura.


## Curl 

Login 
```
curl --location 'http://my-ecs-fargate-cluster-alb-874497135.us-east-2.elb.amazonaws.com/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "carlos", // admin
    "password": "123" // admin 
}'
```

Create tasks
```
curl --location 'http://my-ecs-fargate-cluster-alb-874497135.us-east-2.elb.amazonaws.com/api/tasks' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ3ODk3NjAsImlkIjoyLCJvcmlnX2lhdCI6MTc2NDc4NjE2MCwicm9sZSI6IiIsInVzZXJuYW1lIjoiY2FybG9zIn0.UIRv-Km9JwtCNnUGPTPfWTHkWeI_w3XXAEQ1zZJ50T8' \
--header 'Content-Type: application/json' \
--data '{
    "title": "postmortem session",
    "due_date": "2025-12-08T00:00:00Z"
}'
```

Listar tasks
```
curl --location 'http://my-ecs-fargate-cluster-alb-874497135.us-east-2.elb.amazonaws.com/api/tasks' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ3ODk3NjAsImlkIjoyLCJvcmlnX2lhdCI6MTc2NDc4NjE2MCwicm9sZSI6IiIsInVzZXJuYW1lIjoiY2FybG9zIn0.UIRv-Km9JwtCNnUGPTPfWTHkWeI_w3XXAEQ1zZJ50T8'
```

Obtener una task
```
curl --location 'http://my-ecs-fargate-cluster-alb-874497135.us-east-2.elb.amazonaws.com/api/tasks/7' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ3ODk3NjAsImlkIjoyLCJvcmlnX2lhdCI6MTc2NDc4NjE2MCwicm9sZSI6IiIsInVzZXJuYW1lIjoiY2FybG9zIn0.UIRv-Km9JwtCNnUGPTPfWTHkWeI_w3XXAEQ1zZJ50T8'
```
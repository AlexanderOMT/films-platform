
This repository holds an authentication service using Golang and Docker

As the repository as this file "README.md" file is still in progress

## Define Environment Variables for Running PostgreSQL with Docker
```env
POSTGRES_USER=""
POSTGRES_PASSWORD=""
POSTGRES_DB=""
EXPOSE_PORT=8000 
HOST=0.0.0.0

## Architecture
(...)

## API documentation
(...)

## Folder structure
(...)

Sanitize the input and the field validation:
https://gorm.io/docs/security.html

sudo docker exec -it films-database psql -U user -d films_platform -f ./tmp/schema.sql

sudo docker exec -it films-service sh -c "sh ./db/init.sh"

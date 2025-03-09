
This repository holds an authentication service using Golang and Docker

As the repository as this file "README.md" file is still in progress

## Define Environment Variables for Running PostgreSQL with Docker
```env
POSTGRES_USER=""
POSTGRES_PASSWORD=""
POSTGRES_DB=""
EXPOSE_PORT=8000 // In the API source code is optional, but docker cannot build if not specified because is used in docker compose file 

## Define Optional Environment Variables for Running PostgreSQL with Docker
HOST=0.0.0.0

## Architecture
(...)

## Folder structure
(...)

consider: Should we keep the connection up all the time? Or just we want to query to DB ?
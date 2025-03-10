
This repository holds an authentication service using Golang and Docker

## Define Environment Variables for Running PostgreSQL with Docker
On the one hand, there are credentials variables related to database, on the other hand, there are variables to set the API address and the secret key for signing the JWT Token (singning method: HMAC-SHA256).
```env
POSTGRES_USER= 
POSTGRES_PASSWORD=

POSTGRES_DB=films-database

SECRET_KEY=""
EXPOSE_PORT=8000 
HOST=0.0.0.0
```

## Architecture

The architucture proposed for the implementation is an Hexagonal architucture pattern, which each layer should have a specific role:

1. Router Layer (router, routes)
2. Handler Layer (Controllers - auth_handler, film_handler)
3. Middleware Layer (auth.go)
4. Use Case Layer (auth_service, film_service, etc.)
5. Infrastructure Layer (database, repo implementations)
6. Domain Layer (film, user, repositories)

## API documentation

## Notes:
- Requiere username and password validation
- Replace the token for the JWT

## User Endpoints (Unprotected Routes)

| **Method** | **Endpoint**   | **Description**       | **Response**             | **JWT Authentication** | **Parameters**            | **Example**                                                                                                                                     |
|------------|----------------|-----------------------|--------------------------|------------------------|---------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| GET        | /users         | Get all users          | List of users            | No                     | None                      | `curl -X GET http://localhost:8000/users -H "Content-Type: application/json"`                                                                  |
| POST       | /register      | Register a new user    | User registration        | No                     | `username` and `password` | `curl -X POST http://localhost:8000/register -H "Content-Type: application/json" -d '{"username": "Testuser123", "password": "Password123"}'`  |
| POST       | /login         | User login             | JWT Token                | No                     | `username` and `password` | `curl -X POST http://localhost:8000/login -H "Content-Type: application/json" -d '{"username": "Testuser123", "password": "Password123"}'`     |

## Film Endpoints (Protected Routes)

| **Method** | **Endpoint**   | **Description**                       | **Response**               | **JWT Authentication** | **Example**                                                                                                                  |
|------------|----------------|---------------------------------------|----------------------------|------------------------|------------------------------------------------------------------------------------------------------------------------------|
| POST       | /film          | Create a new film                     | Created Film               | Yes                    | `curl -X POST http://localhost:8000/film -H "Content-Type: application/json" -H "Authorization: token" -d '{"Title": "Title", "Director": "Director", "Release": "2000-01-01T00:00:00Z"}'` |
| GET        | /films         | Get all films                         | List of Films              | Yes                    | `curl -X GET http://localhost:8000/films -H "Content-Type: application/json" -H "Authorization: token"`                       |
| PATCH      | /film          | Update a specific film                | Updated Film               | Yes                    | `curl -X PATCH "http://localhost:8000/film?title=Title" -H "Content-Type: application/json" -d '{"Director": "newDirector"}' -H "Authorization: token"` |
| PUT        | /film          | Replace an existing film              | Replaced Film              | Yes                    | `curl -X PUT "http://localhost:8000/film?title=Title" -H "Content-Type: application/json" -d '{"Title": "newTitle", "Director": "newTitle", "Release": "2000-01-01T00:00:00Z"}' -H "Authorization: token"` |
| DELETE     | /film          | Delete a film                         | Film deleted               | Yes                    | `curl -X DELETE http://localhost:8000/film?title=newTitle -H "Content-Type: application/json" -H "Authorization: token"`      |

## Starting the API

First, clone the repository to your local machine using the method of your choice, e.g.:
```bash
git clone https://github.com/AlexanderOMT/films-platform.git
git clone git@github.com:AlexanderOMT/films-platform.git
```

Once you have cloned the repository, to start the containers defined in docker-compose.yml, run:
```bash
sudo docker compose -f docker-compose.yml build
sudo docker compose -f docker-compose.yml up
```
If you want to clean, rebuild, and restart the containers, these commands are useful and are part of a full restart workflow.
```bash
sudo docker compose -f docker-compose.yml down
sudo docker compose -f docker-compose.yml build
```

### Run SQL to Initialize the Database

First of all, you need to execute the SQL schema to create the tables on the `films_platform` database, use the following command:
```bash
sudo docker exec -it films-database psql -U user -d films_platform -f ./tmp/schema.sql
```
Secondly, you want to fill some sample data using a very simple script
```bash
sudo docker exec -it films-service sh -c "sh ./db/init.sh"
```

If you need to drop (or use `truncate` for cleaning the table) the films and users tables, run the following commands:
```bash
sudo docker exec -it films-database psql -U user -d films_platform -c "DROP TABLE films;"
sudo docker exec -it films-database psql -U user -d films_platform -c "DROP TABLE users;"
```
## Sanitize the Input and Field Validation

To follow best security practices, it's recommended to sanitize inputs and validate fields when working with databases. Using GORM can implement many of them, more details can be found [here](https://gorm.io/docs/security.html).

To validate and apply constraints, it's used json go-playground/validator [docs](https://github.com/go-playground/validator) to validate the json fields when are required for the routes and gorilla/schemma [docs](https://github.com/gorilla/schema) to validate the query paramaters. [Here can find some style patterns](https://www.reddit.com/r/golang/comments/rzvla7/stylepattern_for_go_database_dto_validation_struct/)

## Git conventions
Git usage in this repository is based on [this guide] (https://slicer.readthedocs.io/en/latest/developer_guide/style_guide.html). A brief summary is the followoing:

- **BUG**: A change made to fix a runtime issue  
  _(e.g., exception, or incorrect result)_

- **COMP**: A fix for a compilation issue, error, or warning  

- **ENH**: New functionality added to the project  

- **STYLE**: A change that does not impact the logic or execution of the code  
  _(e.g., improving coding style, comments in source code)_

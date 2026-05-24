User registration and login service. It uses PostgreSQL for storage and JWT for authentication tokens.
For now only serves to survey-forms

Overview
--------

- HTTP server built with `chi`
- PostgreSQL database with automatic schema initialization on startup
- JWT tokens with `email`, `user_id`, and `role` claims
- Email-domain allowlist for registrations
- Password hashing with bcrypt

-------------

environment variables list

- `DATABASE_URL` - PostgreSQL connection string. Must start with `postgres://` or `host=`
- `PORT` - HTTP listen address, for example `:8081`
- `JWT_SECRET` - HMAC secret used to sign and validate tokens
- `JWT_ISSUER` - expected JWT issuer
- `JWT_AUDIENCE` - expected JWT audience
- `ALLOWED_EMAIL_DOMAINS` - comma-separated list of allowed registration domains

Example `.env` values are provided in [.env.example](.env.example).

Docker
------

The repository includes [docker-compose.yml](docker-compose.yml) for PostgreSQL only. It creates a `users_service` database with user `users_app` and password `users_pass` on port `5432`.

```bash
docker compose up -d
```

------

All endpoints accept and return JSON unless noted otherwise.

### Health

- `GET /health`
- Returns `OK`

### Register

- `POST /register`
- Body:

```json
{
	"email": "user@example.com",
	"password": "strongPassword123"
}
```

- Validation rules:
	- Email must be valid and its domain must be listed in `ALLOWED_EMAIL_DOMAINS`
	- Password must be 8 to 72 characters long
	- Password must contain ASCII characters only

Successful registration returns HTTP 200 with an empty body.

### Login

- `POST /login`
- Body:

```json
{
	"email": "user@example.com",
	"password": "strongPassword123"
}
```

- Response example:

```json
{
	"message": "logged in",
	"token": "<jwt>",
	"email": "user@example.com",
	"user_id": "<uuid>",
	"role": "user",
	"expires": "2026-05-13T10:00:00Z"
}
```

The token is valid for 12 hours and is signed with `JWT_SECRET`.

Authentication
--------------

Send the JWT in the `Authorization` header when calling protected routes introduced in future revisions:

```bash
Authorization: Bearer <token>
```

Current storage model
---------------------

On startup the service creates the `users` table if it does not already exist. The table includes:

- `id` as a UUID primary key
- `email` as a unique field
- `password_hash` for bcrypt hashes
- `role` with a default value of `user`
- `created_at` timestamp

Project structure
-----------------

- `main.go` - application entrypoint and route setup
- `internal/config` - environment loading and parsing
- `internal/handlers` - HTTP handlers for login and registration
- `internal/repository` - PostgreSQL access and schema creation
- `internal/dto` - request payload types
- `internal/validations` - email and password validation helpers
- `internal/models` - domain model placeholder
- `pkg/auth` - JWT creation, validation, and password hashing

Testing
-------

Run the test suite with:

```bash
go test ./...
```

Notes
-----

- The current HTTP surface only exposes `GET /health`, `POST /register`, and `POST /login`.
- `main.go` loads `.env` automatically if present.
- The service exits on startup if required configuration is missing.

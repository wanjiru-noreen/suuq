# Suuq Architecture

## Authentication Configuration

The backend creates the `users` table during startup and stores the SQLite database at the path set by `DB_PATH`. Docker Compose uses `/data/suuq.db` on the persistent `suuq-data` volume.

Set `AUTH_SECRET` to a strong secret in deployed environments. If it is not set, the backend uses a development-only fallback and should not be exposed publicly.

### Authentication endpoints

```text
POST /api/auth/register  Create a user and return a bearer token
POST /api/auth/login     Authenticate a user and return a bearer token
GET  /api/auth/me        Return the authenticated user
```

Send the token returned by registration or login with protected requests:

```text
Authorization: Bearer <token>
```

Example registration request:

```bash
curl -X POST http://localhost:8080/api/auth/register \
	-H "Content-Type: application/json" \
	-d '{"name":"Amina","email":"amina@example.com","password":"password123"}'
```

Passwords must contain at least eight characters. Email addresses are normalized to lowercase and duplicate addresses are rejected.

# Suuq Architecture

## Authentication Configuration

The backend creates the `users` table during startup and stores the SQLite database at the path set by `DB_PATH`. Docker Compose uses `/data/suuq.db` on the persistent `suuq-data` volume.

`AUTH_SECRET` is required in all environments and must contain at least 32 non-padding bytes. Generate a random value with `openssl rand -hex 32`. Startup fails for a missing or short secret. There is no fallback. Changing the secret invalidates all existing tokens.

Authentication middleware lives in `backend/internal/auth/auth.go`; register protected routes with `authService.Middleware`. Future business handlers must additionally enforce ownership using `auth.UserID(r.Context())`.

[SQLite requires foreign keys to be enabled per connection](https://www.sqlite.org/foreignkeys.html). SQLite connections enable foreign keys and a five-second busy timeout through the driver connection string. The pool uses one connection to serialize writes. Startup migrations run atomically; the current create-if-absent scripts do not upgrade existing table definitions. Introduce versioned migrations when changing the schema.

## Deployment

- Set `AUTH_SECRET` through the deployment secret manager or an uncommitted `.env` file for Compose. The Go backend automatically loads the first `.env` found in its working directory or its parent directory, supporting local runs from `backend/`. Existing environment variables take precedence, including explicitly empty values. Missing files are allowed; unreadable or malformed files stop startup without logging their contents. Relative `DB_PATH` values are resolved against the working directory.
- Terminate HTTPS at a trusted reverse proxy before exposing the frontend. The supplied Compose configuration is a local HTTP setup, not a complete public deployment.
- The backend port binds to loopback on the host. Route public API requests through Nginx, which limits login and registration to five requests per minute per client IP with a burst of five and returns HTTP 429 for excess requests. Direct access to port 8080 bypasses this limit. When placing another proxy before Nginx, configure trusted real-IP handling so limits apply to clients rather than the proxy.
- Keep the database volume private and configure SQLite-aware backups with a tested restore procedure.
- Use maintained, patched Go and container images and perform dependency vulnerability scanning before release.
- The HTTP server bounds request/header time, limits headers, and allows ten seconds for graceful shutdown on SIGINT/SIGTERM. Authentication JSON bodies are capped at 16 KiB, and authentication responses disable caching.
- Tokens expire after 24 hours. Individual session revocation and password reset are not implemented; account-wide logout currently requires a future session design.


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

Passwords must contain 8–72 bytes (the bcrypt input limit). Email addresses are normalized to lowercase and duplicate addresses are rejected.

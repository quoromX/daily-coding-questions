# DaoForge Backend

Go API for DaoForge with standard-library routing, in-memory demo persistence, seeded problems, auth tokens, dashboard metrics, submissions, and Judge0 integration hooks.

Run from this directory:

```bash
env GOCACHE=/tmp/go-cache go run ./cmd/api
```

Then check:

```bash
curl http://localhost:8080/healthz
```

Core routes:

- `POST /v1/auth/register`
- `POST /v1/auth/login`
- `GET /v1/me`
- `GET /v1/me/dashboard`
- `GET /v1/problems`
- `GET /v1/problems/{slug}`
- `POST /v1/problems/{slug}/run`
- `POST /v1/problems/{slug}/submit`
- `GET /v1/submissions`

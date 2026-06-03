# DaoForge

DaoForge is a Dao and cultivation themed coding challenge platform inspired by algorithm practice platforms and progression-driven learning. The app uses Next.js for the frontend, Go for the backend API, Postgres for persistence, and a self-hosted Judge0 instance for code execution.

Start with the editable docs:

- [Project Plan](docs/PROJECT_PLAN.md)
- [Style Guide](docs/STYLE_GUIDE.md)
- [Judge0 Setup](docs/JUDGE0_SETUP.md)

The repository has been seeded with the intended monorepo shape so implementation can begin after the plan is approved.

## Current Implementation

The MVP scaffold now includes:

- A Go API with auth, seeded problems, dashboard data, submissions, and Judge0 execution hooks.
- A Postgres migration describing the production schema.
- A Next.js app with cultivation-themed pages for landing, login, register, dashboard, manual catalog, problem workbench, submissions, profile, settings, and admin manual drafts.
- Docker Compose for local Postgres, backend, and frontend development.

## Local Development

Backend:

```bash
cd backend
env GOCACHE=/tmp/go-cache go run ./cmd/api
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

Full local stack:

```bash
docker compose up
```

The frontend currently uses local demo data so the UI is usable immediately. The Go API is ready for the frontend to switch from local demo data to `NEXT_PUBLIC_API_URL` calls when you want the next integration pass.

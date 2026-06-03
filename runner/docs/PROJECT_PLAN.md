# DaoForge Project Plan

## Product Direction

DaoForge is a code algorithm platform in the spirit of LeetCode, Codewars, and Boot.dev, but with its own themed identity: a Dao and cultivation sect where users refine algorithmic qi by solving manuals and passing breakthrough attempts.

The working theme is **DaoForge: The Path of Algorithmic Cultivation**. It combines cultivation progression, sect-library organization, and calm ink-wash UI. The product should feel memorable and focused enough for daily coding practice.

## Core Users

- Learners who want structured algorithm practice.
- Interview prep users who want measurable progress.
- Competitive programmers who want speed, accuracy, and history.
- Future community users who may author challenges, compare rankings, and join clans/guilds.

## MVP Goals

- Account creation, login, logout, and session refresh.
- User dashboard with solved problems, streak, rank, recent attempts, language usage, and difficulty spread.
- Problem catalog with filters by difficulty, tags, status, and search.
- Problem detail page with prompt, examples, constraints, hints, discussion placeholder, and code editor.
- Code execution through self-hosted Judge0.
- Submission result display with pass/fail, runtime, memory, stdout, stderr, compile errors, and per-test feedback.
- Submission history per problem and user.
- Admin seed path for creating initial problems and test cases.

## Proposed Tech Stack

- Frontend: Next.js App Router, TypeScript, Tailwind CSS, shadcn-style primitives where useful, Monaco Editor for code editing.
- Backend: Go, chi or Fiber for HTTP routing, pgx for Postgres, goose or golang-migrate for migrations, sqlc if we want typed SQL.
- Database: Postgres.
- Code execution: Judge0 CE or Extra CE self-hosted on DigitalOcean.
- Auth: Email/password with Argon2id password hashing, JWT access tokens, refresh token rotation stored server-side.
- Deployment: Docker Compose for local development; production can later split app, API, Postgres, and Judge0 across droplets/services.

## Repository Shape

```text
.
├── backend/
│   ├── cmd/api/                 # API entrypoint
│   ├── internal/auth/           # registration, login, sessions, passwords
│   ├── internal/config/         # env loading and validation
│   ├── internal/database/       # db pool, queries, repositories
│   ├── internal/judge/          # Judge0 client and result mapping
│   ├── internal/metrics/        # solved counts, streaks, language stats
│   ├── internal/platform/       # shared HTTP middleware/errors
│   ├── internal/problems/       # problem catalog and test cases
│   ├── internal/submissions/    # run/submit lifecycle
│   ├── internal/users/          # user profiles and rank progression
│   └── migrations/              # SQL migrations
├── frontend/                    # Next.js application
├── infra/judge0/                # deployment notes/config overrides
├── docs/                        # product/design/setup docs
└── scripts/                     # local automation
```

## Backend API Plan

### Auth

- `POST /v1/auth/register`
- `POST /v1/auth/login`
- `POST /v1/auth/refresh`
- `POST /v1/auth/logout`
- `GET /v1/me`

### Problems

- `GET /v1/problems`
- `GET /v1/problems/{slug}`
- `GET /v1/problems/{slug}/starter-code`
- `POST /v1/admin/problems`
- `PUT /v1/admin/problems/{id}`

### Code Execution

- `POST /v1/problems/{slug}/run`
- `POST /v1/problems/{slug}/submit`
- `GET /v1/submissions/{id}`
- `GET /v1/problems/{slug}/submissions`
- `GET /v1/me/submissions`

### Dashboard

- `GET /v1/me/dashboard`
- `GET /v1/me/stats`
- `GET /v1/me/activity`

## Data Model Draft

- `users`: account, display name, email, password hash, rank, xp, created time.
- `refresh_tokens`: hashed token, user id, expiry, revocation metadata.
- `problems`: title, slug, difficulty, status, body, constraints, editorial status.
- `problem_tags`: tag names and relationships.
- `starter_code`: language id, language name, template.
- `test_cases`: visible and hidden tests, stdin, expected output, weighting.
- `submissions`: user, problem, language, source, status, score, runtime, memory.
- `submission_results`: per-test status, stdout, stderr, compile output, Judge0 token.
- `daily_activity`: solved counts and streak support.

## Frontend Pages

- `/`: authenticated users go to dashboard; guests see the product entry and login/register actions.
- `/register`: account creation.
- `/login`: login.
- `/dashboard`: metrics, solved history, streak, recent attempts, recommended next trials.
- `/problems`: searchable challenge board.
- `/problems/[slug]`: problem prompt, editor, run/submit panel, test output.
- `/submissions`: personal submission archive.
- `/profile/[handle]`: public progress profile.
- `/settings`: account and editor preferences.
- `/admin/problems`: later-stage problem authoring.

## Theme Concept

The visual world should be a **cultivation sect library**: misty paper surfaces, pine and jade navigation, gold progress accents, soft depth, ink-dark code surfaces, and realm badges. It should reference Dao/cultivation themes through naming, hierarchy, badges, and progression without becoming decorative clutter.

Suggested progression language:

- Novice: Outer Disciple
- Easy track: Foundation Realm
- Medium track: Qi Condensation
- Hard track: Core Formation
- Streaks: Meditation Streak
- XP: Spirit Stones
- Ranking: Realm
- Problem sets: Manuals
- Dashboard: Sect Hall

## Implementation Phases

### Phase 1: Foundation

- Create Next.js app and Go API scaffolds.
- Add Docker Compose for app, API, Postgres, and Redis if needed.
- Add migrations for users, problems, submissions, and dashboard stats.
- Add auth and basic protected routes.

### Phase 2: Problem Solving Loop

- Add problem catalog and problem detail UI.
- Add Monaco editor and language selector.
- Add Judge0 client in Go.
- Implement run and submit flows.
- Store submission results.

### Phase 3: Dashboard and Progression

- Add solved metrics, streaks, XP, ranks, and recent activity.
- Add submission archive.
- Add profile page.

### Phase 4: Authoring and Polish

- Add admin problem creation.
- Add seed challenges.
- Add editorial/hints support.
- Add refined mobile states, loading states, empty states, and error states.

## Decisions To Confirm Before Build

- Final product name: `DaoForge` or another cultivation name.
- Theme balance: more serene Dao sect, more fantasy xianxia, or a clean hybrid.
- Go router preference: `chi`, `Fiber`, or standard library.
- SQL preference: handwritten pgx queries, `sqlc`, or ORM.
- Judge0 flavor: CE or Extra CE.
- Whether submissions should run all hidden tests synchronously for MVP or use async polling/webhooks from the start.

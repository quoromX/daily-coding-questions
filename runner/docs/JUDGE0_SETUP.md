# Judge0 Self-Hosted Setup and App Integration

This guide is for running Judge0 on a DigitalOcean Droplet and connecting it to DaoForge through the Go backend.

Sources checked on June 2, 2026:

- Judge0 docs: https://docs.judge0.com/products/judge0/get_started/self_hosting/
- Judge0 GitHub: https://github.com/judge0/judge0
- Judge0 releases: https://github.com/judge0/judge0/releases
- DigitalOcean Docker 1-Click docs: https://docs.digitalocean.com/products/marketplace/catalog/docker/
- Docker Compose docs: https://docs.docker.com/compose/

## Recommended Deployment Shape

For the MVP, keep Judge0 separate from the app/API server:

- Droplet A: DaoForge frontend, Go API, Postgres, reverse proxy.
- Droplet B: Judge0 CE or Extra CE.

This keeps untrusted code execution isolated from user data and the main application. If budget is tight during development, everything can run locally or on one private test droplet, but production should isolate Judge0.

## DigitalOcean Droplet

Recommended starting point:

- Ubuntu 22.04 or newer.
- Docker 1-Click Droplet, because DigitalOcean says it includes Docker Engine and Docker Compose.
- At least 2 vCPU / 4GB RAM for experimentation.
- Prefer 4 vCPU / 8GB RAM or higher if multiple users may submit concurrently.
- Add swap only as a safety net; do not rely on swap for normal judging.

After creating the droplet:

```bash
ssh root@YOUR_DROPLET_IP
docker version
docker compose version
```

## Install Judge0

Use the latest stable Judge0 CE release from the official GitHub releases page. At the time of checking, public release snippets still point to the v1.13.1 archive pattern.

Example:

```bash
apt update
apt install -y unzip wget
wget https://github.com/judge0/judge0/releases/download/v1.13.1/judge0-v1.13.1.zip
unzip judge0-v1.13.1.zip
cd judge0-v1.13.1
docker compose up -d db redis
sleep 10
docker compose up -d
```

Then test:

```bash
curl http://YOUR_DROPLET_IP/languages
```

For Extra CE, use the matching `-extra` release archive if you need the expanded language set.

## HTTPS and Reverse Proxy

For production, put Judge0 behind HTTPS and restrict who can call it.

Suggested approach:

- Point `judge.yourdomain.com` to the Judge0 droplet.
- Use Caddy, Traefik, or Nginx with Let's Encrypt.
- Allow inbound HTTP/HTTPS only from the public internet if the API needs to call by domain, or restrict to the app server IP if possible.
- Do not expose database or Redis ports.

If the official Judge0 HTTPS archive is used, edit its compose variables for:

- `VIRTUAL_HOST`
- `LETSENCRYPT_HOST`
- `LETSENCRYPT_EMAIL`

## Security Notes

Judge0 executes untrusted code. Treat it as hostile infrastructure.

- Keep Judge0 on its own droplet or private network.
- Patch the host regularly.
- Do not store app secrets on the Judge0 server.
- Rate-limit submissions in the Go API before sending to Judge0.
- Store only necessary source/submission data.
- Keep Judge0 unavailable from browsers if possible; frontend should call the Go API, and the Go API should call Judge0.
- Consider adding an API gateway or reverse proxy auth in front of Judge0 if the self-hosted edition is reachable publicly.

## App Environment

In DaoForge `.env`:

```env
JUDGE0_BASE_URL=https://judge.yourdomain.com
JUDGE0_AUTH_TOKEN=
JUDGE0_WEBHOOK_SECRET=change_me_webhook_secret
```

If Judge0 is only reachable over a private network, use the private URL/IP:

```env
JUDGE0_BASE_URL=http://10.0.0.5
```

## Backend Integration Plan

The frontend never calls Judge0 directly. The flow should be:

1. User submits code from Next.js.
2. Next.js calls the Go API.
3. Go API validates auth, problem slug, language, source length, and rate limits.
4. Go API loads visible or hidden test cases from Postgres.
5. Go API creates Judge0 submissions.
6. Go API polls Judge0 or receives webhook callbacks.
7. Go API stores normalized results in Postgres.
8. Frontend reads the stored result from the Go API.

## Judge0 Request Shape

For a single run:

```http
POST /submissions?wait=true
Content-Type: application/json

{
  "language_id": 71,
  "source_code": "print(input())",
  "stdin": "hello",
  "expected_output": "hello"
}
```

For hidden tests, the backend can submit one Judge0 request per test case for MVP. Later, optimize with batching or additional files when needed.

## Result Mapping

Store Judge0 details but expose a product-friendly status:

| Judge0 Result | DaoForge Status |
| --- | --- |
| Accepted | Passed |
| Wrong Answer | Failed |
| Time Limit Exceeded | Timeout |
| Compilation Error | Compile Error |
| Runtime Error | Runtime Error |
| Internal Error | Judge Error |

## MVP Polling Strategy

Start with synchronous `wait=true` for `Run Code` on visible tests. For `Submit`, prefer async:

- Create submissions.
- Save local submission with `pending` status.
- Poll Judge0 from the API for a short window.
- If still pending, let frontend poll `GET /v1/submissions/{id}`.

Webhooks can be added once the core loop works.

## Local Development

For local work, either:

- Run Judge0 locally with Docker Compose.
- Point `JUDGE0_BASE_URL` to the DigitalOcean Judge0 URL.

The Go API should fail fast on startup if `JUDGE0_BASE_URL` is missing in environments where submissions are enabled.

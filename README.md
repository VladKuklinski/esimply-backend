# eSIMply Backend

REST API backend for the eSIMply iOS app. Serves eSIM plan data for European destinations and proxies AI chat requests to Claude.

## Local Setup

**Prerequisites:** Go 1.22+

```bash
# Clone and enter the project
cd esimply-backend

# Set your Anthropic API key
export ANTHROPIC_API_KEY=sk-ant-...

# Run
go run .
```

The server starts on `http://localhost:8080`.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/countries` | All 15 countries |
| GET | `/countries/{id}/plans` | Plans for a country (e.g. `france`, `czech-republic`) |
| POST | `/ai/chat` | AI plan recommendation |

**AI chat example:**
```bash
curl -X POST http://localhost:8080/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "I am going to Paris for 10 days, mostly sightseeing"}'
```

## Railway Deployment

1. Push this repo to GitHub.
2. Create a new project on [Railway](https://railway.app) and connect the repo.
3. Add environment variable: `ANTHROPIC_API_KEY=sk-ant-...`
4. Railway auto-detects Go and deploys. The `PORT` variable is set automatically.

No `Dockerfile` or `Procfile` needed — Railway's nixpacks builder handles Go projects natively.

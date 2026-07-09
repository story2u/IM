# Opportunity IM Assistant Backend

FastAPI backend for Telegram and WeCom opportunity detection, human review, and after-hours AI replies.

## Run locally with Docker

```bash
cd backend
cp .env.example .env
docker compose build
docker compose up -d postgres redis
docker compose run --rm migrate
docker compose run --rm api python scripts/seed_demo.py
docker compose up api celery_worker celery_beat
```

API docs: <http://localhost:8000/docs>

Admin endpoints require:

```http
Authorization: Bearer change-me
```

## Main frontend-facing endpoints

- `GET /api/v1/opportunities`
- `GET /api/v1/opportunities/{id}`
- `GET /api/v1/messages?opportunity_id={id}`
- `POST /api/v1/opportunities/{id}/manual-reply`
- `POST /api/v1/opportunities/{id}/ai-draft`
- `GET /api/v1/templates`
- `GET /api/v1/configs/work-mode`
- `GET /api/v1/stats/summary`

## Webhooks

- `POST /api/v1/webhooks/telegram`
- `GET /api/v1/webhooks/wecom`
- `POST /api/v1/webhooks/wecom`

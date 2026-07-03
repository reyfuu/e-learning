# Setup Learn Platform

Semua kode ada di folder **`apps/`** (satu-satunya folder proyek di root repo).

## Struktur

```
learn-platform/          # root repo (hanya apps/ + README singkat)
└── apps/
    ├── backend/         # Go API — Echo, GORM, PostgreSQL
    ├── frontend/        # React + Vite
    ├── docker-compose.yml
    ├── Makefile
    ├── docs/
    └── scripts/
```

## PostgreSQL

| Item     | Nilai            |
| -------- | ---------------- |
| User     | `audi`           |
| Password | `090393`         |
| Database | `learn_platform` |

```
postgres://audi:090393@localhost:5432/learn_platform?sslmode=disable
```

```bash
cd apps
make db-setup-local   # PostgreSQL di mesin lokal
make db-docker        # PostgreSQL via Docker
make db-psql          # shell: psql -U audi -d postgres
```

Atau langsung:

```bash
psql -U audi -d postgres
```

Untuk database aplikasi (`learn_platform`):

```bash
psql -U audi -d learn_platform
```

## Menjalankan

```bash
cd apps
make setup
make backend-dev    # http://localhost:8080/health
make frontend-dev   # http://localhost:3000
```

## Environment

| File | Variabel |
| ---- | -------- |
| `backend/.env` | `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGIN` |
| `frontend/.env.local` | `VITE_API_URL` |

## API

| Method | Path |
| ------ | ---- |
| GET | `/health` |
| POST | `/api/v1/auth/register` |
| POST | `/api/v1/auth/login` |
| GET | `/api/v1/users/me` |

# LearnHub

LearnHub adalah platform e-learning. Workspace menggunakan pnpm untuk aplikasi web Next.js; API tetap memakai Go modules.

## Struktur

    apps/web/   Next.js web app (Vercel)
    apps/api/   Go API (container hosting)
    docs/       BRD, FRD, TRD, PRD

## Menjalankan lokal

    pnpm install
    pnpm dev       # web di http://localhost:3000
    pnpm api:dev   # API di http://localhost:8080; memerlukan PostgreSQL

Untuk menjalankan keduanya, gunakan pnpm dev:all. Konfigurasikan apps/api/.env dari .env.example dan arahkan database lokal.

## Build dan deploy

    pnpm build
    pnpm lint

Atur Vercel Root Directory ke apps/web. Deploy API Go sebagai container terpisah dan set URL API pada environment web. Lihat DESIGN.md, PRD_eLearning.md, dan dokumen di docs/.

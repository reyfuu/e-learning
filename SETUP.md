# Setup LearnHub

## Struktur

    learn-2/
    ├── apps/web/   # Next.js, dikelola pnpm dan dideploy ke Vercel
    ├── apps/api/   # Go API, dikelola Go modules dan dideploy sebagai container
    ├── docs/       # Dokumen produk dan teknis
    └── package.json / pnpm-workspace.yaml

## Prasyarat

- Node.js versi yang didukung Next.js dan pnpm 10
- Go 1.22+
- PostgreSQL sesuai konfigurasi backend

## Instalasi dan development

    pnpm install
    cp apps/api/.env.example apps/api/.env
    pnpm dev
    pnpm api:dev

Sesuaikan DATABASE_URL dan JWT_SECRET di apps/api/.env. Jalankan web dan API di terminal terpisah. pnpm dev:all menjalankan keduanya bersama-sama. API memerlukan PostgreSQL aktif.

## Vercel

Hubungkan repository ke Vercel, pilih Root Directory apps/web, dan set environment URL API untuk setiap environment. Go API tidak berjalan sebagai proses Go jangka panjang di Vercel; deploy API sebagai container terpisah.

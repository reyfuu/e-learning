# 📐 System Design Document — LearnHub

**Versi:** 1.0.0
**Status:** Berdasarkan PRD_eLearning.md

Dokumen ini merinci arsitektur teknis, stack teknologi, skema database, dan desain API untuk platform LearnHub.

## 1. Arsitektur Sistem

Platform LearnHub menggunakan arsitektur modular yang memisahkan frontend dan backend service, berkomunikasi melalui RESTful API.

```text
┌─────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                         │
│   Next.js 14 (App Router) + Tailwind CSS + shadcn/ui   │
│   React Query · Zustand · React Hook Form               │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTPS / REST + WebSocket
┌──────────────────────▼──────────────────────────────────┐
│                    API GATEWAY                          │
│              Nginx (Reverse Proxy + TLS)                │
└──────┬───────────────┬────────────────┬─────────────────┘
       │               │                │
┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
│  Auth       │ │  Course     │ │  Payment    │
│  Service    │ │  Service    │ │  Service    │
│  (Go)       │ │  (Go)       │ │  (Go)       │
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │               │                │
       └───────────────▼────────────────┘
                       │
       ┌───────────────▼────────────────┐
       │         PostgreSQL 16          │
       │  (Primary + Read Replicas)     │
       └────────────────────────────────┘
                       │
       ┌───────────────▼────────────────┐
       │    Supporting Services         │
       │  Redis (Cache + Sessions)      │
       │  MinIO / S3 (Media Storage)    │
       │  SMTP / Mailgun (Email)        │
       └────────────────────────────────┘
```

## 2. Tech Stack

### Frontend
- **Framework:** Next.js 14 (App Router)
- **Styling:** Tailwind CSS 3.x + shadcn/ui
- **State Management:** Zustand (Local/Global state), TanStack Query v5 (Server state)
- **Form Handling:** React Hook Form + Zod
- **Editor & Video:** TipTap (Rich Text), Video.js / HTML5 Video Player

### Backend
- **Bahasa:** Go 1.22
- **Framework HTTP:** Gin / Echo
- **Database / ORM:** PostgreSQL 16 dengan sqlc + pgx
- **Authentication:** JWT (golang-jwt) + bcrypt
- **Task Queue:** Asynq (Redis-based)

## 3. Skema Database (Utama)

- **users:** Menyimpan kredensial dan profil user (learner, instructor, admin).
- **courses:** Data kursus yang dibuat oleh instruktur.
- **sections & lessons:** Kurikulum kursus (materi video, teks, kuis).
- **enrollments:** Relasi pendaftaran user pada suatu kursus.
- **lesson_progress:** Tracking posisi dan status penyelesaian pelajaran per user.
- **transactions:** Pencatatan pembayaran dan fee platform.
- **reviews, categories, certificates:** Entitas pendukung lainnya.

*(Untuk skema SQL lengkap, lihat PRD_eLearning.md Section 6).*

## 4. Struktur API (REST)

- **Auth API (`/api/v1/auth/*`):** Register, Login, Refresh token, Reset password, Google OAuth.
- **Course API (`/api/v1/courses/*`):** List kursus, detail, pembuatan materi, ulasan.
- **Enrollment API (`/api/v1/enrollments/*` & `/api/v1/lessons/*`):** Pendaftaran kursus dan update progres belajar.
- **Payment API (`/api/v1/payments/*`):** Checkout transaksi dan webhook dari gateway pembayaran.
- **Admin API (`/api/v1/admin/*`):** Dashboard analitik, moderasi kursus, manajemen user.

## 5. Non-Functional Requirements (NFR)
- **Performa:** P95 response time ≤ 200ms, Web TTI ≤ 3 detik.
- **Ketersediaan & Skalabilitas:** 99.5% uptime, stateless backend, redis caching, mendukung 10.000 concurrent users (fase awal).
- **Keamanan:** HTTPS, JWT token expiry pendek, SQLi/XSS prevention.

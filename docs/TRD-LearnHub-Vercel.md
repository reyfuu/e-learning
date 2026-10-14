# TRD — LearnHub

**Versi:** 1.3 · **Tanggal:** 15 Oktober 2026 · **Status:** Draft target teknis

## 1. Prinsip arsitektur

- **Target:** Next.js App Router menjadi web app yang berjalan di Vercel; versi framework yang benar-benar dipakai dicatat terpisah dari target versi pada PRD induk.
- **Target:** Go 1.22 + Echo menjadi API modular; worker/background jobs berjalan sebagai proses container terpisah.
- Vercel tidak diperlakukan sebagai host proses Go yang selalu hidup. Aplikasi Next.js berbicara ke API Go melalui HTTPS.
- **Target:** PostgreSQL adalah sumber kebenaran dan akses data mengikuti PRD/DESIGN (sqlc + pgx). Implementasi saat ini memakai GORM; ini adalah gap arsitektur yang perlu diputuskan, bukan dianggap telah sesuai.
- File besar/video disimpan di object storage dan didistribusikan melalui CDN; jangan mengalirkan upload video besar melalui function frontend/API biasa.
- Mulai dengan modular monolith, bukan microservices, sampai batas skalabilitas dan kepemilikan tim terbukti.

## 2. Topologi deployment

```text
Browser
  ├── HTTPS → Vercel: Next.js App Router, ISR/SSR, static assets, preview deployments
  └── HTTPS → Go API container: REST /api/v1, auth, course, progress, payment webhook
                    ├── PostgreSQL managed (pooler/connection limits)
                    ├── Redis managed (optional cache/queue/session use)
                    ├── Object storage + CDN (video, PDF, thumbnails)
                    ├── Email provider
                    └── Payment provider
Go worker container ← queue (transcode orchestration, email, certificate, reports)
```

Vercel rewrites/proxy dapat digunakan untuk domain dan routing, tetapi tidak menggantikan deployment API. Gunakan subdomain seperti `www.example.com` dan `api.example.com` atau rewrite terkontrol. Atur CORS hanya untuk origin production dan preview yang disetujui.

## 3. Komponen aplikasi

### Web

- Next.js App Router; Server Components untuk halaman katalog/detail yang dapat di-cache; Client Components hanya untuk interaksi.
- TanStack Query untuk server state interaktif; Zustand untuk state lokal; React Hook Form + Zod untuk form.
- UI: Tailwind + shadcn/ui, dukungan ID/EN, metadata/SEO, aksesibilitas WCAG 2.1 AA.
- Gunakan route handlers hanya untuk kebutuhan web/BFF yang ringan. Aturan domain, pembayaran dan akses course tetap menjadi tanggung jawab API Go.

### API dan worker

- Go modular dengan lapisan handler → service/usecase → repository/query.
- REST JSON dengan versi `/api/v1`, OpenAPI sebagai kontrak, pagination/filter yang tervalidasi.
- Middleware: request ID, structured logging, recovery, CORS, rate limit, auth, timeout dan security headers.
- Worker idempoten untuk pekerjaan lama; transactional outbox atau mekanisme setara untuk menjaga konsistensi event dan database.

### Data

- PostgreSQL managed, migration versioned, backup dan PITR sesuai kemampuan provider.
- Entitas inti mengikuti FRD; tambahan penting: versioning silabus/course, entitlements, attempts, assets, audit events.
- Indeks pada foreign key, slug unik, status/tanggal pencarian, enrollment unik per user-course, progress unik per user-lesson, dan transaksi provider reference.
- Connection pool disetel sesuai limit instance serverless/container; jangan membuka koneksi baru per request tanpa pool/pooler.

### Baseline implementasi repositori

Pemeriksaan manifest dan source pada 15 Oktober 2026 menemukan kondisi berikut. Ini adalah inventaris implementasi saat ini, bukan target arsitektur:

| Area | Kondisi saat ini | Gap terhadap target |
|---|---|---|
| Web | `apps/web/package.json` menetapkan Next.js 16.2.9, React 19.2.4, dan Tailwind CSS 4; dependensi juga mencakup TanStack Query, Zustand, React Hook Form, dan Zod. | Berbeda dari versi Next.js 14 dan Tailwind 3 pada PRD induk. Versi target perlu diselaraskan dengan keputusan produk sebelum deployment produksi. |
| API | `apps/api/go.mod` menetapkan Go 1.25.0 dan Echo 4.15.4; struktur kode memisahkan handler, service, repository, dan model. | PRD induk menetapkan Go 1.22; versi runtime/dependensi perlu keputusan eksplisit. |
| Data | PostgreSQL diakses melalui GORM; `pkg/database/database.go` menjalankan `AutoMigrate` untuk model yang terdaftar. | Belum memakai sqlc/pgx secara langsung atau migrasi SQL berversi seperti target PRD/DESIGN. |
| Rute API | Ada `/health`; auth register/login/refresh; rute user; dan rute course, catalog, enrollment, status, serta module di bawah `/api/v1`. | Belum ditemukan rute implementasi untuk logout, verifikasi/reset email, Google OAuth, payments/webhooks, admin/moderation, quiz, atau certificate. |
| Layanan pendukung | Ada konfigurasi PostgreSQL, JWT, dan middleware Echo. | Worker/queue, Redis, object storage/CDN, email provider, payment provider, dan OpenAPI belum ditemukan pada implementasi yang diperiksa. |

Daftar rute target di bawah adalah kontrak yang direncanakan FRD/TRD dan tidak boleh dianggap sudah tersedia sampai implementasi serta verifikasi endpoint dilakukan. `PRD_eLearning.md` dan `DESIGN.md` tetap menjadi acuan target yang diberikan proyek; perbedaan versi dan ORM di atas perlu diselesaikan sebagai keputusan teknis sebelum dianggap sebagai stack final.

### Runner kode interaktif

- Echo API mengautentikasi learner, memvalidasi ID/versi latihan, mencatat attempt, dan mengantrekan eksekusi; API tidak pernah menjalankan source code learner secara langsung.
- Runner memakai executor terisolasi yang ephemeral untuk setiap attempt. Target pilot: Go dan JavaScript/TypeScript.
- Isolasi wajib: tanpa akses jaringan keluar, secret, database LearnHub, filesystem host, container socket, atau privilege host; batas CPU, memori, wall-clock, proses, output, ukuran payload, dan concurrency diberlakukan.
- Gunakan microVM atau layanan code-execution khusus yang menyediakan batas isolasi kuat; container biasa tanpa sandbox tambahan tidak cukup untuk kode tak tepercaya.
- Simpan versi image/runtime, versi soal, hasil test, durasi, batas resource, status, dan hint unlock. Simpan source code hanya sesuai kebijakan retensi; jangan masukkan source ke log.
- Soal memakai hidden tests yang dievaluasi di runner; response hanya mengembalikan ringkasan hasil dan kategori kegagalan, bukan isi hidden tests.
- Lima hint disimpan sebagai konten bertingkat dan dibuka berurutan. Batas attempt dan rate limit mencegah penyalahgunaan runner.

## 4. Kontrak API target

Prefix `/api/v1`.

| Domain | Rute contoh |
|---|---|
| Auth | `POST /auth/register`, `/auth/login`, `/auth/refresh`, `/auth/logout` |
| Catalog | `GET /tracks`, `GET /courses`, `GET /courses/{slug}`, `GET /courses/{id}/curriculum` |
| Authoring | `POST /courses`, `PATCH /courses/{id}`, `POST /courses/{id}/submit`, section/lesson CRUD |
| Learning | `POST /courses/{id}/enroll`, `GET /me/enrollments`, `PUT /lessons/{id}/progress`, quiz attempts |
| Payment | `POST /payments/checkout`, `GET /payments/history`, `POST /payments/webhooks/{provider}` |
| Admin | user/course moderation, categories, tracks, reporting, audit |

Semua perubahan menggunakan validasi payload dan otorisasi server-side. Operasi yang bisa diulang menerima idempotency key bila relevan. Webhook memakai verifikasi signature provider dan pencatatan event unik.

## 5. Autentikasi dan keamanan

- Gunakan cookie `HttpOnly`, `Secure`, `SameSite` untuk sesi web bila pola deployment memungkinkan; token refresh disimpan aman dan dirotasi. Hindari menyimpan bearer/refresh token jangka panjang di localStorage.
- Password memakai algoritme password hashing yang sesuai dan parameter terkalibrasi; secret dikelola melalui environment/secrets manager.
- Terapkan CSRF protection untuk autentikasi berbasis cookie, CSP, sanitasi konten rich text, rate limit auth, validasi upload, signed URL berumur pendek dan pemeriksaan ownership.
- Pisahkan secret preview dan production; batasi akses database dan storage dengan least privilege.
- Audit perubahan role, moderasi, transaksi/refund, payout dan akses admin.
- Enkripsi transport; backup terenkripsi; kebijakan retensi dan penghapusan data ditetapkan sebelum launch.

## 6. Deploy, CI/CD, dan konfigurasi

- Vercel menghubungkan branch preview dan branch production; cek build dan review preview sebelum promote.
- Go API/worker dibangun sebagai container dan dirilis ke provider container terpisah; migration dijalankan sebagai langkah rilis terkendali, kompatibel dengan versi aplikasi yang berjalan.
- Environment: development, preview/staging, production; database dan credentials harus terpisah.
- Pipeline: format/lint, pemeriksaan dependensi dan secret, validasi kontrak, build web, build image Go, migration check, lalu deployment bertahap.
- Domain, TLS, redirect, callback OAuth, CORS dan webhook URL dikonfigurasi per environment.
- File `.env.example` berisi nama variabel dan contoh nonrahasia; tidak berisi kredensial.

## 7. Observability dan reliability

- Structured logs dengan request/correlation ID; jangan log password, token, data pembayaran sensitif, atau isi privat lesson.
- Metrik: latensi p50/p95/p99, error rate, throughput, DB pool saturation, queue age, webhook retry, playback start, enrollment dan completion.
- Health/readiness endpoint untuk API dan worker; alert berbasis SLO yang disepakati.
- Retry dengan backoff dan dead-letter handling untuk tugas asinkron; timeout eksplisit dan circuit breaking ke provider eksternal.
- Cache hanya untuk data publik yang aman; invalidasi cache saat publish/update. Jangan cache response personal dengan key yang tidak mencakup identitas/entitlement.

## 8. Batas performa dan kapasitas

Target PRD p95 API ≤200 ms dan uptime 99.5% diperlakukan sebagai sasaran yang perlu diukur dan diuji, bukan jaminan awal. Pengecualian untuk operasi lama diproses asinkron. Gunakan CDN untuk aset dan video, pagination untuk daftar, query plan/index review, dan batas upload. Kapasitas 10.000 concurrent users memerlukan load test serta sizing provider tersendiri.

## 9. Keputusan terbuka sebelum production

- Provider container untuk Go API/worker, PostgreSQL, Redis, object storage/CDN.
- Payment gateway (Midtrans/Xendit), komisi, refund, payout, pajak dan mata uang.
- OAuth/email provider, video transcoding/DRM, dan SLA dukungan.
- Target SLO, RPO/RTO, retensi, data residency dan proses respons insiden.
- Perlu tidaknya 2FA, Redis queue, search engine khusus, dan multi-region setelah MVP.

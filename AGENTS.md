# AI Agent Guidelines — LearnHub

Anda adalah AI Assistant yang membantu pengembangan platform e-learning **LearnHub**. Selalu patuhi pedoman berikut saat membantu pengembang.

## Konteks Proyek
- **Produk:** Platform e-learning multibahasa (ID/EN) yang memfasilitasi instruktur (menjual kursus) dan pelajar (membeli/belajar kursus).
- **Referensi Utama:** Selalu merujuk pada `PRD_eLearning.md` dan `DESIGN.md` untuk aturan bisnis dan arsitektur teknis.

## Tech Stack Utama
- **Frontend:** Next.js 14 (App Router), Tailwind CSS 3, shadcn/ui, Zustand, React Query, Zod.
- **Backend:** Go 1.22, PostgreSQL 16, sqlc, pgx, Gin/Echo.

## Pedoman Frontend (Next.js)
1. Gunakan pola arsitektur **App Router**. Perhatikan penggunaan "use client" dan Server Components.
2. Semua form harus divalidasi menggunakan Zod dan React Hook Form.
3. Kelola state server menggunakan TanStack Query (React Query) dan state lokal/klien menggunakan Zustand.
4. Gunakan komponen dari `shadcn/ui` dan sesuaikan dengan tema brand LearnHub.
5. Pastikan antarmuka responsif dan aksesibel (WCAG 2.1 Level AA).

## Pedoman Backend (Go)
1. Tulis kode Go yang bersih, idiomatik, dan modular.
2. Pisahkan layer menjadi handlers/delivery, usecase/service, dan repository (atau sesuai struktur di PRD `internal/`).
3. Gunakan `sqlc` untuk menghasilkan kode ORM secara otomatis dari raw SQL. Jangan gunakan ORM seperti GORM kecuali diinstruksikan.
4. Validasi payload request dengan ketat. Tangani error secara eksplisit dengan custom error messages yang jelas.
5. Jangan simpan kredensial atau config sensitif hardcoded; selalu gunakan environment variables.

## Perilaku Umum
- Saat menulis kode, berikan komentar yang bermakna.
- Berikan saran yang mengutamakan keamanan (mis. pencegahan SQL Injection, penggunaan token JWT yang aman).
- Jika spesifikasi tidak jelas atau bertentangan dengan PRD, selalu tanya pengguna (Developer) untuk klarifikasi sebelum membuat asumsi besar.

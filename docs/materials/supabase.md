# Materi: Supabase untuk Backend-as-a-Service Berbasis PostgreSQL

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun aplikasi modern menggunakan platform backend Supabase: memahami arsitektur berbasis PostgreSQL murni, mengamankan akses data langsung dari browser menggunakan **Row Level Security (RLS)** dan SQL Policies, mengintegrasikan Supabase Auth (Email, OAuth, Magic Link), memanfaatkan Storage untuk file media, berlangganan data real-time via WebSockets, serta menulis Edge Functions berbasis Deno.

## Prasyarat

- Memahami konsep dasar PostgreSQL dan SQL (tabel, foreign key, select/insert).
- Memahami JavaScript / TypeScript dasar.
- Memiliki akun Supabase Cloud atau instance lokal menggunakan Supabase CLI.

## 1. Dasar: filosofi "Open Source Firebase Alternative"

Supabase adalah platform Backend-as-a-Service (BaaS) yang menyediakan fitur instan untuk database, autentikasi, penyimpanan file, dan fungsi serverless. Berbeda dengan platform NoSQL tertutup (seperti Firebase Firestore), Supabase dibangun 100% di atas fondasi standar industri: **PostgreSQL**.

### Komponen Inti Supabase:
- **PostgreSQL Database:** Database relasional penuh dengan hak akses superuser (bisa menambah ekstensi seperti pgvector, PostGIS, pgcrypto).
- **PostgREST:** Web server mandiri yang secara otomatis mengubah skema PostgreSQL menjadi RESTful API instan.
- **GoTrue:** Layanan API autentikasi berbasis JWT untuk mengelola user, sign-in token, dan OAuth.
- **Realtime Server:** Server berbasis Elixir Phoenix yang mendengarkan perubahan PostgreSQL Write-Ahead Log (WAL) dan memancarkannya via WebSocket ke browser klien.
- **Storage:** Penyimpanan file terintegrasi dengan S3 kompatibel yang diatur oleh hak akses PostgreSQL RLS.

## 2. Praktik dasar: Supabase Client SDK dan Row Level Security (RLS)

```sql
-- 1. Skema Tabel di PostgreSQL Supabase
CREATE TABLE profiles (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    avatar_url TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID NOT NULL REFERENCES auth.users(id),
    title TEXT NOT NULL,
    price INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Mengaktifkan Row Level Security (RLS)
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE courses ENABLE ROW LEVEL SECURITY;

-- 3. Membuat Policies Keamanan
-- Siapa pun boleh melihat profil pengguna
CREATE POLICY "Public profiles are viewable by everyone" 
ON profiles FOR SELECT 
USING (true);

-- Pengguna hanya boleh memperbarui profil mereka sendiri
CREATE POLICY "Users can update own profile" 
ON profiles FOR UPDATE 
USING (auth.uid() = id);

-- Publik hanya boleh melihat kursus yang sudah dipublikasikan
CREATE POLICY "Public courses are viewable by everyone" 
ON courses FOR SELECT 
USING (is_published = true);

-- Instruktur hanya boleh mengedit kursus buatan sendiri
CREATE POLICY "Instructors can edit own courses" 
ON courses FOR ALL 
USING (auth.uid() = instructor_id);
```

Integrasi di Frontend dengan TypeScript (`@supabase/supabase-js`):
```typescript
import { createClient } from '@supabase/supabase-js';

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL!;
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!;

export const supabase = createClient(supabaseUrl, supabaseAnonKey);

// 1. Autentikasi Pengguna
async function signIn(email: string, password: string) {
  const { data, error } = await supabase.auth.signInWithPassword({ email, password });
  if (error) console.error('Login error:', error.message);
  return data.user;
}

// 2. Query Data Langsung dari Browser (Aman karena RLS aktif di DB)
async function getPublishedCourses() {
  const { data, error } = await supabase
    .from('courses')
    .select('id, title, price, profiles(full_name)')
    .order('created_at', { ascending: false });

  if (error) throw error;
  return data;
}
```

## 3. Fitur utama dan pola arsitektur

### Realtime Subscriptions via WebSockets
Menerima pembaruan data seketika tanpa polling:

```typescript
// Berlangganan event INSERT pada tabel courses
const channel = supabase
  .channel('public:courses')
  .on(
    'postgres_changes',
    { event: 'INSERT', schema: 'public', table: 'courses' },
    (payload) => {
      console.log('Kursus baru ditambahkan secara realtime:', payload.new);
    }
  )
  .subscribe();
```

### Database Triggers untuk Sinkronisasi `auth.users` ke `public.profiles`
Otomatis buat profil saat ada user mendaftar:

```sql
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS trigger AS $$
BEGIN
  INSERT INTO public.profiles (id, full_name, avatar_url)
  VALUES (new.id, new.raw_user_meta_data->>'full_name', new.raw_user_meta_data->>'avatar_url');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();
```

## 4. Praktik kerja: mencegah kebocoran data akibat lupa mengaktifkan RLS

### Skenario
Developer membuat tabel baru `course_transactions` untuk mencatat nominal pembayaran siswa. Endpoint PostgREST publik mengizinkan user anonim menarik seluruh transaksi siswa lain karena RLS belum diaktifkan pada tabel tersebut.

### Tugas
Aktifkan RLS, buat policy selektif berbasis fungsi helper `auth.uid()`, dan uji akses menggunakan Service Role Key vs Anon Key.

```sql
-- 1. Wajib aktifkan RLS
ALTER TABLE course_transactions ENABLE ROW LEVEL SECURITY;

-- 2. Batasi pembacaan hanya untuk pemilik transaksi atau admin
CREATE POLICY "Students see own transactions" 
ON course_transactions FOR SELECT 
USING (
  auth.uid() = student_id 
  OR 
  EXISTS (SELECT 1 FROM profiles WHERE id = auth.uid() AND role = 'admin')
);

-- 3. Mencegah user meng-insert transaksi fiktif tanpa otentikasi
CREATE POLICY "Students insert own transactions"
ON course_transactions FOR INSERT
WITH CHECK (auth.uid() = student_id);
```

### Lima hint bertahap
1. Tabel PostgreSQL baru di Supabase tidak otomatis memiliki RLS aktif; selalu jalankan `ALTER TABLE ... ENABLE ROW LEVEL SECURITY`.
2. Jika RLS aktif tetapi belum ada policy yang dibuat, PostgreSQL menolak (*deny all*) seluruh akses dari Anon Key.
3. Gunakan `auth.uid()` di dalam ekspresi policy SQL untuk mendapatkan ID user yang saat ini diautentikasi dari token JWT.
4. Jangan pernah membocorkan `service_role_key` ke aplikasi frontend browser/mobile; kunci ini memotong (*bypass*) semua aturan RLS.
5. Gunakan `SECURITY DEFINER` pada fungsi PostgreSQL dengan hati-hati untuk mencegah eskalasi hak istimewa (*privilege escalation*).

## 5. Lanjutan: Supabase Edge Functions (Deno) dan Vector Search (`pgvector`)

- **Edge Functions:** Jalankan serverless TypeScript code secara terdistribusi global via Deno runtime untuk menangani webhook payment gateway (seperti Midtrans/Stripe).
- **pgvector & AI Embeddings:** Simpan embedding teks kurikulum di kolom `vector(1536)` untuk pencarian semantik cerdas (*Semantic Vector Search*) langsung di PostgreSQL.
- **Supabase CLI & Local Development:** Jalankan `supabase start` untuk menjalankan seluruh ekosistem Supabase di Docker lokal secara offline.

## Pemeriksaan pemahaman

1. Apa peran Row Level Security (RLS) dalam mengamankan pemanggilan database langsung dari browser?
2. Apa perbedaan fungsional antara `anon_key` dan `service_role_key` di Supabase?
3. Bagaimana mekanisme Supabase Realtime mendeteksi perubahan data di tabel PostgreSQL?
4. Mengapa trigger `handle_new_user` harus dideklarasikan dengan `SECURITY DEFINER`?

**Kunci ringkas:** RLS memfilter baris database yang diizinkan untuk di-SELECT, UPDATE, atau DELETE langsung di level engine SQL berdasarkan klausa policy dan token JWT user; `anon_key` adalah kunci publik yang mematuhi aturan RLS sedangkan `service_role_key` adalah master key rahasia untuk backend yang mem-bypass seluruh aturan RLS; Supabase Realtime mendengarkan Logical Replication Stream dari Write-Ahead Log (WAL) Postgres lalu menyiarkan event via Phoenix WebSockets; `SECURITY DEFINER` diperlukan agar fungsi trigger berjalan dengan hak akses superuser saat menyisipkan baris ke tabel publik yang dilindungi RLS.

## Referensi

- [Supabase Official Documentation](https://supabase.com/docs)
- [PostgreSQL Row Level Security Documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [Supabase GitHub Architecture](https://github.com/supabase/supabase)

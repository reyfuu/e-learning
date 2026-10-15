# Materi: Next.js App Router dan Full-Stack Modern Web

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun aplikasi web modern menggunakan Next.js App Router: memahami batasan Server vs Client Components, mengelola data fetching dengan caching & revalidasi, memanfaatkan Server Actions untuk mutasi data, serta mengoptimalkan SEO, performa web vital, dan deployment.

## Prasyarat

- Memahami React (komponen, JSX, props, state, hooks dasar).
- Memahami TypeScript dasar dan konsep HTTP (cookies, headers, REST API).
- Node.js versi 18.17+ atau yang lebih baru.

## 1. Dasar: arsitektur Next.js App Router

Next.js mengadopsi paradigma React Server Components (RSC). Secara default di direktori `app/`, semua komponen adalah **Server Components**. Komponen ini dirender di server, tidak menyertakan JavaScript ke bundel browser, dan memiliki akses langsung ke database atau file system server.

### Batasan Penting:
- **Server Components:** Default. Tidak bisa memakai hooks interaktif (`useState`, `useEffect`) atau browser API (`window`, `localStorage`).
- **Client Components:** Ditandai dengan direktif `"use client"` di baris paling atas file. Digunakan untuk komponen interaktif, formulir dengan state lokal, event listeners, atau hooks React standar.
- **File System Routing:** Folder menentukan URL path. File khusus: `page.tsx` (tampilan route), `layout.tsx` (pembungkus halaman bersama), `loading.tsx` (fallback UI streaming Suspense), `error.tsx` (Error Boundary), dan `route.ts` (API endpoint HTTP).

## 2. Praktik dasar: Server Component Fetching & Server Actions

```tsx
// app/courses/page.tsx (Server Component)
import { Suspense } from 'react';
import Link from 'next/link';

interface Course {
  id: string;
  title: string;
  instructor: string;
  price: number;
}

// Data fetching langsung di komponen server tanpa useEffect
async function getCourses(): Promise<Course[]> {
  const res = await fetch('https://api.learnhub.id/v1/courses', {
    next: { revalidate: 60 } // Incremental Static Regeneration (ISR) tiap 60 detik
  });
  if (!res.ok) throw new Error('Gagal memuat kursus');
  return res.json();
}

export default async function CoursesPage() {
  const courses = await getCourses();

  return (
    <main className="p-8 max-w-5xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Katalog Kursus LearnHub</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {courses.map((course) => (
          <div key={course.id} className="border p-4 rounded-lg shadow-sm">
            <h2 className="text-xl font-semibold">{course.title}</h2>
            <p className="text-gray-600">Instruktur: {course.instructor}</p>
            <p className="mt-2 font-bold">Rp {course.price.toLocaleString('id-ID')}</p>
            <Link
              href={`/courses/${course.id}`}
              className="mt-4 inline-block bg-blue-600 text-white px-4 py-2 rounded"
            >
              Lihat Detail
            </Link>
          </div>
        ))}
      </div>
    </main>
  );
}
```

## 3. Fitur utama dan pola arsitektur

### Server Actions untuk Mutasi Data Type-Safe
Server Actions memungkinkan eksekusi fungsi server langsung dari form atau interaksi tanpa membuat endpoint REST API manual:

```tsx
// app/courses/new/actions.ts
'use server';

import { revalidatePath } from 'next/cache';
import { redirect } from 'next/navigation';
import { z } from 'zod';

const CourseSchema = z.object({
  title: z.string().min(5, 'Judul minimal 5 karakter'),
  price: z.coerce.number().min(0, 'Harga tidak boleh negatif'),
});

export async function createCourseAction(formData: FormData) {
  const rawData = {
    title: formData.get('title'),
    price: formData.get('price'),
  };

  const validated = CourseSchema.safeParse(rawData);
  if (!validated.success) {
    return { error: validated.error.flatten().fieldErrors };
  }

  // Simpan ke database (misal Prisma / Drizzle / pgx)
  console.log('Menyimpan ke DB:', validated.data);

  // Revalidasi cache halaman kursus
  revalidatePath('/courses');
  redirect('/courses');
}
```

## 4. Praktik kerja: optimasi streaming UI dengan Suspense & Skeleton

### Skenario
Halaman kursus memiliki dua komponen: profil kursus (cepat, <50ms) dan ulasan siswa dari API pihak ketiga yang lambat (1500ms). Seluruh halaman menjadi lambat dimuat jika keduanya diblokir bersamaan.

### Tugas
Pisahkan komponen ulasan siswa ke dalam Suspense boundary terisolasi dengan skeleton placeholder, sehingga halaman utama terbuka instan dan bagian ulasan streaming secara bertahap.

```tsx
// app/courses/[id]/page.tsx
import { Suspense } from 'react';

async function CourseDetails({ id }: { id: string }) {
  // Query cepat lokal
  return <div className="p-4 bg-white rounded"><h1>Kursus Fullstack Go</h1></div>;
}

async function StudentReviews({ id }: { id: string }) {
  // Query lambat
  await new Promise((r) => setTimeout(r, 2000));
  return <div>List Ulasan Pelajar: 5.0 (Bintang 5)</div>;
}

function ReviewsSkeleton() {
  return <div className="animate-pulse bg-gray-200 h-24 rounded my-4" />;
}

export default function CourseDetailPage({ params }: { params: { id: string } }) {
  return (
    <div className="space-y-6">
      <CourseDetails id={params.id} />
      
      <Suspense fallback={<ReviewsSkeleton />}>
        <StudentReviews id={params.id} />
      </Suspense>
    </div>
  );
}
```

### Lima hint bertahap
1. Jangan membuat satu fungsi async raksasa yang menunggu semua data selesai.
2. Tempatkan komponen yang lambat di dalam file atau sub-komponen async terpisah.
3. Bungkus komponen lambat dengan `<Suspense fallback={<YourSkeleton />}>`.
4. Browser akan menerima HTML awal secara instan, lalu Next.js mengalirkan (streaming) bagian Suspense begitu selesai di server.
5. Gunakan `loading.tsx` untuk loading state di level rute penuh.

## 5. Lanjutan: Caching, Middleware, dan Metadata SEO

- **Next.js Caching Layers:**
  1. *Request Memoization:* Deduplikasi pemanggilan `fetch` yang sama dalam satu render tree.
  2. *Data Cache:* Cache persisten lintas request pengguna (`fetch(..., { cache: 'force-cache' | 'no-store' })`).
  3. *Full Route Cache:* Cache HTML statis pada build time atau ISR.
  4. *Router Cache:* Cache client-side in-memory saat pengguna menavigasi rute.
- **Dynamic Metadata (SEO):** Ekspor fungsi `generateMetadata({ params })` untuk OpenGraph, Twitter cards, dan dynamic title yang ramah mesin pencari.
- **Middleware:** Gunakan `middleware.ts` di root proyek untuk inspeksi cookie sesi, autentikasi redirect, dan rewrite header sebelum request menyentuh server route.

## Pemeriksaan pemahaman

1. Apa perbedaan mendasar antara Server Components dan Client Components?
2. Kapan sebaiknya menggunakan direktif `"use client"`?
3. Bagaimana `revalidatePath` dan `revalidateTag` bekerja pada sistem caching Next.js?
4. Mengapa Server Actions lebih aman untuk mutasi data dibanding endpoint publik tanpa validasi?

**Kunci ringkas:** Server Components dirender eksklusif di server tanpa beban JS di browser, sedangkan Client Components menyertakan JS untuk interaktivitas; gunakan `"use client"` hanya saat membutuhkan state lokal, event listeners, atau browser API; `revalidatePath` menghapus cache rute tertentu agar query data baru diproses ulang; Server Actions berjalan di server, dapat membaca cookies/session secara aman, dan mengeliminasi boilerplate endpoint API terpisah.

## Referensi

- [Next.js Official Documentation](https://nextjs.org/docs)
- [React Server Components Guide](https://react.dev/reference/rsc/server-components)
- [Vercel Deployment & Optimization Guidelines](https://vercel.com/docs)

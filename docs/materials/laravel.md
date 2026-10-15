# Materi: Laravel Framework untuk REST API dan Aplikasi Enterprise

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat mengembangkan aplikasi web dan backend RESTful API dengan ekosistem Laravel: memahami struktur arsitektur MVC dan Service Container, memanfaatkan Eloquent ORM secara optimal (mencegah masalah N+1), mengimplementasikan Form Requests dan Policies untuk validasi & autorisasi, serta memanfaatkan Database Migrations, Queues, dan Events.

## Prasyarat

- Memahami PHP 8.2+ (OOP, namespaces, typing, PDO).
- Memahami konsep basis data relasional (foreign keys, indexing).
- Memiliki Composer, PHP CLI, dan database PostgreSQL/MySQL lokal.

## 1. Dasar: arsitektur dan filosofi Laravel

Laravel mengikuti pola MVC (*Model-View-Controller*) dengan filosofi pengembangan yang elegan (*developer happiness*). Inti dari seluruh framework adalah **IoC Service Container** yang mengelola dependency injection dan binding layanan sistem.

### Komponen Kunci:
- **Routing & Middleware:** Memetakan URI HTTP ke controller action serta memfilter request (autentikasi, rate limit, CORS).
- **Eloquent ORM:** Active Record implementation yang kaya fitur untuk memodelkan tabel relasional menjadi objek PHP.
- **Form Request:** Kelas khusus untuk mengekstraksi aturan validasi dan otorisasi dari controller agar controller tetap ramping (*skinny controller*).
- **Artisan CLI:** Command-line tool bawaan untuk menjalankan migrasi, membuat boilerplate file, dan memantau queue.

## 2. Praktik dasar: REST API Controller, Form Request, dan Eloquent

```php
<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

class StoreCourseRequest extends FormRequest
{
    public function authorize(): bool
    {
        // Hanya user dengan role 'instructor' yang boleh membuat kursus
        return $this->user()?->role === 'instructor';
    }

    public function rules(): array
    {
        return [
            'title' => ['required', 'string', 'min:5', 'max:255'],
            'description' => ['required', 'string'],
            'price' => ['required', 'integer', 'min:0'],
            'category_id' => ['required', 'exists:categories,id'],
        ];
    }
}
```

```php
<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\StoreCourseRequest;
use App\Http\Resources\CourseResource;
use App\Models\Course;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class CourseController extends Controller
{
    public function index(): AnonymousResourceCollection
    {
        // Menggunakan eager loading 'category' dan 'instructor' untuk mencegah masalah N+1
        $courses = Course::with(['category', 'instructor'])
            ->where('is_published', true)
            ->latest()
            ->paginate(15);

        return CourseResource::collection($courses);
    }

    public function store(StoreCourseRequest $request): JsonResponse
    {
        $course = $request->user()->courses()->create($request->validated());

        return (new CourseResource($course->load('category')))
            ->response()
            ->setStatusCode(217);
    }
}
```

## 3. Fitur utama dan pola arsitektur

### Relasi Eloquent dan Eager Loading
Masalah umum performa di Laravel adalah **N+1 Query Problem**, di mana sistem mengeksekusi 1 query untuk mengambil daftar data utama ditambah N query tambahan untuk setiap relasi anak.

```php
// Buruk (N+1 Problem): 1 + 100 queries
$courses = Course::all();
foreach ($courses as $course) {
    echo $course->instructor->name; // Memicu query baru per iterasi
}

// Baik (Eager Loading): Hanya 2 queries
$courses = Course::with('instructor')->get();
foreach ($courses as $course) {
    echo $course->instructor->name; // Diambil dari memory
}
```

### Background Jobs dan Queues
Untuk tugas berat seperti mengirim email sertifikat atau transcoding video materi:

```php
<?php

namespace App\Jobs;

use App\Models\Enrollment;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Mail;

class SendCertificateEmailJob implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    public function __construct(public Enrollment $enrollment) {}

    public function handle(): void
    {
        // Logika pembuatan PDF sertifikat dan pengiriman email
        Mail::to($this->enrollment->user->email)->send(new \App\Mail\CertificateMail($this->enrollment));
    }
}
```

## 4. Praktik kerja: optimasi transaksi keuangan dengan Database Transaction & Pessimistic Locking

### Skenario
Pada saat *flash sale*, kuota diskon terbatas untuk 50 pembeli pertama. Dua pengguna melakukan checkout bersamaan pada kuota terakhir (sisa 1), menyebabkan *race condition* di mana kedua transaksi berhasil dan kuota menjadi negatif (-1).

### Tugas
Terapkan Database Transaction dengan *Pessimistic Locking* (`lockForUpdate()`) di dalam Repository atau Service class untuk memastikan integritas kuota secara atomik.

```php
<?php

namespace App\Services;

use App\Models\Course;
use App\Models\Enrollment;
use App\Models\User;
use Illuminate\Support\Facades\DB;
use Exception;

class EnrollmentService
{
    public function enrollWithQuota(User $user, int $courseId): Enrollment
    {
        return DB::transaction(function () use ($user, $courseId) {
            // Lock baris database dengan SELECT ... FOR UPDATE
            $course = Course::where('id', $courseId)
                ->lockForUpdate()
                ->firstOrFail();

            if ($course->remaining_quota <= 0) {
                throw new Exception('Mohon maaf, kuota pendaftaran kursus ini sudah habis.');
            }

            // Kurangi kuota secara aman
            $course->decrement('remaining_quota');

            // Catat pendaftaran
            return Enrollment::create([
                'user_id' => $user->id,
                'course_id' => $course->id,
                'amount_paid' => $course->discount_price ?? $course->price,
                'status' => 'active',
            ]);
        });
    }
}
```

### Lima hint bertahap
1. Jangan memeriksa sisa kuota dan menguranginya dalam dua operasi query terpisah tanpa isolasi transaksi.
2. `DB::transaction` otomatis melakukan `COMMIT` jika closure selesai tanpa error dan `ROLLBACK` jika ada `Exception`.
3. Gunakan `lockForUpdate()` pada model Eloquent untuk mengunci row sampai transaksi selesai.
4. Pastikan engine database mendukung transaksi (seperti InnoDB di MySQL atau standar PostgreSQL).
5. Buat stress test sederhana menggunakan tool seperti k6 atau Apache Benchmark untuk memvalidasi proteksi race condition.

## 5. Lanjutan: caching, API Resources, dan otentikasi Sanctum

- **Laravel Sanctum:** Solusi otentikasi ringan berbasis cookie SPA atau API token untuk mobile app.
- **Cache Tags & Redis:** Simpan hasil kalkulasi atau data kurikulum statis ke dalam Redis cache (`Cache::tags(['courses'])->remember(...)`).
- **Laravel Pulse & Telescope:** Tool observability real-time untuk memantau waktu respons query lambat, job failure, dan penggunaan memori.

## Pemeriksaan pemahaman

1. Apa itu masalah N+1 pada ORM dan bagaimana cara mengatasinya di Laravel?
2. Apa kegunaan `FormRequest` dibanding menulis validasi langsung di dalam controller?
3. Mengapa `lockForUpdate()` diperlukan saat menangani transaksi kuota terbatas?
4. Apa perbedaan mendasar antara model sinkron dan asinkron (Queue) dalam Laravel?

**Kunci ringkas:** N+1 terjadi saat loop memicu query database baru untuk setiap entitas relasi, diatasi dengan Eager Loading (`with()`); FormRequest memisahkan validasi dan otorisasi dari business logic controller sehingga kode tetap bersih dan mudah diuji; `lockForUpdate()` mengunci baris data di level database sehingga proses konkuren lain harus mengantre, mencegah race condition; Model sinkron memblokir HTTP response sampai eksekusi selesai, sedangkan Queue mendelegasikan tugas berat ke worker background sehingga response pengguna instan.

## Referensi

- [Laravel Official Documentation](https://laravel.com/docs)
- [Laracasts - PHP and Laravel Learning Platform](https://laracasts.com/)
- [Laravel Beyond CRUD by Spatie](https://spatie.be/)

# Materi: Modern PHP 8.x dan Rekayasa Backend Berbasis Standar

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat menulis kode PHP modern (PHP 8.2+) dengan standar industri: memahami sistem tipe yang ketat (*strict types*), Object-Oriented Programming (OOP) tingkat lanjut, penanganan error dan exception, manajemen paket dengan Composer, implementasi standar PSR (PHP Standard Recommendations), serta interaksi basis data yang aman melalui PDO.

## Prasyarat

- Memahami konsep dasar pemrograman dan siklus request-response HTTP.
- Memiliki PHP 8.2+ dan Composer terpasang di terminal lokal.
- Memiliki server database (PostgreSQL/MySQL) atau SQLite.

## 1. Dasar: evolusi PHP modern

PHP 8.x menghadirkan performa tinggi (JIT compiler) dan sistem pengetikan statis yang kuat, menjadikannya bahasa backend yang matang dan aman, jauh dari reputasi lama PHP berbasis skrip prosedural.

### Fitur Kunci PHP 8+:
- **Strict Typing:** `declare(strict_types=1);` mewajibkan validasi tipe argumen dan nilai kembalian.
- **Constructor Property Promotion:** Mengurangi *boilerplate* kelas dengan mendeklarasikan properti langsung di parameter konstruktor.
- **Match Expression:** Ekspresi percabangan ekspresif yang memeriksa tipe secara ketat dan mengembalikan nilai.
- **Readonly Properties & Classes:** Properti atau kelas yang nilainya tidak dapat diubah setelah inisialisasi.
- **Union & Intersection Types:** Fleksibilitas tipe kombinasi (`int|float`, `Countable&Iterator`).

## 2. Praktik dasar: OOP Modern, Validasi Tipe, dan Match

```php
<?php

declare(strict_types=1);

namespace LearnHub\Domain;

enum EnrollmentStatus: string
{
    case PENDING = 'pending';
    case ACTIVE = 'active';
    case COMPLETED = 'completed';
    case CANCELLED = 'cancelled';
}

final readonly class Student
{
    public function __construct(
        public string $id,
        public string $name,
        public string $email,
        public \DateTimeImmutable $registeredAt
    ) {}
}

final class CourseEnrollment
{
    public function __construct(
        public readonly Student $student,
        public readonly string $courseId,
        public EnrollmentStatus $status = EnrollmentStatus::PENDING,
        private ?\DateTimeImmutable $completedAt = null
    ) {}

    public function activate(): void
    {
        if ($this->status !== EnrollmentStatus::PENDING) {
            throw new \DomainException('Hanya pendaftaran berstatus pending yang dapat diaktifkan.');
        }
        $this->status = EnrollmentStatus::ACTIVE;
    }

    public function getLabel(): string
    {
        return match ($this->status) {
            EnrollmentStatus::PENDING => 'Menunggu Pembayaran',
            EnrollmentStatus::ACTIVE => 'Sedang Belajar',
            EnrollmentStatus::COMPLETED => 'Selesai Belajar',
            EnrollmentStatus::CANCELLED => 'Dibatalkan',
        };
    }
}
```

## 3. Fitur utama dan pola arsitektur

### Interaksi Database Aman dengan PDO dan Prepared Statements
Mencegah serangan SQL Injection dengan *prepared statement*:

```php
<?php

declare(strict_types=1);

namespace LearnHub\Infrastructure;

use PDO;
use PDOException;

final class CourseRepository
{
    public function __construct(private PDO $pdo)
    {
        $this->pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        $this->pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);
    }

    public function findById(string $id): ?array
    {
        $stmt = $this->pdo->prepare('SELECT id, title, price, is_published FROM courses WHERE id = :id LIMIT 1');
        $stmt->execute(['id' => $id]);
        $result = $stmt->fetch();

        return $result ?: null;
    }

    public function create(string $id, string $title, int $price): void
    {
        $stmt = $this->pdo->prepare(
            'INSERT INTO courses (id, title, price, is_published, created_at) VALUES (:id, :title, :price, false, NOW())'
        );
        $stmt->execute([
            'id' => $id,
            'title' => $title,
            'price' => $price,
        ]);
    }
}
```

### PSR (PHP Standards Recommendations)
- **PSR-1 & PSR-12:** Standar gaya penulisan kode (*coding style guide*).
- **PSR-4:** Standar autoloader class berbasis namespace direktori.
- **PSR-7 & PSR-15:** Abstraksi pesan HTTP dan middleware web request.

## 4. Praktik kerja: implementasi Dependency Injection dan custom Exception Handler

### Skenario
Aplikasi mengalami kebocoran error trace database mentah ke browser pengguna saat koneksi database terputus, mengekspos informasi kredensial sensitif.

### Tugas
Buat wrapper exception global dan container dependency injection sederhana yang menangkap semua error fatal, mencatat log ke file aman, dan mengembalikan response JSON `500 Server Error` yang bersih bagi klien.

```php
<?php

declare(strict_types=1);

namespace LearnHub\Core;

use Throwable;

final class ErrorHandler
{
    public static function register(): void
    {
        set_exception_handler([self::class, 'handleException']);
    }

    public static function handleException(Throwable $e): void
    {
        // 1. Catat log error teknis secara aman
        error_log(sprintf(
            "[%s] Unhandled Exception: %s in %s:%d\nStack trace:\n%s",
            date('Y-m-d H:i:s'),
            $e->getMessage(),
            $e->getFile(),
            $e->getLine(),
            $e->getTraceAsString()
        ));

        // 2. Bersihkan output buffer
        if (ob_get_length()) {
            ob_clean();
        }

        // 3. Kembalikan response JSON generik ke pengguna
        http_response_code(500);
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode([
            'success' => false,
            'message' => 'Terjadi kesalahan internal server. Tim teknis kami telah menerima notifikasi.',
            'code' => 500
        ], JSON_THROW_ON_ERROR);
        exit;
    }
}
```

### Lima hint bertahap
1. Gunakan `set_exception_handler` untuk menangkap uncaught exceptions di PHP.
2. Jangan pernah menyetel `display_errors = On` pada server *production*.
3. Pisahkan antara pesan log internal yang detail dengan pesan respons publik ke user.
4. Tangani `Throwable` agar dapat menangkap baik `Exception` maupun fatal `Error` (seperti `TypeError` atau `ParseError`).
5. Pastikan header HTTP dikirim sebelum konten body respons.

## 5. Lanjutan: Composer autoloading, testing, dan performa

- **Composer Autoloading:** Definisikan mapping `"psr-4": { "LearnHub\\": "src/" }` pada `composer.json` dan jalankan `composer dump-autoload -o` untuk optimasi di production.
- **Unit Testing:** Gunakan PHPUnit atau Pest untuk menulis automated tests bagi logika domain dan repository.
- **OPcache:** Aktifkan ekstensi `opcache` dan `opcache.preload` di server Linux untuk meng-compile bytecode ke shared memory demi lonjakan throughput request per detik.

## Pemeriksaan pemahaman

1. Apa dampak menambahkan `declare(strict_types=1);` pada bagian awal file PHP?
2. Bagaimana `match` expression lebih unggul dibanding pernyataan `switch` konvensional?
3. Mengapa Prepared Statements di PDO sepenuhnya mencegah SQL Injection?
4. Apa kegunaan standar PSR-4 dalam ekosistem Composer?

**Kunci ringkas:** `strict_types=1` memaksa pengecekan tipe argumen dan return value tanpa konversi implisit; `match` menggunakan perbandingan identitas ketat (`===`), mengembalikan nilai langsung, dan melempar `UnhandledMatchError` jika tidak ada cabang yang cocok; Prepared Statement memisahkan struktur query SQL dari data parameter sehingga input pengguna tidak pernah dieksekusi sebagai perintah SQL; PSR-4 mengotomatisasi pemetaan namespace kelas ke struktur folder fisik file `.php`.

## Referensi

- [PHP Official Documentation](https://www.php.net/manual/en/)
- [PHP: The Right Way](https://phptherightway.com/)
- [PHP-FIG Standards (PSR)](https://www.php-fig.org/psr/)

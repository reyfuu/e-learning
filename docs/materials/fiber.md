# Materi: Fiber Framework Berbasis Fasthttp untuk Throughput Maksimal

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun REST API dengan kecepatan maksimal menggunakan Go Fiber (v2/v3): memahami arsitektur berbasis `fasthttp` (bukan standard library `net/http`), sintaks routing yang mirip dengan Express.js, manajemen memori nol-alokasi (*zero memory allocation*), penggunaan middleware, serta penanganan jebakan siklus hidup pointer (*buffer immutability pitfalls*).

## Prasyarat

- Memahami bahasa Go (struct, pointer, slicing, concurrency).
- Mengerti perbedaan antara TCP connection multiplexing dan HTTP request handling.
- Go 1.21+ terpasang di sistem.

## 1. Dasar: arsitektur Fiber dan Fasthttp

Fiber terinspirasi oleh Express.js dalam hal kemudahan sintaks, namun dibangun di atas mesin jaringan **Fasthttp** (bukan `net/http` bawaan Go). Fasthttp dirancang khusus untuk meminimalkan alokasi memori heap hingga mendekati nol melalui penggunaan ulang buffer memori (*memory pooling*).

### Konsep Kunci:
- **Zero Allocation Philosophy:** String dan slice byte yang diterima dari request tidak disalin secara otomatis; memorinya didaur ulang segera setelah handler selesai.
- **Express-like Syntax:** Metode chaining `app.Get()`, `app.Post()`, `c.Params()`, `c.Query()`, dan `c.JSON()`.
- **Pre-fork Support:** Memanfaatkan multi-core CPU dengan menjalankan beberapa proses worker pada port yang sama menggunakan socket flag `SO_REUSEPORT`.

## 2. Praktik dasar: setup Fiber Server dan Routing

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Course struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       false, // Set true di Linux production untuk multi-process fork
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "LearnHub-Fiber",
		AppName:       "LearnHub API v1",
	})

	// Pasang middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://learnhub.id",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// API Routing
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	api.Post("/courses", handleCreateCourse)
	api.Get("/courses/:id", handleGetCourse)

	log.Fatal(app.Listen(":3000"))
}

func handleCreateCourse(c *fiber.Ctx) error {
	course := new(Course)

	if err := c.BodyParser(course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal membaca body JSON",
		})
	}

	course.ID = "crs_fib_88"

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    course,
	})
}

func handleGetCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"id":    id,
		"title": "Ultra Fast Go with Fiber",
	})
}
```

## 3. Fitur utama dan pola arsitektur

### Jebakan Memori Zero Allocation: `c.Immutable()` vs Salinan Buffer
Karena Fasthttp menggunakan buffer yang sama berulang kali (*buffer reuse*), string yang dikembalikan oleh `c.Params()`, `c.Query()`, atau `c.Body()` **tidak aman** digunakan di dalam goroutine konkuren tanpa disalin:

```go
// BAHAYA (Bisa menghasilkan data korup / garbage string):
func BadHandler(c *fiber.Ctx) error {
	param := c.Params("id") // Menunjuk ke temporary byte buffer
	go func() {
		// Ketika goroutine ini berjalan, buffer memori param mungkin sudah ditimpa request baru!
		println(param)
	}()
	return c.SendStatus(fiber.StatusOK)
}

// BENAR (Menyalin string / mengaktifkan immutability):
func GoodHandler(c *fiber.Ctx) error {
	// Salin nilai secara eksplisit menggunakan utils.CopyString atau string([]byte(...))
	param := string([]byte(c.Params("id")))
	go func() {
		println(param) // Aman!
	}()
	return c.SendStatus(fiber.StatusOK)
}

// Atau aktifkan Immutable di konfigurasi:
// fiber.New(fiber.Config{ Immutable: true })
```

## 4. Praktik kerja: custom rate limiter middleware berbasis token bucket

### Skenario
Endpoint pengiriman submission tugas sering diserang bot dengan ribuan spam request per detik, menyebabkan lonjakan CPU pada server database.

### Tugas
Konfigurasikan middleware `limiter` resmi Fiber dengan pembatasan 10 request per menit per alamat IP klien.

```go
package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func setupRateLimiter(app *fiber.App) {
	app.Use("/api/v1/courses/submit", limiter.New(limiter.Config{
		Max:        10,              // Maksimal 10 request
		Expiration: 1 * time.Minute, // Dalam jendela waktu 1 menit
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Batasi berdasarkan IP klien
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit terlampaui",
				"message": "Anda telah mencapai batas maksimum request. Silakan coba kembali dalam 1 menit.",
			})
		},
	}))
}
```

### Lima hint bertahap
1. Pasang middleware limiter secara spesifik pada rute yang rawan serangan atau secara global pada grup API.
2. Gunakan `c.IP()` untuk mendapatkan IP asli klien (pastikan memvalidasi header proxy jika berada di balik reverse proxy Nginx).
3. `LimitReached` memungkinkan kustomisasi response body dan header retry-after.
4. Di lingkungan multi-server / cluster, gunakan penyimpanan storage Redis untuk limiter (`limiter.Config{ Storage: redisStorage }`).
5. Jangan menyetel limit terlalu ketat pada endpoint aset statis.

## 5. Lanjutan: WebSocket, Monitor dashboard, dan Fiber Prefork

- **Fiber Monitor:** Middleware visualisasi metrik server real-time (penggunaan CPU, RAM, jumlah koneksi aktif) via `/metrics` (`app.Get("/dashboard", monitor.New())`).
- **Prefork Mode:** Menggandakan throughput pada server multicore Linux dengan memanfaatkan syscall `SO_REUSEPORT` kernel.
- **WebSocket:** Mendukung komunikasi dua arah latensi ultra-rendah dengan package `fiber/websocket/v2`.

## Pemeriksaan pemahaman

1. Mengapa Fiber memiliki performa alokasi memori yang jauh lebih rendah dibanding framework berbasis `net/http`?
2. Apa bahaya menggunakan nilai `c.Params()` di dalam background goroutine pada konfigurasi default Fiber?
3. Apa fungsi opsi `Prefork: true` pada konfigurasi Fiber?
4. Kapan sebaiknya kita memilih standard library `net/http` dibanding Fiber?

**Kunci ringkas:** Fiber memakai Fasthttp yang mendaur ulang buffer byte alih-alih mengalokasikan objek baru di heap pada setiap request; String param merujuk langsung ke buffer yang dapat ditimpa oleh koneksi berikutnya sehingga data menjadi rusak di goroutine; `Prefork` menjalankan proses terpisah per core CPU di Linux untuk menghindari lock contention antar-thread; `net/http` lebih dianjurkan jika membutuhkan kompatibilitas HTTP/2 penuh, HTTP/3, atau integrasi luas dengan ekosistem middleware standar Go.

## Referensi

- [Fiber Official Documentation](https://docs.gofiber.io/)
- [Fasthttp GitHub Repository](https://github.com/valyala/fasthttp)
- [Fiber Awesome Collection](https://github.com/gofiber/awesome-fiber)

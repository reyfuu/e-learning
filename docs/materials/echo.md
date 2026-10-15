# Materi: Echo Framework untuk REST API Cepat dan Ringan di Go

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun API HTTP berkinerja tinggi menggunakan Labstack Echo v4: menguasai abstraksi `echo.Context`, routing berbasis Radix tree yang dioptimalkan, request binding & custom validation, pembuatan middleware custom, penanganan error terpusat (`HTTPErrorHandler`), serta pengamanan API produksi.

## Prasyarat

- Memahami bahasa Go (struct, pointer, interface, standard library `net/http`).
- Memahami konsep RESTful API dan JSON payload.
- Go 1.21+ terpasang di sistem.

## 1. Dasar: arsitektur Echo Framework

Echo adalah framework web minimalis untuk Go yang berfokus pada kecepatan eksekusi, penggunaan memori minimal, dan kenyamanan pengembang. Inti dari Echo adalah objek `echo.Context` yang merangkum request HTTP, response writer, path parameters, query parameters, dan shared handler state.

### Keunggulan Echo:
- **Router Cepat:** Menggunakan Radix tree yang sangat efisien dan mampu memproses rute dinamis tanpa alokasi memori berlebih.
- **Middleware Ekstensif:** Menyediakan middleware bawaan (Logger, Recover, CORS, Secure, Rate Limiter).
- **Data Binding & Validasi:** Kemampuan otomatis mem-binding JSON, form, XML, atau query param ke dalam struct Go.
- **Centralized Error Handling:** Format respons error yang dapat diseragamkan dengan satu fungsi handler.

## 2. Praktik dasar: setup Server, Struct Binding, dan Validator

```go
package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Struct Validator kustom mengintegrasikan go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type CreateCoursePayload struct {
	Title       string  `json:"title" validate:"required,min=5,max=100"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"gte=0"`
}

type CourseResponse struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

func main() {
	e := echo.New()

	// 1. Registrasi validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// 2. Middlewares bawaan
	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) // Cegah panic mematikan server
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://learnhub.id"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// 3. Routing
	api := e.Group("/api/v1")
	api.POST("/courses", handleCreateCourse)
	api.GET("/courses/:id", handleGetCourse)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleCreateCourse(c echo.Context) error {
	var payload CreateCoursePayload

	// Bind data dari request body JSON ke struct
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Format JSON tidak valid")
	}

	// Validasi nilai struct
	if err := c.Validate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Simulasikan penyimpanan
	response := CourseResponse{
		ID:    "crs_echo_01",
		Title: payload.Title,
		Price: payload.Price,
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"data":    response,
	})
}

func handleGetCourse(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, echo.Map{
		"id":    id,
		"title": "Belajar Echo Framework",
	})
}
```

## 3. Fitur utama dan pola arsitektur

### Custom Middleware di Echo
Middleware di Echo memiliki signature `func(echo.HandlerFunc) echo.HandlerFunc`:

```go
func CustomAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("X-API-KEY")
		if apiKey != "secret-token-learnhub" {
			return echo.NewHTTPError(http.StatusUnauthorized, "API Key tidak valid")
		}
		
		// Simpan informasi ke context untuk handler berikutnya
		c.Set("client_role", "admin")
		return next(c)
	}
}
```

### Centralized HTTP Error Handler
Seragamkan format error JSON untuk seluruh aplikasi:

```go
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message interface{} = "Terjadi kesalahan internal pada server"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, echo.Map{
				"success": false,
				"error": echo.Map{
					"code":    code,
					"message": message,
				},
			})
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

// Di main():
// e.HTTPErrorHandler = CustomHTTPErrorHandler
```

## 4. Praktik kerja: implementasi Graceful Shutdown pada Echo Server

### Skenario
Saat deployment kontainer Docker baru di Kubernetes, koneksi HTTP yang sedang mengunggah file materi terputus mendadak (*connection reset by peer*) karena server Go langsung berhenti sebelum menyelesaikan request.

### Tugas
Terapkan pola *Graceful Shutdown* menggunakan `os/signal` dan `c.Server.Shutdown(ctx)` dengan batas waktu toleransi 10 detik.

```go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func runServer() {
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Jalankan server di goroutine terpisah
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Mematikan server karena error:", err)
		}
	}()

	// Tangkap sinyal terminasi OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Menerima sinyal shutdown, menutup server dengan aman...")

	// Beri batas waktu 10 detik bagi request yang sedang aktif
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Gagal shutdown graceful:", err)
	}

	e.Logger.Info("Server berhasil berhenti dengan bersih.")
}
```

### Lima hint bertahap
1. `e.Start()` adalah fungsi blocking; jalankan di dalam `go func()` agar eksekusi berlanjut ke listener sinyal.
2. Gunakan `signal.Notify` dengan sinyal `os.Interrupt` dan `syscall.SIGTERM`.
3. Saat `e.Shutdown(ctx)` dipanggil, Echo berhenti menerima request baru dan menunggu koneksi aktif selesai.
4. Gunakan `context.WithTimeout` agar shutdown tidak menggantung selamanya jika ada request yang macet.
5. Pastikan mengabaikan error `http.ErrServerClosed` karena error ini normal saat shutdown dipanggil.

## 5. Lanjutan: profiling, Sub-domain routing, dan template rendering

- **Sub-domain Routing:** Echo mendukung routing berdasarkan Host header (`e.Host("api.learnhub.id").Group(...)`).
- **Gzip Compression:** Aktifkan `middleware.Gzip()` untuk mengompres payload respons JSON besar.
- **Testing:** Gunakan paket `net/http/httptest` bersama `echo.NewContext` untuk menulis unit test handler yang cepat tanpa membuka soket port jaringan fisik.

## Pemeriksaan pemahaman

1. Apa peran `echo.Context` dalam request lifecycle framework Echo?
2. Mengapa middleware `middleware.Recover()` wajib dipasang di aplikasi produksi?
3. Bagaimana cara memvalidasi struct request menggunakan library pihak ketiga di Echo?
4. Apa yang terjadi jika sinyal SIGTERM tidak ditangani dengan Graceful Shutdown?

**Kunci ringkas:** `echo.Context` menyediakan akses ke data request, response builder, parameter URL, dan penyimpanan variabel antar-middleware; `middleware.Recover()` menangkap panic di dalam goroutine request dan mengubahnya menjadi response 500 sehingga seluruh proses server tidak mati; Integrasi dilakukan dengan meng-assign struct yang mengimplementasikan method `Validate(i interface{}) error` ke `e.Validator`; Tanpa graceful shutdown, koneksi TCP langsung diputus paksa sehingga klien menerima network error atau data transaksi tersimpan separuh (*corrupted state*).

## Referensi

- [Echo Official Guide & Documentation](https://echo.labstack.com/)
- [Echo GitHub Repository](https://github.com/labstack/echo)
- [Go HTTP Package Standard Library](https://pkg.go.dev/net/http)

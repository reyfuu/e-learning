# Materi: Gin Web Framework untuk API Berkecepatan Tinggi di Go

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang API RESTful menggunakan Gin Gonic: menguasai routing berbasis Radix Tree dan Route Groups, binding payload (`ShouldBindJSON`, `ShouldBindUri`), integrasi middleware, pengelolaan log terstruktur, serta pengorganisasian kode dengan Clean Architecture.

## Prasyarat

- Memahami bahasa pemrograman Go (tipe data, pointers, structs, interfaces).
- Memahami konsep dasar REST API dan HTTP status codes.
- Go 1.21+ terpasang di sistem.

## 1. Dasar: performa dan kesederhanaan Gin

Gin adalah salah satu framework web paling populer di ekosistem Go. Dibangun dengan fokus pada performa ekstrem (menggunakan HttpRouter yang dioptimalkan), Gin menawarkan produktivitas tinggi dengan API yang ringkas dan footprint memori yang minim.

### Fitur Kunci:
- **Fast HTTP Router:** Berbasis Radix tree tanpa alokasi memori berlebih, mendukung parameter URL (`/users/:id`) dan wildcard (`/files/*filepath`).
- **Built-in Validator:** Terintegrasi langsung dengan `go-playground/validator` melalui tag struct `binding:"required,email"`.
- **Engine & Context:** `gin.Engine` bertindak sebagai multiplexer, dan `*gin.Context` membawa seluruh konteks request/response.
- **Middleware Chaining:** Menggunakan method `c.Next()` dan `c.Abort()` untuk mengendalikan rantai middleware.

## 2. Praktik dasar: setup Router, Model Binding, dan Validation

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCourseInput struct {
	Title       string  `json:"title" binding:"required,min=5,max=100"`
	Category    string  `json:"category" binding:"required,oneof=frontend backend devops mobile"`
	Price       float64 `json:"price" binding:"gte=0"`
	Instructor  string  `json:"instructor" binding:"required"`
}

type CourseResponse struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Category   string  `json:"category"`
	Price      float64 `json:"price"`
	Instructor string  `json:"instructor"`
}

func main() {
	// Buat router tanpa middleware bawaan default (atau gunakan gin.Default())
	router := gin.New()

	// Pasang middleware standar
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Route grouping
	v1 := router.Group("/api/v1")
	{
		courses := v1.Group("/courses")
		{
			courses.GET("", listCoursesHandler)
			courses.GET("/:id", getCourseHandler)
			courses.POST("", createCourseHandler)
		}
	}

	router.Run(":8080")
}

func createCourseHandler(c *gin.Context) {
	var input CreateCourseInput

	// ShouldBindJSON memeriksa header Content-Type dan melakukan validasi tag
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	response := CourseResponse{
		ID:         "crs_gin_101",
		Title:      input.Title,
		Category:   input.Category,
		Price:      input.Price,
		Instructor: input.Instructor,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    response,
	})
}

func getCourseHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"title": "Mastering Go with Gin",
	})
}

func listCoursesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []string{"Go Concurrency", "Gin Web Framework"},
	})
}
```

## 3. Fitur utama dan pola arsitektur

### Middleware Pipeline: `c.Next()` vs `c.Abort()`
Middleware di Gin dapat memotong eksekusi sebelum atau sesudah handler:

```go
func ResponseTimerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Jalankan handler berikutnya
		c.Next()

		// Kode di bawah dieksekusi setelah response selesai dibuat
		latency := time.Since(start)
		c.Header("X-Response-Time", latency.String())
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token otentikasi wajib disertakan"})
			c.Abort() // Mencegah handler selanjutnya dipanggil!
			return
		}
		c.Set("user_id", "usr_999")
		c.Next()
	}
}
```

## 4. Praktik kerja: penanganan context concurrency yang aman

### Skenario
Developer ingin mencatat log aktivitas audit ke database di background menggunakan `go func()` di dalam handler. Aplikasi mengalami panic acak dan crash dengan error *data race / nil pointer dereference* pada objek `*gin.Context`.

### Tugas
Pahami bahwa `*gin.Context` didaur ulang (*reused*) oleh sync.Pool bawaan Gin setelah request selesai. Buat salinan konteks menggunakan `c.Copy()` sebelum dilempar ke goroutine background.

```go
func handleAuditLog(c *gin.Context) {
	// Salin konteks sebelum masuk ke goroutine
	cCp := c.Copy()

	go func() {
		// Gunakan cCp, JANGAN gunakan c asli!
		time.Sleep(1 * time.Second) // Simulasi operasi lambat
		log.Printf("[AUDIT] Path: %s, User: %s", cCp.Request.URL.Path, cCp.GetString("user_id"))
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Permintaan berhasil diproses"})
}
```

### Lima hint bertahap
1. Objek `*gin.Context` bersifat mutable dan tidak thread-safe.
2. Begitu handler selesai dan response dikirim, Gin me-reset `*gin.Context` dan mengembalikannya ke pool memori.
3. Mengakses pointer `c` asli di dalam goroutine konkuren menyebabkan race condition parah.
4. Method `c.Copy()` membuat duplikat konteks yang aman (*read-only state*) untuk goroutine asynchronous.
5. Jika goroutine membutuhkan batas waktu atau koneksi DB, turunkan `context.Context` dari `c.Request.Context()`.

## 5. Lanjutan: Production Mode, Testing, dan Clean Architecture

- **Gin Mode:** Di lingkungan produksi, selalu setel `gin.SetMode(gin.ReleaseMode)` untuk mematikan log debug yang memakan CPU.
- **Integration Testing:** Gunakan `httptest.NewRecorder()` dan `router.ServeHTTP(w, req)` untuk menguji endpoint Gin secara terisolasi.
- **Clean Architecture:** Pisahkan handler Gin dari business logic repository/service agar handler hanya bertanggung jawab pada parsing HTTP dan serialisasi JSON.

## Pemeriksaan pemahaman

1. Apa akibatnya jika memanggil `c.JSON()` setelah `c.Abort()` pada middleware?
2. Mengapa tidak boleh menggunakan pointer `*gin.Context` langsung di dalam goroutine background?
3. Apa perbedaan fungsi `c.BindJSON()` dan `c.ShouldBindJSON()`?
4. Bagaimana Gin mengoptimalkan alokasi memori objek konteks request?

**Kunci ringkas:** `c.Abort()` menghentikan pemanggilan handler berikutnya, namun pemanggilan `c.JSON()` di dalam middleware yang sama tetap mengirim respons ke klien; objek konteks didaur ulang oleh sync.Pool sehingga pointer asli akan mengalami race condition jika dibaca di background thread; `BindJSON` otomatis mengirim respons error 400 jika binding gagal sedangkan `ShouldBindJSON` mengembalikan error Go biasa untuk penanganan fleksibel; Gin menggunakan `sync.Pool` untuk menggunakan kembali memori struct Context tanpa membebani Garbage Collector.

## Referensi

- [Gin Web Framework Documentation](https://gin-gonic.com/docs/)
- [Gin Repository on GitHub](https://github.com/gin-gonic/gin)
- [Go HTTP Routing Benchmarks](https://github.com/julienschmidt/go-http-routing-benchmark)

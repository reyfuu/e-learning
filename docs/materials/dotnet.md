# Materi: ASP.NET Core 8 Web API dan Ekosistem .NET Modern

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang backend API kelas enterprise menggunakan ASP.NET Core 8: memahami pipeline HTTP Middleware, menguasai Dependency Injection bawaan (Transient, Scoped, Singleton), memanfaatkan Entity Framework Core (EF Core 8) untuk persistensi data relasional, mengimplementasikan otentikasi JWT Bearer, serta mengoptimalkan performa aplikasi dengan Minimal APIs dan Response Caching.

## Prasyarat

- Memahami bahasa pemrograman C# (OOP, async/await, LINQ, records).
- Memahami konsep RESTful API dan database relasional (PostgreSQL/SQL Server).
- .NET SDK 8+ terpasang di sistem.

## 1. Dasar: arsitektur ASP.NET Core dan HTTP Pipeline

ASP.NET Core adalah framework web cross-platform dan open-source dengan performa tertinggi dalam berbagai benchmark TechEmpower. Aplikasi ASP.NET Core dikonfigurasi melalui objek `WebApplication` terpadu yang menyatukan Dependency Injection Container dan HTTP Request Pipeline.

### Tiga Tingkatan Lifecycle Dependency Injection:
- **Transient (`AddTransient`):** Dibuat baru setiap kali diminta (*fresh instance*). Cocok untuk service ringan tanpa state.
- **Scoped (`AddScoped`):** Dibuat satu kali per siklus request HTTP. Standar wajib untuk EF Core `DbContext`.
- **Singleton (`AddSingleton`):** Dibuat satu kali saat aplikasi startup dan dipakai bersama oleh seluruh request. Cocok untuk caching in-memory.

## 2. Praktik dasar: Minimal APIs, EF Core, dan DTO Validation

```csharp
using Microsoft.EntityFrameworkCore;
using System.ComponentModel.DataAnnotations;

var builder = WebApplication.CreateBuilder(args);

// 1. Registrasi Database Context (Scoped)
builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

// 2. Registrasi Swagger OpenAPI
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

// 3. Routing Minimal API
var coursesApi = app.MapGroup("/api/v1/courses").WithTags("Courses");

coursesApi.MapGet("/", async (AppDbContext db) =>
{
    var courses = await db.Courses.AsNoTracking().ToListAsync();
    return Results.Ok(courses);
});

coursesApi.MapGet("/{id:guid}", async (Guid id, AppDbContext db) =>
{
    var course = await db.Courses.FindAsync(id);
    return course is not null ? Results.Ok(course) : Results.NotFound();
});

coursesApi.MapPost("/", async (CreateCourseDto dto, AppDbContext db) =>
{
    var course = new Course
    {
        Id = Guid.NewGuid(),
        Title = dto.Title,
        Price = dto.Price,
        CreatedAt = DateTime.UtcNow
    };

    db.Courses.Add(course);
    await db.SaveChangesAsync();

    return Results.Created($"/api/v1/courses/{course.Id}", course);
});

app.Run();

// Definisi Model & DTO
public class Course
{
    public Guid Id { get; set; }
    public required string Title { get; set; }
    public decimal Price { get; set; }
    public DateTime CreatedAt { get; set; }
}

public record CreateCourseDto(
    [Required, MinLength(5)] string Title,
    [Range(0, 100_000_000)] decimal Price
);

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
    public DbSet<Course> Courses => Set<Course>();
}
```

## 3. Fitur utama dan pola arsitektur

### Optimasi Query EF Core: `AsNoTracking()`
Secara default, EF Core melacak perubahan (*Change Tracker*) setiap entitas yang ditarik dari database. Untuk endpoint baca-saja (*read-only*), nonaktifkan tracking untuk melipatgandakan kecepatan dan memotong penggunaan RAM:

```csharp
// Mengabaikan alokasi change tracking snapshot
var courses = await db.Courses
    .AsNoTracking()
    .Where(c => c.Price > 100_000m)
    .ToListAsync();
```

### Custom Middleware Pipeline
Middleware dapat mengintersepsi request sebelum dan sesudah controller:

```csharp
public class RequestTimingMiddleware(RequestDelegate next, ILogger<RequestTimingMiddleware> logger)
{
    public async Task InvokeAsync(HttpContext context)
    {
        var watch = System.Diagnostics.Stopwatch.StartNew();

        await next(context); // Lanjutkan ke middleware berikutnya

        watch.Stop();
        logger.LogInformation("Request {Path} memakan waktu {ElapsedMs} ms",
            context.Request.Path, watch.ElapsedMilliseconds);
    }
}
```

## 4. Praktik kerja: implementasi Global Exception Handling Middleware

### Skenario
Terjadi error koneksi database atau unhandled exception di dalam kode service, menyebabkan ASP.NET Core menampilkan stack trace lengkap yang memuat struktur internal sistem kepada klien publik.

### Tugas
Terapkan interface bawaan .NET 8 `IExceptionHandler` untuk menangkap unhandled exception dan menghasilkan response standar `ProblemDetails` (RFC 7807).

```csharp
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Mvc;

public class GlobalExceptionHandler(ILogger<GlobalExceptionHandler> logger) : IExceptionHandler
{
    public async ValueTask<bool> TryHandleAsync(
        HttpContext httpContext,
        Exception exception,
        CancellationToken cancellationToken)
    {
        logger.LogError(exception, "Unhandled exception: {Message}", exception.Message);

        var problemDetails = new ProblemDetails
        {
            Status = StatusCodes.Status500InternalServerError,
            Title = "Server Error",
            Detail = "Terjadi kesalahan pada sistem internal kami. Tim pengembang telah diberitahu.",
            Instance = httpContext.Request.Path
        };

        httpContext.Response.StatusCode = StatusCodes.Status500InternalServerError;
        await httpContext.Response.WriteAsJsonAsync(problemDetails, cancellationToken);

        return true; // Menandakan bahwa exception telah berhasil ditangani
    }
}

// Di Program.cs:
// builder.Services.AddExceptionHandler<GlobalExceptionHandler>();
// builder.Services.AddProblemDetails();
// ...
// app.UseExceptionHandler();
```

### Lima hint bertahap
1. Di .NET 8, pendekatan paling bersih adalah mengimplementasikan `IExceptionHandler`.
2. Jangan pernah membocorkan `exception.StackTrace` ke klien di lingkungan produksi.
3. Gunakan standar `ProblemDetails` agar struktur response error dikenali oleh pustaka HTTP klien modern.
4. Pastikan memanggil `app.UseExceptionHandler()` di pipeline sebelum middleware routing atau endpoint.
5. Return `true` dari method `TryHandleAsync` untuk menandakan bahwa exception telah selesai diproses.

## 5. Lanjutan: Health Checks, Rate Limiting, dan Kestrel Tuning

- **ASP.NET Core Rate Limiting:** Manfaatkan middleware pembatasan laju bawaan di `Microsoft.AspNetCore.RateLimiting` (Fixed Window, Sliding Window, Token Bucket, Concurrency).
- **Health Checks:** Daftarkan `builder.Services.AddHealthChecks()` dan pantau kesiapan database untuk orkestrasi Kubernetes (`/healthz`).
- **Kestrel Web Server:** Server HTTP multi-platform bawaan .NET yang sangat cepat, mendukung protokol HTTP/2 dan HTTP/3 (QUIC) secara native.

## Pemeriksaan pemahaman

1. Apa bahaya menyuntikkan (inject) service bertipe Scoped ke dalam service bertipe Singleton di ASP.NET Core?
2. Kapan kita wajib menambahkan `.AsNoTracking()` pada query Entity Framework Core?
3. Apa keunggulan Minimal APIs dibanding arsitektur Controller konvensional?
4. Bagaimana fungsi `IExceptionHandler` di .NET 8 menstandardisasi respons error?

**Kunci ringkas:** Menyuntikkan Scoped ke dalam Singleton memicu masalah *captive dependency*, di mana instance scoped hidup selamanya dan dapat menyebabkan kebocoran memori atau koneksi database yang tidak pernah ditutup; `.AsNoTracking()` wajib digunakan pada query pembacaan data tanpa intensi pembaruan (update) agar memori change tracker tidak terbebani; Minimal APIs menghilangkan overhead controller, reflection, dan boilerplate sehingga throughput request per detik lebih tinggi dan kode lebih ringkas; `IExceptionHandler` menangkap exception global secara terpusat dan mengembalikan struktur RFC 7807 ProblemDetails tanpa duplikasi blok try/catch di tiap endpoint.

## Referensi

- [ASP.NET Core Official Documentation (Microsoft Learn)](https://learn.microsoft.com/en-us/aspnet/core/)
- [Entity Framework Core Documentation](https://learn.microsoft.com/en-us/ef/core/)
- [Architecture Principles in ASP.NET Core](https://learn.microsoft.com/en-us/dotnet/architecture/modern-web-apps-azure/)

# Materi: C# 12 dan Pemrograman Berorientasi Objek Modern

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat menguasai sintaks modern C# (C# 11 & 12): memahami sistem pengetikan strongly-typed, fitur Record dan Structs (Value Types vs Reference Types), LINQ (Language Integrated Query), Pattern Matching tingkat lanjut, pemrograman asynchronous dengan `async` / `await` (`Task` / `ValueTask`), serta pengelolaan memori dan Garbage Collection di lingkungan CLR.

## Prasyarat

- Memahami konsep dasar logika pemrograman dan OOP (Class, Inheritance, Polymorphism).
- Memiliki .NET SDK 8+ terpasang di sistem (`dotnet --version`).
- Memahami dasar struktur data (Array, List, Dictionary).

## 1. Dasar: arsitektur Common Language Runtime (CLR)

C# adalah bahasa modern, aman, dan berorientasi objek yang dikembangkan oleh Microsoft dan berjalan di atas runtime **CLR (Common Language Runtime)**. C# mengompilasi kode sumber menjadi IL (Intermediate Language) yang kemudian diubah menjadi instruksi mesin melalui JIT (Just-In-Time) compiler.

### Fitur Kunci C# 12:
- **Primary Constructors:** Mendeklarasikan parameter konstruktor langsung di deklarasi kelas atau struct.
- **Collection Expressions:** Sintaks inisialisasi koleksi terpadu menggunakan tanda kurung siku `[1, 2, 3]`.
- **Nullable Reference Types (`#nullable enable`):** Peringatan compile-time untuk mencegah bug `NullReferenceException`.
- **Pattern Matching Ekspresif:** Penggunaan switch expression dan relasional pattern (`>= 80`).

## 2. Praktik dasar: Record, Primary Constructor, dan LINQ

```csharp
namespace LearnHub.Core;

// 1. Record immutable dengan Primary Constructor
public record Course(
    string Id,
    string Title,
    decimal Price,
    string Category,
    bool IsPublished
);

public class CourseAnalytics
{
    public static void Main()
    {
        // 2. Collection Expression C# 12
        List<Course> courses = [
            new("c1", "Pemrograman C# 12", 250_000m, "Backend", true),
            new("c2", "ASP.NET Core Web API", 400_000m, "Backend", true),
            new("c3", "WPF Desktop Apps", 150_000m, "Desktop", false),
            new("c4", "Game Dev dengan Unity", 350_000m, "Game", true)
        ];

        // 3. Pemrosesan Data Ekspresif dengan LINQ
        var backendCourses = courses
            .Where(c => c.IsPublished && c.Category == "Backend")
            .OrderByDescending(c => c.Price)
            .Select(c => new { c.Title, FormattedPrice = $"Rp {c.Price:N0}" })
            .ToList();

        Console.WriteLine("Kursus Backend Unggulan:");
        foreach (var item in backendCourses)
        {
            Console.WriteLine($"- {item.Title} ({item.FormattedPrice})");
        }
    }
}
```

## 3. Fitur utama dan pola arsitektur

### Pattern Matching dan Switch Expression
Mengubah logika percabangan kompleks menjadi ekspresi yang sangat bersih:

```csharp
public static class DiscountCalculator
{
    public static decimal CalculateDiscount(decimal price, int userTier) => userTier switch
    {
        1 => price * 0.05m, // Tier Basic: diskon 5%
        2 => price * 0.15m, // Tier Pro: diskon 15%
        >= 3 when price > 300_000m => price * 0.25m, // VIP dengan syarat harga
        _ => 0.0m // Default: tidak ada diskon
    };
}
```

### Value Types (Struct) vs Reference Types (Class)
- **Class (Reference Type):** Dialokasikan di memori Heap; dikelola oleh Garbage Collector; disalin berdasarkan referensi pointer.
- **Struct / Record Struct (Value Type):** Dialokasikan di Stack (jika sebagai variabel lokal); disalin berdasarkan nilai (*copy by value*); tidak membebani Garbage Collector.
- **`readonly ref struct`:** Struktur memori khusus (seperti `Span<T>` dan `ReadOnlySpan<T>`) untuk pemrosesan buffer tanpa alokasi heap (*zero allocation*).

## 4. Praktik kerja: optimasi alokasi string dengan `ReadOnlySpan<char>`

### Skenario
Aplikasi memproses jutaan baris string log transaksi berformat `DATE|USER_ID|COURSE_ID|AMOUNT`. Penggunaan metode konvensional `string.Split('|')` mengalokasikan jutaan array dan string baru di Heap, memicu lag akibat siklus Garbage Collector yang sering berjalan (*GC pause*).

### Tugas
Gunakan `ReadOnlySpan<char>` untuk mem-parsing token substring tanpa membuat alokasi objek string baru sama sekali.

```csharp
using System;

public class LogParser
{
    public static void ParseLogLine(string logLine)
    {
        ReadOnlySpan<char> span = logLine.AsSpan();

        int firstPipe = span.IndexOf('|');
        ReadOnlySpan<char> date = span[..firstPipe];

        ReadOnlySpan<char> remaining = span[(firstPipe + 1)..];
        int secondPipe = remaining.IndexOf('|');
        ReadOnlySpan<char> userId = remaining[..secondPipe];

        remaining = remaining[(secondPipe + 1)..];
        int thirdPipe = remaining.IndexOf('|');
        ReadOnlySpan<char> courseId = remaining[..thirdPipe];

        ReadOnlySpan<char> amountSpan = remaining[(thirdPipe + 1)..];
        decimal amount = decimal.Parse(amountSpan);

        // Operasi di atas dilakukan langsung di stack tanpa alokasi heap baru!
        Console.WriteLine($"Parsed User: {userId.ToString()}, Course: {courseId.ToString()}, Amount: {amount:C}");
    }
}
```

### Lima hint bertahap
1. `string.Split()` selalu mengalokasikan array string baru di heap memori.
2. `ReadOnlySpan<T>` adalah representasi pointer memori kontinu di stack (*ref struct*).
3. Mengiris span (`span[..index]`) hanya mengubah pointer offset dan panjang tanpa menyalin byte data.
4. Karena `ref struct` hanya boleh hidup di stack, Span tidak dapat disimpan sebagai field di kelas biasa.
5. Gunakan `Span` untuk tugas-tugas intensif pemrosesan teks, parsing biner, dan serialisasi data.

## 5. Lanjutan: Asynchronous Task vs ValueTask, dan GC Generations

- **ValueTask<T>:** Alternatif hemat alokasi untuk method async yang sering kali selesai secara sinkron (misalnya dari cache lokal) tanpa mengalokasikan objek `Task` baru.
- **Generational Garbage Collection:** CLR membagi Heap menjadi Generasi 0 (objek berumur pendek), Generasi 1 (buffer transisi), dan Generasi 2 (objek berumur panjang seperti singleton cache).
- **IDisposable & `using` declaration:** Memastikan resource unmanaged (file handle, socket koneksi database) dibersihkan secara otomatis.

## Pemeriksaan pemahaman

1. Apa perbedaan utama antara `class` dan `struct` dalam hal alokasi memori di C#?
2. Mengapa `ReadOnlySpan<char>` lebih efisien dibanding metode `string.Substring()`?
3. Apa fungsi keyword `record` dibanding `class` biasa?
4. Kapan kita sebaiknya mengembalikan `ValueTask<T>` daripada `Task<T>`?

**Kunci ringkas:** `class` dialokasikan di heap dan dikelola oleh Garbage Collector sedangkan `struct` dialokasikan di stack tanpa overhead GC; `ReadOnlySpan` adalah slice memori tanpa alokasi objek heap baru sedangkan `Substring()` selalu mengalokasikan objek string baru; `record` menyediakan perbandingan berbasis nilai (*value equality*), deklarasi immutable ringkas, dan copy syntax (`with`); `ValueTask<T>` dianjurkan ketika hasil operasi async sebagian besar sudah tersedia secara sinkron sehingga terhindar dari alokasi objek Task di heap.

## Referensi

- [Official C# Documentation (Microsoft Learn)](https://learn.microsoft.com/en-us/dotnet/csharp/)
- [C# Language Design & Specification](https://github.com/dotnet/csharplang)
- [Pro .NET Memory Management by Konrad Kokosa](https://prodotnetmemory.com/)

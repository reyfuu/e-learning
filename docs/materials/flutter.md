# Materi: Flutter dan Dart untuk Aplikasi Multiplatform Modern

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang aplikasi lintas platform (Android, iOS, Web, Desktop) menggunakan Flutter dan Dart: memahami arsitektur pohon widget (*Widget Tree*), membedakan `StatelessWidget` dan `StatefulWidget`, menguasai manajemen state aplikasi (Provider / BLoC / Riverpod), menangani networking dan JSON serialization, serta menerapkan optimasi rendering 60/120 FPS.

## Prasyarat

- Memahami konsep dasar Object-Oriented Programming (Class, Inheritance, Interfaces).
- Memiliki Flutter SDK dan Dart terpasang di sistem (`flutter doctor`).
- Android Studio / VS Code dengan ekstensi Flutter terpasang.

## 1. Dasar: arsitektur engine Flutter dan bahasa Dart

Flutter berbeda dari framework hybrid lainnya (seperti React Native yang memakai bridge native). Flutter menggambar setiap piksel UI secara langsung menggunakan engine grafis berperforma tinggi (**Impeller** / Skia) langsung ke kanvas GPU.

### Konsep Inti Flutter & Dart:
- **"Everything is a Widget":** Elemen struktural (Button, Text), tata letak (Row, Column, Stack), dan animasi adalah widget.
- **Tiga Lapisan Pohon (Three Trees):**
  1. *Widget Tree:* Konfigurasi UI yang immutable dan murah dibuat ulang.
  2. *Element Tree:* Pengelola lifecycle yang menghubungkan widget dengan objek render.
  3. *RenderObject Tree:* Objek visual berat yang menangani penghitungan ukuran (*layout*), posisi, dan proses penggambaran (*painting*) piksel.
- **Sound Null Safety Dart:** Mencegah error null pointer crash sejak masa kompilasi (*compile-time*).

## 2. Praktik dasar: StatelessWidget, StatefulWidget, dan HTTP Fetching

```dart
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

// Model Kursus
class Course {
  final String id;
  final String title;
  final int price;

  const Course({required this.id, required this.title, required this.price});

  factory Course.fromJson(Map<String, dynamic> json) {
    return Course(
      id: json['id'] as String,
      title: json['title'] as String,
      price: json['price'] as int,
    );
  }
}

class CourseListPage extends StatefulWidget {
  const CourseListPage({super.key});

  @override
  State<CourseListPage> createState() => _CourseListPageState();
}

class _CourseListPageState extends State<CourseListPage> {
  late Future<List<Course>> _coursesFuture;

  @override
  void initState() {
    super.initState();
    _coursesFuture = fetchCourses();
  }

  Future<List<Course>> fetchCourses() async {
    final response = await http.get(Uri.parse('https://api.learnhub.id/v1/courses'));
    if (response.statusCode == 200) {
      final List<dynamic> body = jsonDecode(response.body);
      return body.map((item) => Course.fromJson(item)).toList();
    } else {
      throw Exception('Gagal memuat katalog kursus');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Katalog Belajar LearnHub')),
      body: FutureBuilder<List<Course>>(
        future: _coursesFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Terjadi error: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('Belum ada kursus tersedia'));
          }

          final courses = snapshot.data!;
          return ListView.builder(
            itemCount: courses.length,
            itemBuilder: (context, index) {
              final course = courses[index];
              return Card(
                margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: ListTile(
                  title: Text(course.title, style: const TextStyle(fontWeight: FontWeight.bold)),
                  subtitle: Text('Rp ${course.price.toString()}'),
                  trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                  onTap: () {
                    // Navigasi ke detail
                  },
                ),
              );
            },
          );
        },
      ),
    );
  }
}
```

## 3. Fitur utama dan pola arsitektur

### Pola Arsitektur BLoC (Business Logic Component)
BLoC memisahkan logika UI dari business logic menggunakan aliran data Reactive (`Stream` & `Sink`):

```text
UI (Events) → [BLoC] → (States) → UI (Re-renders)
```

### Navigasi Deklaratif: GoRouter
Untuk mengelola deep linking dan rute web URL yang rapi:

```dart
final GoRouter _router = GoRouter(
  routes: [
    GoRoute(path: '/', builder: (context, state) => const HomePage()),
    GoRoute(
      path: '/courses/:id',
      builder: (context, state) => CourseDetailPage(id: state.pathParameters['id']!),
    ),
  ],
);
```

## 4. Praktik kerja: optimasi jank dan rebuilding widget berlebih

### Skenario
Daftar modul belajar mengalami patah-patah (*UI jank / frame drop*) saat di-scroll. Tool Flutter DevTools Performance Overlay menampilkan grafik merah (render time > 16.6ms per frame).

### Tugas
Identifikasi penyebab lag (penggunaan `ListView` tanpa builder, pembuatan widget inline tanpa `const`, dan pemanggilan `setState` di level paling atas Scaffold) dan lakukan refactoring.

```dart
// SEBELUM (Buruk - me-render seluruh widget dan memicu jank):
// ListView(children: items.map((i) => ExpensiveWidget(i)).toList())

// SESUDAH (Optimal - lazy rendering hanya untuk item yang terlihat di layar):
ListView.builder(
  itemCount: items.length,
  itemExtent: 72.0, // Memberikan tinggi tetap agar engine tidak perlu mengukur ulang layout
  cacheExtent: 200.0, // Cache pra-render sedikit di luar viewport
  itemBuilder: (context, index) {
    return ModuleItemTile(
      key: ValueKey(items[index].id),
      item: items[index],
    );
  },
)

// Gunakan const constructor di setiap tempat yang memungkinkan
class ModuleItemTile extends StatelessWidget {
  final ModuleItem item;
  const ModuleItemTile({super.key, required this.item});

  @override
  Widget build(BuildContext context) {
    return const ListTile(
      leading: Icon(Icons.check_circle_outline), // const widget tidak di-rebuild!
      title: Text('Modul Pembelajaran'),
    );
  }
}
```

### Lima hint bertahap
1. Gunakan konstruktor `const` pada widget statis agar Flutter dapat menggunakan ulang instance widget yang sama di memori.
2. Selalu gunakan `ListView.builder` atau `CustomScrollView` alih-alih `ListView(children: [...])` untuk dataset yang tidak terbatas.
3. Tetapkan `itemExtent` jika tinggi setiap baris item seragam untuk menghemat komputasi layout pass.
4. Jangan pernah menempatkan logika berat atau panggilan HTTP di dalam method `build()`.
5. Isolasi pemanggilan `setState()` ke widget sekecil mungkin (*leaf node*) agar tidak merender ulang seluruh halaman.

## 5. Lanjutan: Platform Channels, Isolates, dan Profiling

- **Isolates:** Eksekusi parsing JSON berukuran raksasa atau enkripsi berat di thread Dart Isolate terpisah menggunakan fungsi `compute()` agar UI thread tidak macet.
- **Platform Channels (MethodChannel):** Jembatan komunikasi untuk memanggil fungsi native platform Android (Kotlin/Java) atau iOS (Swift/Objective-C).
- **Impeller Engine:** Rendering engine generasi terbaru Flutter yang mengeliminasi masalah kompilasi shader jank pada iOS dan Android.

## Pemeriksaan pemahaman

1. Apa perbedaan tanggung jawab antara Widget Tree, Element Tree, dan RenderObject Tree?
2. Mengapa menambahkan keyword `const` pada widget membantu performa rendering Flutter?
3. Kapan fungsi `compute()` atau Isolate perlu dipanggil?
4. Apa resiko memanggil fungsi async di dalam method `build()` widget?

**Kunci ringkas:** Widget Tree adalah konfigurasi deklaratif, Element Tree adalah penaut siklus hidup, dan RenderObject Tree adalah objek yang menghitung koordinat layout dan menggambar piksel; `const` menginstruksikan compiler untuk me-reuse satu instance objek di memori sehingga Flutter tidak perlu merekonstruksi widget saat re-render; `compute()` wajib digunakan saat memproses komputasi CPU berat (seperti decoding file besar) agar UI thread 60fps tidak mengalami lag/jank; Memanggil async di dalam `build()` akan memicu pemanggilan fungsi berulang-ulang tanpa kendali setiap kali widget di-rebuild.

## Referensi

- [Flutter Official Documentation](https://docs.flutter.dev/)
- [Dart Language Tour](https://dart.dev/guides/language/language-tour)
- [Flutter Architecture & State Management](https://flutter.dev/docs/development/data-and-backend/state-mgmt/intro)

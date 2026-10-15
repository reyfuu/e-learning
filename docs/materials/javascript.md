# Materi: JavaScript Modern (ES6+ dan Runtime Browser/Node.js)

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat memahami dasar hingga lanjutan JavaScript: sistem tipe dan cakupan variabel, manipulasi objek/array, mekanisme Event Loop dan asynchronous programming (Promise, async/await), functional & prototype inheritance, serta teknik debugging performa di browser dan server.

## Prasyarat

- Memahami konsep pemrograman dasar (variabel, percabangan, fungsi).
- Memiliki browser modern (Chrome/Firefox) dan Node.js versi 18+ terpasang di sistem lokal.
- Memahami dasar struktur HTML dan alur interaksi web.

## 1. Dasar: ekosistem dan konsep inti JavaScript

JavaScript adalah bahasa berorientasi objek berbasis prototipe dengan pengetikan dinamis (*dynamically typed*), berfitur *single-threaded* non-blocking yang didukung oleh runtime (V8 di Chrome/Node.js).

### Konsep kunci:
- **Deklarasi Variabel:** `let` dan `const` (block-scoped) menggantikan `var` (function-scoped/hoisted).
- **Tipe Data Primitif:** `string`, `number`, `boolean`, `null`, `undefined`, `symbol`, `bigint`.
- **Tipe Referensi:** `Object`, `Array`, `Function`, `Date`, `Map`, `Set`.
- **Eksekusi Asynchronous:** Menggunakan Call Stack, Web APIs (atau Node C++ APIs), Task Queue (Macrotasks), dan Microtask Queue.

## 2. Praktik dasar: sintaks modern ES6+ dan Async/Await

Berikut adalah contoh praktis penggunaan dekonstruksi (*destructuring*), *spread/rest*, arrow functions, serta pemanggilan data asynchronous dengan penanganan error yang benar:

```javascript
// Data transformasi dengan Map & Destructuring
interface User {
  id: number;
  name: string;
  role: string;
}

const users = [
  { id: 1, name: 'Budi Santoso', role: 'student', score: 85 },
  { id: 2, name: 'Siti Aminah', role: 'instructor', score: 95 },
  { id: 3, name: 'Reza Pratama', role: 'student', score: 78 }
];

// 1. Array methods fungsional (filter, map)
const activeStudents = users
  .filter(({ role }) => role === 'student')
  .map(({ id, name, score }) => ({ id, name, isPassed: score >= 80 }));

console.log('Siswa Lulus:', activeStudents);

// 2. Asynchronous fetch dengan async/await dan Promise.all
async function fetchCourseData(courseId) {
  try {
    const [detailsResponse, reviewsResponse] = await Promise.all([
      fetch(`https://api.learnhub.id/v1/courses/${courseId}`),
      fetch(`https://api.learnhub.id/v1/courses/${courseId}/reviews`)
    ]);

    if (!detailsResponse.ok || !reviewsResponse.ok) {
      throw new Error(`Gagal mengambil data kursus: ${detailsResponse.status}`);
    }

    const details = await detailsResponse.json();
    const reviews = await reviewsResponse.json();

    return { ...details, reviews };
  } catch (error) {
    console.error('Fetch error:', error.message);
    throw error;
  }
}
```

## 3. Fitur utama dan pola arsitektur

### Event Loop: Microtasks vs Macrotasks
JavaScript memprioritaskan Microtask Queue (`Promise.then`, `queueMicrotask`, `process.nextTick` di Node) sebelum memproses Macrotask Queue (`setTimeout`, `setInterval`, I/O).

```javascript
console.log('1: Synchronous');

setTimeout(() => {
  console.log('4: Macrotask (setTimeout)');
}, 0);

Promise.resolve().then(() => {
  console.log('3: Microtask (Promise)');
});

console.log('2: Synchronous');

// Urutan output: 1 -> 2 -> 3 -> 4
```

### Closures dan Prototypes
Closure memungkinkan fungsi mengakses variabel dari scope induknya bahkan setelah fungsi induk selesai dieksekusi:

```javascript
function createCounter(initialValue = 0) {
  let count = initialValue; // State privat terisolasi
  return {
    increment: () => ++count,
    decrement: () => --count,
    getValue: () => count
  };
}

const counter = createCounter(10);
console.log(counter.increment()); // 11
console.log(counter.getValue());  // 11
```

## 4. Praktik kerja: optimasi memory leak dan debouncing input

### Skenario
Halaman pencarian kursus di LearnHub mengalami penurunan performa drastis dan freeze saat pengguna mengetik kata kunci dengan cepat di input box. Setiap ketukan keyboard langsung menembak API backend.

### Tugas
Buat utilitas `debounce` murni dan integrasikan ke input event listener agar request hanya dikirim 300ms setelah pengguna selesai mengetik, serta lakukan pembersihan (*cleanup*) listener untuk mencegah memory leak.

### Langkah implementasi:
1. Buat fungsi higher-order `debounce(fn, delay)` yang menampung identifier timeout dalam closure.
2. Setiap kali dipanggil, batalkan timer sebelumnya dengan `clearTimeout`.
3. Pasang handler pada event `input` element HTML.
4. Pastikan `AbortController` digunakan untuk membatalkan HTTP fetch yang masih berjalan jika query baru dikirim.

```javascript
function debounce(fn, delayMs = 300) {
  let timerId;
  return function (...args) {
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => {
      fn.apply(this, args);
    }, delayMs);
  };
}

// Integrasi aman dengan AbortController
let currentAbortController = null;

async function searchCourses(query) {
  if (currentAbortController) {
    currentAbortController.abort(); // Batalkan request sebelumnya
  }
  currentAbortController = new AbortController();

  try {
    const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`, {
      signal: currentAbortController.signal
    });
    const results = await res.json();
    renderCourseResults(results);
  } catch (err) {
    if (err.name !== 'AbortError') {
      console.error('Search failed:', err);
    }
  }
}

const searchInput = document.getElementById('search-input');
const debouncedSearch = debounce((e) => searchCourses(e.target.value), 350);
searchInput.addEventListener('input', debouncedSearch);
```

### Lima hint bertahap
1. Event `input` terpicu pada setiap penekanan tombol; tanpa penundaan, browser mengirim puluhan HTTP request per detik.
2. Timer ID harus disimpan di dalam closure agar tidak bocor ke scope global.
3. Fungsi `clearTimeout` menghapus timer yang tertunda sebelum timer baru dibuat.
4. Gunakan `AbortController` jika ingin menghindari *race condition* response jaringan lama menimpa response baru.
5. Bersihkan event listener saat komponen atau elemen di-unmount dari DOM.

## 5. Lanjutan: performa, keamanan, dan best practices

- **Immutability:** Hindari mutasi langsung pada array dan objek state (`Object.freeze`, structured clone `structuredClone(obj)`).
- **Web Workers:** Pindahkan komputasi berat (seperti enkripsi atau parsing data biner) keluar dari main UI thread ke background worker.
- **Keamanan XSS:** Jangan pernah memasukkan input pengguna langsung ke `innerHTML`. Gunakan `textContent` atau sanitasi dengan library terpercaya.
- **Strict Mode:** Selalu gunakan `"use strict";` (atau module ES6 yang otomatis strict) untuk mendeteksi variabel tak terdeklarasi dan menghindari silent errors.

## Pemeriksaan pemahaman

1. Apa perbedaan urutan eksekusi antara Microtask (Promise) dan Macrotask (setTimeout)?
2. Bagaimana cara kerja closure dalam mempertahankan variabel privat?
3. Mengapa `==` sebaiknya dihindari dan selalu menggunakan `===`?
4. Apa manfaat `AbortController` pada operasi `fetch` di aplikasi web?

**Kunci ringkas:** Microtasks dieksekusi segera setelah stack kosong sebelum event loop mengambil macrotask berikutnya; closure menyimpan referensi scope leksikal fungsi induk; `===` memeriksa nilai dan tipe data tanpa *type coercion* yang tak terduga; `AbortController` memungkinkan pembatalan request jaringan yang tidak lagi relevan sehingga menghemat bandwidth dan mencegah *race condition*.

## Referensi

- [MDN Web Docs - JavaScript Guide](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
- [JavaScript.info - Modern JavaScript Tutorial](https://javascript.info/)
- [V8 JavaScript Engine Documentation](https://v8.dev/)

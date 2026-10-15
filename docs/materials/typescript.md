# Materi: TypeScript untuk Aplikasi Skala Produksi

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat memanfaatkan sistem tipe statis TypeScript untuk mendeteksi bug pada compile-time, mendefinisikan interface dan custom types yang ekspresif, memanfaatkan Generics, menerapkan Type Narrowing, serta merancang arsitektur type-safe di aplikasi frontend maupun backend.

## Prasyarat

- Memahami JavaScript ES6+ (sintaks objek, array, fungsi, dan async/await).
- Memiliki Node.js dan npm/pnpm terpasang.
- Mengerti konsep dasar kompilasi kode dan file `tsconfig.json`.

## 1. Dasar: mengapa TypeScript?

TypeScript adalah *superset* terketik dari JavaScript yang dikompilasi menjadi JavaScript murni. Sistem tipenya bersifat struktural (*structural subtyping* atau duck typing), yang berarti kesesuaian tipe ditentukan oleh bentuk (*shape*) data, bukan namanya saja.

### Konsep kunci:
- **Tipe Primitif & Khusus:** `string`, `number`, `boolean`, `any` (hindari), `unknown` (tipe aman untuk data yang belum diketahui), `never` (kondisi yang tidak mungkin terjadi).
- **Interface vs Type Alias:** `interface` mendukung deklarasi bertahap (*merging*) cocok untuk kontrak objek/OOP, sedangkan `type` mendukung union, intersection, dan tipe kompleks.
- **Strict Mode:** Mengaktifkan `"strict": true` di `tsconfig.json` mencakup `noImplicitAny`, `strictNullChecks`, dan validasi ketat lainnya.

## 2. Praktik dasar: Type, Interface, dan Generics

```typescript
// 1. Interface dan Type Alias
export type Role = 'admin' | 'instructor' | 'student';

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  createdAt: Date;
  metadata?: Record<string, unknown>; // Optional property
}

// 2. Generic API Response Wrapper
export interface ApiResponse<TData> {
  success: boolean;
  data: TData;
  error?: string;
  statusCode: number;
}

// Fungsi Generic dengan constraint
export async function fetchData<T>(url: string): Promise<ApiResponse<T>> {
  const response = await fetch(url);
  if (!response.ok) {
    return {
      success: false,
      data: null as unknown as T,
      error: `HTTP Error: ${response.status}`,
      statusCode: response.status
    };
  }
  const data: T = await response.json();
  return {
    success: true,
    data,
    statusCode: response.status
  };
}

// Penggunaan nyata
interface Course {
  id: string;
  title: string;
  price: number;
}

async function loadCourse(id: string) {
  const result = await fetchData<Course>(`/api/courses/${id}`);
  if (result.success) {
    console.log(result.data.title, result.data.price);
  }
}
```

## 3. Fitur utama dan pola arsitektur

### Discriminated Unions (Tagged Unions)
Teknik krusial dalam TypeScript untuk memodelkan state aplikasi tanpa ambiguitas:

```typescript
type AsyncState<T> =
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'success'; data: T }
  | { status: 'error'; error: Error };

function renderUI(state: AsyncState<User[]>) {
  switch (state.status) {
    case 'idle':
      return 'Silakan mulai pencarian';
    case 'loading':
      return 'Memuat data...';
    case 'success':
      // TypeScript secara otomatis tahu bahwa 'data' ada pada cabang ini
      return `Ditemukan ${state.data.length} pengguna`;
    case 'error':
      return `Terjadi kesalahan: ${state.error.message}`;
  }
}
```

### Utility Types
- `Partial<T>`: Membuat semua properti menjadi opsional.
- `Required<T>`: Membuat semua properti wajib.
- `Pick<T, K>`: Mengambil sebagian properti dari interface.
- `Omit<T, K>`: Menghilangkan properti tertentu dari interface.
- `Record<K, T>`: Membentuk tipe objek key-value pasangan tipe tertentu.

## 4. Praktik kerja: validasi tipe API eksternal dan Type Guard

### Skenario
Aplikasi sering mengalami runtime crash karena backend pihak ketiga mengembalikan payload dengan tipe tidak konsisten (misalnya nilai harga kadang berupa `string`, kadang `number`, atau `null`).

### Tugas
Buat *Custom Type Guard* (*User-Defined Type Guard*) dan validasi runtime agar TypeScript memastikan objek aman sebelum diproses logika bisnis.

```typescript
interface VerifiedCourse {
  id: string;
  title: string;
  price: number;
  published: boolean;
}

// Custom Type Guard menggunakan predikat tipe 'is'
function isVerifiedCourse(obj: unknown): obj is VerifiedCourse {
  if (typeof obj !== 'object' || obj === null) {
    return false;
  }

  const record = obj as Record<string, unknown>;

  return (
    typeof record.id === 'string' &&
    typeof record.title === 'string' &&
    typeof record.price === 'number' &&
    !isNaN(record.price) &&
    typeof record.published === 'boolean'
  );
}

// Pemakaian
async function processExternalCourse(input: unknown) {
  if (isVerifiedCourse(input)) {
    // Di dalam blok ini, TypeScript menjamin input adalah VerifiedCourse
    console.log(`Kursus Valid: ${input.title} - Rp ${input.price.toLocaleString('id-ID')}`);
  } else {
    throw new Error('Payload kursus korup atau tidak memenuhi spesifikasi');
  }
}
```

### Lima hint bertahap
1. Jangan gunakan `any`; gunakan `unknown` untuk data yang belum divalidasi.
2. Predikat tipe menggunakan sintaks `parameterName is TargetType` pada return value fungsi.
3. Selalu periksa `typeof obj === 'object' && obj !== null` sebelum mengakses properti objek.
4. Untuk skenario aplikasi besar, integrasikan library schema validation seperti Zod atau ArkType yang dapat menginferensi tipe TypeScript secara otomatis.
5. Gunakan `satisfies` operator jika ingin memverifikasi bahwa tipe cocok tanpa mengubah tipe aslinya menjadi union yang lebar.

## 5. Lanjutan: Conditional Types dan Mapped Types

```typescript
// Mapped Type: Mengubah semua properti menjadi fungsi getter
type Getters<T> = {
  [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K];
};

// Conditional Type: Mengekstrak tipe dari Promise
type AwaitedType<T> = T extends Promise<infer U> ? U : T;
```

## Pemeriksaan pemahaman

1. Mengapa tipe `unknown` lebih direkomendasikan daripada `any`?
2. Apa fungsi Discriminated Union dalam pengelolaan state UI?
3. Kapan sebaiknya menggunakan `interface` dibanding `type` alias?
4. Bagaimana cara kerja operator `keyof` dan `typeof` dalam TypeScript?

**Kunci ringkas:** `unknown` menuntut pengecekan tipe eksplisit sebelum dapat dimanipulasi sedangkan `any` mematikan semua type checking; Discriminated Union menghilangkan kondisi state yang mustahil (*impossible states*) melalui satu properti pembeda; `interface` lebih unggul untuk OOP dan deklarasi terbuka (merging), sedangkan `type` lebih fleksibel untuk primitive unions dan mapped types; `typeof` mengekstrak tipe dari nilai variabel yang ada, sedangkan `keyof` menghasilkan union dari nama properti objek.

## Referensi

- [Official TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Total TypeScript by Matt Pocock](https://www.totaltypescript.com/)
- [TypeScript Deep Dive](https://basarat.gitbook.io/typescript/)

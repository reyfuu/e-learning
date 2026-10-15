# Materi: Vue.js 3 dan Ekosistem Composition API

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun antarmuka web reaktif dengan Vue 3: menguasai Composition API (`ref`, `reactive`, `computed`, `watch`), membuat reusable *composables*, mengelola state aplikasi dengan Pinia, menangani routing dengan Vue Router, serta menerapkan strategi optimasi rendering.

## Prasyarat

- Memahami JavaScript modern (ES6 modules, arrow functions, destructuring).
- Memahami dasar DOM, event handling, dan CSS.
- Node.js versi 18+ terpasang di sistem.

## 1. Dasar: reaktivitas Vue 3 dan Composition API

Vue 3 menggunakan JavaScript `Proxy` untuk sistem reaktivitasnya. Pendekatan **Composition API** (menggunakan `<script setup>`) memungkinkan pengelompokan logika berdasarkan fitur (*feature-based logic*), berbeda dengan Options API yang memisahkan data, methods, dan computed ke blok terpisah.

### Konsep Kunci:
- `ref()`: Membungkus nilai primitif atau objek ke dalam ref object. Diakses via `.value` di script, otomatis di-*unwrap* di template.
- `reactive()`: Membungkus objek JavaScript langsung ke Proxy object (tidak memerlukan `.value`, tetapi tidak boleh di-destructure tanpa `toRefs`).
- `computed()`: Nilai turunan ter-cache yang hanya dihitung ulang jika dependensinya berubah.
- `watch()` / `watchEffect()`: Menjalankan side-effect saat variabel reaktif berubah.

## 2. Praktik dasar: Single File Component (SFC) dengan `<script setup>`

```vue
<!-- components/CourseCard.vue -->
<script setup lang="ts">
import { ref, computed } from 'vue';

interface Props {
  title: string;
  price: number;
  discountPercentage?: number;
}

const props = withDefaults(defineProps<Props>(), {
  discountPercentage: 0
});

const emit = defineEmits<{
  (e: 'enroll', courseTitle: string): void;
}>();

const isFavorite = ref(false);

const finalPrice = computed(() => {
  if (props.discountPercentage > 0) {
    return props.price - (props.price * props.discountPercentage) / 100;
  }
  return props.price;
});

function toggleFavorite() {
  isFavorite.value = !isFavorite.value;
}

function handleEnroll() {
  emit('enroll', props.title);
}
</script>

<template>
  <div class="border rounded-xl p-5 shadow-sm hover:shadow-md transition">
    <div class="flex justify-between items-start">
      <h3 class="text-lg font-bold">{{ title }}</h3>
      <button @click="toggleFavorite" class="text-xl">
        {{ isFavorite ? '❤️' : '🤍' }}
      </button>
    </div>

    <div class="mt-3">
      <span v-if="discountPercentage > 0" class="text-sm line-through text-gray-400 mr-2">
        Rp {{ price.toLocaleString('id-ID') }}
      </span>
      <span class="text-xl font-bold text-emerald-600">
        Rp {{ finalPrice.toLocaleString('id-ID') }}
      </span>
    </div>

    <button
      @click="handleEnroll"
      class="mt-4 w-full bg-emerald-600 text-white font-medium py-2 rounded-lg hover:bg-emerald-700"
    >
      Daftar Kursus
    </button>
  </div>
</template>
```

## 3. Fitur utama dan pola arsitektur

### Pola Reusable Composables
Composables adalah fungsi konvensi Vue (diawali `use...`) yang memanfaatkan Composition API untuk mengekapsulasi stateful logic:

```typescript
// composables/useFetchCourses.ts
import { ref, onMounted } from 'vue';

export function useFetchCourses() {
  const courses = ref<any[]>([]);
  const isLoading = ref(true);
  const error = ref<string | null>(null);

  async function loadCourses() {
    isLoading.value = true;
    error.value = null;
    try {
      const res = await fetch('/api/courses');
      if (!res.ok) throw new Error('Gagal memuat data');
      courses.value = await res.json();
    } catch (err: any) {
      error.value = err.message || 'Terjadi kesalahan jaringan';
    } finally {
      isLoading.value = false;
    }
  }

  onMounted(loadCourses);

  return { courses, isLoading, error, refetch: loadCourses };
}
```

### State Management dengan Pinia
Pinia adalah store resmi Vue 3 yang ringan dan fully type-safe:

```typescript
// stores/cart.ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useCartStore = defineStore('cart', () => {
  const items = ref<string[]>([]);

  const itemCount = computed(() => items.value.length);

  function addItem(courseId: string) {
    if (!items.value.includes(courseId)) {
      items.value.push(courseId);
    }
  }

  function removeItem(courseId: string) {
    items.value = items.value.filter((id) => id !== courseId);
  }

  return { items, itemCount, addItem, removeItem };
});
```

## 4. Praktik kerja: optimasi form reaktif dan debounce pencarian

### Skenario
Halaman eksplorasi kursus memiliki filter kategori dan pencarian keyword. Input teks yang memicu reaktivitas langsung sering membuat UI berkedip karena request berulang.

### Tugas
Buat composable `useDebouncedRef` yang menunda perubahan state reaktif sampai pengguna berhenti mengetik selama periode waktu tertentu.

```typescript
// composables/useDebouncedRef.ts
import { customRef } from 'vue';

export function useDebouncedRef<T>(value: T, delay = 300) {
  let timeout: ReturnType<typeof setTimeout>;
  return customRef((track, trigger) => {
    return {
      get() {
        track(); // Daftarkan dependensi reaktif
        return value;
      },
      set(newValue: T) {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
          value = newValue;
          trigger(); // Picu update reaktif setelah jeda delay
        }, delay);
      }
    };
  });
}
```

### Lima hint bertahap
1. Gunakan fungsi `customRef` bawaan Vue untuk mengontrol kapan `track()` dan `trigger()` dipanggil.
2. `track()` harus dipanggil pada getter agar Vue mencatat pembaca nilai.
3. `trigger()` harus dipanggil di dalam callback `setTimeout` setelah nilai diperbarui.
4. Jangan lupa `clearTimeout` jika setter dipanggil lagi sebelum timeout selesai.
5. Gunakan composable ini dengan `v-model` pada elemen `<input>` seperti ref biasa.

## 5. Lanjutan: performa, teleports, dan testing

- **Teleport:** Mengirim elemen DOM modal/dialog ke `<body>` tanpa memutus hierarki komponen Vue (`<Teleport to="body">`).
- **KeepAlive:** Menyimpan instance komponen dalam memori saat berpindah tab untuk menjaga input pengguna tetap utuh (`<KeepAlive>`).
- **Vue Test Utils & Vitest:** Uji komponen secara unit test dengan memeriksa emit event dan rendering teks.

## Pemeriksaan pemahaman

1. Kapan harus menggunakan `ref()` dan kapan menggunakan `reactive()`?
2. Apa perbedaan utama antara `watch` dan `watchEffect`?
3. Mengapa Pinia menggantikan Vuex di ekosistem Vue 3?
4. Bagaimana fungsi `customRef` membantu optimasi performa form input?

**Kunci ringkas:** `ref()` cocok untuk semua tipe data (primitif maupun objek) dan lebih aman saat di-destructure, sedangkan `reactive()` hanya untuk objek; `watch` memerlukan sumber eksplisit dan memberi akses old vs new value, sedangkan `watchEffect` otomatis melacak semua dependensi di dalam callback; Pinia lebih sederhana, tidak membutuhkan mutations, dan sepenuhnya type-safe; `customRef` memungkinkan pengaturan timing `trigger()` seperti debouncing secara native.

## Referensi

- [Vue.js Official Documentation](https://vuejs.org/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vue Router Documentation](https://router.vuejs.org/)

# Materi: React Core, Hooks, dan Arsitektur Komponen

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang aplikasi berbasis komponen React dengan pendekatan deklaratif: memahami siklus render dan rekonsiliasi Virtual DOM, menguasai hooks bawaan (`useState`, `useEffect`, `useMemo`, `useCallback`, `useRef`, `useReducer`), merancang custom hooks, serta menerapkan pola optimasi rendering dan state management.

## Prasyarat

- Pemahaman solid tentang JavaScript modern (ES6, closures, destructuring, promises).
- Konsep dasar HTML DOM dan event handling.
- Node.js versi 18+ terpasang di sistem.

## 1. Dasar: model mental React dan Virtual DOM

React bekerja secara deklaratif: antarmuka UI adalah fungsi dari state (`UI = f(state)`). Ketika state berubah, React menjalankan proses render ulang (*re-render*), membandingkan Virtual DOM lama dengan Virtual DOM baru (*reconciliation / Fiber algorithm*), dan hanya memperbarui elemen DOM browser yang benar-benar berubah (*commit phase*).

### Aturan Utama Hooks:
1. Hanya panggil Hooks di level paling atas (*top-level*) fungsi komponen (jangan di dalam loop, percabangan `if`, atau fungsi bersarang).
2. Hanya panggil Hooks dari fungsi komponen React atau Custom Hook.
3. State harus diperlakukan sebagai nilai *immutable* (jangan pernah melakukan `state.push()` atau `state.property = value` secara langsung).

## 2. Praktik dasar: Hooks Esensial dan Form Handling

```tsx
import React, { useState, useEffect, useMemo, useCallback } from 'react';

interface Lesson {
  id: number;
  title: string;
  durationMinutes: number;
  completed: boolean;
}

export function CourseProgressTracker({ initialLessons }: { initialLessons: Lesson[] }) {
  const [lessons, setLessons] = useState<Lesson[]>(initialLessons);
  const [filterQuery, setFilterQuery] = useState('');

  // Toggle status pelajaran tanpa memutasi array asli
  const toggleLesson = useCallback((id: number) => {
    setLessons((prev) =>
      prev.map((item) =>
        item.id === id ? { ...item, completed: !item.completed } : item
      )
    );
  }, []);

  // Perhitungan ter-memoize untuk mencegah re-kalkulasi yang tidak perlu
  const stats = useMemo(() => {
    const total = lessons.length;
    const completedCount = lessons.filter((l) => l.completed).length;
    const percentage = total === 0 ? 0 : Math.round((completedCount / total) * 100);
    return { total, completedCount, percentage };
  }, [lessons]);

  // Filter lessons
  const filteredLessons = useMemo(() => {
    return lessons.filter((l) =>
      l.title.toLowerCase().includes(filterQuery.toLowerCase())
    );
  }, [lessons, filterQuery]);

  return (
    <div className="p-6 bg-white rounded-xl shadow border">
      <div className="mb-4">
        <h2 className="text-xl font-bold">Progres Belajar: {stats.percentage}%</h2>
        <div className="w-full bg-gray-200 h-3 rounded-full mt-2">
          <div
            className="bg-blue-600 h-3 rounded-full transition-all"
            style={{ width: `${stats.percentage}%` }}
          />
        </div>
        <p className="text-sm text-gray-500 mt-1">
          {stats.completedCount} dari {stats.total} modul selesai.
        </p>
      </div>

      <input
        type="text"
        placeholder="Cari materi..."
        value={filterQuery}
        onChange={(e) => setFilterQuery(e.target.value)}
        className="w-full p-2 border rounded mb-4"
      />

      <ul className="space-y-2">
        {filteredLessons.map((lesson) => (
          <li
            key={lesson.id}
            onClick={() => toggleLesson(lesson.id)}
            className={`p-3 rounded border cursor-pointer flex justify-between ${
              lesson.completed ? 'bg-green-50 border-green-300' : 'bg-gray-50'
            }`}
          >
            <span>{lesson.title} ({lesson.durationMinutes} mnt)</span>
            <span>{lesson.completed ? '✅' : '⭕'}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

## 3. Fitur utama dan pola arsitektur

### Pola Custom Hook untuk Abstraksi Logika
Pisahkan logika bisnis dari UI presentasional:

```tsx
// hooks/useLocalStorage.ts
import { useState, useEffect } from 'react';

export function useLocalStorage<T>(key: string, initialValue: T) {
  const [storedValue, setStoredValue] = useState<T>(() => {
    if (typeof window === 'undefined') return initialValue;
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.warn(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  });

  useEffect(() => {
    try {
      window.localStorage.setItem(key, JSON.stringify(storedValue));
    } catch (error) {
      console.warn(`Error setting localStorage key "${key}":`, error);
    }
  }, [key, storedValue]);

  return [storedValue, setStoredValue] as const;
}
```

### Context API vs External Store (Zustand)
Gunakan Context API hanya untuk data yang jarang berubah (tema gelap/terang, locale bahasa). Untuk state sering berubah dengan performa tinggi, gunakan store eksternal seperti Zustand yang mendukung *selective subscription*:

```typescript
// store/useAuthStore.ts
import { create } from 'zustand';

interface AuthState {
  user: { id: string; name: string } | null;
  setUser: (user: { id: string; name: string } | null) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  setUser: (user) => set({ user }),
}));
```

## 4. Praktik kerja: mencegah unnecessary re-renders dan cascading updates

### Skenario
Daftar 1.000 pertanyaan kuis sering mengalami lag input saat pengguna memilih radio button jawaban. Setiap pilihan memicu render ulang seluruh daftar 1.000 komponen kuis anak.

### Tugas
Terapkan `React.memo` dengan perbandingan props yang tepat, serta pastikan callback fungsi dipertahankan referensinya menggunakan `useCallback` agar hanya item yang berubah saja yang dirender ulang.

```tsx
import React, { memo } from 'react';

interface QuestionItemProps {
  id: number;
  text: string;
  isSelected: boolean;
  onSelect: (id: number) => void;
}

// Dibungkus React.memo agar tidak re-render jika isSelected dan onSelect tetap sama
export const QuestionItem = memo(function QuestionItem({
  id,
  text,
  isSelected,
  onSelect
}: QuestionItemProps) {
  console.log(`Render Question ${id}`);
  return (
    <div
      onClick={() => onSelect(id)}
      className={`p-3 border rounded ${isSelected ? 'bg-indigo-50 border-indigo-500' : ''}`}
    >
      <input type="radio" checked={isSelected} readOnly className="mr-2" />
      {text}
    </div>
  );
});
```

### Lima hint bertahap
1. Buat profil render menggunakan React DevTools Profiler untuk melihat komponen mana yang memakan waktu render terlama.
2. `React.memo` melakukan perbandingan *shallow* pada props; jika props berupa objek atau fungsi inline baru, memo akan gagal.
3. Bungkus fungsi event handler di komponen induk dengan `useCallback`.
4. Jika daftar sangat panjang (ribuan baris), pertimbangkan virtualisasi DOM menggunakan `@tanstack/react-virtual`.
5. Hindari memindahkan state lokal ke level global jika state tersebut hanya relevan di satu komponen anak.

## 5. Lanjutan: Concurrent React, Suspense, dan Error Boundaries

- **Error Boundaries:** Komponen kelas (`componentDidCatch` / `getDerivedStateFromError`) untuk menangkap error rendering dan menampilkan UI fallback yang elegan tanpa merusak seluruh layar aplikasi.
- **useTransition & useDeferredValue:** Pisahkan pembaruan state mendesak (ketikan input) dari pembaruan berprioritas rendah (render filter berat) agar browser tidak macet.
- **Strict Mode:** Menjalankan efek dua kali pada mode *development* untuk membantu mendeteksi side-effects yang tidak dibersihkan dengan benar (*cleanup functions*).

## Pemeriksaan pemahaman

1. Apa perbedaan fase Render dan fase Commit pada siklus eksekusi React?
2. Mengapa mutasi langsung pada array state menyebabkan bug rendering di React?
3. Kapan sebaiknya menggunakan `useCallback` dan kapan penggunaannya justru menambah overhead?
4. Bagaimana mekanisme kerja `React.memo` dalam membatasi re-render komponen anak?

**Kunci ringkas:** Fase render menghitung perbandingan Virtual DOM tanpa efek visual, sedangkan commit menerapkan perubahan tersebut ke DOM fisik browser; mutasi langsung tidak mengubah referensi objek sehingga perbandingan kesetaraan (`Object.is`) menganggap state tidak berubah; `useCallback` bermanfaat saat fungsi dioper ke komponen anak ter-memoize atau array dependensi useEffect, namun tidak diperlukan untuk komponen DOM primitif biasa; `React.memo` mengecek kesetaraan *shallow* dari props lama dan baru sebelum memutuskan merender ulang.

## Referensi

- [React Official Documentation](https://react.dev/)
- [Overreacted by Dan Abramov](https://overreacted.io/)
- [TanStack Query & Virtual Docs](https://tanstack.com/)

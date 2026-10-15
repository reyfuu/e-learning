# Materi: Firebase Platform dan Cloud Firestore NoSQL

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang backend berbasis serverless menggunakan Google Firebase: mengintegrasikan Firebase Authentication, menguasai Cloud Firestore NoSQL database (Documents & Collections, Realtime Listeners, Subcollections), menulis Firebase Security Rules yang ketat untuk otorisasi akses, mengimplementasikan Cloud Storage untuk upload file, serta memahami Cloud Functions untuk komputasi backend terisolasi.

## Prasyarat

- Memahami JavaScript / TypeScript modern (Promise, async/await).
- Mengerti konsep NoSQL dan arsitektur aplikasi Single Page Application / Mobile.
- Memiliki akun Google dan akses ke Firebase Console.

## 1. Dasar: ekosistem Firebase dan Cloud Firestore

Firebase adalah platform pengembangan aplikasi komprehensif milik Google. Di pusat ekosistemnya terdapat **Cloud Firestore**, database dokumen NoSQL yang dirancang untuk penskalaan otomatis, sinkronisasi data real-time, dan dukungan operasi luring (*offline sync*).

### Komponen Utama Firebase:
- **Firebase Authentication:** Manajemen identitas pengguna siap pakai (Email/Password, Google Sign-In, Apple, nomor telepon).
- **Cloud Firestore:** Database dokumen fleksibel dan terdistribusi global yang mendukung query kompleks dan snapshot listener real-time.
- **Cloud Storage for Firebase:** Penyimpanan objek biner (gambar, video, PDF) dengan integrasi Google Cloud Storage.
- **Cloud Functions:** Serverless microservice berbasis Node.js/Python yang terpicu secara otomatis oleh event database atau HTTP call.
- **Security Rules:** Logika deklaratif yang dieksekusi di server Google untuk mengontrol izin baca/tulis langsung dari perangkat klien.

## 2. Praktik dasar: Firestore CRUD dan Realtime Listeners (SDK v9+ Modular)

```typescript
import { initializeApp } from 'firebase/app';
import {
  getFirestore,
  collection,
  doc,
  setDoc,
  getDocs,
  query,
  where,
  orderBy,
  onSnapshot
} from 'firebase/firestore';

const firebaseConfig = {
  apiKey: "AIzaSy...",
  authDomain: "learnhub-app.firebaseapp.com",
  projectId: "learnhub-app",
  storageBucket: "learnhub-app.appspot.com",
};

// Inisialisasi Firebase App & Firestore
const app = initializeApp(firebaseConfig);
const db = getFirestore(app);

// 1. Menambahkan Dokumen Kursus Baru
async function addCourse(courseId: string, title: string, price: number, instructorId: string) {
  const courseRef = doc(db, 'courses', courseId);
  await setDoc(courseRef, {
    title,
    price,
    instructorId,
    isPublished: true,
    createdAt: new Date()
  });
  console.log('Kursus berhasil disimpan!');
}

// 2. Berlangganan Realtime Updates (onSnapshot)
function listenToPublishedCourses(onUpdate: (courses: any[]) => void) {
  const q = query(
    collection(db, 'courses'),
    where('isPublished', '==', true),
    orderBy('createdAt', 'desc')
  );

  // Listener aktif: dipanggil seketika saat data di database berubah
  const unsubscribe = onSnapshot(q, (snapshot) => {
    const courses = snapshot.docs.map((d) => ({ id: d.id, ...d.data() }));
    onUpdate(courses);
  }, (error) => {
    console.error('Realtime listener error:', error);
  });

  return unsubscribe; // Panggil fungsi ini saat komponen unmount untuk stop listener
}
```

## 3. Fitur utama dan pola arsitektur

### Firebase Security Rules: Menjaga Keamanan Langsung dari Klien
Karena browser dan aplikasi mobile berkomunikasi langsung dengan Firestore tanpa perantara server konvensional, **Security Rules adalah satu-satunya dinding pertahanan**:

```javascript
// firestore.rules
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {

    // Helper functions
    function isAuthenticated() {
      return request.auth != null;
    }
    
    function isOwner(userId) {
      return isAuthenticated() && request.auth.uid == userId;
    }

    // Aturan koleksi users
    match /users/{userId} {
      allow read: if isAuthenticated();
      allow write: if isOwner(userId);
    }

    // Aturan koleksi courses
    match /courses/{courseId} {
      // Siapa pun (termasuk anonim) boleh membaca kursus yang dipublikasikan
      allow read: if resource.data.isPublished == true;

      // Hanya user terotentikasi yang boleh membuat kursus baru
      allow create: if isAuthenticated() 
        && request.resource.data.instructorId == request.auth.uid
        && request.resource.data.price >= 0;

      // Hanya instruktur pemilik yang boleh mengubah atau menghapus
      allow update, delete: if isAuthenticated() 
        && resource.data.instructorId == request.auth.uid;
    }
  }
}
```

### Cloud Functions untuk Operasi Terproteksi
Jangan mengeksekusi operasi kritis (seperti verifikasi webhook pembayaran) di klien; delegasikan ke Cloud Function:

```typescript
// functions/src/index.ts
import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';

admin.initializeApp();

export const onPaymentWebhook = functions.https.onRequest(async (req, res) => {
  const { orderId, status } = req.body;
  if (status === 'SUCCESS') {
    await admin.firestore().collection('orders').doc(orderId).update({
      paymentStatus: 'PAID',
      paidAt: admin.firestore.FieldValue.serverTimestamp()
    });
  }
  res.status(200).send('Webhook processed');
});
```

## 4. Praktik kerja: optimasi tagihan dan mencegah infinite read loop

### Skenario
Biaya tagihan Firestore melonjak ribuan dolar dalam 2 hari karena tim pengembang tidak sengaja memasang `onSnapshot` di dalam hook `useEffect` React tanpa array dependensi, menyebabkan listener dipasang berulang kali pada setiap render dan menghasilkan jutaan operasi pembacaan dokumen.

### Tugas
Perbaiki *lifecycle* listener, pastikan fungsi `unsubscribe()` selalu dipanggil pada pembersihan efek, serta gunakan teknik agregasi data (*denormalization*) untuk mengurangi frekuensi pembacaan dokumen berulang.

```tsx
import { useEffect, useState } from 'react';

export function CourseList() {
  const [courses, setCourses] = useState([]);

  useEffect(() => {
    // 1. Pasang listener hanya satu kali saat mount
    const unsubscribe = listenToPublishedCourses((data) => {
      setCourses(data);
    });

    // 2. WAJIB: Kembalikan fungsi cleanup agar listener dilepas saat unmount
    return () => {
      unsubscribe();
    };
  }, []); // Array dependensi kosong menjamin efek hanya berjalan sekali

  return <div>{courses.length} Kursus Aktif</div>;
}
```

### Lima hint bertahap
1. Firestore mengenakan biaya berdasarkan jumlah dokumen yang dibaca (*reads*), ditulis (*writes*), dan dihapus (*deletes*).
2. Setiap kali listener `onSnapshot` pertama kali aktif, ia membaca seluruh dokumen yang cocok dengan kueri.
3. Selalu simpan fungsi pembatalan (*unsubscribe*) yang dikembalikan oleh `onSnapshot` dan panggil saat komponen dilepas.
4. Jangan pernah menaruh pemanggilan kueri di dalam loop; gunakan operator `whereIn` (maksimal 30 ID) atau desain subkoleksi.
5. Manfaatkan Firestore Cache lokal (`cacheSizeBytes: CACHE_SIZE_UNLIMITED`) untuk menyajikan data yang tidak berubah tanpa biaya read tambahan.

## 5. Lanjutan: Firestore Indexes, App Check, dan Backup Otomatis

- **Composite Indexes:** Kueri Firestore dengan filter ganda atau kombinasi `where` dan `orderBy` memerlukan Composite Index; Firestore CLI menyediakan link otomatis untuk membuatnya di konsol.
- **Firebase App Check:** Melindungi backend Firestore dan Storage dari serangan bot dan request ilegal dengan memvalidasi keaslian aplikasi via Google Play Integrity atau Apple DeviceCheck.
- **Ekspor Otomatis ke BigQuery:** Streaming data Firestore secara otomatis ke Google BigQuery untuk analitik data dan reporting tanpa membebani performa database operasional.

## Pemeriksaan pemahaman

1. Apa resiko keamanan terbesar jika aplikasi memanggil Cloud Firestore langsung dari aplikasi mobile tanpa Security Rules?
2. Bagaimana cara kerja perhitungan biaya (billing) di Cloud Firestore?
3. Apa fungsi mengembalikan callback `unsubscribe()` dari `onSnapshot`?
4. Mengapa operasi pembayaran harus dieksekusi di Cloud Functions dan bukan di SDK klien?

**Kunci ringkas:** Tanpa Security Rules, siapa pun yang memiliki API key dapat membaca, menimpa, atau menghapus seluruh database secara langsung melalui REST endpoint; Firestore menagih biaya per jumlah operasi read, write, dan delete dokumen serta konsumsi penyimpanan data; Callback unsubscribe menghentikan koneksi listener WebSocket agar aplikasi tidak terus menerima update dan tidak memicu biaya read yang tak terkontrol; Cloud Functions berjalan di lingkungan server yang aman (*admin privileges*) sehingga kode dan kunci rahasia API payment gateway tidak dapat dimanipulasi oleh klien.

## Referensi

- [Firebase Official Documentation](https://firebase.google.com/docs)
- [Cloud Firestore Security Rules Guide](https://firebase.google.com/docs/firestore/security/get-started)
- [Firebase Architecture and Best Practices](https://firebase.google.com/docs/guides)

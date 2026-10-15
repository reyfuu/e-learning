# Materi: Express.js untuk RESTful API dan Middleware Pipeline

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang backend RESTful API berbasis Node.js menggunakan Express.js: memahami alur middleware pipeline, menangani validasi request body, mengelola error handling terpusat, mengimplementasikan otentikasi JWT dengan aman, serta menerapkan praktik hardening keamanan HTTP (CORS, Helmet, Rate Limiter).

## Prasyarat

- Memahami JavaScript modern (ES Modules, Async/Await, Destructuring).
- Memahami konsep dasar protokol HTTP (method, status code, headers, JSON body).
- Node.js versi 18+ terpasang di komputer lokal.

## 1. Dasar: siklus request-response dan middleware

Express.js adalah framework minimalis dan fleksibel untuk Node.js. Inti arsitektur Express dibangun di atas konsep **Middleware Pipeline**: fungsi-fungsi berantai yang menerima objek `req` (Request), `res` (Response), dan `next` (fungsi pemicu middleware berikutnya).

```text
Incoming Request → [CORS] → [Body Parser] → [Auth Guard] → [Route Handler] → [Error Handler] → Response
```

### Karakteristik Middleware:
- Dapat memodifikasi objek `req` atau `res` (misal menambahkan `req.user`).
- Menjalankan kode tertentu (logging, validasi, otentikasi).
- Menghentikan siklus request dengan mengirim respon (`res.json(...)`).
- Memanggil `next()` untuk meneruskan kontrol ke middleware berikutnya.

## 2. Praktik dasar: setup Express Server dengan Zod Validation

```typescript
// server.ts
import express, { Request, Response, NextFunction } from 'express';
import helmet from 'helmet';
import cors from 'cors';
import { z } from 'zod';

const app = express();

// 1. Security & Parsing Middlewares
app.use(helmet());
app.use(cors({ origin: 'https://learnhub.id' }));
app.use(express.json({ limit: '10kb' })); // Batasi ukuran payload untuk cegah DoS

// 2. Skema Validasi Request dengan Zod
const CreateCourseSchema = z.object({
  title: z.string().min(5, 'Judul minimal 5 karakter'),
  price: z.number().nonnegative('Harga tidak boleh negatif'),
  category: z.enum(['frontend', 'backend', 'devops', 'mobile'])
});

// 3. Reusable Validation Middleware
function validateBody(schema: z.ZodSchema) {
  return (req: Request, res: Response, next: NextFunction) => {
    const result = schema.safeParse(req.body);
    if (!result.success) {
      return res.status(400).json({
        success: false,
        message: 'Validasi input gagal',
        errors: result.error.flatten().fieldErrors
      });
    }
    req.body = result.data;
    next();
  };
}

// 4. Route Handler
app.post('/api/v1/courses', validateBody(CreateCourseSchema), (req: Request, res: Response) => {
  const { title, price, category } = req.body;
  // Simpan ke database...
  res.status(201).json({
    success: true,
    data: { id: 'crs_123', title, price, category }
  });
});

// 5. Global Error Handling Middleware (wajib 4 parameter)
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  console.error('[Unhandled Error]:', err.stack);
  res.status(500).json({
    success: false,
    message: 'Internal server error'
  });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server aktif di http://localhost:${PORT}`);
});
```

## 3. Fitur utama dan pola arsitektur

### Pola Router Modular
Bagi rute ke dalam modul-modul terpisah untuk menjaga keterbacaan:

```typescript
// routes/courseRoutes.ts
import { Router } from 'express';
const router = Router();

router.get('/', (req, res) => res.json({ courses: [] }));
router.get('/:id', (req, res) => res.json({ courseId: req.params.id }));

export default router;

// Di file utama:
// app.use('/api/v1/courses', courseRouter);
```

### Async Error Handling Tanpa Try-Catch Berulang
Mulai Express 5 (atau dengan wrapper di Express 4), promise rejection harus diarahkan ke error handler:

```typescript
const asyncHandler = (fn: Function) => (req: Request, res: Response, next: NextFunction) => {
  Promise.resolve(fn(req, res, next)).catch(next);
};

router.get('/data', asyncHandler(async (req: Request, res: Response) => {
  const data = await fetchDatabaseData();
  res.json(data);
}));
```

## 4. Praktik kerja: implementasi JWT Authentication & Role-Based Access Control (RBAC)

### Skenario
Endpoint untuk menghapus modul kursus dapat diakses oleh siapa saja karena tidak adanya verifikasi token JWT dan otorisasi role instruktur.

### Tugas
Buat middleware `authenticateToken` yang memverifikasi access token di header `Authorization: Bearer <token>`, serta middleware `authorizeRole` yang memeriksa role pengguna.

```typescript
import jwt from 'jsonwebtoken';

interface JwtPayload {
  userId: string;
  role: 'admin' | 'instructor' | 'student';
}

declare global {
  namespace Express {
    interface Request {
      user?: JwtPayload;
    }
  }
}

export function authenticateToken(req: Request, res: Response, next: NextFunction) {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];

  if (!token) {
    return res.status(401).json({ success: false, message: 'Access token tidak ditemukan' });
  }

  jwt.verify(token, process.env.JWT_SECRET || 'secret_key', (err, decoded) => {
    if (err) {
      return res.status(403).json({ success: false, message: 'Token kadaluarsa atau tidak valid' });
    }
    req.user = decoded as JwtPayload;
    next();
  });
}

export function authorizeRole(allowedRoles: string[]) {
  return (req: Request, res: Response, next: NextFunction) => {
    if (!req.user || !allowedRoles.includes(req.user.role)) {
      return res.status(403).json({
        success: false,
        message: 'Akses ditolak: Anda tidak memiliki izin untuk tindakan ini'
      });
    }
    next();
  };
}

// Penggunaan di rute penghapusan
router.delete('/:id', authenticateToken, authorizeRole(['admin', 'instructor']), (req, res) => {
  res.json({ success: true, message: 'Modul berhasil dihapus' });
});
```

### Lima hint bertahap
1. Ambil token dari header `Authorization` dengan pola `Bearer <token>`.
2. Kembalikan kode status HTTP `401 Unauthorized` jika token tidak ada, dan `403 Forbidden` jika token rusak/role tidak cocok.
3. Jangan pernah menyimpan rahasia JWT (*secret key*) di kode sumber; gunakan environment variables.
4. Perluas tipe `Express.Request` dengan TypeScript `declare global` agar `req.user` dikenali tanpa error linting.
5. Error handler global harus selalu memiliki 4 parameter: `(err, req, res, next)`.

## 5. Lanjutan: rate limiting, compression, dan graceful shutdown

- **express-rate-limit:** Batasi request dari IP yang sama untuk mencegah serangan brute-force login.
- **Graceful Shutdown:** Tangani sinyal OS `SIGTERM` dan `SIGINT` untuk menutup koneksi database dan menyelesaikan request yang sedang berjalan sebelum proses mati.

```typescript
process.on('SIGTERM', () => {
  console.log('SIGTERM signal received: closing HTTP server');
  server.close(() => {
    console.log('HTTP server closed, closing DB pool');
    db.pool.end();
  });
});
```

## Pemeriksaan pemahaman

1. Mengapa urutan pendaftaran middleware di Express sangat krusial?
2. Bagaimana Express membedakan antara middleware biasa dan error-handling middleware?
3. Apa perbedaan status code 401 dan 403 dalam konteks otentikasi REST API?
4. Mengapa penting membatasi ukuran body parser (`express.json({ limit: '10kb' })`)?

**Kunci ringkas:** Express mengeksekusi middleware secara sekuensial dari atas ke bawah; error-handling middleware secara spesifik memiliki arity 4 parameter `(err, req, res, next)`; 401 berarti belum terautentikasi (belum login/token salah), sedangkan 403 berarti sudah terautentikasi namun tidak memiliki wewenang role; pembatasan body size mencegah serangan Denial-of-Service akibat pengiriman payload berukuran raksasa yang menghabiskan RAM server.

## Referensi

- [Express.js Official Documentation](https://expressjs.com/)
- [Node.js Best Practices Repository](https://github.com/goldbergyoni/nodebestpractices)
- [OWASP REST Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html)

# Materi: AdonisJS untuk TypeScript-First Full-Stack MVC

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun aplikasi backend modern dengan AdonisJS (v6): memahami arsitektur berbasis TypeScript murni, menguasai routing dan HTTP controllers, menggunakan Lucid ORM untuk relasi dan migrasi database, menerapkan validasi menggunakan VineJS, serta mengamankan API dengan sistem autentikasi dan Bouncer authorization.

## Prasyarat

- Memahami TypeScript tingkat menengah (tipikal objek, interface, async/await).
- Mengerti konsep arsitektur MVC (*Model-View-Controller*) dan basis data SQL.
- Node.js versi 20+ terpasang di sistem.

## 1. Dasar: filosofi AdonisJS v6

AdonisJS adalah framework Node.js berfitur lengkap (*batteries-included*) yang dirancang khusus untuk TypeScript. Berbeda dengan framework minimalis, AdonisJS menyediakan standar terpadu untuk ORM, autentikasi, mailing, validasi, dan CLI (Ace), mirip dengan pengalaman Laravel atau Ruby on Rails di dunia JavaScript/TypeScript.

### Karakteristik Kunci:
- **TypeScript First:** Dukungan tipe statis tanpa konfigurasi rumit, termasuk type-safe routing dan query builder.
- **Lucid ORM:** Active Record ORM yang elegan dengan query builder bertenaga Knex.
- **VineJS:** Schema validator tercepat di ekosistem JavaScript, terintegrasi langsung dengan AdonisJS.
- **IoC Container:** Manajemen injeksi dependensi yang bersih dan modular.

## 2. Praktik dasar: Lucid Model, VineJS Validator, dan Controller

```typescript
// app/models/course.ts
import { DateTime } from 'luxon'
import { BaseModel, column, belongsTo } from '@adonisjs/lucid/orm'
import type { BelongsTo } from '@adonisjs/lucid/types/relations'
import User from '#models/user'

export default class Course extends BaseModel {
  @column({ isPrimary: true })
  declare id: number

  @column()
  declare title: string

  @column()
  declare description: string

  @column()
  declare price: number

  @column()
  declare instructorId: number

  @belongsTo(() => User, { foreignKey: 'instructorId' })
  declare instructor: BelongsTo<typeof User>

  @column.dateTime({ autoCreate: true })
  declare createdAt: DateTime

  @column.dateTime({ autoCreate: true, autoUpdate: true })
  declare updatedAt: DateTime
}
```

```typescript
// app/validators/course.ts
import vine from '@vinejs/vine'

export const createCourseValidator = vine.compile(
  vine.object({
    title: vine.string().trim().minLength(5).maxLength(255),
    description: vine.string().trim(),
    price: vine.number().min(0),
  })
)
```

```typescript
// app/controllers/courses_controller.ts
import type { HttpContext } from '@adonisjs/core/http'
import Course from '#models/course'
import { createCourseValidator } from '#validators/course'

export default class CoursesController {
  public async index({ response }: HttpContext) {
    const courses = await Course.query().preload('instructor').orderBy('created_at', 'desc')
    return response.ok(courses)
  }

  public async store({ request, response, auth }: HttpContext) {
    const user = auth.getUserOrFail()
    const payload = await request.validateUsing(createCourseValidator)

    const course = await Course.create({
      ...payload,
      instructorId: user.id,
    })

    return response.created(course)
  }

  public async show({ params, response }: HttpContext) {
    const course = await Course.findOrFail(params.id)
    await course.load('instructor')
    return response.ok(course)
  }
}
```

## 3. Fitur utama dan pola arsitektur

### Routing Type-Safe
AdonisJS v6 menggunakan sub-path imports `#models`, `#controllers`, dan `#validators` yang diatur via `package.json`:

```typescript
// start/routes.ts
import router from '@adonisjs/core/services/router'
import { middleware } from '#start/kernel'

const CoursesController = () => import('#controllers/courses_controller')

router
  .group(() => {
    router.get('/courses', [CoursesController, 'index'])
    router.get('/courses/:id', [CoursesController, 'show'])

    // Rute yang memerlukan otentikasi
    router
      .group(() => {
        router.post('/courses', [CoursesController, 'store'])
      })
      .use(middleware.auth())
  })
  .prefix('/api/v1')
```

### Otorisasi Akses dengan Bouncer
AdonisJS menyediakan Bouncer untuk mengontrol otorisasi hak akses (*policies & abilities*):

```typescript
// app/abilities/main.ts
import { Bouncer } from '@adonisjs/bouncer'
import User from '#models/user'
import Course from '#models/course'

export const editCourse = Bouncer.ability((user: User, course: Course) => {
  return user.id === course.instructorId || user.role === 'admin'
})
```

## 4. Praktik kerja: migrasi skema dan pagination pada dataset besar

### Skenario
Halaman eksplorasi kursus melambat ketika jumlah kursus mencapai ribuan, dan API mengembalikan seluruh data sekaligus tanpa pembatasan batas pagination.

### Tugas
Terapkan Database Migration untuk mengindeks kolom pencarian dan gunakan pagination bawaan Lucid ORM pada endpoint katalog.

```typescript
// database/migrations/1710000000000_create_courses_table.ts
import { BaseSchema } from '@adonisjs/lucid/schema'

export default class extends BaseSchema {
  protected tableName = 'courses'

  async up() {
    this.schema.createTable(this.tableName, (table) => {
      table.increments('id')
      table.string('title').notNullable().index()
      table.text('description').nullable()
      table.integer('price').unsigned().notNullable()
      table.integer('instructor_id').unsigned().references('users.id').onDelete('CASCADE')
      
      table.timestamp('created_at', { useTz: true }).notNullable()
      table.timestamp('updated_at', { useTz: true }).notNullable()
    })
  }

  async down() {
    this.schema.dropTable(this.tableName)
  }
}
```

Controller dengan pagination:
```typescript
public async index({ request, response }: HttpContext) {
  const page = request.input('page', 1)
  const limit = request.input('limit', 20)

  const courses = await Course.query()
    .preload('instructor')
    .paginate(page, limit)

  return response.ok(courses)
}
```

### Lima hint bertahap
1. Buat migrasi menggunakan perintah `node ace make:migration courses`.
2. Selalu tambahkan indeks (`.index()`) pada kolom yang sering digunakan untuk filter atau pencarian (`title`, `category_id`).
3. Gunakan method `.paginate(page, limit)` pada query Lucid untuk otomatis menghasilkan metadata pagination (total, perPage, currentPage).
4. Gunakan `auth.getUserOrFail()` untuk menjamin ketersediaan data pengguna yang telah terotentikasi.
5. Konfigurasikan guard auth di `config/auth.ts` untuk memilih antara session cookie atau access token.

## 5. Lanjutan: event emitter, queue, dan health check

- **Ace CLI:** Utilitas command-line tangguh untuk custom tasks, database seeding, dan repl debugging (`node ace repl`).
- **Health Checks:** Paket `@adonisjs/health-checks` bawaan untuk monitoring status koneksi database, disk space, dan memory.
- **Edge Template Engine:** Jika membangun full-stack SSR konvensional, AdonisJS menyertakan engine template Edge yang sangat cepat dan aman.

## Pemeriksaan pemahaman

1. Apa keunggulan AdonisJS dibanding merangkai Express.js dari nol untuk proyek besar?
2. Bagaimana cara VineJS meningkatkan keamanan dan kecepatan validasi data di AdonisJS?
3. Apa fungsi Bouncer dalam arsitektur aplikasi AdonisJS?
4. Mengapa AdonisJS v6 menggunakan standar *subpath imports* (`#models/*`, `#controllers/*`)?

**Kunci ringkas:** AdonisJS menyediakan struktur terpadu (*batteries-included*) dengan standar ORM, auth, dan validasi tanpa perlunya merakit pustaka pihak ketiga secara manual; VineJS mengompilasi skema menjadi kode JavaScript murni dengan optimasi JIT sehingga jauh lebih cepat; Bouncer mengisolasi aturan otorisasi pengguna secara eksplisit; *subpath imports* adalah standar Node.js modern yang mengeliminasi masalah *relative path hell* (`../../..`) dan fully type-safe.

## Referensi

- [AdonisJS Official Documentation](https://docs.adonisjs.com/)
- [Lucid ORM Documentation](https://lucid.adonisjs.com/)
- [VineJS Documentation](https://vinejs.dev/)

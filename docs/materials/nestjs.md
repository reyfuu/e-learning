# Materi: NestJS untuk Arsitektur Backend Skala Enterprise

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat membangun aplikasi backend Node.js yang modular dan terstruktur rapi menggunakan NestJS: memahami Dependency Injection (DI) dan arsitektur berbasis Modul, mengimplementasikan Controllers, Providers/Services, DTO Validation dengan Pipes, Authentication/Authorization Guards, serta Interceptors dan Exception Filters.

## Prasyarat

- Memahami TypeScript tingkat menengah (Decorators, Types, Interfaces, Generics).
- Memahami konsep dasar REST API dan arsitektur berorientasi objek (OOP).
- Node.js versi 18+ dan `@nestjs/cli` terpasang secara lokal.

## 1. Dasar: arsitektur modular NestJS

NestJS terinspirasi oleh Angular, memadukan paradigma Object-Oriented Programming (OOP), Functional Programming (FP), dan Functional Reactive Programming (FRP). NestJS menggunakan Express (atau Fastify) di bawah tenda, namun memberikan struktur arsitektur yang baku dan *opinionated*.

### Blok Pembangun Utama:
- **Modules (`@Module`):** Mengelompokkan komponen terkait (Controllers, Providers) dalam satu batas domain fungsional.
- **Controllers (`@Controller`):** Menangani incoming HTTP requests dan memetakan rute ke method handler.
- **Providers / Services (`@Injectable`):** Kelas yang berisi logika bisnis dan dapat diinjeksi melalui Dependency Injection (DI).
- **Pipes (`@Injectable` implements `PipeTransform`):** Mentransformasi atau memvalidasi payload request sebelum menyentuh handler.
- **Guards (`@Injectable` implements `CanActivate`):** Menentukan apakah suatu request diizinkan mengakses handler (Autentikasi & RBAC).

## 2. Praktik dasar: DTO, Controller, dan Service dengan Dependency Injection

```typescript
// courses/dto/create-course.dto.ts
import { IsString, IsNotEmpty, IsNumber, Min, IsEnum } from 'class-validator';

export enum CourseLevel {
  BEGINNER = 'beginner',
  INTERMEDIATE = 'intermediate',
  ADVANCED = 'advanced',
}

export class CreateCourseDto {
  @IsString()
  @IsNotEmpty({ message: 'Judul kursus tidak boleh kosong' })
  title: string;

  @IsNumber()
  @Min(0, { message: 'Harga minimal adalah 0' })
  price: number;

  @IsEnum(CourseLevel, { message: 'Level kursus tidak valid' })
  level: CourseLevel;
}
```

```typescript
// courses/courses.service.ts
import { Injectable, NotFoundException } from '@nestjs/common';
import { CreateCourseDto } from './dto/create-course.dto';

export interface Course {
  id: string;
  title: string;
  price: number;
  level: string;
  createdAt: Date;
}

@Injectable()
export class CoursesService {
  private courses: Course[] = [];

  create(dto: CreateCourseDto): Course {
    const newCourse: Course = {
      id: `crs_${Date.now()}`,
      ...dto,
      createdAt: new Date(),
    };
    this.courses.push(newCourse);
    return newCourse;
  }

  findAll(): Course[] {
    return this.courses;
  }

  findById(id: string): Course {
    const course = this.courses.find((c) => c.id === id);
    if (!course) {
      throw new NotFoundException(`Kursus dengan ID ${id} tidak ditemukan`);
    }
    return course;
  }
}
```

```typescript
// courses/courses.controller.ts
import { Controller, Get, Post, Body, Param, UsePipes, ValidationPipe } from '@nestjs/common';
import { CoursesService } from './courses.service';
import { CreateCourseDto } from './dto/create-course.dto';

@Controller('courses')
export class CoursesController {
  constructor(private readonly coursesService: CoursesService) {}

  @Post()
  @UsePipes(new ValidationPipe({ whitelist: true, forbidNonWhitelisted: true }))
  create(@Body() createCourseDto: CreateCourseDto) {
    return this.coursesService.create(createCourseDto);
  }

  @Get()
  findAll() {
    return this.coursesService.findAll();
  }

  @Get(':id')
  findById(@Param('id') id: string) {
    return this.coursesService.findById(id);
  }
}
```

## 3. Fitur utama dan pola arsitektur

### Alur Eksekusi Request Lifecycle di NestJS
Setiap request HTTP melalui urutan pipeline terstandarisasi:
```text
Request → Middleware → Guards → Interceptors (Pre) → Pipes → Controller Handler → Interceptors (Post) → Exception Filters → Response
```

### Dependency Injection Hierarkis
NestJS mengelola siklus hidup (*scope*) objek secara otomatis:
- **DEFAULT (Singleton):** Instance tunggal dibagikan ke seluruh aplikasi (sangat efisien memori).
- **REQUEST:** Instance baru dibuat untuk setiap request masuk (berguna untuk multi-tenancy).
- **TRANSIENT:** Instance baru dibuat setiap kali diinjeksi ke provider lain.

## 4. Praktik kerja: implementasi Custom Auth Guard dan Roles Decorator

### Skenario
Endpoint untuk mempublikasikan kurikulum harus dilindungi. Hanya pengguna dengan role `instructor` atau `admin` yang diizinkan memanggil endpoint tersebut, dan payload JWT harus diekstrak ke dalam parameter handler.

### Tugas
Buat decorator kustom `@Roles()` dan Guard `RolesGuard` yang membaca metadata refleksi menggunakan `Reflector`.

```typescript
// auth/decorators/roles.decorator.ts
import { SetMetadata } from '@nestjs/common';

export const ROLES_KEY = 'roles';
export const Roles = (...roles: string[]) => SetMetadata(ROLES_KEY, roles);
```

```typescript
// auth/guards/roles.guard.ts
import { Injectable, CanActivate, ExecutionContext, ForbiddenException } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { ROLES_KEY } from '../decorators/roles.decorator';

@Injectable()
export class RolesGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  canActivate(context: ExecutionContext): boolean {
    const requiredRoles = this.reflector.getAllAndOverride<string[]>(ROLES_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);

    if (!requiredRoles) {
      return true; // Endpoint publik/tidak ada batasan role
    }

    const { user } = context.switchToHttp().getRequest();
    if (!user || !requiredRoles.includes(user.role)) {
      throw new ForbiddenException('Akses ditolak: role Anda tidak mencukupi');
    }

    return true;
  }
}
```

### Lima hint bertahap
1. Buat decorator khusus untuk menempelkan array role yang diizinkan ke handler method.
2. Manfaatkan `Reflector` dari `@nestjs/core` untuk membaca metadata dekorator di dalam Guard.
3. Gunakan `context.switchToHttp().getRequest()` untuk mengakses objek request standar Express.
4. Pastikan `AuthGuard` (yang menempelkan `req.user`) dieksekusi sebelum `RolesGuard`.
5. Daftarkan `ValidationPipe` secara global di `main.ts` menggunakan `app.useGlobalPipes()`.

## 5. Lanjutan: Microservices, TypeORM/Prisma, dan Swagger OpenAPI

- **OpenAPI / Swagger:** Integrasikan `@nestjs/swagger` untuk menghasilkan dokumentasi API interaktif secara otomatis dari kelas DTO.
- **Microservices:** NestJS mendukung transport layer modular seperti Redis pub/sub, RabbitMQ, Kafka, dan gRPC tanpa mengubah kode controller dasar.
- **Exception Filters:** Tangani format error respons agar seragam di seluruh aplikasi dengan `@Catch()`.

## Pemeriksaan pemahaman

1. Apa perbedaan tanggung jawab antara Pipes dan Guards di NestJS?
2. Mengapa Dependency Injection di NestJS default-nya bertipe Singleton?
3. Apa fungsi properti `whitelist: true` pada `ValidationPipe`?
4. Bagaimana urutan eksekusi antara Middleware, Guard, dan Interceptor?

**Kunci ringkas:** Pipes bertugas memvalidasi dan mentransformasi data input, sedangkan Guards bertugas menentukan izin akses (*canActivate*); Singleton menghemat alokasi memori dan waktu komputasi karena instance hanya diinisialisasi sekali saat bootstrap; `whitelist: true` secara otomatis menghapus properti asing yang tidak terdaftar di kelas DTO; urutan eksekusinya adalah: Middleware → Guard → Interceptor (pre) → Pipe → Handler → Interceptor (post).

## Referensi

- [NestJS Official Documentation](https://docs.nestjs.com/)
- [NestJS Best Practices and Patterns](https://github.com/jmcdo29/testing-nestjs)
- [Class Validator Documentation](https://github.com/typestack/class-validator)

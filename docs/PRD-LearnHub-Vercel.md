# PRD — LearnHub: Katalog Silabus Teknologi

**Versi:** 1.3 · **Tanggal:** 15 Oktober 2026 · **Status:** Draft kebutuhan produk target

## 1. Visi dan ringkasan produk

LearnHub membantu pelajar membangun kompetensi teknologi melalui jalur belajar yang berurutan, mulai dari Basic hingga Advanced. Instruktur menyajikan materi sebagai kursus dengan lesson, demonstrasi, latihan, proyek dan asesmen. Platform mendukung Indonesia/Inggris dan bisnis kursus gratis maupun berbayar.

Dokumen ini adalah PRD turunan untuk inisiatif katalog silabus teknologi, bukan pengganti [PRD produk induk](../PRD_eLearning.md). Persyaratan di sini menjelaskan keadaan yang ingin dicapai; implementasi saat ini lebih terbatas dan dirangkum di [baseline TRD](TRD-LearnHub-Vercel.md#baseline-implementasi-repositori). Stack target pada PRD induk dan DESIGN juga belum sepenuhnya cocok dengan manifest serta kode saat ini; gap tersebut dicatat di TRD untuk keputusan eksplisit.

Frontend Next.js di-deploy pada Vercel. Backend Go yang telah ada tetap berjalan sebagai API/container terpisah. Rincian arsitektur ada di [TRD](TRD-LearnHub-Vercel.md); kebutuhan bisnis dan fungsional ada di [BRD](BRD-LearnHub-Vercel.md) dan [FRD](FRD-LearnHub-Vercel.md).

## 2. Sasaran dan bukan sasaran

### Sasaran

- Temukan jalur teknologi menurut kategori, prasyarat, tingkat dan bahasa.
- Setiap teknologi dalam katalog memiliki capaian Basic dan Advanced serta praktik yang dapat dinilai.
- Learner dapat enrollment, belajar, melanjutkan progres, mengerjakan asesmen dan menerima sertifikat.
- Instruktur dapat menyusun dan mengajukan materi untuk moderasi.
- Admin dapat mengelola versi silabus dan katalog.

### Bukan sasaran MVP

- Menyediakan seluruh 40+ jalur sebagai konten lesson lengkap pada hari pertama.
- Menjalankan semua kode learner secara bebas di server LearnHub.
- Menggantikan IDE lokal/cloud lab atau sertifikasi resmi vendor.
- Aplikasi native, live class, dan microservices.

## 3. Persona dan kebutuhan

- **Learner pemula/karier:** ingin urutan jelas, prasyarat, contoh praktik, feedback, dan bukti kemajuan.
- **Learner berpengalaman:** ingin materi Advanced seperti arsitektur, keamanan, performa, observability dan integrasi.
- **Instructor/subject matter expert:** ingin authoring yang mudah, versi materi, kuis/proyek dan analitik.
- **Admin akademik:** ingin mengendalikan kualitas, konsistensi level dan publikasi.

## 4. Taksonomi dan pola silabus

Setiap jalur teknologi memiliki: tujuan, target audiens, prasyarat, estimasi jam yang ditetapkan setelah produksi, versi teknologi/tool, bahasa materi, capaian pembelajaran, modul Basic, modul Advanced, proyek portofolio, rubrik asesmen, dan referensi resmi. Basic dan Advanced adalah tingkat capaian; modul intermediate dapat ditambahkan sebagai jembatan bila diperlukan.

**Basic** mencakup instalasi/setup, konsep inti, sintaks/pola dasar, debugging, praktik kecil, dan dasar keamanan/kualitas.

**Advanced** mencakup desain sistem/pola lanjut, integrasi, pengujian, keamanan, performa, observability/operasi, trade-off, serta proyek produksi yang sesuai bidang. Materi Advanced tidak mengasumsikan topik spesialis sebagai prasyarat tanpa penjelasan.

## 5. Katalog silabus

Topik berikut adalah cetak biru kurikulum. Setiap butir menjadi rangkaian lesson, demo, latihan, dan asesmen; bukan klaim bahwa semua materi sudah diproduksi.

### A. Bahasa dan framework frontend

| Teknologi | Basic | Advanced |
|---|---|---|
| JavaScript | runtime, variabel/tipe, operator, kontrol alur, fungsi, array/object, DOM, module, async dasar | event loop, promise/concurrency, closure, prototype, memory/performance, testing, security browser, desain package |
| TypeScript | anotasi, inference, interface/type, union, function, generics dasar, konfigurasi | conditional/mapped types, narrowing lanjutan, type-safe API, declaration/module, strict migration, type-level design dan build performance |
| React | JSX, komponen, props/state, event, list/form, hooks dasar, routing dan fetch | state architecture, custom hooks, concurrent rendering concepts, server/client boundaries, performance, testing, accessibility dan error handling |
| Vue | template, reactivity, component, props/events, directives, Composition API dasar, router | composables, state architecture, rendering/performance, SSR, testing, TypeScript, plugin/design system dan security |
| Inertia.js | konsep server-driven SPA, adapter, routes, page props, forms dan validation | shared data, partial reload, authorization, asset/versioning, SSR, pagination, testing dan integrasi framework |
| Next.js | App Router, routes/layouts, Server/Client Components, data fetching, forms, metadata, deploy | caching/revalidation, rendering strategy, streaming, server actions/API boundaries, auth, security, performance, observability dan Vercel deployment |
| Nuxt.js | Vue/Nuxt setup, pages/layouts, components, composables, data fetch, runtime config | SSR/SSG/ISR, Nitro server routes, caching, modules, auth/security, performance, testing dan deployment targets |

### B. Bahasa dan framework backend

| Teknologi | Basic | Advanced |
|---|---|---|
| PHP | syntax, types, functions, arrays, OOP, Composer, HTTP dan exceptions | modern PHP typing, attributes, performance, security, async/runtime options, package design dan testing |
| Laravel | routing, controllers, Blade/API, Eloquent dasar, migration, validation, auth | service/container, queues/events, policies, API resources, query/performance, testing, security dan deployment |
| Express.js | Node runtime, routes, middleware, request/response, validation, error handling | layering, async errors, auth, security headers, rate limits, observability, performance, testing dan graceful shutdown |
| NestJS | modules, controllers, providers, DI, DTO validation, REST, guards | interceptors/pipes, CQRS patterns, queues, microservice transports, OpenAPI, testing, security dan observability |
| AdonisJS | project structure, routing, controllers, Lucid ORM, validation, auth dasar | authorization, queues/events, API architecture, database performance, testing, security, deployment dan observability |
| Java | syntax, types, OOP, collections, exceptions, generics, I/O dan build tools | concurrency, JVM/memory, streams, design patterns, profiling, security, modular design dan testing |
| Spring Boot | project setup, REST, DI, configuration, persistence, validation, basic security | Spring Security, transactions, testing, messaging, observability, resilience, performance, container deployment dan API versioning |
| Go | syntax, types, functions, structs/interfaces, errors, modules, HTTP dan context | concurrency, cancellation, profiling, memory, API design, security, database patterns, observability dan graceful shutdown |
| Echo | routing, handler/context, binding dan validation, middleware, REST API, error handling | middleware/auth, request lifecycle, testing, rate limits, observability, performance dan deployment aman |
| Gin | routing, handlers, binding/validation, middleware, JSON API, errors | auth/RBAC, context, testing, rate limiting, versioning, performance, observability dan secure deployment |
| Fiber | routing, handlers, middleware, validation, REST, configuration | context/lifecycle pitfalls, auth, testing, profiling, security, concurrency/IO tradeoffs, observability dan deployment |
| Python | syntax, data structures, functions, modules, exceptions, venv, typing dasar | advanced typing, iterators/generators, concurrency, packaging, profiling, security, testing dan architecture |
| Django | models, ORM, admin, views, templates/API, forms, auth dasar | query optimization, permissions, async tasks, caching, security, testing, API design, scaling dan deployment |
| FastAPI | path operations, Pydantic, dependency injection, OpenAPI, async basics | auth/scopes, async database, background jobs, testing, performance, security, observability dan deployment |
| Flask | routes, templates/API, request lifecycle, config, blueprints dan extensions | app factory, auth, ORM patterns, async/task queues, testing, security, performance dan deployment |
| Ruby | syntax, collections, blocks, methods, modules, classes, exceptions, Bundler | metaprogramming, concurrency, memory/performance, testing, packaging, security dan maintainable design |
| Ruby on Rails | MVC, routes, Active Record, migrations, views/API, validations, scaffolding | callbacks/service boundaries, authorization, background jobs, query/performance, security, testing, caching dan deployment |
| C# | syntax, types, OOP, collections, LINQ, exceptions, async/await dasar | generics, concurrency, memory/GC, performance, dependency injection, security, testing dan architecture |
| .NET | solution/project, ASP.NET Core API, DI, configuration, EF Core dasar, logging | middleware, auth, resilient services, performance, observability, testing, containerization, security dan deployment |
| C++ | compilation, types, control flow, functions, classes, STL, pointers/references | RAII, templates, move semantics, concurrency, memory safety, profiling, build systems, testing dan modern C++ architecture |

### C. Mobile

| Teknologi | Basic | Advanced |
|---|---|---|
| Flutter | Dart essentials, widgets, layout, navigation, forms, state dasar, HTTP | state architecture, platform channels, performance, accessibility, testing, offline/sync, security dan release pipeline |

### D. Database dan backend platform

| Teknologi | Basic | Advanced |
|---|---|---|
| PostgreSQL | relational model, SQL CRUD, joins, constraints, indexes dasar, transactions | query plans, indexing, isolation/locking, partitioning, backup/restore, security, replication dan tuning |
| MySQL | schema, SQL CRUD, joins, constraints, indexes, transactions | optimizer/EXPLAIN, InnoDB, locking, replication, backup/recovery, security dan performance tuning |
| MariaDB | SQL/schema, CRUD, joins, constraints, indexing dasar, transactions | optimizer, engines, replication/clustering concepts, backup/recovery, security, performance dan migration compatibility |
| MongoDB | documents/collections, CRUD, schema design, indexes, aggregation dasar | transactions, aggregation optimization, sharding/replica sets, schema evolution, security, backup dan performance |
| Prisma ORM | model relational data dengan schema/contract, koneksi PostgreSQL, typed queries, CRUD, relasi dan migration dasar | introspection, transactions, migration review dan expand/contract, indeks/query performance, akses data multi-tenant dan upgrade antarversi |
| Drizzle ORM | schema TypeScript dan type inference, koneksi PostgreSQL, CRUD, relations API dan migration dengan Drizzle Kit | SQL-like query builder, joins, transactions, migration review, indeks/query performance, debugging SQL dan akses data multi-tenant |
| Supabase | project, Postgres basics, Auth, table/API, Storage, client SDK | Row Level Security, policies, Edge Functions, realtime, migrations, performance, secrets dan production security |
| Oracle Database | relational concepts, SQL, tables, joins, constraints, sequences, transactions | PL/SQL, execution plans, indexing, roles, backup/recovery, RAC/Data Guard concepts dan performance |
| Firebase | project setup, Authentication, Firestore/Realtime Database, Storage, SDK | Security Rules, Cloud Functions, indexing, offline sync, cost/performance, App Check, emulator testing dan migration |

### E. Cloud dan platform engineering

| Teknologi | Basic | Advanced |
|---|---|---|
| AWS | account/IAM, regions, VPC concepts, S3, compute, managed DB, CloudWatch | least privilege, multi-account, autoscaling, networking, resilience, cost controls, IaC, security dan disaster recovery |
| Google Cloud | projects/IAM, regions, Cloud Storage, compute, managed DB, logging | VPC, service accounts, autoscaling, GKE/serverless, resilience, cost, security, observability dan recovery |
| Azure | subscriptions/resource groups, Entra/IAM, storage, compute, managed DB, Monitor | networking, identity, scaling, containers, resilience, cost governance, security, IaC dan disaster recovery |
| Docker | images, containers, Dockerfile, volumes, ports, Compose dan registries | multi-stage/hardened builds, networking, secrets, SBOM/scanning, optimization, orchestration readiness dan CI |
| Kubernetes | pods, deployments, services, config/secrets, probes, namespaces | RBAC, ingress/network policy, autoscaling, storage, rollout, observability, security, operators dan recovery |
| Terraform | providers, resources, variables, outputs, state, modules dasar | remote state/locking, modules, imports/refactors, workspaces, policy, secrets, drift, testing dan multi-environment design |

### F. CI/CD, QA, dan kolaborasi

| Teknologi | Basic | Advanced |
|---|---|---|
| Jenkins | jobs, pipelines, agents, credentials, artifacts, webhook | shared libraries, scalable agents, pipeline security, approvals, caching, observability, backup dan plugin governance |
| GitHub Actions | workflows, triggers, jobs, steps, secrets, artifacts, environments | reusable workflows, OIDC, least privilege, matrix/caching, release promotion, supply-chain security dan self-hosted runners |
| Selenium | WebDriver, locators, waits, browser actions, assertions, page objects | grid/parallel execution, stable test design, cross-browser, CI, accessibility/performance integration dan debugging flakiness |
| Appium | driver setup, capabilities, locators, gestures, waits, mobile test lifecycle | device farms, parallelism, hybrid apps, test architecture, CI, flakiness diagnostics dan secure test data |
| TestRail | projects/suites/cases, runs, results, plans, traceability | custom fields, API integration, automation mapping, reporting, governance, permissions dan release analytics |
| Jira | projects, issues, workflows, boards, filters, permissions | workflow design, automation, JQL, dashboards, integrations, governance, release planning dan reporting |
| Postman | requests, collections, environments, variables, auth, tests dasar | scripting, Newman/CI, contract checks, mocks, monitors, secret handling, reusable test architecture dan team governance |

### G. AI/ML dan agent engineering

| Teknologi | Basic | Advanced |
|---|---|---|
| PyTorch | tensors, autograd, datasets/dataloaders, model, loss, optimizer | custom training loops, distributed training, mixed precision, profiling, reproducibility, deployment dan model security |
| LangChain | model/chat interfaces, prompts, chains/runnables, retrievers, tools dasar | retrieval design, structured output, evaluation, tracing, latency/cost, guardrails, deployment dan provider portability |
| LangGraph | nodes, state, edges, conditional flow, persistence dasar | cycles, checkpointing, human-in-loop, retries, subgraphs, streaming, state/version design, evaluation dan production observability |
| CrewAI | agents, roles, tasks, tools, crews, process dasar | delegation, memory/context, tool security, deterministic workflows, evaluation, tracing, cost control dan production operations |
| QLoRA | transformer/LLM overview, quantization concepts, dataset preparation, LoRA adapters, supervised fine-tuning | memory/compute tradeoffs, data quality, hyperparameter tuning, evaluation, merging/deployment, safety, reproducibility dan licensing |
| MCP | client/server concepts, tools, resources, prompts, transport, local setup | authorization, trust boundaries, validation, secure tool execution, transport/deployment, versioning, observability dan conformance testing |
| AI Agents | model/tool loop, instructions, tools, state, basic retrieval, human approval | planning limits, permissions, memory, evals, failure recovery, observability, prompt injection defense, cost/latency dan safe deployment |

## 6. Persyaratan produk dan UX

- Halaman publik: landing, kategori, technology track, detail kursus, instructor, pencarian dan SEO metadata.
- Learner: dashboard, enrollment, player, outline, progress, notes/bookmarks (opsional), quiz/project, review dan certificate.
- Learner: editor kode interaktif di browser dengan starter code, run/submit, feedback dari test, dan tepat lima hint bertingkat pada latihan pilot.
- Contoh skenario kerja: memperbaiki validasi dan pagination endpoint Go Echo; menangani loading/error/empty state pada integrasi API TypeScript; menulis test regresi untuk bug yang diberikan. Semua prompt memakai data fiktif dan tidak menyalin kode atau soal rahasia perusahaan.
- Instructor: dashboard, course editor, lesson ordering, asset manager, preview, submission dan analytics.
- Admin: user/course moderation, category/track/version manager, transaction overview dan audit.
- Kursus menampilkan level, prasyarat, capaian, bahasa, estimasi waktu, versi teknologi, kurikulum dan metode asesmen.
- UI responsif, aksesibel, ID/EN, state loading/empty/error yang jelas.

## 7. Persyaratan nonfungsional

- Keamanan: HTTPS, RBAC, validasi server, proteksi CSRF/XSS/SQLi, rate limiting, audit, secret management dan signed media URLs.
- Eksekusi kode: source tidak dijalankan di proses API atau frontend; runner terisolasi, sekali pakai, tanpa network/secrets/akses host, dengan limit CPU, memori, waktu, output dan concurrency.
- Performa: ukur Core Web Vitals dan API p95; sasaran awal mengikuti PRD induk p95 API ≤200 ms untuk endpoint umum, dikecualikan pekerjaan asinkron.
- Ketersediaan: sasaran 99.5% ditinjau berdasarkan SLA provider dan praktik operasional.
- Skalabilitas: CDN, pagination, connection pooling, background jobs dan pengujian beban sebelum klaim concurrent users.
- Aksesibilitas: WCAG 2.1 AA sebagai sasaran verifikasi.
- Lokalitas: ID/EN untuk UI dan kemampuan kursus menandai bahasa kontennya.

## 8. Prioritas rilis

### MVP

- Auth, role, catalog/track, course authoring, moderation, free enrollment, lesson text/video, progress, search dasar dan dashboard.
- Pilot konten untuk rumpun web (JavaScript, TypeScript, React, Next.js), backend Go dengan Echo, SQL (PostgreSQL), ORM TypeScript (Prisma dan Drizzle), serta GitHub Actions/Docker. Latihan coding interaktif dimulai pada Go dan JavaScript/TypeScript; ekspansi katalog mengikuti review ahli.
- Vercel frontend dan Go API container, storage/CDN, observability minimum.

### Rilis berikutnya

- Kuis/proyek, sertifikat, reviews, notifikasi, pembayaran, promo dan analytics instruktur.
- Ekspansi seluruh katalog sesuai matriks silabus.
- AI/agent lessons dengan evaluasi dan kebijakan keamanan khusus.

### Kemudian

- Lab interaktif terisolasi, rekomendasi personal, jalur berbasis kompetensi, mobile app dan live learning bila tervalidasi.

## 9. Metrik produk

- Funnel: visit → signup → enrollment → lesson start.
- Aktivasi: learner menyelesaikan lesson pertama dalam periode yang ditentukan.
- Engagement: weekly active learners, lesson completion, quiz/project attempts.
- Outcome: course/track completion, assessment pass rate, learner satisfaction.
- Supply: kursus submitted/approved, waktu moderasi, aktivitas instruktur.
- Bisnis: conversion, gross revenue, refund, payout dan biaya per learner.
- Reliability: web vitals, API latency/error, uptime, queue delay dan media playback startup.

## 10. Kriteria penerimaan rilis awal

- Learner dapat menemukan track, membaca Basic/Advanced, membuat akun, enrollment, menuntaskan lesson dan memulihkan progress.
- Instructor dapat membuat course dengan seksi/lesson dan mengirimkannya untuk moderasi.
- Admin dapat approve/reject dan semua perubahan penting masuk audit log.
- Web berjalan di Vercel dan mengakses Go API melalui HTTPS dengan environment yang dipisahkan.
- Learner dapat menjalankan minimal satu latihan coding Go dan satu JavaScript/TypeScript di browser; runner menguji solusi secara terisolasi dan menampilkan hasil tanpa membocorkan hidden tests.
- Setiap latihan pilot memiliki skenario kerja realistis, rubrik evaluasi, dan tepat lima hint progresif yang dapat dibuka satu per satu.
- Kursus privat dan data learner tidak dapat diakses tanpa entitlement/otorisasi.
- Silabus katalog mencakup seluruh teknologi yang diminta di tabel bagian 5 dan memiliki Basic serta Advanced.

## 11. Ketergantungan dan keputusan yang masih perlu ditetapkan

- Provider hosting Go, database, queue, storage/CDN, email, payment, domain dan analytics.
- Komisi, gateway, refund, payout, format sertifikat, kebijakan konten dan target SLA.
- Ahli reviewer, versi/toolchain yang didukung, durasi kursus, lisensi materi, dan rubrik asesmen.
- Strategi isolasi bila kelak menyediakan code execution/lab.

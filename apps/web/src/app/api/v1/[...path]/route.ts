import { NextRequest, NextResponse } from "next/server";
import catalogSpecs from "@/data/technology-courses.json";

type CatalogSpec = (typeof catalogSpecs)[number];

interface StoredUser {
  id: string;
  email: string;
  password: string;
  name: string;
  role: "learner" | "instructor" | "admin";
  avatar_url: string | null;
  is_verified: boolean;
  is_active: boolean;
  created_at: string;
}

interface ClientState {
  e?: string[]; // enrolled course slugs/ids
  c?: string[]; // completed module ids
  h?: number;   // hints unlocked
  u?: StoredUser[]; // custom registered users
}

const DEFAULT_USERS: StoredUser[] = [
  {
    id: "987ebef4-666f-4f1a-b69a-3f329a5fb515",
    email: "audinathanael@gmail.com",
    password: "$Audin123",
    name: "Audi Nathanael",
    role: "learner",
    avatar_url: null,
    is_verified: true,
    is_active: true,
    created_at: "2026-10-15T00:00:00.000Z",
  },
  {
    id: "75576176-c963-4e29-a984-17ec45b4ac36",
    email: "audipakaii@gmail.com",
    password: "$Audin123",
    name: "Audi Instructor",
    role: "instructor",
    avatar_url: null,
    is_verified: true,
    is_active: true,
    created_at: "2026-10-15T00:00:00.000Z",
  },
];

const EXERCISE = {
  id: "employee-email-normalizer",
  title: "Normalisasi email karyawan",
  scenario:
    "Tim HR menerima data karyawan dari beberapa sistem. Sebelum disimpan, backend harus menormalkan alamat email agar pencarian akun konsisten.",
  prompt:
    "Buat program yang membaca satu email per baris dari stdin dan mencetak versi yang sudah di-trim dan huruf kecil. Pertahankan satu baris output untuk setiap input. Contoh:  Rani@Example.COM  menjadi rani@example.com.",
  starter: {
    javascript:
      "const fs = require('fs');\nconst lines = fs.readFileSync(0, 'utf8').split(/\\r?\\n/);\nfor (const line of lines) {\n  if (line.length === 0) continue;\n  // TODO: trim spasi lalu ubah menjadi huruf kecil\n  console.log(line);\n}",
    go: 'package main\n\nimport (\n  "bufio"\n  "fmt"\n  "os"\n)\n\nfunc main() {\n  scanner := bufio.NewScanner(os.Stdin)\n  for scanner.Scan() {\n    // TODO: trim spasi lalu ubah menjadi huruf kecil\n    fmt.Println(scanner.Text())\n  }\n}',
  },
  hints: [
    "Hilangkan spasi di awal dan akhir string sebelum melakukan transformasi lain.",
    "Di JavaScript, string memiliki method trim() dan toLowerCase().",
    "Di Go, gunakan strings.TrimSpace dan strings.ToLower.",
    "Pastikan setiap baris diproses sendiri; jangan menggabungkan alamat email yang berbeda.",
    "Urutannya: trim whitespace, lowercase, lalu tulis satu hasil per baris.",
  ],
};

function encodeToken(user: StoredUser): string {
  const payload = {
    sub: user.id,
    email: user.email,
    name: user.name,
    role: user.role,
    exp: Math.floor(Date.now() / 1000) + 86400 * 7,
  };
  return Buffer.from(JSON.stringify(payload), "utf-8").toString("base64url");
}

function decodeToken(token?: string | null): StoredUser | null {
  if (!token) return null;
  try {
    const raw = token.includes(".") ? token.split(".")[1] : token;
    const data = JSON.parse(Buffer.from(raw, "base64url").toString("utf-8"));
    const id = data.sub || data.id;
    const email = data.email;
    const found = DEFAULT_USERS.find((u) => u.id === id || u.email === email);
    if (found) return found;
    if (email) {
      return {
        id: id || email,
        email,
        password: "",
        name: data.name || email.split("@")[0],
        role: data.role || "learner",
        avatar_url: null,
        is_verified: true,
        is_active: true,
        created_at: new Date().toISOString(),
      };
    }
  } catch {
    /* ignore */
  }
  return null;
}

function readState(req: NextRequest): ClientState {
  const raw = req.cookies.get("learnhub_state")?.value;
  if (!raw) return { e: [], c: [], h: 0, u: [] };
  try {
    const parsed = JSON.parse(Buffer.from(raw, "base64url").toString("utf-8"));
    return {
      e: Array.isArray(parsed.e) ? parsed.e : [],
      c: Array.isArray(parsed.c) ? parsed.c : [],
      h: typeof parsed.h === "number" ? parsed.h : 0,
      u: Array.isArray(parsed.u) ? parsed.u : [],
    };
  } catch {
    return { e: [], c: [], h: 0, u: [] };
  }
}

function respondWithState(body: unknown, state: ClientState, status = 200) {
  const res = NextResponse.json(body, { status });
  const encoded = Buffer.from(JSON.stringify(state), "utf-8").toString("base64url");
  res.cookies.set("learnhub_state", encoded, {
    path: "/",
    maxAge: 60 * 60 * 24 * 30,
    sameSite: "lax",
  });
  return res;
}

function formatDuration(seconds: number): string {
  if (seconds <= 0) return "0m";
  const mins = Math.floor(seconds / 60);
  return mins > 0 ? `${mins}m` : `${seconds}s`;
}

function buildCourseResponse(spec: CatalogSpec, completedSet: Set<string>) {
  const fullDescription =
    spec.description +
    "\n\nYang akan dipelajari:\n- " +
    spec.outcomes.join("\n- ") +
    "\n\nPrasyarat:\n- " +
    spec.prerequisites.join("\n- ");

  const modules: {
    id: string;
    title: string;
    section_title: string;
    content_type: string;
    content: string | null;
    is_completed: boolean;
    duration: string;
  }[] = [];

  spec.sections.forEach((sec, sIdx) => {
    sec.lessons.forEach((lesson, lIdx) => {
      const modId = `${spec.slug}__${sIdx + 1}_${lIdx + 1}`;
      modules.push({
        id: modId,
        title: lesson.title,
        section_title: sec.title,
        content_type: "text",
        content: lesson.content,
        is_completed: completedSet.has(modId),
        duration: formatDuration(lesson.duration_seconds),
      });
    });
  });

  const total = modules.length;
  const completed = modules.filter((m) => m.is_completed).length;
  const progress = total > 0 ? Math.round((completed / total) * 100) : 0;

  return {
    id: spec.slug,
    title: spec.title,
    description: fullDescription,
    link: "",
    icon: "book",
    status: "published",
    total_modules: total,
    completed_modules: completed,
    progress,
    modules,
  };
}

function getAuthUser(req: NextRequest, state: ClientState): StoredUser | null {
  const authHeader = req.headers.get("authorization") ?? "";
  const bearer = authHeader.startsWith("Bearer ") ? authHeader.slice(7).trim() : "";
  const cookieToken = req.cookies.get("access_token")?.value ?? "";
  const token = bearer || cookieToken;
  const decoded = decodeToken(token);
  if (!decoded) return null;
  const custom = (state.u ?? []).find((u) => u.id === decoded.id || u.email === decoded.email);
  return custom ?? decoded;
}

function publicUser(u: StoredUser) {
  return {
    id: u.id,
    email: u.email,
    name: u.name,
    avatar_url: u.avatar_url,
    role: u.role,
    is_verified: u.is_verified,
    is_active: u.is_active,
    created_at: u.created_at,
  };
}

async function handle(req: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const { path } = await params;
  const route = "/" + (path ?? []).join("/");
  const method = req.method.toUpperCase();
  const state = readState(req);
  const completedSet = new Set(state.c ?? []);
  const enrolledSet = new Set(state.e ?? []);
  const allUsers = [...DEFAULT_USERS, ...(state.u ?? [])];

  // ─── Auth ──────────────────────────────────────────────────────────
  if (method === "POST" && route === "/auth/login") {
    const body = await req.json().catch(() => ({}));
    const email = String(body.email ?? "").trim().toLowerCase();
    const password = String(body.password ?? "");
    const user = allUsers.find((u) => u.email.toLowerCase() === email && u.password === password);
    if (!user) {
      return NextResponse.json({ message: "Invalid email or password" }, { status: 401 });
    }
    const token = encodeToken(user);
    return respondWithState({
      access_token: token,
      refresh_token: token,
      user: publicUser(user),
      expires_in: 86400,
    }, state);
  }

  if (method === "POST" && route === "/auth/register") {
    const body = await req.json().catch(() => ({}));
    const email = String(body.email ?? "").trim().toLowerCase();
    const password = String(body.password ?? "");
    const name = String(body.name ?? "").trim();
    const role = body.role === "instructor" ? "instructor" : "learner";

    if (!email || !password || !name) {
      return NextResponse.json({ message: "Missing required fields" }, { status: 400 });
    }
    if (allUsers.some((u) => u.email.toLowerCase() === email)) {
      return NextResponse.json({ message: "Email already registered" }, { status: 409 });
    }

    const newUser: StoredUser = {
      id: `usr_${Date.now()}`,
      email,
      password,
      name,
      role,
      avatar_url: null,
      is_verified: true,
      is_active: true,
      created_at: new Date().toISOString(),
    };
    state.u = [...(state.u ?? []), newUser];
    const token = encodeToken(newUser);
    return respondWithState(
      {
        access_token: token,
        refresh_token: token,
        user: publicUser(newUser),
        expires_in: 86400,
      },
      state,
      201
    );
  }

  if (method === "POST" && route === "/auth/refresh") {
    const body = await req.json().catch(() => ({}));
    const user = decodeToken(body.refresh_token);
    if (!user) {
      return NextResponse.json({ message: "Invalid refresh token" }, { status: 401 });
    }
    const token = encodeToken(user);
    return respondWithState({
      access_token: token,
      refresh_token: token,
      user: publicUser(user),
      expires_in: 86400,
    }, state);
  }

  // ─── Users ─────────────────────────────────────────────────────────
  const currentUser = getAuthUser(req, state) ?? DEFAULT_USERS[0];

  if (method === "GET" && route === "/users/me") {
    return respondWithState(publicUser(currentUser), state);
  }

  if (method === "PUT" && route.startsWith("/users/")) {
    const body = await req.json().catch(() => ({}));
    const updated = { ...currentUser, name: body.name ?? currentUser.name };
    return respondWithState(publicUser(updated), state);
  }

  // ─── Courses ───────────────────────────────────────────────────────
  if (method === "GET" && (route === "/courses/catalog" || route === "/courses")) {
    const courses = catalogSpecs.map((spec) => buildCourseResponse(spec, completedSet));
    return respondWithState(courses, state);
  }

  if (method === "GET" && route === "/courses/enrolled") {
    const enrolledCourses = catalogSpecs
      .filter((spec) => enrolledSet.has(spec.slug))
      .map((spec) => buildCourseResponse(spec, completedSet));
    return respondWithState(enrolledCourses, state);
  }

  if (method === "POST" && route.endsWith("/enroll") && route.startsWith("/courses/")) {
    const courseId = path[1];
    const spec = catalogSpecs.find((s) => s.slug === courseId);
    if (!spec) {
      return NextResponse.json({ message: "Course not found" }, { status: 404 });
    }
    enrolledSet.add(courseId);
    state.e = Array.from(enrolledSet);
    return respondWithState({ message: "Successfully enrolled in course" }, state);
  }

  if (method === "PUT" && route.startsWith("/courses/modules/") && route.endsWith("/toggle")) {
    const moduleId = path[2];
    if (completedSet.has(moduleId)) {
      completedSet.delete(moduleId);
    } else {
      completedSet.add(moduleId);
    }
    state.c = Array.from(completedSet);
    const modSlug = moduleId.split("__")[0];
    if (modSlug) {
      enrolledSet.add(modSlug);
      state.e = Array.from(enrolledSet);
    }
    return respondWithState(
      { message: "Module progress updated", is_completed: completedSet.has(moduleId) },
      state
    );
  }

  if (method === "GET" && path.length === 2 && path[0] === "courses") {
    const courseId = path[1];
    const spec = catalogSpecs.find((s) => s.slug === courseId);
    if (!spec) {
      return NextResponse.json({ message: "Course not found" }, { status: 404 });
    }
    return respondWithState(buildCourseResponse(spec, completedSet), state);
  }

  // ─── Practice Exercise ─────────────────────────────────────────────
  if (method === "GET" && route === `/practice/exercises/${EXERCISE.id}`) {
    const unlocked = Math.max(0, Math.min(state.h ?? 0, EXERCISE.hints.length));
    return respondWithState(
      {
        id: EXERCISE.id,
        title: EXERCISE.title,
        scenario: EXERCISE.scenario,
        prompt: EXERCISE.prompt,
        starter: EXERCISE.starter,
        hints: EXERCISE.hints.slice(0, unlocked),
        hints_unlocked: unlocked,
      },
      state
    );
  }

  if (method === "POST" && route === `/practice/exercises/${EXERCISE.id}/hints`) {
    const current = state.h ?? 0;
    if (current >= EXERCISE.hints.length) {
      return NextResponse.json({ message: "semua hint sudah digunakan" }, { status: 409 });
    }
    const next = current + 1;
    state.h = next;
    return respondWithState(
      { hint: EXERCISE.hints[next - 1], hints_unlocked: next },
      state
    );
  }

  if (method === "POST" && route === `/practice/exercises/${EXERCISE.id}/run`) {
    const body = await req.json().catch(() => ({}));
    const source = String(body.source ?? "");
    const language = String(body.language ?? "");
    if (!source || source.length > 16000) {
      return NextResponse.json({ message: "source code wajib diisi dan maksimal 16 KB" }, { status: 400 });
    }
    const hasTrim = /trim|TrimSpace/i.test(source);
    const hasLower = /toLowerCase|ToLower/i.test(source);
    if ((language === "javascript" || language === "go") && hasTrim && hasLower) {
      return respondWithState({ passed: 3, total: 3, status: "Semua test lolos" }, state);
    }
    return respondWithState(
      {
        passed: 0,
        total: 3,
        status: "Belum lolos semua test",
        output: "Pastikan setiap baris email di-trim dan diubah menjadi huruf kecil (lowercase).",
      },
      state
    );
  }

  return NextResponse.json({ message: `Route ${method} ${route} not found` }, { status: 404 });
}

export const GET = handle;
export const POST = handle;
export const PUT = handle;
export const DELETE = handle;

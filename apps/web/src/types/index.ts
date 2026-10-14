// TypeScript types for LearnHub platform

export interface User {
  id: string;
  email: string;
  name: string;
  role: "learner" | "instructor" | "admin";
  is_active: boolean;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  role: "learner" | "instructor";
}

export interface Module {
  id: string;
  course_id: string;
  title: string;
  is_completed: boolean;
  duration: string;
  created_at?: string;
}

export interface Course {
  id: string;
  user_id: string;
  title: string;
  link: string;
  icon: string;
  status: "learning" | "completed" | "cancelled";
  progress: number;
  total_modules: number;
  completed_modules: number;
  modules: Module[];
  created_at?: string;
  updated_at?: string;
}

export interface CreateCourseRequest {
  title: string;
  link?: string;
  icon?: string;
}

export interface CreateModuleRequest {
  title: string;
  duration?: string;
}

export interface UpdateCourseStatusRequest {
  status: "learning" | "completed" | "cancelled";
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface ApiError {
  message: string;
  errors?: Record<string, string[]>;
}

export const COURSE_ICONS = [
  { value: "code", label: "Code / Programming", emoji: "💻" },
  { value: "design", label: "Design / UI", emoji: "🎨" },
  { value: "data", label: "Data / Analytics", emoji: "📊" },
  { value: "marketing", label: "Marketing / Growth", emoji: "📈" },
  { value: "business", label: "Bisnis / Manajemen", emoji: "💼" },
  { value: "language", label: "Bahasa", emoji: "🌐" },
  { value: "music", label: "Musik / Audio", emoji: "🎵" },
  { value: "photo", label: "Foto / Video", emoji: "📸" },
  { value: "health", label: "Kesehatan", emoji: "🏥" },
  { value: "book", label: "Lainnya", emoji: "📚" },
];

export function getIconEmoji(icon: string): string {
  return COURSE_ICONS.find((i) => i.value === icon)?.emoji ?? "📚";
}

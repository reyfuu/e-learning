import axios from "axios";
import Cookies from "js-cookie";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: { "Content-Type": "application/json" },
  timeout: 15000,
});

// Attach access token on every request
api.interceptors.request.use((config) => {
  const token = Cookies.get("access_token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Auto-refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config;
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true;
      const refresh = Cookies.get("refresh_token");
      if (refresh) {
        try {
          const { data } = await axios.post(
            `${API_URL}/api/v1/auth/refresh`,
            { refresh_token: refresh }
          );
          const { access_token, refresh_token } = data;
          Cookies.set("access_token", access_token, { expires: 1 });
          Cookies.set("refresh_token", refresh_token, { expires: 7 });
          original.headers.Authorization = `Bearer ${access_token}`;
          return api(original);
        } catch {
          Cookies.remove("access_token");
          Cookies.remove("refresh_token");
          window.location.href = "/login";
        }
      }
    }
    return Promise.reject(error);
  }
);

// ─── Auth ────────────────────────────────────────────────────────────────────
export const authApi = {
  login: (data: { email: string; password: string }) =>
    api.post("/auth/login", data).then((r) => r.data),
  register: (data: { email: string; password: string; name: string; role: "learner" | "instructor" }) =>
    api.post("/auth/register", data).then((r) => r.data),
  refresh: (refresh_token: string) =>
    api.post("/auth/refresh", { refresh_token }).then((r) => r.data),
};

// ─── Users ───────────────────────────────────────────────────────────────────
export const userApi = {
  me: () => api.get("/users/me").then((r) => r.data),
  update: (id: string, data: { name?: string }) =>
    api.put(`/users/${id}`, data).then((r) => r.data),
};

// ─── Courses ─────────────────────────────────────────────────────────────────
export const courseApi = {
  list: () => api.get("/courses").then((r) => r.data),
  catalog: () => api.get("/courses/catalog").then((r) => r.data),
  enrolled: () => api.get("/courses/enrolled").then((r) => r.data),
  get: (id: string) => api.get(`/courses/${id}`).then((r) => r.data),
  create: (data: { title: string; link?: string; icon?: string }) =>
    api.post("/courses", data).then((r) => r.data),
  enroll: (id: string) => api.post(`/courses/${id}/enroll`, {}).then((r) => r.data),
  updateStatus: (id: string, status: string) =>
    api.put(`/courses/${id}/status`, { status }).then((r) => r.data),
  createModule: (courseId: string, data: { title: string; duration?: string }) =>
    api.post(`/courses/${courseId}/modules`, data).then((r) => r.data),
  toggleModule: (moduleId: string) =>
    api.put(`/courses/modules/${moduleId}/toggle`, {}).then((r) => r.data),
};

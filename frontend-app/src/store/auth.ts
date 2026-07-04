import { create } from "zustand";
import { persist } from "zustand/middleware";
import Cookies from "js-cookie";
import { User } from "@/types";

interface AuthState {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
  setUser: (user: User) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      isAuthenticated: false,

      setAuth: (user, accessToken, refreshToken) => {
        Cookies.set("access_token", accessToken, { expires: 1, sameSite: "strict" });
        Cookies.set("refresh_token", refreshToken, { expires: 7, sameSite: "strict" });
        set({ user, accessToken, isAuthenticated: true });
      },

      logout: () => {
        Cookies.remove("access_token");
        Cookies.remove("refresh_token");
        set({ user: null, accessToken: null, isAuthenticated: false });
      },

      setUser: (user) => set({ user }),
    }),
    {
      name: "learnhub-auth",
      partialize: (state) => ({ user: state.user, isAuthenticated: state.isAuthenticated }),
    }
  )
);

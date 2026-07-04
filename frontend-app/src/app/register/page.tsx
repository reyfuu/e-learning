"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { authApi } from "@/lib/api";
import { GraduationCap, ArrowRight, Mail, Lock, Loader2, User, BookOpen, Briefcase } from "lucide-react";

export default function RegisterPage() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<"learner" | "instructor">("learner");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await authApi.register({ name, email, password, role });
      setAuth(response.user, response.access_token, response.refresh_token);
      router.push("/dashboard");
    } catch (err) {
      const message =
        typeof err === "object" && err !== null && "response" in err
          ? // eslint-disable-next-line @typescript-eslint/no-explicit-any
            (err as any).response?.data?.message
          : undefined;
      setError(message || "Registration failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative min-h-screen overflow-hidden bg-slate-950 text-slate-100 selection:bg-indigo-500/30 py-8">
      {/* Dynamic Background */}
      <div className="absolute inset-0 bg-grid opacity-20" />
      <div className="absolute inset-0 noise" />
      <div className="absolute -top-[40%] -right-[20%] h-[80%] w-[60%] rounded-full bg-indigo-600/10 blur-[120px]" />
      <div className="absolute -bottom-[40%] -left-[20%] h-[80%] w-[60%] rounded-full bg-sky-600/10 blur-[120px]" />

      <div className="relative z-10 flex min-h-screen items-center justify-center px-6 sm:px-8">
        <div className="w-full max-w-[420px]">
          {/* Main Card */}
          <div className="relative rounded-[32px] bg-slate-900/60 p-8 sm:p-10 backdrop-blur-xl border border-white/5 shadow-2xl shadow-indigo-950/50">
            {/* Header */}
            <div className="flex flex-col items-center mb-8">
              <div className="relative flex h-16 w-16 items-center justify-center rounded-3xl bg-gradient-to-br from-indigo-500 to-sky-500 text-white shadow-lg shadow-indigo-500/30 mb-6 group">
                <div className="absolute inset-0 rounded-3xl bg-white/20 opacity-0 transition-opacity group-hover:opacity-100" />
                <GraduationCap className="h-8 w-8 relative z-10" strokeWidth={2.5} />
              </div>
              
              <h1 className="text-3xl font-bold tracking-tight text-white mb-2">
                Join LearnHub
              </h1>
              <p className="text-sm text-slate-400 text-center max-w-[280px]">
                Create your account to unlock thousands of premium courses.
              </p>
            </div>

            {/* Error Message */}
            {error && (
              <div className="mb-6 flex items-center gap-3 rounded-2xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-200 animate-in fade-in slide-in-from-top-2">
                <div className="h-1.5 w-1.5 rounded-full bg-red-400 shrink-0" />
                <p>{error}</p>
              </div>
            )}

            {/* Form */}
            <form onSubmit={handleSubmit} className="space-y-5">
              
              {/* Role Selection */}
              <div className="grid grid-cols-2 gap-3 mb-2">
                <button
                  type="button"
                  onClick={() => setRole("learner")}
                  className={`flex flex-col items-center justify-center p-3 rounded-2xl border transition-all ${
                    role === "learner"
                      ? "border-indigo-500 bg-indigo-500/10 text-white shadow-[0_0_15px_rgba(99,102,241,0.15)]"
                      : "border-white/10 bg-black/20 text-slate-400 hover:border-white/20 hover:text-slate-300"
                  }`}
                >
                  <BookOpen className="h-5 w-5 mb-1.5" />
                  <span className="text-xs font-semibold">Student</span>
                </button>
                <button
                  type="button"
                  onClick={() => setRole("instructor")}
                  className={`flex flex-col items-center justify-center p-3 rounded-2xl border transition-all ${
                    role === "instructor"
                      ? "border-sky-500 bg-sky-500/10 text-white shadow-[0_0_15px_rgba(14,165,233,0.15)]"
                      : "border-white/10 bg-black/20 text-slate-400 hover:border-white/20 hover:text-slate-300"
                  }`}
                >
                  <Briefcase className="h-5 w-5 mb-1.5" />
                  <span className="text-xs font-semibold">Lecturer</span>
                </button>
              </div>

              <div className="space-y-1.5">
                <label htmlFor="name" className="text-sm font-medium text-slate-300 ml-1">
                  Full Name
                </label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-500 group-focus-within:text-indigo-400 transition-colors">
                    <User className="h-5 w-5" />
                  </div>
                  <input
                    id="name"
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    placeholder="John Doe"
                    className="w-full rounded-2xl border border-white/10 bg-black/20 pl-11 pr-4 py-3.5 text-slate-100 placeholder:text-slate-500 focus:border-indigo-500/50 focus:bg-indigo-500/5 focus:outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label htmlFor="email" className="text-sm font-medium text-slate-300 ml-1">
                  Email Address
                </label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-500 group-focus-within:text-indigo-400 transition-colors">
                    <Mail className="h-5 w-5" />
                  </div>
                  <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    placeholder="you@example.com"
                    className="w-full rounded-2xl border border-white/10 bg-black/20 pl-11 pr-4 py-3.5 text-slate-100 placeholder:text-slate-500 focus:border-indigo-500/50 focus:bg-indigo-500/5 focus:outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label htmlFor="password" className="text-sm font-medium text-slate-300 ml-1">
                  Password
                </label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-500 group-focus-within:text-indigo-400 transition-colors">
                    <Lock className="h-5 w-5" />
                  </div>
                  <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    minLength={8}
                    placeholder="Min. 8 characters"
                    className="w-full rounded-2xl border border-white/10 bg-black/20 pl-11 pr-4 py-3.5 text-slate-100 placeholder:text-slate-500 focus:border-indigo-500/50 focus:bg-indigo-500/5 focus:outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all"
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="group relative flex w-full items-center justify-center gap-2 rounded-2xl bg-gradient-to-r from-indigo-500 to-sky-500 px-8 py-3.5 text-sm font-semibold text-white shadow-[0_0_20px_rgba(99,102,241,0.3)] transition-all hover:shadow-[0_0_25px_rgba(99,102,241,0.5)] active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-70 disabled:active:scale-100 mt-4"
              >
                {loading ? (
                  <>
                    <Loader2 className="h-5 w-5 animate-spin" />
                    <span>Creating account...</span>
                  </>
                ) : (
                  <>
                    <span>Create account</span>
                    <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
                  </>
                )}
              </button>
            </form>
          </div>
          
          {/* Footer */}
          <p className="mt-8 text-center text-sm text-slate-400">
            Already have an account?{' '}
            <Link 
              href="/login" 
              className="font-semibold text-indigo-400 hover:text-indigo-300 transition-colors relative after:absolute after:bottom-0 after:left-0 after:h-[1px] after:w-full after:origin-bottom-right after:scale-x-0 after:bg-indigo-400 after:transition-transform after:duration-300 hover:after:origin-bottom-left hover:after:scale-x-100"
            >
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}

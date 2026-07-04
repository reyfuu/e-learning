"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { courseApi } from "@/lib/api";
import { ArrowLeft, BookOpen, Loader2, Save } from "lucide-react";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";

export default function AddCoursePage() {
  const router = useRouter();
  const { user } = useAuthStore();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Only instructors should access this
  if (user && user.role !== "instructor") {
    return (
      <div className="flex h-[50vh] flex-col items-center justify-center text-center">
        <h2 className="text-xl font-semibold text-white">Access Denied</h2>
        <p className="text-slate-400 mt-2">Only instructors can add new courses.</p>
        <Link href="/dashboard" className="mt-4 text-indigo-400 hover:text-indigo-300">
          Back to Dashboard
        </Link>
      </div>
    );
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await courseApi.create({ title }); // We can add description if the API accepts it later
      router.push("/dashboard/courses");
    } catch (err) {
      setError("Failed to create course. Please try again.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link 
            href="/dashboard/courses"
            className="p-2 rounded-xl bg-slate-800/50 text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
          >
            <ArrowLeft className="w-5 h-5" />
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-white">Create New Course</h1>
            <p className="text-sm text-slate-400">Fill in the details to publish a new course.</p>
          </div>
        </div>
      </div>

      <div className="bg-slate-900/60 border border-slate-800 backdrop-blur-xl rounded-[24px] p-6 sm:p-8">
        {error && (
          <div className="mb-6 p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-sm text-red-200">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label htmlFor="title" className="text-sm font-medium text-slate-300">
              Course Title
            </label>
            <div className="relative group">
              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-500 group-focus-within:text-indigo-400 transition-colors">
                <BookOpen className="h-5 w-5" />
              </div>
              <input
                id="title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                placeholder="e.g. Master React in 30 Days"
                className="w-full rounded-2xl border border-white/10 bg-black/20 pl-11 pr-4 py-3.5 text-slate-100 placeholder:text-slate-500 focus:border-indigo-500/50 focus:bg-indigo-500/5 focus:outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label htmlFor="description" className="text-sm font-medium text-slate-300">
              Description (Optional)
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="What will students learn in this course?"
              rows={4}
              className="w-full rounded-2xl border border-white/10 bg-black/20 px-4 py-3.5 text-slate-100 placeholder:text-slate-500 focus:border-indigo-500/50 focus:bg-indigo-500/5 focus:outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all resize-none"
            />
          </div>

          <div className="pt-4 flex items-center justify-end gap-3 border-t border-slate-800">
            <Link 
              href="/dashboard/courses"
              className="px-6 py-2.5 rounded-xl text-sm font-medium text-slate-300 hover:text-white hover:bg-slate-800 transition-all"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={loading || !title}
              className="flex items-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-sky-500 px-6 py-2.5 text-sm font-semibold text-white shadow-[0_0_20px_rgba(99,102,241,0.2)] transition-all hover:shadow-[0_0_25px_rgba(99,102,241,0.4)] active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50"
            >
              {loading ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              Create Course
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

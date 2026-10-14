"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { courseApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import {
  ArrowLeft,
  BookOpen,
  CheckCircle2,
  Circle,
  Clock,
  ExternalLink,
  Loader2,
  Plus,
  Layers,
} from "lucide-react";

interface Module {
  id: string;
  title: string;
  is_completed: boolean;
  duration: string;
  link?: string;
}

interface Course {
  id: string;
  title: string;
  link?: string;
  icon?: string;
  status: string;
  progress: number;
  total_modules: number;
  completed_modules: number;
  modules: Module[];
}

export default function CourseDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { user } = useAuthStore();
  const [course, setCourse] = useState<Course | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [addingModule, setAddingModule] = useState(false);
  const [newModuleTitle, setNewModuleTitle] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const [toggling, setToggling] = useState<string | null>(null);

  const fetchCourse = async () => {
    try {
      const data = await courseApi.get(id);
      setCourse(data);
    } catch {
      setError("Could not load course. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCourse();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const handleToggle = async (moduleId: string) => {
    setToggling(moduleId);
    try {
      await courseApi.toggleModule(moduleId);
      await fetchCourse();
    } catch {
      /* ignore */
    } finally {
      setToggling(null);
    }
  };

  const handleAddModule = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newModuleTitle.trim()) return;
    setAddingModule(true);
    try {
      await courseApi.createModule(id, { title: newModuleTitle });
      setNewModuleTitle("");
      setShowAddForm(false);
      await fetchCourse();
    } catch {
      /* ignore */
    } finally {
      setAddingModule(false);
    }
  };

  if (loading) {
    return (
      <div className="flex h-[60vh] items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-400" />
      </div>
    );
  }

  if (error || !course) {
    return (
      <div className="flex h-[60vh] flex-col items-center justify-center text-center gap-4">
        <BookOpen className="w-12 h-12 text-slate-600" />
        <h2 className="text-xl font-semibold text-white">Course Not Found</h2>
        <p className="text-slate-400">{error || "This course could not be loaded."}</p>
        <button
          onClick={() => router.back()}
          className="px-5 py-2 rounded-xl bg-slate-800 text-slate-300 hover:text-white transition-colors"
        >
          Go Back
        </button>
      </div>
    );
  }

  const isInstructor = user?.role === "instructor";

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Back */}
      <Link
        href="/dashboard/courses"
        className="inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors group"
      >
        <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
        Back to Courses
      </Link>

      {/* Hero Card */}
      <div className="relative rounded-[28px] overflow-hidden bg-gradient-to-br from-slate-900 via-indigo-950/40 to-slate-900 border border-slate-800 p-8">
        {/* Glows */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500/10 rounded-full blur-3xl pointer-events-none" />
        <div className="absolute bottom-0 left-0 w-64 h-64 bg-sky-500/10 rounded-full blur-3xl pointer-events-none" />

        <div className="relative z-10 flex flex-col sm:flex-row gap-6 sm:items-start justify-between">
          <div className="flex-1 space-y-3">
            {/* Status badge */}
            <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold border ${
              course.status === "published"
                ? "bg-emerald-500/10 border-emerald-500/30 text-emerald-400"
                : course.status === "draft"
                ? "bg-slate-700/60 border-slate-600/40 text-slate-400"
                : "bg-indigo-500/10 border-indigo-500/30 text-indigo-400"
            }`}>
              <span className="w-1.5 h-1.5 rounded-full bg-current" />
              {course.status.charAt(0).toUpperCase() + course.status.slice(1)}
            </span>

            <h1 className="text-3xl font-bold text-white leading-tight">{course.title}</h1>

            {course.link && (
              <a
                href={course.link}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 text-sm text-indigo-400 hover:text-indigo-300 transition-colors"
              >
                <ExternalLink className="w-4 h-4" />
                Course Resource Link
              </a>
            )}

            {/* Stats */}
            <div className="flex flex-wrap gap-4 pt-2 text-sm text-slate-400">
              <div className="flex items-center gap-1.5">
                <Layers className="w-4 h-4 text-indigo-400" />
                <span>{course.total_modules} Sections</span>
              </div>
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                <span>{course.completed_modules} Completed</span>
              </div>
              <div className="flex items-center gap-1.5">
                <Clock className="w-4 h-4 text-sky-400" />
                <span>{course.progress}% Progress</span>
              </div>
            </div>
          </div>

          {/* Progress ring */}
          <div className="flex items-center justify-center shrink-0">
            <div className="relative w-24 h-24">
              <svg className="w-full h-full -rotate-90" viewBox="0 0 100 100">
                <circle cx="50" cy="50" r="42" fill="none" stroke="rgb(30,41,59)" strokeWidth="10" />
                <circle
                  cx="50" cy="50" r="42" fill="none"
                  stroke="url(#progressGrad)" strokeWidth="10"
                  strokeLinecap="round"
                  strokeDasharray={`${2 * Math.PI * 42}`}
                  strokeDashoffset={`${2 * Math.PI * 42 * (1 - course.progress / 100)}`}
                  className="transition-all duration-1000"
                />
                <defs>
                  <linearGradient id="progressGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" stopColor="#6366f1" />
                    <stop offset="100%" stopColor="#38bdf8" />
                  </linearGradient>
                </defs>
              </svg>
              <div className="absolute inset-0 flex items-center justify-center">
                <span className="text-xl font-bold text-white">{course.progress}%</span>
              </div>
            </div>
          </div>
        </div>

        {/* Progress bar */}
        <div className="relative z-10 mt-6">
          <div className="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-indigo-500 to-sky-400 rounded-full transition-all duration-1000"
              style={{ width: `${course.progress}%` }}
            />
          </div>
        </div>
      </div>

      {/* Sections / Modules */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Sections</h2>
          {isInstructor && (
            <button
              onClick={() => setShowAddForm((v) => !v)}
              className="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 hover:bg-indigo-500/20 transition-colors"
            >
              <Plus className="w-4 h-4" />
              Add Section
            </button>
          )}
        </div>

        {/* Add module form */}
        {showAddForm && isInstructor && (
          <form
            onSubmit={handleAddModule}
            className="flex gap-3 p-4 rounded-2xl bg-slate-900/60 border border-slate-700 backdrop-blur-xl"
          >
            <input
              type="text"
              value={newModuleTitle}
              onChange={(e) => setNewModuleTitle(e.target.value)}
              placeholder="Section title..."
              autoFocus
              className="flex-1 bg-black/20 border border-white/10 rounded-xl px-4 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:border-indigo-500/50 focus:ring-2 focus:ring-indigo-500/10 transition-all"
            />
            <button
              type="submit"
              disabled={addingModule || !newModuleTitle.trim()}
              className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-indigo-500 to-sky-500 text-sm font-semibold text-white disabled:opacity-50"
            >
              {addingModule ? <Loader2 className="w-4 h-4 animate-spin" /> : "Add"}
            </button>
            <button
              type="button"
              onClick={() => setShowAddForm(false)}
              className="px-4 py-2.5 rounded-xl text-sm text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
            >
              Cancel
            </button>
          </form>
        )}

        {course.modules.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 rounded-[24px] border border-dashed border-slate-800 text-center">
            <BookOpen className="w-10 h-10 text-slate-600 mb-3" />
            <p className="text-slate-400">No sections yet.</p>
            {isInstructor && (
              <button
                onClick={() => setShowAddForm(true)}
                className="mt-4 px-5 py-2 rounded-xl text-sm text-indigo-400 border border-indigo-500/20 hover:bg-indigo-500/10 transition-colors"
              >
                Add your first section
              </button>
            )}
          </div>
        ) : (
          <div className="space-y-3">
            {course.modules.map((module, index) => (
              <div
                key={module.id}
                className={`group flex items-center gap-4 p-4 rounded-2xl border transition-all ${
                  module.is_completed
                    ? "bg-emerald-500/5 border-emerald-500/20"
                    : "bg-slate-900/60 border-slate-800 hover:border-slate-700"
                }`}
              >
                {/* Toggle completion */}
                <button
                  onClick={() => handleToggle(module.id)}
                  disabled={toggling === module.id}
                  className="shrink-0 text-slate-500 hover:text-indigo-400 transition-colors"
                  title={module.is_completed ? "Mark as incomplete" : "Mark as complete"}
                >
                  {toggling === module.id ? (
                    <Loader2 className="w-5 h-5 animate-spin text-indigo-400" />
                  ) : module.is_completed ? (
                    <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                  ) : (
                    <Circle className="w-5 h-5" />
                  )}
                </button>

                {/* Index */}
                <span className="shrink-0 text-xs font-bold text-slate-600 w-5 text-center">
                  {String(index + 1).padStart(2, "0")}
                </span>

                {/* Content */}
                <div className="flex-1 min-w-0">
                  <p className={`text-sm font-medium truncate ${
                    module.is_completed ? "text-slate-400 line-through" : "text-slate-100"
                  }`}>
                    {module.title}
                  </p>
                  <div className="flex items-center gap-3 mt-0.5">
                    <span className="text-xs text-slate-500 flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      {module.duration}
                    </span>
                  </div>
                </div>

                {/* External link */}
                {module.link && (
                  <a
                    href={module.link}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="shrink-0 p-2 rounded-lg text-slate-500 hover:text-indigo-400 hover:bg-indigo-500/10 transition-all"
                    title="Open resource"
                  >
                    <ExternalLink className="w-4 h-4" />
                  </a>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

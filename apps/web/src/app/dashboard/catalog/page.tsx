"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { courseApi } from "@/lib/api";
import { BookOpen, Clock, Layers, Search, UserRound, Loader2, CheckCircle2 } from "lucide-react";

interface Course {
  id: string;
  title: string;
  link?: string; // holds instructor name from API
  status: string;
  total_modules: number;
  progress: number;
}

export default function CatalogPage() {
  const [courses, setCourses] = useState<Course[]>([]);
  const [enrolledIds, setEnrolledIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const [enrollingId, setEnrollingId] = useState<string | null>(null);
  const [search, setSearch] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [catalog, enrolled] = await Promise.all([
          courseApi.catalog(),
          courseApi.enrolled(),
        ]);
        setCourses(catalog);
        setEnrolledIds(new Set((enrolled as Course[]).map((c) => c.id)));
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const handleEnroll = async (courseId: string) => {
    setEnrollingId(courseId);
    try {
      await courseApi.enroll(courseId);
      setEnrolledIds((prev) => new Set([...prev, courseId]));
    } catch (err) {
      console.error(err);
    } finally {
      setEnrollingId(null);
    }
  };

  const filtered = courses.filter((c) =>
    c.title.toLowerCase().includes(search.toLowerCase()) ||
    (c.link || "").toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Browse Courses</h1>
          <p className="text-sm text-slate-400 mt-1">
            Discover and enroll in courses created by our lecturers.
          </p>
        </div>

        {/* Search */}
        <div className="relative sm:w-72">
          <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-500">
            <Search className="w-4 h-4" />
          </div>
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search courses..."
            className="w-full rounded-2xl border border-white/10 bg-slate-900/60 pl-10 pr-4 py-2.5 text-sm text-slate-100 placeholder:text-slate-500 focus:border-indigo-500/50 focus:outline-none focus:ring-2 focus:ring-indigo-500/10 transition-all"
          />
        </div>
      </div>

      {/* Grid */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="h-56 rounded-2xl bg-slate-800/40 animate-pulse border border-slate-800" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-24 rounded-[32px] bg-slate-900/40 border border-dashed border-slate-800 text-center">
          <BookOpen className="w-12 h-12 text-slate-600 mb-4" />
          <h3 className="text-xl font-semibold text-white mb-2">No Courses Found</h3>
          <p className="text-slate-400 max-w-xs">
            {search ? `No courses match "${search}".` : "No courses have been created by lecturers yet."}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filtered.map((course) => {
            const isEnrolled = enrolledIds.has(course.id);
            const isEnrolling = enrollingId === course.id;

            return (
              <div
                key={course.id}
                className="group flex flex-col rounded-2xl bg-slate-900/60 border border-slate-800 overflow-hidden hover:border-indigo-500/40 transition-all hover:shadow-[0_0_30px_rgba(99,102,241,0.12)]"
              >
                {/* Thumbnail */}
                <div className="h-36 bg-gradient-to-br from-slate-800 to-indigo-950/40 relative p-5 flex flex-col justify-end overflow-hidden">
                  <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-500/10 rounded-full blur-2xl group-hover:bg-indigo-500/20 transition-colors" />
                  <div className="absolute -bottom-8 -left-8 w-32 h-32 bg-sky-500/10 rounded-full blur-2xl" />

                  {/* Status badge */}
                  <span className={`relative z-10 self-start mb-2 inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[10px] font-semibold border ${
                    course.status === "published"
                      ? "bg-emerald-500/10 border-emerald-500/30 text-emerald-400"
                      : "bg-slate-700/60 border-slate-600/40 text-slate-400"
                  }`}>
                    <span className="w-1 h-1 rounded-full bg-current" />
                    {course.status}
                  </span>

                  <h3 className="text-base font-semibold text-white relative z-10 leading-snug line-clamp-2">
                    {course.title}
                  </h3>
                </div>

                {/* Info */}
                <div className="p-4 flex flex-col flex-1 gap-3">
                  <div className="flex items-center gap-4 text-xs text-slate-400">
                    {course.link && (
                      <div className="flex items-center gap-1.5 truncate">
                        <UserRound className="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                        <span className="truncate">{course.link}</span>
                      </div>
                    )}
                    <div className="flex items-center gap-1.5 shrink-0">
                      <Layers className="w-3.5 h-3.5 text-sky-400" />
                      <span>{course.total_modules} sections</span>
                    </div>
                  </div>

                  <div className="mt-auto flex gap-2">
                    {isEnrolled ? (
                      <Link
                        href={`/dashboard/courses/${course.id}`}
                        className="flex-1 flex items-center justify-center gap-2 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm font-semibold py-2 hover:bg-emerald-500/20 transition-colors"
                      >
                        <CheckCircle2 className="w-4 h-4" />
                        Open Course
                      </Link>
                    ) : (
                      <button
                        onClick={() => handleEnroll(course.id)}
                        disabled={isEnrolling}
                        className="flex-1 flex items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-sky-500 text-white text-sm font-semibold py-2 shadow-[0_0_15px_rgba(99,102,241,0.2)] hover:shadow-[0_0_20px_rgba(99,102,241,0.4)] transition-all active:scale-[0.98] disabled:opacity-60"
                      >
                        {isEnrolling ? (
                          <Loader2 className="w-4 h-4 animate-spin" />
                        ) : (
                          "Enroll Now"
                        )}
                      </button>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

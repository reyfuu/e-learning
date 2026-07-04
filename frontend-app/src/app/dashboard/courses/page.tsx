"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { courseApi } from "@/lib/api";
import { BookOpen, Clock, FileText, Plus, Trophy } from "lucide-react";

export default function CoursesPage() {
  const { user } = useAuthStore();
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [courses, setCourses] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCourses = async () => {
      try {
        let data;
        if (user?.role === "instructor") {
          data = await courseApi.list();
          setCourses(data.courses || data || []);
        } else {
          data = await courseApi.enrolled();
          setCourses(data || []);
        }
      } catch (err) {
        console.error("Failed to fetch courses:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchCourses();
  }, [user?.role]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">
            {user?.role === "instructor" ? "My Courses" : "My Learning"}
          </h1>
          <p className="text-sm text-slate-400 mt-1">
            {user?.role === "instructor"
              ? "Manage the courses you are teaching."
              : "Track your progress and continue learning."}
          </p>
        </div>

        {user?.role === "instructor" && (
          <Link
            href="/dashboard/courses/new"
            className="flex items-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-sky-500 px-5 py-2.5 text-sm font-semibold text-white shadow-[0_0_20px_rgba(99,102,241,0.2)] transition-all hover:shadow-[0_0_25px_rgba(99,102,241,0.4)] active:scale-[0.98]"
          >
            <Plus className="w-4 h-4" />
            Add New Course
          </Link>
        )}
      </div>

      {/* Courses Grid */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-64 rounded-2xl bg-slate-800/40 animate-pulse border border-slate-800" />
          ))}
        </div>
      ) : courses.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 px-4 text-center rounded-[32px] bg-slate-900/40 border border-dashed border-slate-800">
          <div className="w-16 h-16 rounded-full bg-slate-800/80 flex items-center justify-center mb-4">
            <BookOpen className="w-8 h-8 text-slate-500" />
          </div>
          <h3 className="text-xl font-semibold text-white mb-2">No Courses Found</h3>
          <p className="text-slate-400 max-w-sm">
            {user?.role === "instructor"
              ? "You haven't created any courses yet. Start by creating your first course."
              : "You haven't enrolled in any courses yet. Browse our catalog to get started."}
          </p>
          {user?.role === "instructor" && (
            <Link
              href="/dashboard/courses/new"
              className="mt-6 px-6 py-2.5 rounded-xl bg-white/10 hover:bg-white/15 text-white font-medium transition-colors border border-white/5"
            >
              Create First Course
            </Link>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course) => (
            <Link
              key={course.id}
              href={`/dashboard/courses/${course.id}`}
              className="group flex flex-col rounded-2xl bg-slate-900/60 border border-slate-800 overflow-hidden hover:border-indigo-500/50 transition-all hover:shadow-[0_0_30px_rgba(99,102,241,0.15)]"
            >
              {/* Thumbnail Area */}
              <div className="h-40 bg-gradient-to-br from-slate-800 to-slate-900 relative p-6 flex flex-col justify-end overflow-hidden">
                <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-500/10 rounded-full blur-2xl group-hover:bg-indigo-500/20 transition-colors" />
                <div className="absolute -bottom-8 -left-8 w-32 h-32 bg-sky-500/10 rounded-full blur-2xl group-hover:bg-sky-500/20 transition-colors" />
                
                <h3 className="text-lg font-semibold text-white relative z-10 leading-tight">
                  {course.title}
                </h3>
              </div>

              {/* Course Info */}
              <div className="p-5 flex flex-col flex-1 gap-4">
                <div className="flex items-center gap-4 text-sm text-slate-400">
                  <div className="flex items-center gap-1.5">
                    <FileText className="w-4 h-4 text-indigo-400" />
                    <span>{course.total_modules || 0} Modules</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <Clock className="w-4 h-4 text-sky-400" />
                    <span>{course.status || 'Draft'}</span>
                  </div>
                </div>

                {/* Progress bar */}
                <div className="mt-auto space-y-2">
                  <div className="flex items-center justify-between text-xs font-medium">
                    <span className="text-slate-400">Progress</span>
                    <span className="text-white">{course.progress || 0}%</span>
                  </div>
                  <div className="h-2 w-full bg-slate-800 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-gradient-to-r from-indigo-500 to-sky-400 rounded-full transition-all duration-1000 ease-out relative"
                      style={{ width: `${course.progress || 0}%` }}
                    >
                      <div className="absolute top-0 right-0 bottom-0 w-10 bg-white/20 blur-[2px]" />
                    </div>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

"use client";

import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import {
  BookOpen,
  GraduationCap,
  LayoutDashboard,
  LogOut,
  Plus,
  User,
  Compass,
} from "lucide-react";
import { useAuthStore } from "@/store/auth";
import { cn, getInitials } from "@/lib/utils";

interface SidebarProps {
  activePath?: string;
}

export function Sidebar({ activePath: activeProp = "" }: SidebarProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { user, logout } = useAuthStore();

  const activePath = activeProp || pathname;

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  return (
    <aside className="flex flex-col w-64 min-h-screen bg-slate-900/80 border-r border-slate-800 backdrop-blur-xl shrink-0">
      {/* Logo */}
      <div className="flex items-center gap-3 px-6 py-5 border-b border-slate-800">
        <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/30">
          <GraduationCap className="w-5 h-5 text-white" />
        </div>
        <div>
          <span className="text-base font-bold text-slate-100">LearnHub</span>
          <p className="text-[10px] text-slate-500 uppercase tracking-wider">E-Learning</p>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 p-4 space-y-1">
        {/* Common links */}
        {[
          { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
          { href: "/dashboard/courses", label: user?.role === "instructor" ? "My Courses" : "My Learning", icon: BookOpen },
        ].map(({ href, label, icon: Icon }) => {
          const isActive = activePath === href;
          return (
            <Link
              key={href}
              href={href}
              className={cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200",
                isActive
                  ? "bg-indigo-600/20 text-indigo-400 border border-indigo-500/20"
                  : "text-slate-400 hover:text-slate-100 hover:bg-slate-800/60"
              )}
            >
              <Icon className="w-4 h-4 shrink-0" />
              {label}
              {isActive && (
                <span className="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-400" />
              )}
            </Link>
          );
        })}

        {/* Student-only: Browse catalog */}
        {user?.role !== "instructor" && (
          <Link
            href="/dashboard/catalog"
            className={cn(
              "flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200",
              activePath === "/dashboard/catalog"
                ? "bg-indigo-600/20 text-indigo-400 border border-indigo-500/20"
                : "text-slate-400 hover:text-slate-100 hover:bg-slate-800/60"
            )}
          >
            <Compass className="w-4 h-4 shrink-0" />
            Browse Courses
            {activePath === "/dashboard/catalog" && (
              <span className="ml-auto w-1.5 h-1.5 rounded-full bg-indigo-400" />
            )}
          </Link>
        )}

        {/* Instructor-only: Add Course */}
        {user?.role === "instructor" && (
          <div className="pt-4">
            <p className="px-3 mb-2 text-[10px] font-semibold uppercase tracking-wider text-slate-600">
              Actions
            </p>
            <Link
              href="/dashboard/courses/new"
              className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-400 hover:text-slate-100 hover:bg-slate-800/60 transition-all"
            >
              <Plus className="w-4 h-4" />
              Add Course
            </Link>
          </div>
        )}
      </nav>


      {/* User profile */}
      <div className="p-4 border-t border-slate-800">
        <div className="flex items-center gap-3 p-3 rounded-xl bg-slate-800/60 border border-slate-700/50">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-xs font-bold text-white shrink-0">
            {user ? getInitials(user.name) : <User className="w-4 h-4" />}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-slate-100 truncate">{user?.name ?? "User"}</p>
            <p className="text-xs text-slate-500 truncate capitalize">{user?.role || "Student"}</p>
          </div>
          <button
            onClick={handleLogout}
            title="Logout"
            className="text-slate-500 hover:text-red-400 transition-colors p-1 rounded-lg hover:bg-red-400/10"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </aside>
  );
}

"use client";

import { useAuthStore } from "@/store/auth";
import LearnerDashboard from "@/components/dashboard/LearnerDashboard";
import InstructorDashboard from "@/components/dashboard/InstructorDashboard";
import AdminDashboard from "../admin/dashboard/page"; // We can reuse the page if we want, or eventually move it to components

export default function DashboardPage() {
  const { user, isAuthenticated } = useAuthStore();

  if (!isAuthenticated || !user) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <p className="text-lg text-muted-foreground">Silakan login terlebih dahulu...</p>
      </div>
    );
  }

  switch (user.role) {
    case "admin":
      return <AdminDashboard />;
    case "instructor":
      return <InstructorDashboard />;
    case "learner":
    default:
      return <LearnerDashboard />;
  }
}

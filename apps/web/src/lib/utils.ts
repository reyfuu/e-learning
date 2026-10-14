import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(dateString: string): string {
  return new Intl.DateTimeFormat("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(new Date(dateString));
}

export function formatDuration(minutes: number): string {
  if (minutes < 60) return `${minutes} menit`;
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  return m > 0 ? `${h}j ${m}m` : `${h} jam`;
}

export function getInitials(name: string): string {
  return name
    .split(" ")
    .slice(0, 2)
    .map((n) => n[0])
    .join("")
    .toUpperCase();
}

export function getStatusColor(status: string) {
  switch (status) {
    case "completed":
      return "text-emerald-400 bg-emerald-400/10 border-emerald-400/20";
    case "learning":
      return "text-brand-400 bg-brand-400/10 border-brand-400/20";
    case "cancelled":
      return "text-red-400 bg-red-400/10 border-red-400/20";
    default:
      return "text-slate-400 bg-slate-400/10 border-slate-400/20";
  }
}

export function getStatusLabel(status: string) {
  switch (status) {
    case "completed": return "Selesai";
    case "learning": return "Sedang Belajar";
    case "cancelled": return "Dibatalkan";
    default: return status;
  }
}

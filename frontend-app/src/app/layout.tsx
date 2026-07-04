import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Providers } from "@/providers/providers";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "LearnHub — Platform E-Learning Modern",
    template: "%s | LearnHub",
  },
  description:
    "Platform e-learning terbaik untuk belajar skill baru. Akses ratusan kursus dari instruktur berpengalaman kapan saja dan di mana saja.",
  keywords: ["e-learning", "kursus online", "belajar online", "platform belajar"],
  openGraph: {
    title: "LearnHub — Platform E-Learning Modern",
    description: "Belajar skill baru bersama instruktur terbaik Indonesia",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id" className={`${inter.variable} h-full`}>
      <body className="min-h-full antialiased bg-[#0f172a] text-slate-100">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}

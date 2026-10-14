"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

const EXERCISE = "employee-email-normalizer";
type Exercise = { title: string; scenario: string; prompt: string; starter: Record<string, string>; hints_unlocked: number };
type RunResult = { passed: number; total: number; status: string; output?: string };

export default function PracticePage() {
  const [exercise, setExercise] = useState<Exercise | null>(null);
  const [language, setLanguage] = useState("javascript");
  const [source, setSource] = useState("");
  const [hints, setHints] = useState<string[]>([]);
  const [result, setResult] = useState<RunResult | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    api.get(`/practice/exercises/${EXERCISE}`).then(({ data }) => {
      setExercise(data);
      setHints(data.hints ?? []);
      setSource(data.starter.javascript);
    }).catch(() => setError("Latihan gagal dimuat. Silakan muat ulang halaman."));
  }, []);

  const changeLanguage = (next: string) => {
    setLanguage(next);
    setSource(exercise?.starter[next] ?? "");
    setResult(null);
  };

  const unlockHint = async () => {
    setBusy(true); setError("");
    try {
      const { data } = await api.post(`/practice/exercises/${EXERCISE}/hints`);
      setHints((current) => [...current, data.hint]);
      setExercise((current) => current ? { ...current, hints_unlocked: data.hints_unlocked } : current);
    } catch (e: any) {
      setError(e.response?.data?.message ?? "Hint gagal dibuka.");
    } finally { setBusy(false); }
  };

  const runCode = async () => {
    setBusy(true); setError(""); setResult(null);
    try {
      const { data } = await api.post(`/practice/exercises/${EXERCISE}/run`, { language, source }, { timeout: 45000 });
      setResult(data);
    } catch (e: any) {
      setError(e.response?.data?.message ?? "Kode gagal dijalankan.");
    } finally { setBusy(false); }
  };

  if (!exercise) return <main className="mx-auto max-w-5xl p-8 text-slate-300">{error || "Memuat latihan…"}</main>;

  return (
    <main className="mx-auto max-w-6xl space-y-6">
      <header>
        <p className="text-sm font-semibold uppercase tracking-wider text-indigo-400">Latihan kerja · Backend</p>
        <h1 className="mt-2 text-3xl font-bold text-white">{exercise.title}</h1>
      </header>

      <section className="rounded-2xl border border-slate-800 bg-slate-900 p-6">
        <h2 className="font-semibold text-white">Situasi kerja</h2>
        <p className="mt-2 text-slate-300">{exercise.scenario}</p>
        <h3 className="mt-5 font-semibold text-white">Tugas</h3>
        <p className="mt-2 text-slate-300">{exercise.prompt}</p>
      </section>

      <section className="overflow-hidden rounded-2xl border border-slate-800 bg-slate-950">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-slate-800 p-4">
          <label className="flex items-center gap-3 text-sm text-slate-300">Bahasa
            <select aria-label="Bahasa pemrograman" value={language} onChange={(e) => changeLanguage(e.target.value)} className="rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white">
              <option value="javascript">JavaScript (Node.js)</option><option value="go">Go</option>
            </select>
          </label>
          <button onClick={runCode} disabled={busy} className="rounded-lg bg-indigo-500 px-5 py-2 font-semibold text-white hover:bg-indigo-400 disabled:cursor-wait disabled:opacity-60">
            {busy ? "Memproses…" : "Jalankan & periksa"}
          </button>
        </div>
        <label htmlFor="code-editor" className="sr-only">Editor kode</label>
        <textarea id="code-editor" spellCheck={false} aria-label="Editor kode" value={source} onChange={(e) => setSource(e.target.value)} className="min-h-[360px] w-full resize-y bg-[#0b1120] p-5 font-mono text-sm leading-6 text-slate-100 outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500" />
        <div aria-live="polite" className="border-t border-slate-800 p-4">
          {result && <p className={result.passed === result.total ? "font-semibold text-emerald-400" : "font-semibold text-amber-300"}>{result.status} · {result.passed}/{result.total} test lolos{result.output ? ` — ${result.output}` : ""}</p>}
          {error && <p role="alert" className="text-sm text-rose-300">{error}</p>}
          {!error && !result && <p className="text-sm text-slate-500">Kode diperiksa dengan beberapa input uji. Maksimal 16 KB per pengiriman.</p>}
        </div>
      </section>

      <section className="rounded-2xl border border-slate-800 bg-slate-900 p-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div><h2 className="font-semibold text-white">Petunjuk bertahap</h2><p className="mt-1 text-sm text-slate-400">{exercise.hints_unlocked} dari 5 hint telah dibuka.</p></div>
          <button onClick={unlockHint} disabled={busy || exercise.hints_unlocked >= 5} className="rounded-lg border border-indigo-400/40 px-4 py-2 text-sm font-medium text-indigo-300 hover:bg-indigo-500/10 disabled:opacity-50">Buka hint berikutnya</button>
        </div>
        {hints.length > 0 && <ol className="mt-4 list-inside list-decimal space-y-2 text-slate-300">{hints.map((hint, i) => <li key={i}>{hint}</li>)}</ol>}
      </section>
    </main>
  );
}

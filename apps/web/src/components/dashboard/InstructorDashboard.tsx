export default function InstructorDashboard() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-4">Instructor Dashboard</h1>
      <p className="text-muted-foreground">
        Selamat datang, Instruktur! Kelola kursus Anda dan pantau pendapatan di sini.
      </p>
      
      <div className="mt-8 grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Total Pendapatan</h3>
          <p className="text-3xl font-bold text-green-400">Rp 2.5M</p>
        </div>
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Siswa Aktif</h3>
          <p className="text-3xl font-bold text-brand-400">128</p>
        </div>
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Rating Rata-rata</h3>
          <p className="text-3xl font-bold text-yellow-400">4.8</p>
        </div>
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Kursus Aktif</h3>
          <p className="text-3xl font-bold text-accent-400">5</p>
        </div>
      </div>
    </div>
  );
}

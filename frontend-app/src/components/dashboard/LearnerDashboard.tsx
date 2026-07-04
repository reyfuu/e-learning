export default function LearnerDashboard() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-4">Learner Dashboard</h1>
      <p className="text-muted-foreground">
        Selamat datang! Di sini Anda dapat melihat progres belajar dan daftar kursus yang Anda ikuti.
      </p>
      
      <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Kursus Aktif</h3>
          <p className="text-3xl font-bold text-brand-400">3</p>
        </div>
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Sertifikat</h3>
          <p className="text-3xl font-bold text-accent-400">1</p>
        </div>
        <div className="glass p-6 rounded-xl">
          <h3 className="font-semibold text-lg mb-2">Jam Belajar</h3>
          <p className="text-3xl font-bold text-green-400">12j</p>
        </div>
      </div>
    </div>
  );
}

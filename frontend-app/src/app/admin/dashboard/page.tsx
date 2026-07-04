export default function AdminDashboard() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-4">Admin Panel</h1>
      <p className="text-muted-foreground">
        Overview performa platform LearnHub.
      </p>
      
      <div className="mt-8 grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="glass p-6 rounded-xl border-l-4 border-brand-500">
          <h3 className="font-semibold text-lg mb-2">Total Pengguna</h3>
          <p className="text-3xl font-bold">1,240</p>
        </div>
        <div className="glass p-6 rounded-xl border-l-4 border-accent-500">
          <h3 className="font-semibold text-lg mb-2">Total Kursus</h3>
          <p className="text-3xl font-bold">45</p>
        </div>
        <div className="glass p-6 rounded-xl border-l-4 border-green-500">
          <h3 className="font-semibold text-lg mb-2">Revenue Bulan Ini</h3>
          <p className="text-3xl font-bold">Rp 15.2M</p>
        </div>
        <div className="glass p-6 rounded-xl border-l-4 border-red-500">
          <h3 className="font-semibold text-lg mb-2">Kursus Pending Review</h3>
          <p className="text-3xl font-bold">3</p>
        </div>
      </div>
    </div>
  );
}

# Materi: AWS Cloud Architecture dan Rekayasa Layanan Terkelola

**Format:** teks · **Level:** Basic–Advanced · **Estimasi:** 60–90 menit

## Tujuan belajar

Setelah materi ini, pelajar dapat merancang arsitektur cloud yang aman, andal, dan skalabel di Amazon Web Services (AWS): memahami AWS Well-Architected Framework, mengelola keamanan identitas dengan IAM (*Least Privilege*), merancang jaringan privat dengan Amazon VPC (Subnets, Route Tables, NAT Gateways), menyebarkan komputasi (EC2, ECS/EKS, Lambda Serverless), mengelola penyimpanan data (S3, RDS Multi-AZ, DynamoDB), serta memantau operasi dengan CloudWatch.

## Prasyarat

- Memahami konsep dasar jaringan komputer (CIDR, IP address, Subnet, DNS, Port, Firewall).
- Memiliki akun AWS (Free Tier) dan AWS CLI terpasang di sistem lokal.
- Memahami konsep dasar virtualisasi dan server Linux.

## 1. Dasar: arsitektur global AWS dan Shared Responsibility Model

AWS (Amazon Web Services) adalah penyedia platform komputasi awan terbesar di dunia. Infrastruktur global AWS dibangun di atas konsep **Regions** (wilayah geografis fisik di seluruh dunia) dan **Availability Zones (AZ)** (satu atau beberapa data center fisik terisolasi dengan koneksi latensi ultra-rendah di dalam satu Region).

### Model Tanggung Jawab Bersama (Shared Responsibility Model):
- **Keamanan *OF* the Cloud (Tanggung Jawab AWS):** Mengamankan perangkat keras fisik, fasilitas data center, jaringan global, dan software virtualisasi hipervisor.
- **Keamanan *IN* the Cloud (Tanggung Jawab Pengguna):** Mengonfigurasi enkripsi data, firewall (*Security Groups*), manajemen OS tamu (patching EC2), dan hak akses pengguna (*IAM policies*).

## 2. Praktik dasar: IAM Least Privilege dan S3 Bucket Hardening

```json
// Kebijakan IAM (Policy JSON) yang menerapkan prinsip hak akses terendah (Least Privilege):
// Hanya mengizinkan upload materi video ke folder 'uploads/' di bucket spesifik LearnHub
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowListCourseBucket",
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::learnhub-media-production",
      "Condition": {
        "StringLike": {
          "s3:prefix": ["uploads/*"]
        }
      }
    },
    {
      "Sid": "AllowUploadOnlyUnderUploadsFolder",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:PutObjectAcl"
      ],
      "Resource": "arn:aws:s3:::learnhub-media-production/uploads/*"
    }
  ]
}
```

Pemberian hak akses sementara (*Presigned URL*) via AWS SDK (Node.js/TypeScript) untuk upload video tanpa mengekspos kredensial:
```typescript
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';

const s3Client = new S3Client({ region: 'ap-southeast-1' });

export async function generateUploadUrl(fileName: string, contentType: string): Promise<string> {
  const command = new PutObjectCommand({
    Bucket: 'learnhub-media-production',
    Key: `videos/${Date.now()}-${fileName}`,
    ContentType: contentType,
  });

  // URL aman yang hanya valid selama 15 menit
  return await getSignedUrl(s3Client, command, { expiresIn: 900 });
}
```

## 3. Fitur utama dan pola arsitektur

### Desain Jaringan Virtual Private Cloud (VPC)
Pola arsitektur standar enterprise membagi VPC menjadi tiga lapisan subnet lintas minimal 2 Availability Zone:
1. **Public Subnet:** Berisi Application Load Balancer (ALB) dan NAT Gateway. Memiliki rute langsung ke *Internet Gateway (IGW)*.
2. **Private Application Subnet:** Berisi kontainer aplikasi backend (ECS/EC2) tanpa IP publik langsung. Berkomunikasi ke internet hanya melalui NAT Gateway.
3. **Isolated Database Subnet:** Berisi database PostgreSQL/MySQL (RDS Multi-AZ). Tidak memiliki akses internet sama sekali untuk keamanan data absolut.

```text
Internet → [Internet Gateway] 
             ↓
        [Public Subnet] → Application Load Balancer
             ↓
        [Private App Subnet] → Web / API Containers (EC2 / ECS)
             ↓
        [Database Subnet] → Amazon RDS Multi-AZ (PostgreSQL)
```

### Komputasi Serverless dengan AWS Lambda
Mengeksekusi fungsi backend hanya saat ada event/request masuk, dengan penskalaan instan dan biaya $0 saat tidak ada trafik:

```typescript
// lambda/processThumbnail.ts
import { S3Event } from 'aws-lambda';

export const handler = async (event: S3Event) => {
  for (const record of event.Records) {
    const bucket = record.s3.bucket.name;
    const key = decodeURIComponent(record.s3.object.key.replace(/\+/g, ' '));
    console.log(`Memproses thumbnail untuk file baru: ${bucket}/${key}`);
    // Resize gambar dan simpan ke folder /thumbnails
  }
  return { status: 200 };
};
```

## 4. Praktik kerja: mitigasi insiden kebocoran kredensial root dan audit trail

### Skenario
Developer secara tidak sengaja mengunggah Access Key dan Secret Key akun root AWS ke repository GitHub publik. Dalam hitungan menit, bot penyerang menyalakan 50 instance EC2 besar bertipe GPU di Region Frankfurt untuk menambang cryptocurrency (*crypto-mining*), menimbulkan tagihan ribuan dolar.

### Tugas
Terapkan prosedur tanggap darurat (*Incident Response*): cabut (*revoke*) kredensial segera, hapus instance tidak sah, aktifkan AWS CloudTrail dan AWS Budgets, serta wajibkan Multi-Factor Authentication (MFA) pada seluruh akses.

### Langkah Tanggap Darurat:
1. **Cabut Kredensial:** Buka AWS IAM Console menggunakan akun root via browser, segera *Deactivate* dan *Delete* Access Key yang bocor.
2. **Terminasi Resource Tak Dikenal:** Periksa EC2 Console di semua Region aktif dan lakukan *Terminate* pada seluruh instance tak sah.
3. **Investigasi Audit Trail:** Buka **AWS CloudTrail**, filter event `RunInstances` berdasarkan Access Key ID penyerang untuk mengidentifikasi seluruh resource yang sempat dimodifikasi.
4. **Pasang AWS Budgets:** Atur alarm batas anggaran (misal: kirim email peringatan jika tagihan bulanan diproyeksikan melebihi $50).
5. **Kebijakan Pencegahan:** Kunci akun root dan jangan pernah membuat Access Key untuk akun root; buat IAM User atau IAM Identity Center (SSO) dengan izin terbatas.

### Lima hint bertahap
1. Jangan pernah menyimpan `AWS_ACCESS_KEY_ID` atau `AWS_SECRET_ACCESS_KEY` di dalam kode program atau file git.
2. Jika aplikasi berjalan di EC2 atau ECS, gunakan **IAM Roles** alih-alih hardcoded access key (*instance profile*).
3. Gunakan AWS Secrets Manager atau AWS Systems Manager Parameter Store untuk menyimpan rahasia konfigurasi database.
4. AWS CloudTrail mencatat setiap panggilan API yang terjadi di akun AWS Anda secara permanen.
5. S3 Block Public Access harus selalu aktif secara default untuk mencegah kebocoran data sensitif ke publik.

## 5. Lanjutan: Enkripsi KMS, Auto Scaling, dan Well-Architected Pillars

- **AWS KMS (Key Management Service):** Enkripsi data saat istirahat (*Encryption at Rest*) untuk S3, EBS, dan RDS menggunakan master key yang dikelola perangkat keras FIPS 140-2.
- **Auto Scaling Group (ASG):** Otomatis menambah atau mengurangi jumlah instance komputasi berdasarkan metrik utilisasi CPU atau antrean request target.
- **Enam Pilar AWS Well-Architected Framework:**
  1. *Operational Excellence:* Otomasi deployment dan pemantauan.
  2. *Security:* Perlindungan data dan integritas hak akses.
  3. *Reliability:* Kemampuan pemulihan diri dari kegagalan infrastruktur (*Self-healing*).
  4. *Performance Efficiency:* Pemilihan tipe resource yang optimal sesuai beban.
  5. *Cost Optimization:* Penghapusan resource menganggur dan pemilihan kapasitas reserved/spot.
  6. *Sustainability:* Mengurangi dampak lingkungan dari beban komputasi.

## Pemeriksaan pemahaman

1. Apa perbedaan tanggung jawab keamanan pengguna vs penyedia pada AWS Shared Responsibility Model?
2. Mengapa subnet database (RDS) sebaiknya tidak diberikan akses ke Internet Gateway maupun NAT Gateway?
3. Apa risiko utama menggunakan Access Key akun Root untuk keperluan operasional harian?
4. Bagaimana IAM Role pada instance EC2 lebih aman daripada menyimpan Access Key di file `.env`?

**Kunci ringkas:** AWS bertanggung jawab atas keamanan fisik infrastruktur dan software dasar cloud sedangkan pengguna bertanggung jawab atas keamanan OS tamu, konfigurasi firewall, enkripsi data, dan hak akses IAM; Subnet database diisolasi penuh agar tidak dapat diakses atau disusupi dari jaringan internet luar secara langsung; Akun root memiliki wewenang tak terbatas tanpa batasan izin sehingga jika bocor seluruh infrastruktur dapat dihapus atau disalahgunakan tanpa batas; IAM Role menyediakan token kredensial sementara (*temporary credentials*) yang otomatis dirotasi oleh AWS tanpa perlu menyimpan kunci rahasia statis di server.

## Referensi

- [AWS Documentation Home](https://docs.aws.com/)
- [AWS Well-Architected Framework Whitepaper](https://aws.amazon.com/architecture/well-architected/)
- [AWS Security Best Practices in IAM](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html)

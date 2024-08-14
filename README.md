# Loan Management API
API ini digunakan untuk mengelola data pelanggan, pinjaman, jadwal pembayaran, dan pembayaran dalam sistem manajemen pinjaman. Setiap endpoint dilengkapi dengan contoh permintaan dan respons sukses. Perhatikan bahwa respons dapat bervariasi jika terjadi kesalahan.

## Cara Menjalankan Aplikasi di Laptop/PC
Untuk menjalankan aplikasi ini secara lokal, ikuti langkah-langkah berikut:
1. Clone repository:

```sh
git clone git@github.com:jacky-htg/loans.git
```

2. Eksekusi query yang ada di file ./migrations/01_ddl_loans.sql untuk membuat tabel yang diperlukan di database.

3. Salin file .env.example menjadi .env dan isi data-data konfigurasi yang diperlukan.

4. Install dependencies dengan perintah:

```sh
go mod tidy
```

5. Jalankan aplikasi dengan perintah:

```sh
go run main.go
```

## Daftar API
### 1. Create Customer
- Endpoint: POST /customers
- Deskripsi: Menambahkan pelanggan baru ke sistem.
- Request:
```bash
curl --location 'localhost:8080/customers' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-secret-token' \
--data-raw '{
    "name": "Jacky Chan",
    "email": "jacky.chan@gmail.com",
    "phone": "08123456789"
}'
```
- Response:
```json
{
    "customer_id": 1,
    "message": "Customer created successfully"
}
```

### 2. Get Customer by ID
- Endpoint: GET /customers/:customerId
- Deskripsi: Mengambil informasi pelanggan berdasarkan ID.
- Request:
```bash
curl --location 'localhost:8080/customers/1' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "customer_id": 1,
    "name": "Jacky Chan",
    "email": "jacky.chan@gmail.com",
    "phone": "08123456789"
}
```

### 3. Create Loan
- Endpoint: POST /loans
- Deskripsi: Membuat pinjaman baru untuk pelanggan.
- Request:
```bash
curl --location 'localhost:8080/loans' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-secret-token' \
--data '{
    "customer_id": 1,
    "principal_amount": 5000000,
    "term_in_weeks": 50,
    "interest_rate": 10
}'
```
- Response:
```json
{
    "loan_id": 6,
    "message": "Loan created successfully"
}
```

### 4. Get Loan by ID
- Endpoint: GET /loans/:loanId
- Deskripsi: Mengambil informasi pinjaman berdasarkan ID.
- Request:
```bash
curl --location 'localhost:8080/loans/6' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "loan_id": 6,
    "customer_id": 1,
    "principal_amount": 5000000,
    "term_in_weeks": 50,
    "interest_rate": 10
}
```

### 5. Get Payment Schedule
- Endpoint: GET /loans/:loanId/schedule
- Deskripsi: Mengambil jadwal pembayaran untuk pinjaman tertentu.
- Request:
```bash
curl --location 'localhost:8080/loans/6/schedule' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
[
    {"installment_number":1,"due_date":"2024-08-20","amount":110000,"status":"Unpaid"},
    {"installment_number":2,"due_date":"2024-08-27","amount":110000,"status":"Unpaid"},
    {"installment_number":3,"due_date":"2024-09-03","amount":110000,"status":"Unpaid"},
    ...
    {"installment_number":50,"due_date":"2025-07-29","amount":110000,"status":"Unpaid"}
]
```

### 6. Make Payment
- Endpoint: POST /loans/:loanId/payment
- Deskripsi: Melakukan pembayaran cicilan pinjaman.
- Request:
```bash
curl --location 'localhost:8080/loans/6/payment' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your-secret-token' \
--data '{
    "installment_number": 2,
    "date": "2024-08-13",
    "amount": 110000 
}'
```
- Response:
```json
{
    "payment_id": 2,
    "message": "Payment successfully recorded"
}
```

### 7. Get Outstanding Amount
- Endpoint: GET /loans/:loanId/outstanding
- Deskripsi: Mengambil jumlah pinjaman yang belum dibayar.
- Request:
```bash
curl --location 'localhost:8080/loans/6/outstanding' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "loan_id": 6,
    "outstanding_amount": 4800000
}
```

### 8. Check Delinquent Status
- Endpoint: GET /loans/:loanId/delinquent
- Deskripsi: Memeriksa apakah peminjam termasuk delinquent.
- Request:
```bash
curl --location 'localhost:8080/loans/6/delinquent' \
--header 'Authorization: Bearer your-secret-token'
```
- Response:
```json
{
    "is_delinquent": false,
    "loan_id": "6"
}
```

### Catatan
Respons yang ditampilkan di atas adalah respons sukses. Dalam hal terjadi kesalahan, respons akan bervariasi tergantung pada jenis kesalahan yang terjadi, misalnya kesalahan validasi, kesalahan otentikasi, atau kesalahan server.
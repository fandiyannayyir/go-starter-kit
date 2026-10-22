# go-starter-kit

Template layanan Go untuk produk Fandi. Organisasi kode ala jajanin
(`domain → service → repository → handler`, route sebagai data) di atas infra
ala remittance-core (`dig`, `TxManager`, zap, Echo).

## Stack

- **HTTP:** Echo v4 + error handler terpusat (`AppError` → envelope JSON)
- **DI:** `dig` (`internal/di`)
- **DB:** PostgreSQL + GORM (query GORM biasa di balik interface repository)
- **Cache:** Redis
- **Migrasi:** goose (`migrations/`, dijalankan lewat perintah `migrate`)
- **Config:** viper + `.env` (lihat `.env.example`)
- **CLI:** cobra (`server`, `migrate`)

## Menjalankan (lokal)

```bash
cp .env.example .env            # sesuaikan bila perlu
docker compose up -d            # postgres + redis
go run . migrate up             # buat tabel (jalankan dari root project)
go run . server                 # start HTTP server di APP_HOST:APP_PORT
```

## Endpoint

| Method | Path              | Keterangan                         |
|--------|-------------------|------------------------------------|
| GET    | `/health`         | liveness (proses hidup)            |
| GET    | `/ready`          | readiness (ping DB + Redis)        |
| POST   | `/v1/examples`    | buat example (`{name, description}`)|
| GET    | `/v1/examples`    | list (query `page`, `per_page`)    |
| GET    | `/v1/examples/:id`| ambil satu                         |
| PUT    | `/v1/examples/:id`| ubah                               |
| DELETE | `/v1/examples/:id`| hapus                              |

Contoh:

```bash
curl -s localhost:8080/ready
curl -s -X POST localhost:8080/v1/examples \
  -H 'content-type: application/json' \
  -d '{"name":"produk A","description":"contoh"}'
```

## Struktur

```
cmd/                 entrypoint cobra (server, migrate)
configs/             config (viper) + koneksi DB
internal/
  di/                wiring dig
  shared/            dependency, TxManager, logger (zap), cache
  middleware/        RequestID, logger (CORS/Recover pakai bawaan Echo)
  server/            bootstrap Echo + error handler + graceful shutdown
  domain/<ctx>/      entity murni + domain error (tanpa dependensi DB)
  service/           interface repository + DTO + logika bisnis
  repository/        model GORM + implementasi repository
  handler/           handler HTTP + tabel route
migrations/          migrasi goose
pkg/                 apperror, response (Envelope[T]), route (registry)
```

## Membuat domain baru

Salin `internal/domain/example`, `service/example_*`, `repository/example_*`,
`handler/example_handler.go` sebagai pola, ganti nama, lalu daftarkan route-nya
di `cmd/server.go`. Tambah migrasi baru di `migrations/` (`NNNNN_nama.sql`).

## Kloning untuk produk

Ganti nama module `go-starter-kit` di `go.mod` dan seluruh import path dengan
nama produk.

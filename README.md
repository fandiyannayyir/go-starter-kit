# go-starter-kit

Starter kit layanan HTTP berbahasa Go dengan arsitektur berlapis
(`domain → service → repository → handler`, route sebagai data) di atas
Echo + dig + GORM — siap di-clone sebagai fondasi layanan baru.

## Stack

- **HTTP:** Echo v4 + error handler terpusat (`AppError` → envelope JSON)
- **DI:** `dig` (`internal/di`)
- **DB:** PostgreSQL + GORM (query GORM biasa di balik interface repository)
- **Cache:** Redis
- **Migrasi:** goose, **ter-embed ke binary** (`migrations/embed.go`) — perintah
  `migrate` jalan dari direktori mana pun, deploy cukup satu binary
- **Config:** viper + `.env` (lihat `.env.example`)
- **CLI:** cobra (`server`, `migrate`)
- **Lint:** golangci-lint (`.golangci.yml`) — `golangci-lint run`
- **Deploy:** `Dockerfile` multi-stage (binary tunggal + tzdata)

## Menjalankan (lokal)

```bash
cp .env.example .env            # sesuaikan bila perlu
docker compose up -d            # postgres + redis
go run . migrate up             # buat tabel (migrasi ter-embed di binary)
go run . server                 # start HTTP server di APP_HOST:APP_PORT
```

## Testing

```bash
go test ./...          # semua test; integration test memutar Postgres via Docker (testcontainers)
go test -short ./...   # lewati integration test (tanpa Docker)
```

Tiga lapis: unit service (repo palsu, tanpa DB), handler lewat HTTP in-process
(httptest — menguji middleware, validator, router, error handler), dan
integration repository (Postgres asli + migrasi ter-embed).

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
  di/                wiring dig (agregasi Register semua layer)
  shared/            dependency, TxManager, logger (zap), cache, validate
  middleware/        RequestID, logger (CORS/Recover pakai bawaan Echo)
  server/            bootstrap Echo + error handler + graceful shutdown
  router/            satu file routing modular (satu fungsi per domain)
  domain/<ctx>/      entity murni + domain error (tanpa dependensi DB)
  service/<ctx>/     subfolder per domain: interface repo + DTO + logika
  repository/        model GORM + implementasi repository
  handler/           handler HTTP murni (route-nya di router/)
migrations/          migrasi goose (ter-embed ke binary)
pkg/                 apperror, response (Envelope[T]), route (registry)
```

## Membuat domain baru

Salin `internal/domain/example` dan `internal/service/example/` sebagai pola,
lalu buat `repository/<ctx>_repository.go` + `handler/<ctx>_handler.go`.
Terakhir, daftarkan: `Provide` konstruktor di `service/di.go`,
`repository/di.go`, `handler/di.go`; satu fungsi route di
`internal/router/router.go`; dan migrasi baru di `migrations/`
(`NNNNN_nama.sql` — otomatis ikut ter-embed).

## Memakai sebagai template

Ganti nama module `go-starter-kit` di `go.mod` dan seluruh import path dengan
nama project Anda.

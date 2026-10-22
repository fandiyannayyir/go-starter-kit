# go-starter-kit — Rencana Rebuild (LOCKED 2026-10-22)

Template Go **umum untuk semua produk Fandi**. Arah: **organisasi kode ala
jajanin (`donation-core`) + infra ala remittance-core**, DI pakai `dig`.

## Prinsip

- **Organisasi domain = jajanin:** `domain` (model murni) → `repository`
  (GORM impl; interface repository-nya didefinisikan di dalam `service`, idiom
  Go — tanpa package `ports` terpisah) → `service` → `handler` + route-as-data.
- **Infra = remittance-core:** `dig` wiring, `shared.Dependency`,
  `TxManager` (ctx-based), zap + OTel, docker-compose. (gorm-gen TIDAK dipakai —
  repository pakai GORM biasa; type-safety codegen tak sepadan untuk template.)
- **Response + error = jajanin:** `pkg/response.Envelope[T]` + `pkg/apperror`
  + Echo error handler terpusat.
- **Lean:** satu slice contoh end-to-end, bukan belasan domain. Tunda gRPC,
  property-test berat, worker, email, outbond (stub bila perlu).

## Sudah ada (dipertahankan)

- `internal/di/*` — dig container (`GetContainer`), `RegisterDependency`
  (config, logger, database, redis, TxManager).
- `internal/shared/dependency.go` — `Dependency{dig.In}` + `Close()`.
- `internal/shared/transaction.go` — `TxManager.WithTransaction` + `GetDB(ctx)`.
- `internal/shared/log/*`, `internal/shared/cache/*` (zap + redis).
- `configs/config.go` (viper + bindEnvs reflektif), `configs/db/postgres.go`.

## Perlu dibereskan / dibangun

1. **Cleanup:** `.gitignore` + hapus sampah build ter-track (`.gocache/`,
   `.gomodcache/`); konsolidasi config ganda (`configs/` vs
   `internal/config/` — pakai `configs/`, buang `internal/config`).
2. **`pkg/` fondasi:** `apperror` (AppError + New/Wrap + codes), `response`
   (Envelope[T] + OK/Created/Error/ValidationError/Paginated), `route`
   (Registry route-as-data + guard duplicate).
3. **HTTP infra:** `internal/middleware` (RequestID, CORS, Recover, OTel,
   logger), Echo error handler terpusat (map AppError→envelope), server
   bootstrap + graceful shutdown.
4. **Organisasi domain (slice contoh):**
   - `healthcheck` — liveness/readiness (ping DB+cache).
   - `example` CRUD end-to-end: `domain/example` (entity+constants+errors) →
     `service/example_repository.go` (interface repository, di sisi konsumen) →
     `repository/example_{model,repository}.go` (GORM biasa memenuhi interface) →
     `service/example_{service,dto}.go` → `handler/example_handler.go` +
     `handler/routes_*.go` → wiring dig.
5. **cmd cobra:** `main.go` → `root` → `server` (dig+echo+graceful),
   `migrate` (goose), `seed` (stub).
6. **Infra files:** `migrations/` (goose, digenerate dari model),
   `docker-compose.yml` (postgres+redis), `.env.example`, `Makefile` (opsional).

## Keputusan rekonsiliasi (ringkas)

| Aspek | Pilihan | Sumber |
|---|---|---|
| DI | `dig` | remittance |
| Organisasi | domain/repository/service/handler (interface repo di `service`) | jajanin |
| Transaksi | `TxManager` ctx-based | remittance |
| Logger | zap + OTel | remittance |
| Query | GORM biasa di balik interface repository (didefinisikan di `service`) | jajanin |
| Migrasi | goose, digenerate dari model GORM | jajanin |
| cmd | cobra | jajanin |
| Response + error | `pkg/response.Envelope[T]` + `pkg/apperror` | jajanin |
| Config | viper + bindEnvs reflektif (sudah ada) | keduanya |

## TIDAK dibawa

Domain remittance (master/pricing/transaction, backbone/HMAC, `RM_*`,
`remittance-contracts`), domain jajanin (donation/creator/dll), gRPC,
property-test berat, worker/email/outbond penuh (stub saja bila perlu).

## Urutan build & verifikasi

Tiap tahap tutup dengan `go build ./... && go vet ./...`.

1. Cleanup (gitignore, hapus config ganda) → build.
2. `pkg/apperror`, `pkg/response`, `pkg/route` → build (independen, aman).
3. Middleware + Echo error handler + server bootstrap → build.
4. Slice `healthcheck` end-to-end + wiring dig + route → `go run` readiness OK.
5. Slice `example` CRUD end-to-end + migrasi → `go run` + smoke curl.
6. `docker-compose`, `.env.example`, README, Makefile → setup lokal jalan.

## Catatan

- Module name tetap `go-starter-kit`; produk yang meng-clone mengganti path
  module.

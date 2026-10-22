package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	domain "go-starter-kit/internal/domain/example"
	"go-starter-kit/internal/repository"
	"go-starter-kit/migrations"
)

// setupDB starts a throwaway Postgres (testcontainers), runs the embedded goose
// migrations against it and returns a connected *gorm.DB. Requires Docker; use
// `go test -short ./...` to skip.
func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("short mode: integration test needs Docker")
	}

	ctx := context.Background()
	pgc, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pgc.Terminate(context.Background()) })

	dsn, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("sql db: %v", err)
	}
	goose.SetBaseFS(migrations.Embed)
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("goose dialect: %v", err)
	}
	if err := goose.Up(sqlDB, "."); err != nil {
		t.Fatalf("goose up: %v", err)
	}
	return db
}

// TestExampleRepository_CRUD exercises the real SQL/GORM mapping that the
// in-memory fake in the service tests cannot prove.
func TestExampleRepository_CRUD(t *testing.T) {
	repo := repository.NewExampleRepository(setupDB(t))
	ctx := context.Background()

	e := &domain.Example{ID: "it-1", Name: "produk A", Description: "awal"}
	if err := repo.Create(ctx, e); err != nil {
		t.Fatalf("create: %v", err)
	}
	if e.CreatedAt.IsZero() || e.UpdatedAt.IsZero() {
		t.Fatal("create should backfill timestamps")
	}

	got, err := repo.FindByID(ctx, "it-1")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "produk A" || got.Description != "awal" {
		t.Fatalf("find mismatch: %+v", got)
	}

	got.Name = "produk B"
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	items, total, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Name != "produk B" {
		t.Fatalf("list mismatch: total=%d items=%+v", total, items)
	}

	if err := repo.Delete(ctx, "it-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, "it-1"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound after delete, got %v", err)
	}
	if err := repo.Delete(ctx, "it-1"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound on double delete, got %v", err)
	}
}

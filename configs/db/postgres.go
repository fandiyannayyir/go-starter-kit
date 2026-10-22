package db

import (
	"fmt"
	"time"

	"go-starter-kit/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewPostgreDB(conf *configs.Config) *Database {
	tz := conf.DBConfig.Timezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	sslmode := conf.DBConfig.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		conf.DBConfig.Host,
		conf.DBConfig.User,
		conf.DBConfig.Pass,
		conf.DBConfig.Name,
		conf.DBConfig.Port,
		sslmode,
		tz,
	)

	gormConfig := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	if conf.AppConfig.Debug {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	d, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL: %w", err))
	}

	dbSQL, err := d.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get underlying sql.DB: %w", err))
	}

	if conf.DBConfig.MaxIdleConnection > 0 {
		dbSQL.SetMaxIdleConns(conf.DBConfig.MaxIdleConnection)
	} else {
		dbSQL.SetMaxIdleConns(10)
	}
	if conf.DBConfig.MaxOpenConnection > 0 {
		dbSQL.SetMaxOpenConns(conf.DBConfig.MaxOpenConnection)
	} else {
		dbSQL.SetMaxOpenConns(100)
	}
	if conf.DBConfig.MaxConnLifeTime > 0 {
		dbSQL.SetConnMaxLifetime(conf.DBConfig.MaxConnLifeTime)
	} else {
		dbSQL.SetConnMaxLifetime(30 * time.Minute)
	}

	if err := dbSQL.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping PostgreSQL: %w", err))
	}

	return &Database{DB: d}
}

func (d *Database) PingDb() error {
	dbSQL, err := d.DB.DB()
	if err != nil {
		return err
	}
	return dbSQL.Ping()
}

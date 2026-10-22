package configs

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Of *Config = &Config{}

type Config struct {
	AppConfig App   `mapstructure:"app"`
	DBConfig  DB    `mapstructure:"db"`
	Redis     Redis `mapstructure:"redis"`
	JWT       JWT   `mapstructure:"jwt"`
	OTel      OTel  `mapstructure:"otel"`
}

type App struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Version  string `mapstructure:"version"`
	Env      string `mapstructure:"env"`
	TimeZone string `mapstructure:"time_zone"`
	Debug    bool   `mapstructure:"debug"`
}

type DB struct {
	Host              string        `mapstructure:"host"`
	Port              string        `mapstructure:"port"`
	User              string        `mapstructure:"user"`
	Pass              string        `mapstructure:"pass"`
	Name              string        `mapstructure:"name"`
	SSLMode           string        `mapstructure:"sslmode"`
	Timezone          string        `mapstructure:"timezone"`
	MaxIdleConnection int           `mapstructure:"max_idle_connection"`
	MaxOpenConnection int           `mapstructure:"max_open_connection"`
	MaxConnLifeTime   time.Duration `mapstructure:"max_conn_life_time"`
}

type Redis struct {
	Host  string `mapstructure:"host"`
	Port  string `mapstructure:"port"`
	User  string `mapstructure:"user"`
	Pass  string `mapstructure:"pass"`
	Index int    `mapstructure:"index"`
}

type JWT struct {
	Secret     string `mapstructure:"secret"`
	AccessTTL  string `mapstructure:"access_ttl"`
	RefreshTTL string `mapstructure:"refresh_ttl"`
}

type OTel struct {
	CollectorAddr string `mapstructure:"collector_addr"`
	ServiceName   string `mapstructure:"service_name"`
}

var (
	cfg  *Config
	once = sync.Once{}
)

func GetConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load(".env")

		v := viper.New()
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		cfg = &Config{}
		bindEnvs(v, cfg, "")

		if err := v.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
			dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
			)
		}); err != nil {
			panic(fmt.Errorf("failed to unmarshal config: %w", err))
		}

		Of = cfg

		tz := cfg.AppConfig.TimeZone
		if tz == "" {
			tz = "Asia/Jakarta"
		}
		loc, err := time.LoadLocation(tz)
		if err != nil {
			loc, _ = time.LoadLocation("Asia/Jakarta")
		}
		time.Local = loc
	})

	return cfg
}

func bindEnvs(v *viper.Viper, iface interface{}, prefix string) {
	ifaceVal := reflect.ValueOf(iface)
	if ifaceVal.Kind() == reflect.Pointer {
		ifaceVal = ifaceVal.Elem()
	}
	ifaceType := ifaceVal.Type()

	for i := 0; i < ifaceVal.NumField(); i++ {
		field := ifaceType.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" || tag == "-" {
			continue
		}

		envKey := tag
		if prefix != "" {
			envKey = prefix + "." + tag
		}

		if field.Type.Kind() == reflect.Struct {
			bindEnvs(v, ifaceVal.Field(i).Interface(), envKey)
		} else {
			_ = v.BindEnv(envKey)
		}
	}
}

func IsProduction() bool {
	return strings.EqualFold(Of.AppConfig.Env, "prod") ||
		strings.EqualFold(Of.AppConfig.Env, "production")
}

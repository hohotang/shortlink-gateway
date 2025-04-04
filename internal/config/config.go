package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port           int           `mapstructure:"port"`
	Env            string        `mapstructure:"env"`
	ServiceName    string        `mapstructure:"service_name"`
	OTLPEndpoint   string        `mapstructure:"otel_exporter_otlp_endpoint"`
	UseGrpc        bool          `mapstructure:"use_grpc"`
	GrpcServerAddr string        `mapstructure:"grpc_server_addr"`
	GrpcTimeout    time.Duration `mapstructure:"grpc_timeout"` // Timeout in seconds for gRPC calls
}

func Load() *Config {
	v := viper.New()

	v.SetDefault("grpc_timeout", 5*time.Second)

	// 設定 config 檔案來源
	v.SetConfigName("config") // config.yaml
	v.SetConfigType("yaml")
	v.AddConfigPath(".")        // 專案根目錄
	v.AddConfigPath("./config") // 也支援 config/ 資料夾

	// 支援讀取環境變數：PORT、ENV 等（自動覆蓋）
	v.SetEnvPrefix("SHORTLINK") // 例如 SHORTLINK_PORT
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 讀取 config.yaml（如果有）
	if err := v.ReadInConfig(); err != nil {
		log.Printf("⚠️ config.yaml not found: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	return &cfg
}

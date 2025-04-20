package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Port            int           `mapstructure:"port"`
	Env             string        `mapstructure:"env"`
	ServiceName     string        `mapstructure:"service_name"`
	OTLPEndpoint    string        `mapstructure:"otel_exporter_otlp_endpoint"`
	TracesEndpoint  string        `mapstructure:"traces_endpoint"`
	MetricsEndpoint string        `mapstructure:"metrics_endpoint"`
	UseGrpc         bool          `mapstructure:"use_grpc"`
	GrpcServerAddr  string        `mapstructure:"grpc_server_addr"`
	GrpcTimeout     time.Duration `mapstructure:"grpc_timeout"`
}

// Load loads configuration from config.yaml and environment variables
func Load() *Config {
	v := viper.New()

	// Set default values
	v.SetDefault("port", 8080)
	v.SetDefault("env", "development")
	v.SetDefault("service_name", "api-gateway")
	v.SetDefault("otel_exporter_otlp_endpoint", "localhost:4318")
	v.SetDefault("traces_endpoint", "localhost:4318")
	v.SetDefault("metrics_endpoint", "localhost:9090")
	v.SetDefault("use_grpc", true)
	v.SetDefault("grpc_server_addr", "localhost:50051")
	v.SetDefault("grpc_timeout", 5*time.Second)

	// Set configuration file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// Automatically replace dots with underscores for environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load config from environment variables with prefix SHORTLINK_
	v.SetEnvPrefix("SHORTLINK")
	v.AutomaticEnv()

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return &cfg
}

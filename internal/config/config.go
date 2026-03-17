package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	Log      LogConfig      `yaml:"log"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type KafkaConfig struct {
	BrokerAddr string `yaml:"broker_addr"`
	Topic      string `yaml:"topic"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()

	cfg := defaults()
	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode yaml: %w", err)
	}

	if v := os.Getenv("GRPC_PORT"); v != "" {
		if _, err = fmt.Sscanf(v, "%d", &cfg.GRPC.Port); err != nil {
			return nil, fmt.Errorf("config: GRPC_PORT: %w", err)
		}
	}
	if v := os.Getenv("POSTGRES_DSN"); v != "" {
		cfg.Postgres.DSN = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("KAFKA_BROKER_ADDR"); v != "" {
		cfg.Kafka.BrokerAddr = v
	}
	if v := os.Getenv("KAFKA_TOPIC"); v != "" {
		cfg.Kafka.Topic = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
	return cfg, nil
}

func defaults() *Config {
	return &Config{
		GRPC:     GRPCConfig{Port: 50051},
		Postgres: PostgresConfig{DSN: "postgres://tms:tms@localhost:5432/tms?sslmode=disable"},
		Redis:    RedisConfig{Addr: "localhost:6379"},
		Kafka:    KafkaConfig{BrokerAddr: "localhost:9092", Topic: "shipment-events"},
		Log:      LogConfig{Level: "info"},
	}
}

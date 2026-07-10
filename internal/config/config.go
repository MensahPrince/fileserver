package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	GRPCPort    string
	HTTPPort    string
	DBDSN       string
	StorageRoot string
	StorageType string // "disk" or "s3", useful later
}

func Load() (*Config, error) {
	cfg := &Config{
		GRPCPort:    os.Getenv("GRPC_PORT"),
		HTTPPort:    os.Getenv("HTTP_PORT"),
		DBDSN:       os.Getenv("DB_DSN"),
		StorageRoot: os.Getenv("STORAGE_ROOT"),
		StorageType: os.Getenv("STORAGE_TYPE"),
	}

	// GRPC port validation
	if cfg.GRPCPort == "" {
		cfg.GRPCPort = ":50051"
	} else {
		p := strings.TrimPrefix(cfg.GRPCPort, ":")
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 65535 {
			return nil, fmt.Errorf("invalid GRPC_PORT: %q", cfg.GRPCPort)
		}
	}

	// HTTP port validation
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = ":8800"
	} else {
		p := strings.TrimPrefix(cfg.HTTPPort, ":")
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 65535 {
			return nil, fmt.Errorf("invalid HTTP_PORT: %q", cfg.HTTPPort)
		}
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("no DB_DSN found")
	}

	if cfg.StorageRoot == "" {
		fmt.Println("Storage root not found. Defaulting to ./data")
		cfg.StorageRoot = "./data"
	}

	if cfg.StorageType != "disk" && cfg.StorageType != "s3" {
		return nil, fmt.Errorf("invalid STORAGE_TYPE: %q (must be disk or s3)", cfg.StorageType)
	}

	return cfg, nil
}
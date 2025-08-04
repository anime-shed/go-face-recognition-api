package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	Pigo   PigoConfig
	Limits LimitsConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// PigoConfig holds pigo face detection configuration
type PigoConfig struct {
	MinSize       int
	MaxSize       int
	ShiftFactor   float32
	ScaleFactor   float32
	IoUThreshold  float32
	MinConfidence float32
}

// LimitsConfig holds various limits for the application
type LimitsConfig struct {
	MaxImageSize int64
	MaxWidth     int
	MaxHeight    int
	RateLimit    int
	RateBurst    int
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", ":8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
		},
		Pigo: PigoConfig{
			MinSize:       getIntEnv("PIGO_MIN_SIZE", 25),
			MaxSize:       getIntEnv("PIGO_MAX_SIZE", 1000),
			ShiftFactor:   getFloat32Env("PIGO_SHIFT_FACTOR", 0.2),
			ScaleFactor:   getFloat32Env("PIGO_SCALE_FACTOR", 1.1),
			IoUThreshold:  getFloat32Env("PIGO_IOU_THRESHOLD", 0.6),
			MinConfidence: getFloat32Env("PIGO_MIN_CONFIDENCE", 12.0),
		},
		Limits: LimitsConfig{
			MaxImageSize: getInt64Env("MAX_IMAGE_SIZE", 5242880), // 5MB
			MaxWidth:     getIntEnv("MAX_WIDTH", 2000),
			MaxHeight:    getIntEnv("MAX_HEIGHT", 2000),
			RateLimit:    getIntEnv("RATE_LIMIT", 100),
			RateBurst:    getIntEnv("RATE_BURST", 10),
		},
	}
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getFloat32Env(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatValue)
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

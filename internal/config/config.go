package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL				string
	SerpAPIKey				string
	TrendsGeo					string // e.g., "US"
	ResultsPerTrend		int
	TargetDate				time.Time
}

func Load() Config {
	dateStr := os.Getenv("TARGET_DATE") // YYYY-MM-DD
	var targetDate time.Time

	if dateStr != "" {
		parsed, _ := time.Parse("2006-01-02", dateStr)
		targetDate = parsed
	}

	if targetDate.IsZero() {
		targetDate = time.Now()
	}

	return Config {
		DatabaseURL:				getEnv("DATABASE_URL", ""),
		SerpAPIKey:					getEnv("SERPAPI_KEY", ""),
		TrendsGeo:					getEnv("TRENDS_GEO", ""),
		ResultsPerTrend: 		getEnvInt("RESULTS_PER_TREND", 100),
		TargetDate:					targetDate,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return fallback
}
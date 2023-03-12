package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload" // Automatically load ".env" file in root
)

type AppConfig struct {
	LivenessPort                int
	Port                        int
	AppName                     string
	GracefulShutdownSeconds     int
	EnableHTTP2                 bool
	IsDevEnv                    bool
	GoogleTranslateV2APIKey     string
	GoogleTranslateV3ProjectKey string
}

func ApplicationConfig() AppConfig {
	livenessPort := envVarAtoi("LIVENESS_PROBE_PORT")
	port := envVarAtoi("PORT")
	appName := envVarAsStr("APP_NAME")
	gracefulShutdownSeconds := envVarAtoi("GRACEFUL_SHUTDOWN_SECONDS")
	enableHTTP2 := envVarAsBool("ENABLE_HTTP2")
	isDevEnv := envVarAsBool("DEVELOPMENT_MODE")
	googleTranslateV2APIKey := envVarAsStr("GOOGLE_TRANSLATE_V2_API_KEY")
	googleTranslateV3ProjectKey := envVarAsStr("GOOGLE_TRANSLATE_V3_PROJECT_KEY")

	return AppConfig{
		LivenessPort:                livenessPort,
		Port:                        port,
		AppName:                     appName,
		GracefulShutdownSeconds:     gracefulShutdownSeconds,
		EnableHTTP2:                 enableHTTP2,
		IsDevEnv:                    isDevEnv,
		GoogleTranslateV2APIKey:     googleTranslateV2APIKey,
		GoogleTranslateV3ProjectKey: googleTranslateV3ProjectKey,
	}
}

func envVarAtoi(envName string) int {
	valueStr := os.Getenv(envName)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(err)
	}

	return value
}

func envVarAsBool(envName string) bool {
	valueStr := os.Getenv(envName)
	return valueStr == "true"
}

func envVarAsStr(envName string) string {
	valueStr := os.Getenv(envName)
	return valueStr
}

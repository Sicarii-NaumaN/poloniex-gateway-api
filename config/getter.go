// Functions to get different types env params

package config

import (
	"os"
	"strconv"
	"time"

	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
)

const tempErr = "Config parameter %s is incorrect: %s"

func GetConfigString(key configKey) string {
	value := os.Getenv(string(key))
	if value == "" {
		logger.Fatalf(tempErr, key, value)
	}
	return value
}

func GetConfigBool(key configKey) bool {
	valueStr := os.Getenv(string(key))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, err)
	}
	return value
}

func GetConfigInt(key configKey) int {
	valueStr := os.Getenv(string(key))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, valueStr)
	}
	return value
}

func GetConfigDuration(key configKey) time.Duration {
	valueStr := os.Getenv(string(key))
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, valueStr)
	}
	return value
}

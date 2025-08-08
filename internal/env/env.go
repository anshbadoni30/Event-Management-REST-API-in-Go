package env

import (
	"os"
	"strconv"
)

func GetEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intvalue,err:=strconv.Atoi(value);err==nil{
			return intvalue
		}
	}
	return defaultValue
}


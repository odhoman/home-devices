package utils

import (
	"errors"
	"os"
)

func ResolveDefaultDate(date, defaultValue int64) int64 {

	if date == 0 {
		return defaultValue
	}
	return date
}

func GetArrayNodeFromMap(valueMap map[string]interface{}, field string) []string {
	return valueMap[field].([]string)
}

func GetValueProperty(propertyKey string) (string, error) {

	value := os.Getenv(propertyKey)

	if value == "" {
		return "", errors.New("No value found for property" + propertyKey)
	}

	return value, nil
}

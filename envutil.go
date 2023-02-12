package main

import (
	"errors"
	"os"
	"strconv"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func GetEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func GetEnvInt(key string) (int, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetEnvBool(key string) (bool, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func SafeGetEnvStr(key string, defaultValue string) string {
	value, err := GetEnvStr(key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

func SafeGetEnvInt(key string, defaultValue int) int {
	value, err := GetEnvInt(key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

func SafeGetEnvBool(key string, defaultValue bool) bool {
	value, err := GetEnvBool(key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

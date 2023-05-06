package common

import "os"

type EnvVars map[string]string

func GetEnvVars(keys []string) EnvVars {
	result := make(EnvVars)
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			result[key] = value
		}
	}
	return result
}

func MergeEnvVars(envVars ...EnvVars) EnvVars {
	result := make(EnvVars)

	for _, env := range envVars {
		for key, value := range env {
			if key != "" && value != "" {
				result[key] = value
			}
		}
	}

	return result
}

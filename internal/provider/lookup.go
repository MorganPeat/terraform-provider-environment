package provider

import (
	"fmt"
	"os"
)

func lookupEnvironmentVariable(name string) (string, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable '%s' not found", name)
	}

	return v, nil
}

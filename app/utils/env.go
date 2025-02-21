package utils

import (
	"log"
	"os"
)

// ResetEnv - сбрасывает чувствительные переменные окружения.
func ResetEnv(envs ...string) {
	for _, env := range envs {
		if err := os.Unsetenv(env); err != nil {
			log.Printf("[WARN] can't unset env %s, %s", env, err)
		}
	}
}

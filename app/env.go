package main

import (
	"log"
	"os"
)

// resetEnv() - сбрасывает чувствительные переменные окружения.
func resetEnv(envs ...string) {
	for _, env := range envs {
		if err := os.Unsetenv(env); err != nil {
			log.Printf("[WARN] can't unset env %s, %s", env, err)
		}
	}
}

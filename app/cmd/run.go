package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbzer0/go-rest-template/app/utils"
)

// RunCommand реализует запуск сервера как субкоманду.
// Встраиваем Configuration, чтобы флаги были доступны только при вызове команды run.
type RunCommand struct {
	Configuration
	// Version передаётся извне и не парсится из флагов.
	Version string `no-flag:"true"`
}

// Execute запускает сервер с учетом сброса чувствительных переменных,
// инициализации логирования и graceful shutdown.
func (c *RunCommand) Execute(args []string) error {
	// Создаем базовый контекст.
	ctx := context.Background()

	// Сбрасываем чувствительные переменные (например, DSURL и DSDB).
	utils.ResetEnv(c.DSURL, c.DSDB)

	// Инициализируем логирование.
	if err := NewLogCommand(c.Dbg).Execute(ctx); err != nil {
		return err
	}

	// Создаем контекст с отменой для graceful shutdown.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Graceful shutdown: отслеживаем сигналы прерывания и завершения процесса.
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	// Запускаем сервер.
	return NewServerCommand(&c.Configuration, c.Version).Execute(ctx)
}

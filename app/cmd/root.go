package cmd

import (
	"os"

	"github.com/jessevdk/go-flags"
)

// RootCommand является корневой командой и содержит субкоманду run.
type RootCommand struct {
	Run *RunCommand `command:"run" description:"Запуск HTTP-сервера"`
}

// Execute парсит аргументы командной строки и запускает выбранную субкоманду.
func Execute(version string) {
	root := &RootCommand{
		Run: &RunCommand{
			Configuration: Configuration{}, // значения по умолчанию из тегов структуры
			Version:       version,
		},
	}
	parser := flags.NewParser(root, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

type LogCommand struct {
	debug bool
}

func NewLogCommand(debug bool) *LogCommand {
	return &LogCommand{debug: debug}
}

func (c *LogCommand) Execute(ctx context.Context) error {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if c.debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = "DEBUG"
	}

	log.SetOutput(filter)
	return nil
}

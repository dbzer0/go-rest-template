package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbzer0/go-rest-template/app/utils"
	"github.com/jessevdk/go-flags"
)

type RootCommand struct {
	version string
	opts    *Configuration
}

func NewRootCommand(version string) *RootCommand {
	return &RootCommand{
		version: version,
		opts:    &Configuration{},
	}
}

func (c *RootCommand) Execute(ctx context.Context) error {
	if err := c.parseFlags(); err != nil {
		return err
	}

	// сбрасываем чувствительные переменные
	utils.ResetEnv(c.opts.DSURL, c.opts.DSDB)

	if err := NewLogCommand(c.opts.Dbg).Execute(ctx); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	return NewServerCommand(c.opts, c.version).Execute(ctx)
}

func (c *RootCommand) parseFlags() error {
	p := flags.NewParser(c.opts, flags.Default)
	if _, err := p.Parse(); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && errors.Is(flagsErr.Type, flags.ErrHelp) {
			os.Exit(0)
		}
		return err
	}
	return nil
}

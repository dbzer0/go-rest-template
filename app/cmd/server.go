// app/commands/server.go
package cmd

import (
	"context"
	"log"

	"github.com/dbzer0/go-rest-template/app/database"
	"github.com/dbzer0/go-rest-template/app/database/drivers"
	"github.com/dbzer0/go-rest-template/app/manager/test"
)

type ServerCommand struct {
	opts    *Configuration
	version string
}

func NewServerCommand(opts *Configuration, version string) *ServerCommand {
	return &ServerCommand{
		opts:    opts,
		version: version,
	}
}

func (c *ServerCommand) Execute(ctx context.Context) error {
	ds, err := c.setupDatastore(ctx)
	if err != nil {
		return err
	}
	defer ds.Close(context.Background())

	testManager := test.NewManager(ds)
	_ = testManager

	httpCmd := NewHTTPCommand(
		ctx,
		c.opts,
		// testManager,
		c.version,
	)
	return httpCmd.Execute(ctx)
}

func (c *ServerCommand) setupDatastore(ctx context.Context) (drivers.DataStore, error) {
	ds, err := database.Connect(drivers.DataStoreConfig{
		URL:           c.opts.DSURL,
		DataStoreName: c.opts.DSName,
		DataBaseName:  c.opts.DSDB,
	})
	if err != nil {
		return nil, err
	}

	if err = ds.Connect(); err != nil {
		return nil, err
	}

	log.Printf("[INFO] connected to %s", ds.Name())
	return ds, nil
}

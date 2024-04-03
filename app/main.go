package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbzer0/go-rest-template/app/database"
	"github.com/dbzer0/go-rest-template/app/database/drivers"
	"github.com/dbzer0/go-rest-template/app/director"

	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"
)

var version = "unknown"

type configuration struct {
	DSName string `short:"n" long:"ds" env:"DATASTORE" description:"DataStore name (format: mongo/null)" required:"false" default:"mongo"`
	DSDB   string `short:"d" long:"ds-db" env:"DATASTORE_DB" description:"DataStore database name (format: PROJECTNAME)" required:"false" default:"PROJECTNAME"`
	DSURL  string `short:"u" long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL (format: mongodb://localhost:27017)" required:"false" default:"mongodb://localhost:27017"`

	ListenAddr string `short:"l" long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":8080"`
	BasePath   string `long:"base-path" env:"BASE_PATH" description:"base path of the host" required:"false" default:"sso"`
	CertFile   string `short:"c" long:"cert" env:"CERT_FILE" description:"Location of the SSL/TLS cert file" required:"false" default:""`
	KeyFile    string `short:"k" long:"key" env:"KEY_FILE" description:"Location of the SSL/TLS key file" required:"false" default:""`

	Dbg       bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	IsTesting bool `long:"testing" env:"APP_TESTING" description:"testing mode"`
}

func main() {
	fmt.Printf("PROJECTNAME %s\n", version)

	var opts configuration

	// парсинг опций
	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.Parse(); err != nil {
		log.Println("[ERROR] Ошибка парсинга опций:", err)
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && errors.Is(flagsErr.Type, flags.ErrHelp) {
			os.Exit(0)
		}
	}

	setupLog(opts.Dbg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ловим сигнал для graceful termination
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Print("[WARN] interrupt signal")
		cancel()
	}()

	run(ctx, &opts)

	log.Printf("[INFO] process terminated")
}

// setupLog настраивает уровни логирования и вывод логгера в os.Stdout.
func setupLog(dbg bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if dbg {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = "DEBUG"
	}

	log.SetOutput(filter)
}

// run запускает основной цикл программы, стартующий все остальные приложения.
func run(ctx context.Context, opts *configuration) {
	ds, err := database.Connect(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})
	if err != nil {
		log.Printf("[ERROR] cannot connect to datastore %s: %v", opts.DSName, err)
		return
	}
	defer ds.Close(context.Background())

	if err = ds.Connect(); err != nil {
		log.Printf("[ERROR] cannot connect to database %s: %v", ds.Name(), err)
		return
	}
	log.Printf("[INFO] connected to %s", ds.Name())

	serverApp := director.NewServerApp(ctx)

	serverApp.Run()
}

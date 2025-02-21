package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dbzer0/go-rest-template/app/cmd"
)

var version = "unknown"

func main() {
	fmt.Printf("PROJECTNAME %s\n", version)

	rootCmd := cmd.NewRootCommand(version)
	if err := rootCmd.Execute(context.Background()); err != nil {
		log.Fatal(err)
	}
}

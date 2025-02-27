package main

import (
	"fmt"

	"github.com/dbzer0/go-rest-template/app/cmd"
)

var version = "unknown"

func main() {
	fmt.Printf("PROJECTNAME %s\n", version)

	cmd.Execute(version)
}

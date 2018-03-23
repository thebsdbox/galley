package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thebsdbox/galley/pkg/storage"
	"github.com/thebsdbox/galley/pkg/webserver"
)

func envIsEnabled(envVariable string) bool {
	if os.Getenv(envVariable) == "ENABLED" {
		return true
	}
	return false
}

func main() {

	// Command line flags
	webserverFlag := flag.Bool("webserver", envIsEnabled("galley_webserver"), "Enable the API and health check webserver")

	flag.Parse()

	// Check that the project name is the remaining argument, if not print out the errors
	remArgs := flag.Args()
	if len(remArgs) == 0 {
		fmt.Printf("USAGE: %s [options] <path to storage/file> \n\n", filepath.Base(os.Args[0]))
		flag.Usage()
		os.Exit(1)
	}

	// This should be the path to the storage that will be used
	storagePath := remArgs[0]

	if *webserverFlag == true {
		// This starts the API/Health webserver in its own GO Routine
		webserver.StartWebServer()
	}
	storage.DetectStorage(storagePath)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/thebsdbox/galley/pkg/storage"
	"github.com/thebsdbox/galley/pkg/webserver"

	log "github.com/Sirupsen/logrus"
)

func envIsEnabled(envVariable string) bool {
	if os.Getenv(envVariable) == "ENABLED" {
		return true
	}
	return false
}

func envIsNumber(envVariable string, lowNum int, highNum int, defaultNum int) int {
	if os.Getenv(envVariable) != "" {
		envLogLevel, err := strconv.Atoi(os.Getenv(envVariable))
		if err != nil {
			log.Errorf("%v", err)
		}
		if envLogLevel < lowNum && highNum > 5 {
			log.Fatalf("Log level must be between %d and %d", lowNum, highNum)
		}
		return envLogLevel
	}
	return defaultNum
}

func main() {

	// Command line flags
	log.Println("Starting Galley")
	webserverFlag := flag.Bool("webserver", envIsEnabled("galley_webserver"), "Enable the API and health check webserver")
	initFlag := flag.Bool("initStorage", false, "Initialise the StorageDevice")
	forceFlag := flag.Bool("force", false, "Force an operation, CAUTION can cause data-loss")

	logLevel := flag.Int("logging", envIsNumber("galley_logging", 0, 5, 2), "Set logging 0 = none, 5 = debug")

	flag.Parse()

	log.SetLevel(log.Level(*logLevel))

	// Check that the project name is the remaining argument, if not print out the errors
	remArgs := flag.Args()
	if len(remArgs) == 0 {
		fmt.Printf("USAGE: %s [options] <path to storage/file> \n\n", filepath.Base(os.Args[0]))
		flag.Usage()
		log.Fatalln("Failed to start Galley")
	}

	// This should be the path to the storage that will be used
	storagePath := remArgs[0]

	if *initFlag == true {
		header, err := storage.ReadHeader(storagePath)
		if err != nil {
			log.Errorf("%v", err)
		}
		if storage.HeaderMatches(header) == true && *forceFlag != true {
			log.Fatalln("Galley formatted disk, use --force to wipe storage")
		}
		storage.InitialiseDisk(storagePath)
		os.Exit(0)
	}

	if *webserverFlag == true {
		// This starts the API/Health webserver in its own GO Routine
		webserver.StartWebServer()
	}
	err := storage.DetectStorage(storagePath)
	if err != nil {
		log.Errorf("%v", err)
	}
}

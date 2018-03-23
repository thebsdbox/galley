package storage

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

// DetectStorage - This will identity the storage type from the path
func DetectStorage(path string) error {

	// Open the path and determine backing type
	pathInfo, err := os.Stat(path)

	if err != nil {
		return err
	}

	// Error use-cases
	switch pathInfo.Mode() {
	case os.ModeSymlink:
		return fmt.Errorf("Symlinks are not supported backing devices")
	case os.ModeNamedPipe:
		return fmt.Errorf("UNIX Pipes are not supported")
	case os.ModeDir:
		return fmt.Errorf("%s is a Directory and not supported", path)
	case os.ModeSocket:
		return fmt.Errorf("UNIX Sockets are not supported")
	}

	// Check for devices
	if pathInfo.Mode() == os.ModeDevice {

	}

	// Originally was the IsRegular method
	if pathInfo.Mode()&(os.ModeType|os.ModeCharDevice) == 0 {
		// Regular file
		log.Debugln("File Mode starting ")
		if pathInfo.Size() == 0 {
			return fmt.Errorf("File Size is zero")
		}
		header, err := ReadHeader(path)
		if err != nil {
			return err
		}
		log.Debugf("Disk Header: v%s.%s", string(header.Version[0]), string(header.Version[1]))

		if HeaderMatches(header) == true {
			log.Println("Disk header Matches")
		} else {
			log.Warnln("Incorrect Disk Header")
		}

	}
	return nil
}

//InitialiseDisk -
func InitialiseDisk(path string) error {

	err := WriteHeader(path, NewHeader())
	if err != nil {
		return err
	}
	return nil
}

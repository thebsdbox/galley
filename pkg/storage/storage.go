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
		log.Debugln("Device Mode Starting")
	}

	// Originally was the IsRegular method
	if pathInfo.Mode()&(os.ModeType|os.ModeCharDevice) == 0 {
		// Regular file
		log.Debugln("File Mode starting")
		if pathInfo.Size() == 0 {
			return fmt.Errorf("File Size is zero")
		}
		header, err := ReadHeader(path)
		if err != nil {
			return err
		}
		// Debug output about the Disk Image
		log.Debugf("Disk Header Version: v%s.%s", string(header.Version[0]), string(header.Version[1]))
		log.Debugf("Disk UUID: %s", header.UUID.String())
		log.Debugf("Header Size: %d", header.Size)
		log.Debugf("Disk Size: %d", header.DiskSize)
		log.Debugf("Block Count: %d", header.BlockCount)

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

// AddBlock - Add another storage block
func AddBlock(path string) error {
	err := WriteBlock(path)
	if err != nil {
		return err
	}
	return nil
}

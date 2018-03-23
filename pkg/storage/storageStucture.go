package storage

import (
	"encoding/binary"
	"fmt"
)

// The storage Header is present at the beginning of any storage being used by galley
type storageHeader struct {
	uuid            [64]byte // Unique Identifier for the volume
	previousAddress [16]byte // Last used address (IPv4)
	currentAddress  [16]byte // currently used address (IPv4)
}

// ReadHeader - this will attempt to read the storage header from the path passed
func ReadHeader(path string) {
	header := storageHeader{}
	fmt.Printf("Attempting to read %d bytes of header\n", binary.Size(header))
}

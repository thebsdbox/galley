package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

var header = [4]byte{0x70, 0x6f, 0x6f, 0x70}

// DiskHeader - The storage Header is present at the beginning of any storage being used by galley
type DiskHeader struct {
	Header          [4]byte  // Header to determine if the device has been initialized
	Version         [2]byte  // Version control of galley structure
	UUID            [64]byte // Unique Identifier for the volume
	PreviousAddress [16]byte // Last used address (IPv4)
	CurrentAddress  [16]byte // currently used address (IPv4)
}

// HeaderMatches - This is used to ensure that disk has been initialised
func HeaderMatches(readHeader *DiskHeader) bool {
	if readHeader.Header == header {
		return true
	}
	return false
}

// ReadHeader - this will attempt to read the storage header from the path passed
func ReadHeader(path string) (*DiskHeader, error) {
	header := DiskHeader{}
	headerSize := binary.Size(header)
	log.Debugf("Attempting to read %d bytes of header", headerSize)

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	binaryData := make([]byte, headerSize)
	bytesRead, err := file.Read(binaryData)
	if err != nil {
		return nil, err
	}

	log.Debugf("Read %d bytes from file %s", bytesRead, path)

	buffer := bytes.NewBuffer(binaryData)
	err = binary.Read(buffer, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

// WriteHeader - This writes a header to the storage defined at path
func WriteHeader(path string, header *DiskHeader) error {
	headerSize := binary.Size(header)
	fmt.Printf("%v\n", header)

	log.Debugf("Attempting to write %d bytes of header", headerSize)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	if err != nil {
		return err
	}

	var binaryBuffer bytes.Buffer
	binary.Write(&binaryBuffer, binary.BigEndian, header)

	bytesWritten, err := file.WriteAt(binaryBuffer.Bytes(), 0)
	if err != nil {
		return err
	}
	log.Debugf("Written %d bytes to file %s", bytesWritten, path)
	return nil
}

//NewHeader - creates a new default header
func NewHeader() *DiskHeader {
	nh := DiskHeader{}
	nh.Header = header
	nh.Version[0] = 0x30
	nh.Version[1] = 0x31
	return &nh
}

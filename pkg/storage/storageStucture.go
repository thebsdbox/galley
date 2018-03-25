package storage

import (
	"bytes"
	"encoding/binary"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

// Magic for header
var header = [4]byte{0x70, 0x6f, 0x6f, 0x70}

// Magic for block
var block = [2]byte{0x5c, 0x24}

var blockSize uint64

// DiskHeader - The storage Header is present at the beginning of any storage being used by galley
type DiskHeader struct {
	Magic           [4]byte   // Magic bytes to identify a Disk (Header)
	Size            uint8     // Size of the Header
	Version         [2]byte   // Version control of galley structure
	UUID            uuid.UUID // Unique Identifier for the volume
	PreviousAddress [4]byte   // Last used address (IPv4)
	CurrentAddress  [4]byte   // currently used address (IPv4)
	DiskSize        uint64    // Size of the Disk
	BlockCount      uint64    // Amount of blocks provisioned
}

// DiskBlock - Defines a block of Data
type DiskBlock struct {
	Magic         [2]byte   // Magic bytes to identify a block
	Size          uint8     // Size of Header
	BlockPosition uint64    // Position of the block
	AltUUID       uuid.UUID // Alternative Identical UUID block
	BlockSize     uint64    // Size of Block on Disk

}

func init() {
	// Default is set to 512KB
	blockSize = 512000
}

//NewHeader - creates a new default header populated with fixed length values
func NewHeader() *DiskHeader {
	nh := DiskHeader{}
	nh.Size = uint8(binary.Size(nh))
	nh.Magic = header      // Set the Header Identifier
	nh.Version[0] = 0x30   // 0
	nh.Version[1] = 0x31   // 1
	nh.UUID = uuid.NewV1() // Create a UUID for the disk image
	nh.BlockCount = 0      // Set block count to zero
	return &nh
}

//NewBlock - Creates a new block and header
func NewBlock() *DiskBlock {
	nb := DiskBlock{}
	nb.Size = uint8(binary.Size(nb))
	nb.Magic = block // Set the Block Identifier
	return &nb
}

// HeaderMatches - This is used to ensure that disk has been initialised
func HeaderMatches(readHeader *DiskHeader) bool {
	if readHeader.Magic == header {
		return true
	}
	return false
}

// ReadHeader - this will attempt to read the storage header from the path passed
func ReadHeader(path string) (*DiskHeader, error) {
	header := DiskHeader{}
	headerSize := binary.Size(header)

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

	log.Debugf("Read %d header bytes from file %s", bytesRead, path)

	buffer := bytes.NewBuffer(binaryData)
	err = binary.Read(buffer, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

// WriteHeader - This writes a header to the storage defined at path
func WriteHeader(path string, header *DiskHeader) error {
	log.Printf("Initialising Disk [%s]", path)
	headerSize := binary.Size(header)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	if err != nil {
		return err
	}

	// Retrieve size of storage
	pathInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Set size of disk within header
	header.DiskSize = uint64(pathInfo.Size())

	var binaryBuffer bytes.Buffer
	binary.Write(&binaryBuffer, binary.BigEndian, header)

	bytesWritten, err := file.WriteAt(binaryBuffer.Bytes(), 0)
	if err != nil {
		return err
	}
	log.Debugf("Written %d of %d header bytes to file %s", bytesWritten, headerSize, path)
	log.Printf("Initialised Disk ID: %s", header.UUID.String())
	return nil
}

// WriteBlock - Writes a new block to the underlying storage along with header
func WriteBlock(path string) error {
	h, err := ReadHeader(path)
	if err != nil {
		return err
	}

	// Calculate new block starting point Header + (block * blocksize)
	nextBlockAddr := uint64(h.Size) + (h.BlockCount * blockSize)
	if nextBlockAddr > h.DiskSize {
		log.Fatalf("Unable to add additional block as no space left on Disk [%s]", path)
	}

	log.Debugf("Next Block being written at address: %d", nextBlockAddr)

	block := NewBlock()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()
	if err != nil {
		return err
	}

	var binaryBuffer bytes.Buffer
	binary.Write(&binaryBuffer, binary.BigEndian, block)
	bytesWritten, err := file.WriteAt(binaryBuffer.Bytes(), int64(nextBlockAddr))
	if err != nil {
		return err
	}
	log.Debugf("Written %d bytes to file %s", bytesWritten, path)

	h.BlockCount = h.BlockCount + 1
	log.Printf("Added Block %d", h.BlockCount)

	// Update the header with new block count
	err = WriteHeader(path, h)
	if err != nil {
		return err
	}
	return nil
}

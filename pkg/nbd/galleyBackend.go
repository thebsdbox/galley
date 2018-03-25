package nbd

import (
	"golang.org/x/net/context"
)

type GalleyBackend struct {
	volumeID string
}

// WriteAt implements Backend.WriteAt
func (fb *GalleyBackend) WriteAt(ctx context.Context, b []byte, offset int64, fua bool) (int, error) {

	return 0, nil
}

// ReadAt implements Backend.ReadAt
func (gb *GalleyBackend) ReadAt(ctx context.Context, b []byte, offset int64) (int, error) {
	return 0, nil
}

// TrimAt implements Backend.TrimAt
func (gb *GalleyBackend) TrimAt(ctx context.Context, length int, offset int64) (int, error) {
	return length, nil
}

// Flush implements Backend.Flush
func (gb *GalleyBackend) Flush(ctx context.Context) error {
	return nil
}

// Close implements Backend.Close
func (gb *GalleyBackend) Close(ctx context.Context) error {
	return nil
}

// Size implements Backend.Size
func (gb *GalleyBackend) Geometry(ctx context.Context) (uint64, uint64, uint64, uint64, error) {
	return 0, 1, 32 * 1024, 128 * 1024 * 1024, nil
}

// Size implements Backend.HasFua
func (gb *GalleyBackend) HasFua(ctx context.Context) bool {
	return true
}

// Size implements Backend.HasFua
func (gb *GalleyBackend) HasFlush(ctx context.Context) bool {
	return true
}

// Generate a new file backend
func NewGalleyBackend(ctx context.Context, ec *ExportConfig) (Backend, error) {

	return &GalleyBackend{
		volumeID: "file",
	}, nil
}

func init() {
	RegisterBackend("galley", NewGalleyBackend)
}

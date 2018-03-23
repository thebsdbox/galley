package storage

// DetectStorage - This will identity the storage type from the path
func DetectStorage(path string) error {

	// Open the path and determine backing type

	ReadHeader(path)
	return nil
}

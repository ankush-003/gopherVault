package filebased

import (
	"fmt"
	"os"

	"github.com/gopherVault/db"
)

type FileBasedDB struct {
	config *db.DBConfig
}

func NewFileBasedDB(config *db.DBConfig) *FileBasedDB {
	return &FileBasedDB{config: config}
}

func (db *FileBasedDB) SaveToFile(data []byte) error {
	// filepointer with read and write permission (0644 -> permission to read and write)
	// os.O_CREATE -> create the file if it doesn't exist
	// os.O_WRONLY -> write only mode
	// os.O_TRUNC -> truncate the file if it exists
	fp, err := os.OpenFile(db.config.Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		// log.Fatalf("Failed to open file: %v", err)
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fp.Close()

	_, err = fp.Write(data)
	if err != nil {
		// log.Fatalf("Failed to write data to file: %v", err)
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	// flush the data to the file
	return fp.Sync()
}

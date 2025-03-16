package filebased

import (
	"fmt"
	"os"
	"github.com/gopherVault/db"
	"math/rand"
	"strconv"
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

func (db *FileBasedDB) SaveToFileAtomic(data []byte) error {
	// os.O_EXCL -> if the file already exists, return an error

	tmp := fmt.Sprintf("%s.%s.tmp", db.config.Path, randomString())
	fp, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		fp.Close()
		if err != nil {
			os.Remove(tmp) // remove the file if there is an error
		}
	}()

	_, err = fp.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	if err := fp.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %v", err)
	}

	err = os.Rename(tmp, db.config.Path)
	return nil
}

func randomString() string {
	num := rand.Intn(100000)
	return strconv.Itoa(num)
}
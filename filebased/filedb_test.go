package filebased

import (
	"testing"
	"sync"
	"strconv"
	"github.com/gopherVault/db"
)

var (
	dbConfig = &db.DBConfig{
		Path: "test.db",
	}
)

func TestFileBasedDB_SaveToFile(t *testing.T) {
	fileBasedDB := NewFileBasedDB(dbConfig)

	err := fileBasedDB.SaveToFile([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to save data to file: %v", err)
	}
}

func TestParalledSimpleSave(t *testing.T) {
	fileBasedDB := NewFileBasedDB(dbConfig)
	
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			fileBasedDB.SaveToFile([]byte("Hello, World by writer " + strconv.Itoa(i)))
		}()
	}
	wg.Wait()
}

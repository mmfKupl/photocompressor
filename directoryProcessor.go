package photocompressor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func filesProcessor(dir string, callback func(string), bunchSize int8) error {
	// Create a buffered channel to limit the number of concurrent goroutines
	sem := make(chan struct{}, bunchSize)
	var wg sync.WaitGroup

	// Walk through the directory
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has a .json extension
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Acquire a slot in the semaphore
			sem <- struct{}{}
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				// Call the callback function
				callback(p)
				// Release the slot in the semaphore
				<-sem
			}(path)
		}
		return nil
	})

	// Wait for all goroutines to finish
	wg.Wait()
	return err
}

func createDirIfNotExist(dir string) error {
	fmt.Println("Creating directory", dir)

	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

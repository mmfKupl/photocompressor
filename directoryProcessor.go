package photocompressor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func (compressor *PhotoCompressor) filesProcessor(callback func(string)) error {
	// Create a buffered channel to limit the number of concurrent goroutines
	sem := make(chan struct{}, compressor.BunchSize)
	var wg sync.WaitGroup

	totalFiles, err := countFilesInDir(compressor.DirPath)
	if err != nil {
		return err
	}

	// Walk through the directory
	processedFiles := 0
	err = filepath.WalkDir(compressor.DirPath, func(path string, info os.DirEntry, err error) error {
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
				processedFiles++
				compressor.updateLoader(p, processedFiles, totalFiles)
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

func (compressor *PhotoCompressor) updateLoader(inputFile string, processedFiles, totalFiles int) {
	//clearConsole()

	fmt.Printf("\033[1m\033[33mInput directory:\033[0m %s\n", compressor.DirPath)    // Bold yellow text for output file
	fmt.Printf("\033[1m\033[33mOutput directory:\033[0m %s\n", compressor.OutputDir) // Bold yellow text for output file
	fmt.Print("\n")
	fmt.Printf("\033[1m\033[32mProcessed files:\033[0m %d/%d\n", processedFiles, totalFiles) // Bold green text for processed files
	fmt.Printf("\033[1m\033[34mProcessing:\033[0m %s\n", inputFile)                          // Bold blue text for current file
	fmt.Print("\n")
}

func countFilesInDir(dir string) (int, error) {
	var count int
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

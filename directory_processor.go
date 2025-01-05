package photocompressor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	errorsLog   []error
	activeFiles []string
	mu          sync.Mutex
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
		if !info.IsDir() {
			sem <- struct{}{}
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				mu.Lock()
				activeFiles = append(activeFiles, p)
				mu.Unlock()
				callback(p)
				mu.Lock()
				activeFiles = removeFile(activeFiles, p)
				mu.Unlock()
				processedFiles += 1
				compressor.updateLoader(processedFiles, totalFiles)
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
	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

func (compressor *PhotoCompressor) updateLoader(processedFiles, totalFiles int) {
	clearConsole()

	fmt.Printf("\033[1m\033[33mInput directory:\033[0m %s\n", compressor.DirPath)
	fmt.Printf("\033[1m\033[33mOutput directory:\033[0m %s\n", compressor.OutputDir)
	fmt.Print("\n")
	fmt.Printf("\033[1m\033[32mProcessed files:\033[0m %d/%d\n", processedFiles, totalFiles)
	fmt.Print("\n")
	fmt.Println("\u001B[1m\u001B[34mCurrently processing files:\u001B[0m")
	mu.Lock()
	for _, file := range activeFiles {
		fmt.Println(file)
	}
	mu.Unlock()
	fmt.Print("\n")
	printErrors()
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

func logError(err error) {
	errorsLog = append(errorsLog, err)
}

func printErrors() {
	if len(errorsLog) > 0 {
		fmt.Println("Errors occurred during execution:")
		for _, err := range errorsLog {
			fmt.Println(err)
		}
	}
}

func removeFile(files []string, file string) []string {
	for i, f := range files {
		if f == file {
			return append(files[:i], files[i+1:]...)
		}
	}
	return files
}

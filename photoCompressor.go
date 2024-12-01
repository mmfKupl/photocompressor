package photocompressor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PhotoCompressor struct {
	DirPath   string
	BunchSize int8
	OutputDir string
}

func (compressor *PhotoCompressor) Run() error {
	startTime := time.Now()
	defer printDuration(startTime)

	err := createDirIfNotExist(compressor.OutputDir)
	if err != nil {
		return err
	}

	err = compressor.filesProcessor(func(path string) {
		err := compressor.handleFile(path)
		if err != nil {
			logError(fmt.Errorf("Error processing file:", err))
		}
	})
	if err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}

func (compressor *PhotoCompressor) handleFile(path string) error {
	if filepath.Ext(path) == ".json" {
		err := copyFile(path, compressor.OutputDir)
		if err != nil {
			return err
		}
		return nil
	}

	err := compressor.copyPhotoFileWithCompression(path)
	if err != nil {
		return err
	}

	return nil
}

func (compressor *PhotoCompressor) copyPhotoFileWithCompression(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	outputPath := filepath.Join(compressor.OutputDir, strings.TrimSuffix(filepath.Base(path), ext)+ext)

	switch ext {
	case ".jpg", ".jpeg", ".png":
		outputPath = strings.TrimSuffix(outputPath, ext) + ".jpg"
		err := ffmpeg.Input(path).
			Output(outputPath, ffmpeg.KwArgs{"q:v": 25}).
			OverWriteOutput().
			Run()
		if err != nil {
			return fmt.Errorf("error compressing image '%s': %w", path, err)
		}
	case ".heic":
		// Just copy the HEIC file to the output folder
		err := copyFile(path, compressor.OutputDir)
		if err != nil {
			return fmt.Errorf("error copying HEIC file '%s': %w", path, err)
		}
	case ".mp4", ".avi", ".mov", ".mkv":
		outputPath = strings.TrimSuffix(outputPath, ext) + ".mp4"
		err := ffmpeg.Input(path).
			Output(outputPath, ffmpeg.KwArgs{"c:v": "libx264", "crf": 23}).
			OverWriteOutput().
			Run()
		if err != nil {
			return fmt.Errorf("error compressing video '%s': %w", path, err)
		}
	default:
		return fmt.Errorf("unsupported file type '%s': %s", path, ext)
	}

	return nil
}

func copyFile(path string, outputDir string) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer closeFileHandler(inputFile)

	outputPath := filepath.Join(outputDir, filepath.Base(path))
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer closeFileHandler(outputFile)

	_, err = io.Copy(outputFile, inputFile)
	return err
}

func printDuration(startTime time.Time) {
	duration := time.Since(startTime)
	if duration.Minutes() >= 1 {
		fmt.Printf("Total working time: %v minutes\n", duration.Minutes())
	} else if duration.Seconds() >= 1 {
		fmt.Printf("Total working time: %v seconds\n", duration.Seconds())
	} else {
		fmt.Printf("Total working time: %v milliseconds\n", duration.Milliseconds())
	}
}

func closeFileHandler(file *os.File) {
	err := file.Close()
	if err != nil {
		logError(fmt.Errorf("Error closing file:", err))
	}
}

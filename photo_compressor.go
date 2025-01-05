package photocompressor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	exif "photocompressor/pkg/exif_metadata"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PhotoCompressor struct {
	DirPath       string
	BunchSize     int8
	OutputDir     string
	CompressLevel int
}

func (compressor *PhotoCompressor) Run() error {
	startTime := time.Now()
	defer func() {
		printDuration(startTime)
		printDirSize(compressor.DirPath, "Input directory size", "red")
		printDirSize(compressor.OutputDir, "Output directory size", "green")
	}()

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
	if filepath.Ext(path) == ".DS_Store" {
		return nil
	}
	if filepath.Ext(path) == ".json" {
		err := copyFile(path, compressor.OutputDir)
		if err != nil {
			return err
		}
		return nil
	}

	outputPath, err := compressor.copyPhotoFileWithCompression(path)
	if err != nil {
		return err
	}

	err = copyExifMetadata(path, outputPath)
	if err != nil {
		return err
	}

	return nil
}

func copyExifMetadata(sourceFilePath, targetFilePath string) error {
	possibleJsonFile := sourceFilePath + ".json"
	isJsonFileExist := true
	useJsonMetadata := true
	var jsonExifMetadata *exif.ExifMetadata

	if _, err := os.Stat(possibleJsonFile); os.IsNotExist(err) {
		fmt.Printf("Json file not found: %s\n", possibleJsonFile)
		isJsonFileExist = false
		useJsonMetadata = false
	}

	if isJsonFileExist {
		googleTakeoutMetadata, err := exif.GetGoogleTakeoutMetadata(possibleJsonFile)
		if err != nil {
			fmt.Printf("Error getting metadata from json file: %w", err)
			useJsonMetadata = false
		} else {
			jsonExifMetadata = googleTakeoutMetadata.ToExifMetadata()
		}
	}

	sourceExifMetadata, err := exif.GetFileMetadata(sourceFilePath)
	if err != nil {
		fmt.Printf("Error getting metadata from source file: %w", err)
		useJsonMetadata = true
	}

	if sourceExifMetadata != nil && jsonExifMetadata != nil {
		originalTime, _ := sourceExifMetadata.GetOriginalTime()
		jsonTime, _ := jsonExifMetadata.GetOriginalTime()

		if originalTime == jsonTime {
			useJsonMetadata = false
		}
	}

	if useJsonMetadata {
		return exif.AddMetadataToFile(targetFilePath, jsonExifMetadata)
	} else {
		return exif.CloneMetadataToFile(sourceFilePath, targetFilePath)
	}

}

func (compressor *PhotoCompressor) copyPhotoFileWithCompression(path string) (outputPath string, err error) {
	ext := strings.ToLower(filepath.Ext(path))
	outputPath = filepath.Join(compressor.OutputDir, strings.TrimSuffix(filepath.Base(path), ext)+ext)

	switch ext {
	case ".jpg", ".jpeg", ".png":
		outputPath = strings.TrimSuffix(outputPath, ext) + ".jpg"
		err := ffmpeg.Input(path).
			Output(outputPath, ffmpeg.KwArgs{"q:v": compressor.CompressLevel}).
			OverWriteOutput().
			Run()
		if err != nil {
			return "", fmt.Errorf("error compressing image '%s': %w", path, err)
		}
	case ".heic":
		// Just copy the HEIC file to the output folder
		err := copyFile(path, compressor.OutputDir)
		if err != nil {
			return "", fmt.Errorf("error copying HEIC file '%s': %w", path, err)
		}
	case ".mp4", ".avi", ".mov", ".mkv":
		outputPath = strings.TrimSuffix(outputPath, ext) + ".mp4"
		err := ffmpeg.Input(path).
			Output(outputPath, ffmpeg.KwArgs{
				"c:v": "libx264",
				"crf": 28,   // Higher compression
				"r":   30,   // Fixed framerate (e.g., 30 fps)
				"b:v": "1M", // Bitrate (e.g., 1 Mbps)
			}).
			OverWriteOutput().
			Run()
		if err != nil {
			return "", fmt.Errorf("error compressing video '%s': %w", path, err)
		}
	default:
		return "", fmt.Errorf("unsupported file type '%s': %s", path, ext)
	}

	return outputPath, nil
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
		fmt.Printf("Total working time: %.2f minutes\n", float64(duration.Minutes()))
	} else if duration.Seconds() >= 1 {
		fmt.Printf("Total working time: %.2f seconds\n", float64(duration.Seconds()))
	} else {
		fmt.Printf("Total working time: %.2f milliseconds\n", float64(duration.Milliseconds()))
	}
}

func printDirSize(dirPath, label, colorFlag string) {
	var size int64
	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		logError(fmt.Errorf("Error calculating directory size: %w", err))
		return
	}

	var colorCode string
	if colorFlag == "red" {
		colorCode = "\033[1m\033[31m" // Bold red
	} else {
		colorCode = "\033[1m\033[32m" // Bold green
	}

	fmt.Printf("\033[1m\033[33m%s:\033[0m %s%.2f GB\033[0m\n", label, colorCode, float64(size)/(1024*1024*1024))
}

func closeFileHandler(file *os.File) {
	err := file.Close()
	if err != nil {
		logError(fmt.Errorf("Error closing file:", err))
	}
}

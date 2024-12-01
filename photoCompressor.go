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
	endTime := time.Now()

	defer func() {
		endTime = time.Now()
		duration := endTime.Sub(startTime)
		fmt.Printf("Total working time: %.2f seconds\n", duration.Seconds())
	}()

	err := createDirIfNotExist(compressor.OutputDir)
	if err != nil {
		return err
	}

	err = filesProcessor(compressor.DirPath, func(path string) {
		err := compressor.handleFile(path)
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
	}, compressor.BunchSize)
	if err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}

func (compressor *PhotoCompressor) handleFile(path string) error {
	fmt.Println("Processing", path)

	metadata, err := parseMetadata(path)
	if err != nil {
		return err
	}

	if photoMeta, ok := metadata.(*photoMetadata); ok {
		err := compressor.copyPhotoFileWithCompression(photoMeta)
		if err != nil {
			return err
		}
	}

	err = copyFile(path, compressor.OutputDir)
	if err != nil {
		return err
	}

	return nil
}

func (compressor *PhotoCompressor) copyPhotoFileWithCompression(photo *photoMetadata) error {
	path := photo.FilePath
	ext := strings.ToLower(filepath.Ext(path))
	outputPath := filepath.Join(compressor.OutputDir, strings.TrimSuffix(filepath.Base(path), ext)+"-compressed"+ext)

	switch ext {
	case ".jpg", ".jpeg", ".png", ".heic":
		outputPath = strings.TrimSuffix(outputPath, ext) + ".jpg"
		err := ffmpeg.Input(path).
			Output(outputPath, ffmpeg.KwArgs{"q:v": 25}).
			OverWriteOutput().
			Run()
		if err != nil {
			return fmt.Errorf("error compressing image: %w", err)
		}
	case ".mp4", ".avi", ".mov", ".mkv":
		outputPath = strings.TrimSuffix(outputPath, ext) + ".mp4"
		err := ffmpeg.Input(path).
			Output(outputPath, ffmpeg.KwArgs{"c:v": "libx264", "crf": 23}).
			OverWriteOutput().
			Run()
		if err != nil {
			return fmt.Errorf("error compressing video: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
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
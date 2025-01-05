package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"photocompressor"
	"runtime"
)

func main() {
	// Limit the CPU usage to 2 cores
	runtime.GOMAXPROCS(1)

	if !checkFFmpegInstalled() {
		fmt.Println("Error: ffmpeg is not installed.")
		os.Exit(1)
	}

	// Define flags
	inputDir := flag.String("input", "", "Input directory path")
	outputDir := flag.String("output", "", "Output directory path")
	bunchSize := flag.Int("bunch", 5, "Bunch size")
	compressLevel := flag.Int("compressLevel", 23, "Compression level")

	// Parse flags
	flag.Parse()

	// Check if flags are provided, if not ask for input
	if *inputDir == "" {
		fmt.Print("Enter input directory path: ")
		fmt.Scanf("%[^\n]", inputDir)
	}
	if *inputDir == "" {
		fmt.Println("Error: Input directory path is required.")
		os.Exit(1)
	}
	// Check if the input directory exists
	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		fmt.Println("Error: Input directory does not exist.")
		fmt.Println(*inputDir)
		os.Exit(1)
	}

	if *outputDir == "" {
		fmt.Print("Enter output directory path: ")
		fmt.Scanf("%[^\n]", outputDir)
	}
	if *outputDir == "" {
		*outputDir = *inputDir + "-compressed"
	}

	fmt.Println("Input directory:", *inputDir)
	fmt.Println("Output directory:", *outputDir)

	if *bunchSize == 0 {
		*bunchSize = 5
	}

	compressor := photocompressor.PhotoCompressor{
		DirPath:       *inputDir,
		BunchSize:     int8(*bunchSize),
		OutputDir:     *outputDir,
		CompressLevel: *compressLevel,
	}

	err := compressor.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func checkFFmpegInstalled() bool {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	return err == nil
}

package main

import (
	"flag"
	"fmt"
	"os"
	"photocompressor"
)

func main() {
	// Define flags
	inputDir := flag.String("input", "", "Input directory path")
	outputDir := flag.String("output", "", "Output directory path")
	bunchSize := flag.Int("bunch", 5, "Bunch size")

	// Parse flags
	flag.Parse()

	// Check if flags are provided, if not ask for input
	if *inputDir == "" {
		fmt.Print("Enter input directory path: ")
		fmt.Scanln(inputDir)
	}
	if *inputDir == "" {
		fmt.Println("Error: Input directory path is required.")
		os.Exit(1)
	}
	// Check if the input directory exists
	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		fmt.Println("Error: Input directory does not exist.")
		os.Exit(1)
	}

	if *outputDir == "" {
		fmt.Print("Enter output directory path: ")
		fmt.Scanln(outputDir)
	}
	if *outputDir == "" {
		*outputDir = *inputDir + "-compressed"
	}

	if *bunchSize == 0 {
		*bunchSize = 5
	}

	compressor := photocompressor.PhotoCompressor{
		DirPath:   *inputDir,
		BunchSize: int8(*bunchSize),
		OutputDir: *outputDir,
	}

	err := compressor.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

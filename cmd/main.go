package main

import (
	"flag"
	"fmt"
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
	if *outputDir == "" {
		fmt.Print("Enter output directory path: ")
		fmt.Scanln(outputDir)
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

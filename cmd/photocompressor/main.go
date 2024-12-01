package main

import (
	"flag"
	"fmt"
	"os"
	"photocompressor"
)

func main() {
	fmt.Println("INIT")
	// Define flags
	inputDir := flag.String("input", "", "Input directory path")
	outputDir := flag.String("output", "", "Output directory path")
	bunchSize := flag.Int("bunch", 5, "Bunch size")

	// Parse flags
	flag.Parse()

	fmt.Printf("inputDir: %s\n", *inputDir)
	fmt.Printf("outputDir: %s\n", *outputDir)
	fmt.Printf("bunchSize: %s\n\n", *bunchSize)

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
		DirPath:   *inputDir,
		BunchSize: int8(*bunchSize),
		OutputDir: *outputDir,
	}

	err := compressor.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

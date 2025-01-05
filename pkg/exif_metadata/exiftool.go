package exif_metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func GetFileMetadata(filePath string) (ExifMetadata, error) {
	output, err := executeExiftoolCommand(filePath, "-json")
	if err != nil {
		return nil, err
	}

	var fileMetadata []ExifMetadata
	if err := json.Unmarshal(output, &fileMetadata); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON output: %w", err)
	}

	if len(fileMetadata) == 0 {
		return nil, fmt.Errorf("no file info found")
	}

	return fileMetadata[0], nil
}

func CloneMetadataToFile(sourceFilePath, targetFilePath string) error {
	_, err := executeExiftoolCommand("-tagsFromFile", sourceFilePath, "-all:all<=all:all", targetFilePath, "-overwrite_original")
	if err != nil {
		return err
	}
	return nil
}

func AddMetadataToFile(targetFilePath string, metadata *ExifMetadata) error {
	jsonMetadata, err := metadata.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println("JJJ", string(jsonMetadata))

	tempFileName, err := createTempJsonFile(metadata)
	defer func(filePath string) {
		err := removeTempJsonFile(filePath)
		if err != nil {
			fmt.Println("error removing temp file:", err)
		}
	}(tempFileName)

	o, err := executeExiftoolCommand(
		"-json="+tempFileName,
		targetFilePath,
		"-overwrite_original",
	)
	if err != nil {
		return err
	}

	fmt.Println(string(o))

	return nil
}

func executeExiftoolCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("exiftool", args...)
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing exiftool: %w", err)
	}

	return output, nil
}

func createTempJsonFile(metadata *ExifMetadata) (string, error) {
	jsonMetadata, err := metadata.ToJSON()
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "metadata*.json")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}

	if _, err := tempFile.Write(jsonMetadata); err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func removeTempJsonFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error deleting temp file: %w", err)
	}

	return nil
}

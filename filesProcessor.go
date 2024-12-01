package photocompressor

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type timeInfo struct {
	Timestamp string `json:"timestamp"`
}

type photoMetadata struct {
	Title          string   `json:"title"`
	CreationTime   timeInfo `json:"creationTime"`
	PhotoTakenTime timeInfo `json:"photoTakenTime"`
	FilePath       string
}

type albumMetadata struct {
	Title string   `json:"title"`
	Date  timeInfo `json:"date"`
}

func parseMetadata(path string) (interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var photoMeta photoMetadata
	var albumMeta albumMetadata

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&photoMeta); err == nil {
		if photoMeta.CreationTime.Timestamp != "" {
			photoMeta.FilePath = strings.TrimSuffix(path, ".json")
			return &photoMeta, nil
		}
	}

	file.Seek(0, 0) // Reset file pointer to the beginning
	decoder = json.NewDecoder(file)
	if err := decoder.Decode(&albumMeta); err == nil {
		if albumMeta.Date.Timestamp != "" {
			return &albumMeta, nil
		}
	}

	return nil, errors.New("unknown metadata format")
}

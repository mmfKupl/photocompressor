package exif_metadata

import (
	"encoding/json"
	"fmt"
	"github.com/KingAkeem/go-dms/dms"
	"os"
	"strconv"
	"time"
)

func GetGoogleTakeoutMetadata(filePath string) (GoogleTakeoutMetadata, error) {
	var metadata GoogleTakeoutMetadata

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return metadata, fmt.Errorf("error reading file: %w", err)
	}

	if err := json.Unmarshal(fileContent, &metadata); err != nil {
		return metadata, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return metadata, nil
}

func (g *GoogleTakeoutMetadata) ToExifMetadata() *ExifMetadata {
	dmsCoords := g.getDMS()
	exifMetadata := ExifMetadata{
		"Title":              g.Title,
		"Description":        g.Description,
		"DateTimeOriginal":   g.getPhotoTakenTime().Format(ExifDateTimeFormat),
		"OffsetTimeOriginal": "+00:00",
		"GPSAltitude":        g.getFormatedAltitude(),
		"GPSLatitude":        dmsCoords.Latitude.String(),
		"GPSLongitude":       dmsCoords.Longitude.String(),
	}

	return &exifMetadata
}

func (g *GoogleTakeoutMetadata) getPhotoTakenTime() time.Time {
	photoTakenTime, err := strconv.ParseInt(g.PhotoTakenTime.Timestamp, 10, 64)
	if err != nil {
		photoTakenTime = 0
	}

	return time.Unix(photoTakenTime, 0)
}

func (g *GoogleTakeoutMetadata) getDMS() *dms.DMS {
	dmsCoords, err := dms.NewDMS(dms.DecimalDegrees{
		Latitude:  g.GeoDataExif.Latitude,
		Longitude: g.GeoDataExif.Longitude,
	})

	if err != nil {
		return nil
	}

	return dmsCoords
}

func (g *GoogleTakeoutMetadata) getFormatedAltitude() string {
	return fmt.Sprintf("%.1f m", g.GeoDataExif.Altitude)
}

package exif_metadata

type CreationTime struct {
	Timestamp string `json:"timestamp"`
	Formatted string `json:"formatted"`
}

type PhotoTakenTime struct {
	Timestamp string `json:"timestamp"`
	Formatted string `json:"formatted"`
}

type GeoDataExif struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Altitude      float64 `json:"altitude"`
	LatitudeSpan  float64 `json:"latitudeSpan"`
	LongitudeSpan float64 `json:"longitudeSpan"`
}

type GeoData struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Altitude      float64 `json:"altitude"`
	LatitudeSpan  float64 `json:"latitudeSpan"`
	LongitudeSpan float64 `json:"longitudeSpan"`
}

type GoogleTakeoutMetadata struct {
	CreationTime   CreationTime   `json:"creationTime"`
	PhotoTakenTime PhotoTakenTime `json:"photoTakenTime"`
	GeoDataExif    GeoDataExif    `json:"geoDataExif"`
	GeoData        GeoData        `json:"geoData"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
}

const (
	ExifDateFormat     = "2006:01:02"
	ExifTimeFormat     = "15:04:05"
	ExifDateTimeFormat = ExifDateFormat + " " + ExifTimeFormat
)

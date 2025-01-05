package exif_metadata

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExifMetadata map[string]any

func (em *ExifMetadata) GetByKey(k string) interface{} {
	if em == nil {
		return nil
	}
	v, found := (*em)[k]
	if !found {
		return nil
	}

	return v
}

func (em *ExifMetadata) GetOriginalTime() (*time.Time, error) {
	if em == nil {
		return nil, fmt.Errorf("metadata is nil")
	}
	dateTimeOriginal := em.GetByKey("DateTimeOriginal").(string)
	offsetTime := em.GetByKey("OffsetTime").(string)

	localTime, err := time.Parse(ExifDateTimeFormat, dateTimeOriginal)
	if err != nil {
		fmt.Println("error parsing DateTimeOriginal:", err)
		return nil, err
	}

	offset, err := time.ParseDuration(offsetTime[:3] + "h" + offsetTime[4:] + "m")
	if err != nil {
		fmt.Println("error parsing OffsetTime:", err)
		return nil, err
	}

	utcTime := localTime.Add(-offset)

	return &utcTime, nil
}

func (em *ExifMetadata) ToJSON() ([]byte, error) {
	// Преобразуем карту в JSON строку
	jsonData, err := json.Marshal(em)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshalling JSON: %w", err)
	}

	return jsonData, nil
}

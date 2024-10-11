package utils

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type AudioTrack struct {
	Type         string `json:"type"`
	Label        string `json:"label"`
	TrackId      string `json:"trackId"`
	FileName     string `json:"fileName"`
	StartTimeStr string `json:"startTime"`
	StartTime    time.Time
	EndTimeStr   string `json:"endTime"`
}

func CreateFromFile(path string) ([]AudioTrack, error) {
	// Открываем файл
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	var tracks []AudioTrack

	// Декодируем JSON из файла в структуру
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tracks)
	if err != nil {
		log.Printf("Ошибка при декодировании JSON:", err)
		return nil, err
	}

	for i := range tracks {
		layout := "2006-01-02T15:04:05.000Z"
		tracks[i].StartTime, err = time.Parse(layout, tracks[i].StartTimeStr)
		// log.Printf("layout %s, StartTimeStr %s, StartTime %f", layout, track.StartTimeStr, track.StartTime.UnixMilli())
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
		}
		// track.StartTime= starttime.UnixMilli()

	}

	return tracks, nil
}

// Функция для вычисления задержки в секундах
func CalculateOffset(videoStartTime, trackStartTime time.Time) float64 {
	return trackStartTime.Sub(videoStartTime).Seconds()
}

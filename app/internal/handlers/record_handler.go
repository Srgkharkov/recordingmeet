package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Srgkharkov/recordingmeet/internal/meet"
	"github.com/Srgkharkov/recordingmeet/internal/utils"
)

func HandleRecordRequest(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain")
	// w.Header().Set("Transfer-Encoding", "chunked")
	// w.Header().Set("Cache-Control", "no-cache")

	// flusher, ok := w.(http.Flusher)
	// if !ok {
	// 	http.Error(w, "Streaming not supported", http.StatusInternalServerError)
	// 	return
	// }
	// flusher.Flush()

	link := r.URL.Query().Get("link")
	if link == "" {
		http.Error(w, "Link parameter is required", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if link == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	record, err := meet.NewRecordByLink(link, userID, time.Now().UnixMilli())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing link: %v", err)
		return
	}

	switch record.Service {
	case "GM":

		record.Status = "CREATED"

		// Кодируем объект Record в JSON
		recordJSON, err := json.Marshal(record)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error encoding record JSON: %v", err)
			// record.log.Printf("ошибка кодирования объекта Record в JSON: %v", err)
			// return fmt.Errorf("ошибка кодирования объекта Record в JSON: %v", err)
		}

		file, err := os.Create(filepath.Join(record.DirName, "record.json"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error creating record.json\n %v", err)
			// record.log.Printf("ошибка кодирования объекта Record в JSON: %v", err)
			// return fmt.Errorf("ошибка кодирования объекта Record в JSON: %v", err)
			return
		}
		file.Write(recordJSON)
		file.Close()

		err = utils.RunRecorder(record.ID, record.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error running recording container: %v", err)
			return
		}

		fmt.Fprintf(w, "Start recording, ID:%s", record.ID)
		log.Printf("Start recording, ID:%s", record.ID)
	//go exec recorder

	// ch := make(chan error)
	// go meet.RecordGoogleMeet(ch, record)
	// err := <-ch

	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "Error processing link: %v", err)
	// } else {

	// 	// Кодируем объект Record в JSON
	// 	recordJSON, err := json.Marshal(record)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "Error encoding JSON: %v", err)
	// 		// record.log.Printf("ошибка кодирования объекта Record в JSON: %v", err)
	// 		// return fmt.Errorf("ошибка кодирования объекта Record в JSON: %v", err)
	// 	}

	// 	// Устанавливаем заголовок Content-Type для JSON-ответа
	// 	w.Header().Set("Content-Type", "application/json")

	// 	// Пишем []byte в ответ
	// 	_, err = w.Write(recordJSON)
	// 	if err != nil {
	// 		http.Error(w, "Unable to write response", http.StatusInternalServerError)
	// 		return
	// 	}

	// }
	// for msg := range ch {
	// 	fmt.Fprintf(w, "%s\n", msg)
	// 	// flusher.Flush()
	// }
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported platform: %s\n", record.Service)
	}
}

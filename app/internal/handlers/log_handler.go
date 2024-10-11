package handlers

import (
	// "archive/zip"
	// "fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Srgkharkov/recordingmeet/internal/utils"
)

func HandleLogRequest(w http.ResponseWriter, r *http.Request) {
	// Получаем название директории из параметра запроса
	dir := r.URL.Query().Get("recordsid")
	if dir == "" {
		http.Error(w, "Необходимо указать директорию через параметр recordsid", http.StatusBadRequest)
		return
	}

	recordsDirName, err := utils.GetRecordsDir()
	if err != nil {
		http.Error(w, "Ошибка при чтении директории с записями", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(recordsDirName, dir, "record.log")

	// Открываем файл log.log
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Не удалось открыть лог файл", http.StatusInternalServerError)
		log.Printf("Ошибка открытия файла: %v", err)
		return
	}
	defer file.Close()

	// Устанавливаем заголовок ответа, что мы возвращаем текст
	w.Header().Set("Content-Type", "text/plain")

	// Читаем содержимое файла и отправляем его в ответ
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Не удалось прочитать лог файл", http.StatusInternalServerError)
		log.Printf("Ошибка при чтении файла: %v", err)
	}

}

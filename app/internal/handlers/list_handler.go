package handlers

import (
	// "encoding/json"
	"fmt"
	// "log"
	"net/http"
	"os"
	"path/filepath"

	// "time"

	// "github.com/Srgkharkov/recordingmeet/internal/meet"
	"github.com/Srgkharkov/recordingmeet/internal/utils"
)

func HandleListRequest(w http.ResponseWriter, r *http.Request) {
	recordsDirName, err := utils.GetRecordsDir()
	if err != nil {
		http.Error(w, "Ошибка при чтении директории с записями", http.StatusInternalServerError)
		return
	}

	// Список для хранения всех директорий
	directories := []string{}

	// Используем filepath.Walk для рекурсивного обхода директорий
	err = filepath.Walk(recordsDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Если это директория, добавляем в список
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Ошибка при обходе директорий", http.StatusInternalServerError)
		return
	}

	// Возвращаем список директорий в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v\n", directories)
}

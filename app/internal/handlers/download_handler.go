package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Srgkharkov/recordingmeet/internal/utils"
)

func HandleDownloadRequest(w http.ResponseWriter, r *http.Request) {
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

	// Проверяем, существует ли директория
	info, err := os.Stat(filepath.Join(recordsDirName, dir))
	if os.IsNotExist(err) {
		http.Error(w, "Директория архива не существует, возможно неверно указан идентификатор архива", http.StatusNotFound)
		return
	}
	if !info.IsDir() {
		http.Error(w, "Указанный путь не является директорией", http.StatusBadRequest)
		return
	}

	// Проверяем наличие файлов в директории
	files, err := os.ReadDir(filepath.Join(recordsDirName, dir))
	if err != nil {
		http.Error(w, "Ошибка при чтении директории", http.StatusInternalServerError)
		return
	}
	if len(files) == 0 {
		http.Error(w, "В директории нет файлов", http.StatusNotFound)
		return
	}

	isExistTimelineFile := false

	// Проверяем наличие файла timeline.json
	for _, file := range files {
		if file.Name() == "timeline.json" {
			isExistTimelineFile = true
			break
		}
	}

	if !isExistTimelineFile {
		http.Error(w, "В директории нет файла timeline.json.\nВозможно запись файлов еще не закончена или произошла ошибка при записи", http.StatusNotFound)
		return
	}

	// Создаем архив в памяти
	archiveName := dir + ".zip"
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", archiveName))
	w.Header().Set("Content-Type", "application/zip")
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Добавляем файлы в архив
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(filepath.Join(recordsDirName, dir), file.Name())
			err := AddFileToZip(zipWriter, filePath)
			if err != nil {
				http.Error(w, "Ошибка при добавлении файла в архив", http.StatusInternalServerError)
				return
			}
		}
	}
}

// Вспомогательная функция для добавления файлов в архив
func AddFileToZip(zipWriter *zip.Writer, filePath string) error {
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Получаем информацию о файле
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Создаем запись в архиве для текущего файла
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate // Используем сжатие

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Копируем содержимое файла в архив
	_, err = io.Copy(writer, file)
	return err
}
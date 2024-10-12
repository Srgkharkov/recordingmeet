package main

import (
	"encoding/json"
	// "fmt"
	// "io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Srgkharkov/recordingmeet/internal/meet"
	"github.com/Srgkharkov/recordingmeet/internal/utils"
)

var recorder_id = os.Getenv("RECORDER_ID")

func main() {

	if recorder_id == "" {
		log.Fatal("ENV RECORDER_ID is not set")
		return
	}

	recordsDirName, err := utils.GetFullpathDirRecord(recorder_id)
	if err != nil {
		log.Fatalf("Ошибка при чтении директории с записями. Error:%s", err)
		return
	}

	// Открываем (или создаем) файл для записи логов
	filelog, err := os.OpenFile(filepath.Join(recordsDirName, "record.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла лога: %v", err)
	}
	defer filelog.Close()

	// Перенаправляем стандартный логгер на запись в файл
	log.SetOutput(filelog)

	record, err := meet.ReadFromFile(filepath.Join(recordsDirName, "record.json"))
	if err != nil {
		log.Fatalf("Ошибка при чтении файла record.json: %v", err)
	}

	// // Открываем файл
	// file, err := os.Open(filepath.Join(recordsDirName, "record.json"))
	// if err != nil {
	// 	log.Fatalf("Ошибка открытия файла:", err)
	// 	return
	// }
	// defer file.Close()

	// var record meet.Record

	// // Декодируем JSON из файла в структуру
	// decoder := json.NewDecoder(file)
	// err = decoder.Decode(&record)
	// if err != nil {
	// 	log.Fatalf("Ошибка при декодировании JSON:", err)
	// 	return
	// }

	// log.Println("Decoding JSON complete.")

	// record.DirName = recordsDirName

	// record.SetLogger(log.Default())

	// // Форматируем JSON с отступами
	// formattedJSON, err := json.MarshalIndent(record, "", "  ") // отступ 2 пробела
	// if err != nil {
	// 	log.Fatalf("Ошибка при форматировании JSON: %v", err)
	// }

	// log.Print(formattedJSON)

	if record.Status != "CREATED" {
		log.Fatalf("Параметр status должен иметь значение CREATED")
	}

	switch record.Service {
	case "GM":

		record.Status = "PREPARING"

		// Кодируем объект Record в JSON
		recordJSON, err := json.Marshal(record)
		if err != nil {
			log.Fatalf("Error encoding record JSON: %v", err)
		}

		// Запись обновлённого JSON в файл
		err = os.WriteFile(filepath.Join(recordsDirName, "record.json"), recordJSON, 0644)
		if err != nil {
			log.Fatalf("Ошибка при записи JSON в файл: %v", err)
		}

		record.Status = "RECORDING"

		// Кодируем объект Record в JSON
		recordJSON, err = json.Marshal(record)
		if err != nil {
			log.Fatalf("Error encoding record JSON: %v", err)
		}

		// Запись обновлённого JSON в файл
		err = os.WriteFile(filepath.Join(recordsDirName, "record.json"), recordJSON, 0644)
		if err != nil {
			log.Fatalf("Ошибка при записи JSON в файл: %v", err)
		}

		meet.RecordGoogleMeet(record)

		record.Status = "RECORDED"

		// Запись обновлённого JSON в файл
		err = os.WriteFile(filepath.Join(recordsDirName, "record.json"), recordJSON, 0644)
		if err != nil {
			log.Fatalf("Ошибка при записи JSON в файл: %v", err)
		}
		// file, err := os.Create(filepath.Join(record.DirName, "record.json"))
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintf(w, "Error creating record.json\n %v", err)
		// 	// record.log.Printf("ошибка кодирования объекта Record в JSON: %v", err)
		// 	// return fmt.Errorf("ошибка кодирования объекта Record в JSON: %v", err)
		// 	return
		// }
		// file.Write(recordJSON)
		// file.Close()

		// fmt.Fprintf(w, "Start recording, ID:%s", record.ID)
		// log.Printf("Start recording, ID:%s", record.ID)

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
		// w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(w, "Unsupported platform: %s\n", record.Service)
	}
	// rec.

	// os
	//
}

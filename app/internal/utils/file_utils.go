package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func saveScreenshoot(bytes *[]byte, path string) error {
	err := os.WriteFile(path, *bytes, 0644)
	if err != nil {
		return err
	}

	return nil

}

func savePage(data *string, path string) error {
	// Сохраняем HTML в файл
	err := os.WriteFile(path, []byte(*data), os.ModePerm)
	if err != nil {
		log.Printf("Ошибка при записи html файла:%v\n", err)
		return err
	}

	log.Println("Запись в файл успешно завершена")
	return nil
}

func GetRecordsDir() (string, error) {
	// Получаем текущую рабочую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении рабочей директории:", err)
		return "", err
	}

	// Переходим на один уровень выше
	parentDir := filepath.Dir(currentDir)

	// Создаём путь для директории "records"
	recordsDir := filepath.Join(parentDir, "records")

	// Проверяем, существует ли директория
	if info, err := os.Stat(recordsDir); err == nil {
		// Если ошибок нет, проверяем, является ли это директорией
		if info.IsDir() {
			fmt.Println("Директория существует:", recordsDir)
		} else {
			fmt.Println("Это не директория:", recordsDir)
			return "", fmt.Errorf("Это не директория:%s", recordsDir)
		}
	} else if os.IsNotExist(err) {
		// Если директория не существует
		fmt.Println("Директория не существует:", recordsDir)
		return "", err
	} else {
		// В случае других ошибок
		fmt.Println("Ошибка при проверке директории:", err)
		return "", err
	}
	return recordsDir, nil
}

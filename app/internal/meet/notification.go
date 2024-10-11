package meet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	// "github.com/Srgkharkov/recordingmeet/internal/meet"
)

// NotifyAfterExecution отправляет уведомление с объектом Record по адресу из переменной окружения
func NotifyAfterExecution(record *Record) error {
	// Получаем URL из переменной окружения
	notificationURL := os.Getenv("NOTIFICATION_URL")
	if notificationURL == "" {
		log.Printf("переменная окружения NOTIFICATION_URL не задана")
		return fmt.Errorf("переменная окружения NOTIFICATION_URL не задана")
	}

	// Кодируем объект Record в JSON
	recordJSON, err := json.Marshal(record)
	if err != nil {
		log.Printf("ошибка кодирования объекта Record в JSON: %v", err)
		return fmt.Errorf("ошибка кодирования объекта Record в JSON: %v", err)
	}

	// Отправляем POST-запрос с JSON
	resp, err := http.Post(notificationURL, "application/json", bytes.NewBuffer(recordJSON))
	if err != nil {
		log.Printf("ошибка при отправке POST-запроса: %v", err)
		return fmt.Errorf("ошибка при отправке POST-запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем успешность запроса
	if resp.StatusCode != http.StatusOK {
		log.Printf("не удалось отправить уведомление по адресу %s, код ответа: %d", notificationURL, resp.StatusCode)
		return fmt.Errorf("не удалось отправить уведомление, код ответа: %d", resp.StatusCode)
	}

	log.Println("Уведомление успешно отправлено!")
	return nil
}

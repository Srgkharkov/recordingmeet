package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	// "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	// "os"
)

func RunRecorder(recorderID string, containerName string) error {
	// Создание нового клиента Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// Параметры конфигурации контейнера
	config := &container.Config{
		Image: "go-recorder", // Убедись, что образ уже существует (собран и доступен)
		Env:   []string{fmt.Sprintf("RECORDER_ID=%s", recorderID)},
	}

	// recordsDirName, err := GetRecordsDir()
	// if err != nil {
	// 	// http.Error(w, "Ошибка при чтении директории с записями", http.StatusInternalServerError)
	// 	return err
	// }

	hostConfig := &container.HostConfig{
		// AutoRemove: true,
		Binds: []string{fmt.Sprintf("%s:/records", os.Getenv("RECORDS_DIR"))}, // Пример монтирования томов
		// Binds: []string{"/home/sergei/recordingmeet/records:/records"}, // Пример монтирования томов
	}
	log.Println(hostConfig.Binds)

	// Создание контейнера
	resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	// Запуск контейнера
	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	log.Println("Recorder container started with ID:", resp.ID)
	return nil
}

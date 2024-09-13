package main

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/domsnapshot"
	"github.com/chromedp/chromedp"

	cu "github.com/Davincible/chromedp-undetected"
)

type MeetService struct {
	ShortName     string
	ID            string
	link          string
	dirName       string
	parentDirName string
}

// ParseLink — функция для разбора ссылки и создания переменной типа MeetService
func ParseLinkAndCreateDir(link string) (*MeetService, error) {
	// Разбор ссылки
	parsedURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	// Получаем домен (хост)
	host := parsedURL.Host

	// Получаем путь (например, /j/88330051248)
	pathURL := parsedURL.Path

	parentDirName, err := GetRecordsDir()
	if err != nil {
		// get path of the working directory
		parentDirName, err = os.Getwd()
		if err != nil {
			return nil, err
		}

	}

	// Проверяем, содержит ли ссылка meet.google.com
	if host == "meet.google.com" {
		ms := MeetService{
			ShortName:     "GM", // Google Meet
			ID:            strings.Join(strings.Split(pathURL, "/"), ""),
			link:          link,
			parentDirName: parentDirName,
		}
		ms.dirName = fmt.Sprintf("%s_%s_%d", ms.ShortName, ms.ID, time.Now().Unix())
		recordDir := path.Join(parentDirName, ms.dirName)

		err = os.Mkdir(recordDir, 0755)
		if err != nil {
			// Возвращаем клиенту ответ о том, что не смогли создать директорию
			return nil, err
		}
		return &ms, nil
	}

	// Проверяем, содержит ли ссылка zoom.us
	if strings.Contains(host, "zoom.us") {
		ms := MeetService{
			ShortName:     "ZOOM", // Google Meet
			ID:            strings.Split(pathURL, "/")[1],
			link:          link,
			parentDirName: parentDirName,
		}
		ms.dirName = fmt.Sprintf("%s_%s_%d", ms.ShortName, ms.ID, time.Now().Unix())
		recordDir := path.Join(parentDirName, ms.dirName)

		err = os.Mkdir(recordDir, 0755)
		if err != nil {
			// Возвращаем клиенту ответ о том, что не смогли создать директорию
			return nil, err
		}
		return &ms, nil
	}

	return nil, fmt.Errorf("неизвестный сервис")
}

func main() {
	http.HandleFunc("/record", handleRecordRequest)
	http.HandleFunc("/download", handleDownloadRequest)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// link := "https://meet.google.com/mvg-jivq-qdb"

	// ms, err := ParseLink(link)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = os.Mkdir(ms.dirName, 0755)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // get path of the working directory
	// workingDirectoryPath, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Путь к папке с расширением
	// extensionPath := "./tabCapture-recorder"

	// ID вашего расширения
	// extensionID := "khpccidebljocmbjapckddndjoefhpjo"

	// // New creates a new context for use with chromedp. With this context
	// // you can use chromedp as you normally would.
	// ctx, cancel, err := cu.New(cu.NewConfig(
	// 	// Remove this if you want to see a browser window.
	// 	cu.WithHeadless(), //Требуется xvfb
	// 	// cu.WithExtensions(extensionPath),
	// 	// cu.WithNoSandbox(true),
	// 	// cu.WithChromeFlags(
	// 	// 	chromedp.Flag("disable-popup-blocking", true), // Отключаем блокировку всплывающих окон
	// 	// 	// chromedp.Flag("profile.default_content_setting_values.automatic_downloads", 1), // Разрешение на автозагрузку нескольких файлов
	// 	// 	chromedp.Flag("download.default_directory", fmt.Sprintf("/home/sergei/recordingmeet/app/")), // Указываем директорию загрузки
	// 	// 	chromedp.Flag("download.prompt_for_download", false),                                        // Отключаем запрос на подтверждение загрузки
	// 	// ),

	// 	// If the webelement is not found within 10 minuties, timeout.
	// 	cu.WithTimeout(10*time.Minute),
	// ))
	// if err != nil {
	// 	panic(err)
	// }
	// defer cancel()

	// // Селектор кнопки без камеры и микрофона
	// buttonSelectorWOCamMic := `button[jsname="IbE0S"]`

	// // Селектор поля ввода имени
	// inputSelector := `input#c16, input#c17`

	// // Селектор кнопки Join
	// buttonSelectorJoin := `button[jsname="Qx7uuf"]`

	// // Селектор кнопки Завершения звонка
	// buttonSelectorEnd := `button[aria-label="Покинуть видеовстречу"]`

	// // Селектор для элемента с количеством участников
	// participantCountSelector := "div.uGOf1d"

	// // Переменная для хранения значения количества участников
	// var participantCount string
	// // buttonSelectorInExtension := `share-audio-button`

	// // // Отправка сообщения расширению 'khpccidebljocmbjapckddndjoefhpjo',
	// // extensionID := "khpccidebljocmbjapckddndjoefhpjo" // замените на ваш реальный ID расширения
	// // message := fmt.Sprintf(`
	// // 	console.log('chromedp.Evaluate');
	// // 	// Отправка сообщения в фоновый скрипт для начала захвата
	// // 	chrome.runtime.sendMessage('%s', { action: 'startCapture2' }, (response) => {
	// // 		if (chrome.runtime.lastError) {
	// // 			console.error('Error sending message:', chrome.runtime.lastError);
	// // 		} else if (response && response.status === 'success') {
	// // 			console.log('Capture started successfully');
	// // 		} else {
	// // 			console.error('Error starting capture:', response.error);
	// // 		}
	// // 	});
	// // `, extensionID)

	// // message := fmt.Sprintf(`document.querySelector('button[id="%s"]').click();`, "button34232432432")
	// // message2 := fmt.Sprintf(`document.querySelector('button[id="%s"]').click();`, "button342324324322")
	// // message2 := `// Отправка сообщения из popup.js или content script
	// // 	chrome.runtime.sendMessage({ greeting: "hello from button click" }, (response) => {
	// // 		if (response) {
	// // 			console.log('Response from background:', response.farewell);
	// // 		} else {
	// // 			console.error('No response from background script');
	// // 		}
	// // 	});
	// //  	// chrome.runtime.sendMessage({action: 'startCapture'}, (response) => {
	// //  	// 	if (response.status === 'success') {
	// //  	// 		console.log('Capture started successfully');
	// //  	// 	} else {
	// //  	// 		console.error('Error starting capture:', response.error);
	// //  	// 	}
	// //  	// });
	// // 	`

	// //
	// if err := chromedp.Run(ctx,
	// 	// Check if we pass anti-bot measures.
	// 	chromedp.Navigate("https://meet.google.com/mvg-jivq-qdb"),
	// 	browser.
	// 		SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
	// 		WithDownloadPath(fmt.Sprintf("%s/%s", workingDirectoryPath, ms.dirName)). //fmt.Sprintf("/home/sergei/recordingmeet/app/")).
	// 		WithEventsEnabled(true),
	// 	chromedp.WaitVisible(buttonSelectorWOCamMic), // Ожидание видимости кнопки
	// 	chromedp.Sleep(200*time.Millisecond),
	// 	chromedp.Click(buttonSelectorWOCamMic), // Клик по кнопке
	// 	chromedp.WaitVisible(inputSelector),    // Ожидание видимости поля ввода
	// 	chromedp.Sleep(2000*time.Millisecond),
	// 	chromedp.Clear(inputSelector),                    // Очистка поля ввода (если нужно)
	// 	chromedp.SendKeys(inputSelector, "BotRecording"), // Ввод текста в поле
	// 	chromedp.WaitEnabled(buttonSelectorJoin),         // Ожидание, пока кнопка станет активной
	// 	chromedp.Sleep(200*time.Millisecond),
	// 	chromedp.Click(buttonSelectorJoin), // Клик по кнопке
	// 	// chromedp.Sleep(15*time.Second),
	// 	chromedp.WaitVisible(buttonSelectorEnd),
	// 	chromedp.Sleep(1*time.Second),
	// 	chromedp.ActionFunc(func(ctx context.Context) error {
	// 		// Ваш JavaScript код для запуска на странице
	// 		return chromedp.Evaluate(mediarecorderjs, nil).Do(ctx)
	// 	}),
	// 	chromedp.Sleep(10*time.Second),
	// 	chromedp.ActionFunc(func(ctx context.Context) error {
	// 		for {
	// 			// Чтение значения количества участников
	// 			err := chromedp.Text(participantCountSelector, &participantCount, chromedp.NodeVisible).Do(ctx)
	// 			if err != nil {
	// 				log.Printf("Не удалось получить количество участников: %v", err)
	// 				continue
	// 			}

	// 			// Преобразование значения в целое число
	// 			count, err := strconv.Atoi(participantCount)
	// 			if err != nil {
	// 				log.Printf("Не удалось преобразовать количество участников в число: %v", err)
	// 				continue
	// 			}

	// 			// Если участников меньше 2, кликаем по кнопке выхода
	// 			if count < 2 {
	// 				log.Println("Участников меньше 2, покидаем встречу.")
	// 				err := chromedp.Click(buttonSelectorEnd).Do(ctx)
	// 				if err != nil {
	// 					log.Printf("Не удалось кликнуть по кнопке выхода: %v", err)
	// 				}
	// 				break
	// 			}

	// 			// Ожидание перед следующей проверкой
	// 			time.Sleep(5 * time.Second)
	// 		}
	// 		return nil
	// 	}),
	// 	chromedp.Sleep(15*time.Second),

	// 	// chromedp.Evaluate(message, nil),
	// 	// chromedp.Sleep(2*time.Second),
	// 	// chromedp.Evaluate(message2, nil),
	// 	// chromedp.Sleep(15*time.Second),
	// 	// chromedp.Sleep(1*time.Second),
	// 	// chromedp.Evaluate(message2, nil),
	// 	// chromedp.Evaluate(`// Отправка сообщения из popup.js или content script
	// 	// 	chrome.runtime.sendMessage({action: 'startCapture'}, (response) => {
	// 	// 		if (response.status === 'success') {
	// 	// 			console.log('Capture started successfully');
	// 	// 		} else {
	// 	// 			console.error('Error starting capture:', response.error);
	// 	// 		}
	// 	// 	});`, nil),
	// 	// chromedp.Navigate(fmt.Sprintf("chrome-extension://%s/home.html", extensionID)),
	// 	// chromedp.Sleep(1*time.Second),
	// 	// chromedp.WaitVisible(buttonSelectorInExtension),
	// 	// chromedp.Click(buttonSelectorInExtension),
	// 	// chromedp.Sleep(100*time.Second),
	// ); err != nil {
	// 	panic(err)
	// }

	fmt.Println("Undetected!")

	// // Открытие popup расширения и нажатие кнопок
	// err = chromedp.Run(ctx,
	// 	chromedp.Navigate("chrome-extension://khpccidebljocmbjapckddndjoefhpjo/home.html"),
	// 	chromedp.WaitVisible("#start-recording-button"),
	// 	chromedp.Click("#start-recording-button"),
	// 	chromedp.Sleep(5*time.Second), // Запись длится 5 секунд
	// 	chromedp.Click("#stop-recording-button"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Запись успешно выполнена")

	// // Дайте время для завершения действия после клика
	// time.Sleep(1 * time.Second)

	// fmt.Println("Undetected2!")

	// if err := chromedp.Run(ctx,
	// 	chromedp.Evaluate(message, nil),
	// ); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Message sent to extension!")

	// // Дайте время для завершения действия после клика
	// time.Sleep(50 * time.Second)

	// // Откройте popup расширения, используя JavaScript
	// // Это потребует, чтобы иконка расширения была видима на панели инструментов
	// if err := chromedp.Run(ctx,
	// 	chromedp.Evaluate(`chrome.runtime.sendMessage("`+extensionID+`", { action: "open_popup" });`, nil),
	// ); err != nil {
	// 	panic(err)
	// }

	// // Дайте время, чтобы popup открылся
	// time.Sleep(5 * time.Second)

	// // Перейдите в popup-окно расширения, если это возможно
	// if err := chromedp.Run(ctx,
	// 	chromedp.Click(buttonSelectorInExtension),
	// ); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Button clicked in extension popup!")
	// time.Sleep(50 * time.Second)

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

func handleRecordRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	link := r.URL.Query().Get("link")
	ms, err := ParseLinkAndCreateDir(link)
	if err != nil {
		// Возвращаем клиенту ответ о том, что не смогли обработать ссылку
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad link: %s\nError: %s\n", link, err)
		return
	}

	switch ms.ShortName {
	case "GM":
		ch := make(chan string)
		go recGM(ch, ms)
		for message := range ch { // Получаем данные из канала до его закрытия
			fmt.Fprintf(w, "%s\n", message)
			flusher.Flush()
		}
		break
	// case "ZOOM":
	// 	// Возвращаем клиенту ответ о том, что запись началась
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprintf(w, "Recording started for meeting: %s", ms.dirName)
	// 	break
	default:
		// Возвращаем клиенту ответ о том, что запись началась
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported platform: %s\n", ms.ShortName)
		break
	}
}

func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	// Получаем название директории из параметра запроса
	dir := r.URL.Query().Get("recordsid")
	if dir == "" {
		http.Error(w, "Необходимо указать директорию через параметр recordsid", http.StatusBadRequest)
		return
	}

	recordsDirName, err := GetRecordsDir()
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
			err := addFileToZip(zipWriter, filePath)
			if err != nil {
				http.Error(w, "Ошибка при добавлении файла в архив", http.StatusInternalServerError)
				return
			}
		}
	}
}

// Вспомогательная функция для добавления файлов в архив
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
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

// func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
// 	recarchive := r.URL.Query().Get("recarchive")
// 	// get path of the working directory
// 	// workingDirectoryPath, err := os.Getwd()
// 	// if err != nil {
// 	// 	// Возвращаем клиенту ответ о том, что не смогли получить путь к рабочей директории
// 	// 	w.WriteHeader(http.StatusInternalServerError)
// 	// 	fmt.Fprintf(w, "Can`t get working directory path, error: %s", err)
// 	// }
// 	src, err := os.Stat(fmt.Sprintf("./%s", recarchive))
// 	if err != nil || !src.IsDir() {
// 		// Если файл не найден, возвращаем 404
// 		http.Error(w, "File not found", http.StatusNotFound)
// 		return
// 	}

// 	// http.ServeFile(w, r, filePath)
// }

func recGM(ch chan string, ms *MeetService) {
	ch <- fmt.Sprintf("Archive ID:%s", ms.dirName)
	// func recGM(w *http.ResponseWriter, ms *MeetService) (int, error) {
	// New creates a new context for use with chromedp. With this context
	// you can use chromedp as you normally would.
	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		cu.WithHeadless(), //Требуется xvfb

		// If the webelement is not found within 10 minuties, timeout.
		cu.WithTimeout(5*time.Minute),
	))
	if err != nil {
		panic(err)
	}
	defer cancel()

	// Селектор кнопки без камеры и микрофона
	buttonSelectorWOCamMic := `button[jsname="IbE0S"]`

	// Селектор поля ввода имени
	inputSelector := `input[aria-label="Укажите свое имя"], input#c15, input#c16, input#c17`

	// Селектор кнопки Join
	buttonSelectorJoin := `button[jsname="Qx7uuf"]`

	// Селектор кнопки Завершения звонка
	buttonSelectorEnd := `button[aria-label="Покинуть видеовстречу"]`

	// Селектор для элемента с количеством участников
	participantCountSelector := "div.uGOf1d"

	// Переменная для хранения значения количества участников
	var participantCount string

	// // Запуск Chromedp в горутине
	// go func() {
	// var screenshotBuffer []byte
	// defer cancel() // Освобождаем ресурсы, когда горутина завершится
	if err := chromedp.Run(ctx,

		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath(path.Join(ms.parentDirName, ms.dirName)). //fmt.Sprintf("/home/sergei/recordingmeet/app/")).
			WithEventsEnabled(true),

		chromedp.Navigate(ms.link),

		runWithTimeout(
			ch,
			"Waiting for body visibility",
			10*time.Second,
			chromedp.WaitVisible(`body`, chromedp.ByQuery),
		),

		runWithTimeout(
			ch,
			"Waiting for buttonSelectorWOCamMic visibility",
			10*time.Second,
			chromedp.WaitVisible(buttonSelectorWOCamMic),
		), // Ожидание видимости кнопки

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(buttonSelectorWOCamMic), // Клик по кнопке

		runWithTimeout(
			ch,
			"Waiting for inputSelector visibility",
			10*time.Second,
			chromedp.WaitVisible(inputSelector),
		), // Ожидание видимости поля ввода

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Clear(inputSelector),                    // Очистка поля ввода (если нужно)
		chromedp.SendKeys(inputSelector, "BotRecording"), // Ввод текста в поле

		runWithTimeout(
			ch,
			"Waiting for buttonSelectorJoin visibility",
			10*time.Second,
			chromedp.WaitEnabled(buttonSelectorJoin),
		), // Ожидание, пока кнопка станет активной

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(buttonSelectorJoin), // Клик по кнопке

		runWithTimeout(
			ch,
			"Waiting for buttonSelectorEnd visibility",
			20*time.Second,
			chromedp.WaitVisible(buttonSelectorEnd),
		),

		chromedp.ActionFunc(func(ctx context.Context) error {
			ch <- fmt.Sprintf("Connection to the meeting was completed successfully.")
			close(ch)
			// Ваш JavaScript код для запуска на странице
			return chromedp.Evaluate(mediarecorderjs, nil).Do(ctx)
		}),
		chromedp.Sleep(10*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				// Чтение значения количества участников
				err := chromedp.Text(participantCountSelector, &participantCount, chromedp.NodeVisible).Do(ctx)
				if err != nil {
					log.Printf("Не удалось получить количество участников: %v", err)
					break
				}

				// Преобразование значения в целое число
				count, err := strconv.Atoi(participantCount)
				if err != nil {
					log.Printf("Не удалось преобразовать количество участников в число: %v", err)
					break
				}

				// Если участников меньше 2, кликаем по кнопке выхода
				if count < 2 {
					log.Println("Участников меньше 2, покидаем встречу.")
					err := chromedp.Click(buttonSelectorEnd).Do(ctx)
					if err != nil {
						log.Printf("Не удалось кликнуть по кнопке выхода: %v", err)
					}
					break
				}

				// Ожидание перед следующей проверкой
				time.Sleep(5 * time.Second)
			}
			return nil
		}),
		chromedp.Sleep(10*time.Second),
	); err != nil {
		ch <- fmt.Sprintf("Error in chromedp.Run: %s", err)
		close(ch)
		log.Printf("Error in chromedp.Run: %s", err)
		return
	}

	// }()
	return

}

// Вспомогательная функция для создания контекста с таймаутом для набора действий
func runWithTimeout(ch chan string, message string, timeout time.Duration, action chromedp.Action) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		log.Printf("Starting action: %s", message)
		// Создаём новый контекст с таймаутом для определённого набора действий
		newctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err := action.Do(newctx)
		if err != nil {
			ch <- fmt.Sprintf("Error during action: %s, error: %v\n", message, err)
			log.Printf("Error during action: %s, error: %v\n", message, err)

			var screenshotBuffer []byte
			chromedp.CaptureScreenshot(&screenshotBuffer).Do(ctx)
			filename := fmt.Sprintf("%d", time.Now().Unix())
			err := saveScreenshoot(&screenshotBuffer, fmt.Sprintf("../log/screenshots/%s.png", filename))
			if err != nil {
				log.Printf("Can`t save screenshot\n")
			} else {
				log.Printf("Saved screenshot:%s\n", filename)
			}

			var htmlSnapshot []string
			domsnapshot.CaptureSnapshot(htmlSnapshot)
			err = savePage(htmlSnapshot, fmt.Sprintf("../log/screenshots/%s.html", filename))
		} else {
			log.Printf("Successfully finished action: %s\n", message)
		}
		return err
	}
}

// // Вспомогательная функция для создания контекста с таймаутом для набора действий
// func sendMessage(ch chan string, message string, timeout time.Duration, action chromedp.Action) chromedp.ActionFunc {
// 	return func(ctx context.Context) error {
// 		log.Printf("Starting action: %s", message)
// 		// Создаём новый контекст с таймаутом для определённого набора действий
// 		newctx, cancel := context.WithTimeout(ctx, timeout)
// 		defer cancel()

// 		err := action.Do(newctx)
// 		if err != nil {
// 			log.Printf("Error during action: %s, error: %v\n", message, err)

// 			var screenshotBuffer []byte
// 			chromedp.CaptureScreenshot(&screenshotBuffer).Do(ctx)
// 			filename := fmt.Sprintf("%d", time.Now().Unix())
// 			err := saveScreenshoot(&screenshotBuffer, fmt.Sprintf("../log/screenshots/%s.png", filename))
// 			if err != nil {
// 				log.Printf("Can`t save screenshot\n")
// 			} else {
// 				log.Printf("Saved screenshot:%s\n", filename)
// 			}

// 			var htmlSnapshot []string
// 			domsnapshot.CaptureSnapshot(htmlSnapshot)
// 			err = savePage(htmlSnapshot, fmt.Sprintf("../log/screenshots/%s.html", filename))
// 		} else {
// 			log.Printf("Successfully finished action: %s\n", message)
// 		}
// 		return err
// 	}
// }

func saveScreenshoot(bytes *[]byte, path string) error {
	err := os.WriteFile(path, *bytes, 0644)
	if err != nil {
		return err
	}

	return nil

}

func savePage(data []string, path string) error {
	// Создаем файл для записи
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Не удалось создать файл: %v", err)
	}
	defer file.Close()

	// Создаем writer для буферизованной записи
	writer := bufio.NewWriter(file)

	// Записываем строки в файл
	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Ошибка при записи строки: %v", err)
		}
	}

	// Не забудьте сбросить данные из буфера в файл
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Ошибка при завершении записи в файл: %v", err)
	}

	log.Println("Запись в файл успешно завершена")
	return nil
}
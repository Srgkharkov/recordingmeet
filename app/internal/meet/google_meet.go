package meet

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/Srgkharkov/recordingmeet/internal/utils"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/cdproto/browser"

	// "github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// var cmd *exec.Cmd

// RecordGoogleMeet starts the recording process for Google Meet.
// func RecordGoogleMeet(ch chan error, rec *Record) {
func RecordGoogleMeet(rec *Record) error {
	// defer rec.CloseFile()
	defer NotifyAfterExecution(rec)
	// log.Printf("Archive ID:%s\nConnecting Google Meet: %s\n", rec.ID, rec.Link)
	// ch <- fmt.Sprintf("Archive ID:%s", rec.ID)
	log.Printf("Archive ID:%s\nConnecting Google Meet: %s\n", rec.ID, rec.Link)
	// New creates a new context for use with chromedp. With this context
	// you can use chromedp as you normally would.
	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		// cu.WithHeadless(), //Required xvfb
		cu.WithChromeFlags(chromedp.Flag("window-size", "1920,1080")),
		// cu.WithChromeFlags(chromedp.Flag("window-position", "0,0")),
		// cu.WithChromeFlags(chromedp.Flag("start-maximized", true)),
		cu.WithChromeFlags(chromedp.Flag("kiosk", true)),
		// cu.WithChromeFlags(chromedp.Flag("start-fullscreen", true)),
		// cu.WithChromeFlags(chromedp.Flag("start-maximized", true)),

		// If the webelement is not found within 121 minuties, timeout.
		cu.WithTimeout(121*time.Minute),
	))
	if err != nil {
		panic(err)
	}
	defer cancel()

	// Селектор кнопки без камеры и микрофона
	buttonSelectorWOCamMic := `button[jsname="IbE0S"]`

	// Селектор поля ввода имени
	inputSelector := `input[aria-label="Your name"], input[aria-label="Укажите свое имя"], input#c15, input#c16, input#c17`

	// Селектор кнопки Join
	buttonSelectorJoin := `button[jsname="Qx7uuf"]`

	// Селектор кнопки Завершения звонка
	buttonSelectorEnd := `button[aria-label="Покинуть видеовстречу"], button[aria-label="Leave call"]`

	// Селектор для элемента с количеством участников
	participantCountSelector := "div.uGOf1d"

	// Переменная для хранения значения количества участников
	var participantCount string

	// listenForNetworkEvent(ctx, rec)

	if err := chromedp.Run(ctx,

		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath(rec.DirName). //fmt.Sprintf("/home/sergei/recordingmeet/app/")).
			WithEventsEnabled(true),

		chromedp.Navigate(rec.Link),

		// startRecording(cmd),

		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	// ch <- fmt.Sprintf("Connection to the meeting was completed successfully.")
		// 	log.Printf("Connection to the meeting was completed successfully.")
		// 	// ch <- nil
		// 	// close(ch)
		// 	// Ваш JavaScript код для запуска на странице
		// 	return chromedp.Evaluate(utils.Mediarecorderjs, nil).Do(ctx)
		// }),
		runWithTimeout(
			// ch,
			"Waiting for body visibility",
			10*time.Second,
			chromedp.WaitVisible(`body`, chromedp.ByQuery),
			rec,
		),

		// startRecording(),

		tryClosePopup(
			ctx,
			buttonSelectorWOCamMic,
			3*time.Second,
		), // Ожидание видимости кнопки

		chromedp.Sleep(100*time.Millisecond),
		// chromedp.Click(buttonSelectorWOCamMic), // Клик по кнопке

		runWithTimeout(
			// ch,
			"Waiting for inputSelector visibility",
			10*time.Second,
			chromedp.WaitVisible(inputSelector),
			rec,
		), // Ожидание видимости поля ввода

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Clear(inputSelector),                    // Очистка поля ввода (если нужно)
		chromedp.SendKeys(inputSelector, "BotRecording"), // Ввод текста в поле

		runWithTimeout(
			// ch,
			"Waiting for buttonSelectorJoin visibility",
			10*time.Second,
			chromedp.WaitEnabled(buttonSelectorJoin),
			rec,
		), // Ожидание, пока кнопка станет активной

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(buttonSelectorJoin), // Клик по кнопке

		runWithTimeout(
			// ch,
			"Waiting for buttonSelectorEnd visibility",
			60*time.Second,
			chromedp.WaitVisible(buttonSelectorEnd),
			rec,
		),
	); err != nil {
		// ch <- fmt.Sprintf("Error in chromedp.Run: %s", err)
		log.Printf("Error in chromedp.Run: %s", err)
		// close(ch)
		// log.Printf("Error in chromedp.Run: %s", err)
		return err
	}

	cmd := exec.Command("ffmpeg", "-video_size", "1920x1080", "-framerate", "25", "-f", "x11grab", "-i", ":99", filepath.Join(rec.DirName, "video.mp4"))
	// cmd := exec.Command("ffmpeg", "-video_size", "1920x1080", "-framerate", "25", "-f", "x11grab", "-i", ":99", "-f", "alsa", "-i", "hw:0,0", "/records/output.mp4")

	// Запуск команды
	if err := cmd.Start(); err != nil {
		log.Printf("Не удалось начать запись: %v", err)
		return err
	}

	log.Println("Запись началась...")

	videoStartTime := time.Now().Add(time.Millisecond * 100)

	if err := chromedp.Run(ctx,
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	// ch <- fmt.Sprintf("Connection to the meeting was completed successfully.")
		// 	log.Printf("Connection to the meeting was completed successfully.")
		// 	// ch <- nil
		// 	// close(ch)
		// 	// Ваш JavaScript код для запуска на странице
		// 	return chromedp.Evaluate(utils.Mediarecorderjs, nil).Do(ctx)
		// }),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(utils.Mediarecorderjs, nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// ch <- fmt.Sprintf("Connection to the meeting was completed successfully.")
			// log.Printf("Connection to the meeting was completed successfully.")
			// ch <- nil
			// close(ch)
			// Ваш JavaScript код для запуска на странице
			return chromedp.Sleep(10 * time.Second).Do(ctx)
		}),
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
					log.Printf("Участников меньше 2, покидаем встречу.")
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
		// stopRecording(cmd),
		// chromedp.Sleep(10*time.Second),
	); err != nil {
		// ch <- fmt.Sprintf("Error in chromedp.Run: %s", err)
		log.Printf("Error in chromedp.Run: %s", err)
		// close(ch)
		// log.Printf("Error in chromedp.Run: %s", err)
		return err
	}

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		log.Printf("Не удалось остановить запись: %v", err)
		return err
	}
	log.Println("Запись завершена")

	time.Sleep(1 * time.Second)

	tracks, err := utils.CreateFromFile(filepath.Join(rec.DirName, "timeline.json"))
	if err != nil {
		log.Printf("Ошибка при обработке файла timeline.json")
		return err
	}

	// Начинаем строить команду FFmpeg
	args := []string{"-i", filepath.Join(rec.DirName, "video.mp4")}

	// Добавляем каждый аудиотрек с учетом задержки
	for _, track := range tracks {
		log.Printf("layout %s, StartTimeStr %s", videoStartTime.UnixMilli(), track.StartTime.UnixMilli())
		offset := utils.CalculateOffset(videoStartTime, track.StartTime)
		log.Printf("layout %s, StartTimeStr %s, StartTime %lf", videoStartTime.UnixMilli(), track.StartTime.UnixMilli(), offset)
		// Если есть задержка, добавляем её через itsoffset
		if offset > 0 {
			args = append(args, "-itsoffset", fmt.Sprintf("%.3f", offset))
		}
		// Добавляем аудиотрек
		args = append(args, "-i", filepath.Join(rec.DirName, track.FileName))
	}

	// argsaudio := make([]string, len(args))
	// copy(argsaudio, args[2:])

	// Добавляем фильтры и маппинг видео и аудио для выхода
	args = append(args, "-c:v", "copy", "-c:a", "aac", filepath.Join(rec.DirName, "record.mp4"))

	// filtercomplexstr := "-filter_complex '"
	// // argsaudio = append(argsaudio, "-filter_complex '")
	// for i, _ := range tracks {
	// 	filtercomplexstr += fmt.Sprintf("[%d:a]", i)
	// }
	// filtercomplexstr += "amix=inputs=3:duration=longest'"
	// argsaudio = append(argsaudio, "-c:a", "libmp3lame", "-q:a", "2", filepath.Join(rec.DirName, "output_without_video.mp3"))

	// Печатаем сформированную команду
	fmt.Println("FFmpeg command:", args)

	// Выполняем команду
	cmd = exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command: %v\nOutput: %s", err, output)
	}
	fmt.Println(string(output))
	// fmt.Println(argsaudio)

	// Выполняем команду
	// cmd2 := exec.Command("ffmpeg", "-i", filepath.Join(rec.DirName, "output.mp4"), "-vn", "-acodec", "copy", filepath.Join(rec.DirName, "output_audio_only.mp3"))
	// output2, err := cmd2.CombinedOutput()
	// cmdaudio := exec.Command("ffmpeg", argsaudio...)
	// outputaudio, err := cmdaudio.CombinedOutput()
	// if err != nil {
	// 	log.Fatalf("Error executing command for audiofile: %v\nOutput: %s", err, outputaudio)
	// }
	// fmt.Println(string(outputaudio))

	cmdaudio := exec.Command("ffmpeg", "-i", filepath.Join(rec.DirName, "record.mp4"), "-vn", "-c:a", "libmp3lame", "-q:a", "2", filepath.Join(rec.DirName, "record.mp3"))
	outputaudio, err := cmdaudio.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command for audiofile: %v\nOutput: %s", err, outputaudio)
	}
	fmt.Println(string(outputaudio))

	// time.Sleep(10 * time.Second)

	// cmd = exec.Command("ffmpeg")

	return nil

}

// Вспомогательная функция для создания контекста с таймаутом для набора действий
// func runWithTimeout(ch chan error, message string, timeout time.Duration, action chromedp.Action, rec *Record) chromedp.ActionFunc {
func runWithTimeout(message string, timeout time.Duration, action chromedp.Action, rec *Record) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		log.Printf("Starting action: %s", message)
		// Создаём новый контекст с таймаутом для определённого набора действий
		newctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err := action.Do(newctx)
		if err != nil {
			// ch <- fmt.Errorf("Error during action: %s, error: %v\n", message, err)
			// close(ch)
			log.Printf("Error during action: %s, error: %v\n", message, err)

			var screenshotBuffer []byte
			chromedp.CaptureScreenshot(&screenshotBuffer).Do(ctx)
			filename := fmt.Sprintf("%d", time.Now().Unix())
			err := utils.SaveScreenshoot(&screenshotBuffer, fmt.Sprintf("%s/%s.png", rec.DirName, filename))
			if err != nil {
				log.Printf("Can`t save screenshot\n")
			} else {
				log.Printf("Saved screenshot:%s\n", filename)
			}

			var htmlContent string
			chromedp.OuterHTML("html", &htmlContent).Do(ctx)
			// domsnapshot.CaptureSnapshot(htmlSnapshot)
			err = utils.SavePage(&htmlContent, fmt.Sprintf("%s/%s.html", rec.DirName, filename))
			if err != nil {
				log.Printf("Can`t save page\n")
			}
		} else {
			log.Printf("Successfully finished action: %s\n", message)
		}
		return err
	}
}

// Функция для ожидания элемента и клика по нему с тайм-аутом
func tryClosePopup(ctx context.Context, selector string, timeout time.Duration) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		// Создаем контекст с тайм-аутом
		cctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Пробуем дождаться видимости элемента в течение заданного времени
		err := chromedp.Run(cctx,
			chromedp.WaitVisible(selector),
			chromedp.Click(selector),
		)

		// Если элемент не найден или произошла ошибка — просто возвращаем nil
		if err != nil && cctx.Err() == context.DeadlineExceeded {
			log.Println("Всплывающее окно не появилось.")
			return nil
		}

		if err != nil {
			log.Printf("Неизвестная ошибка при попытке закрытия попапа, %s\n", err)
			return err
		}

		log.Printf("Closed popup Without camera and microphone\n")

		return nil
	}
}

func listenForNetworkEvent(ctx context.Context, rec *Record) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			log.Printf("* console.%s call:\n", ev.Type)
			for _, arg := range ev.Args {
				log.Printf("%s - %s\n", arg.Type, arg.Value)
				if fmt.Sprintf("%s", arg.Value) == "\"Recording started for track\"" {
					rec.StreamCount++
				}
			}
			// case *network.EventResponseReceived:
			//     resp := ev.Response
			//     if len(resp.Headers) != 0 {
			//         log.Printf("received headers: %s", resp.Headers)
			//     }
		}
		// other needed network Event
	})
}

// Вспомогательная функция для создания контекста с таймаутом для набора действий
// func runWithTimeout(ch chan error, message string, timeout time.Duration, action chromedp.Action, rec *Record) chromedp.ActionFunc {
// func startRecording() chromedp.ActionFunc {
// 	return func(ctx context.Context) error {
// 		// Команда для захвата экрана и звука
// 		cmd = exec.Command("ffmpeg", "-video_size", "1920x1080", "-framerate", "25", "-f", "x11grab", "-i", ":99", "-f", "alsa", "-i", "default", "/records/output.mp4")

// 		// Запуск команды
// 		if err := cmd.Start(); err != nil {
// 			log.Printf("Не удалось начать запись: %v", err)
// 			return err
// 		}

// 		log.Println("Запись началась...")

// 		// time.Sleep(10 * time.Second)
// 		// if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
// 		// 	log.Printf("Не удалось остановить запись: %v", err)
// 		// 	return err
// 		// }
// 		// fmt.Println("Запись завершена")

// 		return nil
// 	}
// }

// func stopRecording(cmd *exec.Cmd) chromedp.ActionFunc {
// 	return func(ctx context.Context) error {
// 		// Отправляем сигнал прерывания (Ctrl+C) для корректного завершения ffmpeg
// 		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
// 			log.Printf("Не удалось остановить запись: %v", err)
// 			return err
// 		}
// 		fmt.Println("Запись завершена")

// 		// // Команда для захвата экрана и звука
// 		// 		cmd := exec.Command("ffmpeg", "-video_size", "1280x720", "-framerate", "25", "-f", "x11grab", "-i", ":0.0", "-f", "pulse", "-ac", "2", "-i", "default", "output.mkv")

// 		// 		// Запуск команды
// 		// 		if err := cmd.Start(); err != nil {
// 		// 			log.Printf("Не удалось начать запись: %v", err)
// 		// 			return err
// 		// 		}
// 		// 		log.Println("Запись началась...")
// 		return nil
// 	}
// }

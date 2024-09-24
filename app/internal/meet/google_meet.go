package meet

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Srgkharkov/recordingmeet/internal/utils"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/cdproto/browser"

	// "github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// RecordGoogleMeet starts the recording process for Google Meet.
func RecordGoogleMeet(ch chan error, rec *Record) {
	defer rec.CloseFile()
	defer NotifyAfterExecution(rec)
	// log.Printf("Archive ID:%s\nConnecting Google Meet: %s\n", rec.ID, rec.Link)
	// ch <- fmt.Sprintf("Archive ID:%s", rec.ID)
	rec.log.Printf("Archive ID:%s\nConnecting Google Meet: %s\n", rec.ID, rec.Link)
	// New creates a new context for use with chromedp. With this context
	// you can use chromedp as you normally would.
	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		cu.WithHeadless(), //Required xvfb

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

	listenForNetworkEvent(ctx, rec)

	if err := chromedp.Run(ctx,

		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath(rec.DirName). //fmt.Sprintf("/home/sergei/recordingmeet/app/")).
			WithEventsEnabled(true),

		chromedp.Navigate(rec.Link),

		runWithTimeout(
			ch,
			"Waiting for body visibility",
			10*time.Second,
			chromedp.WaitVisible(`body`, chromedp.ByQuery),
			rec,
		),

		tryClosePopup(
			ctx,
			buttonSelectorWOCamMic,
			3*time.Second,
		), // Ожидание видимости кнопки

		chromedp.Sleep(100*time.Millisecond),
		// chromedp.Click(buttonSelectorWOCamMic), // Клик по кнопке

		runWithTimeout(
			ch,
			"Waiting for inputSelector visibility",
			10*time.Second,
			chromedp.WaitVisible(inputSelector),
			rec,
		), // Ожидание видимости поля ввода

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Clear(inputSelector),                    // Очистка поля ввода (если нужно)
		chromedp.SendKeys(inputSelector, "BotRecording"), // Ввод текста в поле

		runWithTimeout(
			ch,
			"Waiting for buttonSelectorJoin visibility",
			10*time.Second,
			chromedp.WaitEnabled(buttonSelectorJoin),
			rec,
		), // Ожидание, пока кнопка станет активной

		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(buttonSelectorJoin), // Клик по кнопке

		runWithTimeout(
			ch,
			"Waiting for buttonSelectorEnd visibility",
			60*time.Second,
			chromedp.WaitVisible(buttonSelectorEnd),
			rec,
		),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// ch <- fmt.Sprintf("Connection to the meeting was completed successfully.")
			rec.log.Printf("Connection to the meeting was completed successfully.")
			// ch <- nil
			// close(ch)
			// Ваш JavaScript код для запуска на странице
			return chromedp.Evaluate(utils.Mediarecorderjs, nil).Do(ctx)
		}),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// ch <- fmt.Sprintf("Connection to the meeting was completed successfully.")
			// rec.log.Printf("Connection to the meeting was completed successfully.")
			ch <- nil
			close(ch)
			// Ваш JavaScript код для запуска на странице
			return chromedp.Sleep(10 * time.Second).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				// Чтение значения количества участников
				err := chromedp.Text(participantCountSelector, &participantCount, chromedp.NodeVisible).Do(ctx)
				if err != nil {
					rec.log.Printf("Не удалось получить количество участников: %v", err)
					break
				}

				// Преобразование значения в целое число
				count, err := strconv.Atoi(participantCount)
				if err != nil {
					rec.log.Printf("Не удалось преобразовать количество участников в число: %v", err)
					break
				}

				// Если участников меньше 2, кликаем по кнопке выхода
				if count < 2 {
					rec.log.Printf("Участников меньше 2, покидаем встречу.")
					err := chromedp.Click(buttonSelectorEnd).Do(ctx)
					if err != nil {
						rec.log.Printf("Не удалось кликнуть по кнопке выхода: %v", err)
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
		// ch <- fmt.Sprintf("Error in chromedp.Run: %s", err)
		rec.log.Printf("Error in chromedp.Run: %s", err)
		// close(ch)
		// log.Printf("Error in chromedp.Run: %s", err)
		return
	}

	return

}

// Вспомогательная функция для создания контекста с таймаутом для набора действий
func runWithTimeout(ch chan error, message string, timeout time.Duration, action chromedp.Action, rec *Record) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		rec.log.Printf("Starting action: %s", message)
		// Создаём новый контекст с таймаутом для определённого набора действий
		newctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err := action.Do(newctx)
		if err != nil {
			ch <- fmt.Errorf("Error during action: %s, error: %v\n", message, err)
			close(ch)
			rec.log.Printf("Error during action: %s, error: %v\n", message, err)

			var screenshotBuffer []byte
			chromedp.CaptureScreenshot(&screenshotBuffer).Do(ctx)
			filename := fmt.Sprintf("%d", time.Now().Unix())
			err := utils.SaveScreenshoot(&screenshotBuffer, fmt.Sprintf("%s/%s.png", rec.DirName, filename))
			if err != nil {
				rec.log.Printf("Can`t save screenshot\n")
			} else {
				rec.log.Printf("Saved screenshot:%s\n", filename)
			}

			var htmlContent string
			chromedp.OuterHTML("html", &htmlContent).Do(ctx)
			// domsnapshot.CaptureSnapshot(htmlSnapshot)
			err = utils.SavePage(&htmlContent, fmt.Sprintf("%s/%s.html", rec.DirName, filename))
		} else {
			rec.log.Printf("Successfully finished action: %s\n", message)
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
			rec.log.Printf("* console.%s call:\n", ev.Type)
			for _, arg := range ev.Args {
				rec.log.Printf("%s - %s\n", arg.Type, arg.Value)
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

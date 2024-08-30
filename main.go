package main

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"

	cu "github.com/Davincible/chromedp-undetected"
)

func main() {
	// Путь к папке с расширением
	extensionPath := "./tabCapture-recorder"

	// ID вашего расширения
	// extensionID := "khpccidebljocmbjapckddndjoefhpjo"

	// New creates a new context for use with chromedp. With this context
	// you can use chromedp as you normally would.
	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		// cu.WithHeadless(),
		cu.WithExtensions(extensionPath),
		cu.WithNoSandbox(true),

		// If the webelement is not found within 10 minuties, timeout.
		cu.WithTimeout(10*time.Minute),
	))
	if err != nil {
		panic(err)
	}
	defer cancel()

	// Селектор элемента <button>
	buttonSelector := `button[jsname="IbE0S"]`

	// Селектор поля ввода имени
	inputSelector := `input#c16, input#c17`

	// Селектор кнопки
	buttonSelectorJoin := `button[jsname="Qx7uuf"]`

	// buttonSelectorInExtension := `share-audio-button`

	// // Отправка сообщения расширению 'khpccidebljocmbjapckddndjoefhpjo',
	// extensionID := "khpccidebljocmbjapckddndjoefhpjo" // замените на ваш реальный ID расширения
	// message := fmt.Sprintf(`
	// 	console.log('chromedp.Evaluate');
	// 	// Отправка сообщения в фоновый скрипт для начала захвата
	// 	chrome.runtime.sendMessage('%s', { action: 'startCapture2' }, (response) => {
	// 		if (chrome.runtime.lastError) {
	// 			console.error('Error sending message:', chrome.runtime.lastError);
	// 		} else if (response && response.status === 'success') {
	// 			console.log('Capture started successfully');
	// 		} else {
	// 			console.error('Error starting capture:', response.error);
	// 		}
	// 	});
	// `, extensionID)

	// message := fmt.Sprintf(`document.querySelector('button[id="%s"]').click();`, "button34232432432")
	// message2 := fmt.Sprintf(`document.querySelector('button[id="%s"]').click();`, "button342324324322")
	// message2 := `// Отправка сообщения из popup.js или content script
	// 	chrome.runtime.sendMessage({ greeting: "hello from button click" }, (response) => {
	// 		if (response) {
	// 			console.log('Response from background:', response.farewell);
	// 		} else {
	// 			console.error('No response from background script');
	// 		}
	// 	});
	//  	// chrome.runtime.sendMessage({action: 'startCapture'}, (response) => {
	//  	// 	if (response.status === 'success') {
	//  	// 		console.log('Capture started successfully');
	//  	// 	} else {
	//  	// 		console.error('Error starting capture:', response.error);
	//  	// 	}
	//  	// });
	// 	`

	if err := chromedp.Run(ctx,
		// Check if we pass anti-bot measures.
		chromedp.Navigate("https://meet.google.com/mvg-jivq-qdb"),
		chromedp.WaitVisible(buttonSelector), // Ожидание видимости кнопки
		chromedp.Sleep(200*time.Millisecond),
		chromedp.Click(buttonSelector),      // Клик по кнопке
		chromedp.WaitVisible(inputSelector), // Ожидание видимости поля ввода
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.Clear(inputSelector),            // Очистка поля ввода (если нужно)
		chromedp.SendKeys(inputSelector, "fdg"),  // Ввод текста в поле
		chromedp.WaitEnabled(buttonSelectorJoin), // Ожидание, пока кнопка станет активной
		chromedp.Sleep(200*time.Millisecond),
		chromedp.Click(buttonSelectorJoin), // Клик по кнопке
		chromedp.Sleep(15*time.Second),
		// chromedp.Evaluate(message, nil),
		// chromedp.Sleep(2*time.Second),
		// chromedp.Evaluate(message2, nil),
		// chromedp.Sleep(15*time.Second),
		// chromedp.Sleep(1*time.Second),
		// chromedp.Evaluate(message2, nil),
		// chromedp.Evaluate(`// Отправка сообщения из popup.js или content script
		// 	chrome.runtime.sendMessage({action: 'startCapture'}, (response) => {
		// 		if (response.status === 'success') {
		// 			console.log('Capture started successfully');
		// 		} else {
		// 			console.error('Error starting capture:', response.error);
		// 		}
		// 	});`, nil),
		// chromedp.Navigate(fmt.Sprintf("chrome-extension://%s/home.html", extensionID)),
		// chromedp.Sleep(1*time.Second),
		// chromedp.WaitVisible(buttonSelectorInExtension),
		// chromedp.Click(buttonSelectorInExtension),
		// chromedp.Sleep(100*time.Second),
	); err != nil {
		panic(err)
	}

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

	fmt.Println("Message sent to extension!")

	// // Дайте время для завершения действия после клика
	time.Sleep(50 * time.Second)

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
	time.Sleep(50 * time.Second)

}

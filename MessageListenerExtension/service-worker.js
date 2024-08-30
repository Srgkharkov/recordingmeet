// service-worker.js
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    console.log(sender.tab ? "from a content script:" + sender.tab.url : "from the extension");

    if (request.action === "hello from content script" || request.action === "hello from button click") {
        console.log("Received message:", request.action);
        sendResponse({ success: "goodbye" });
        return true; // Добавляем эту строку для асинхронного ответа
    } else if (request.action === "startRecording") {
        console.log("Received message:", request.action);
        captureTabAudioVideo();
//        chrome.tabCapture.capture({ audio: true, video: true }, (stream) => {
//            if (chrome.runtime.lastError || !stream) {
//                console.error('Ошибка при захвате аудио и видео:', chrome.runtime.lastError);
//                sendResponse({ success: "false" });
//                return;
//            }
//
//            const options = { mimeType: 'video/webm; codecs=vp8,opus' };
//            recorder = new MediaRecorder(stream, options);
//
//            recorder.ondataavailable = (e) => {
//                if (e.data.size > 0) {
//                    chunks.push(e.data);
//                }
//            };
//
//            recorder.onstop = () => {
//                const blob = new Blob(chunks, { type: 'video/webm' });
//                chunks = [];
//                saveToFile(blob, "recording.webm");
//            };
//
//            recorder.start();
        sendResponse({ success: "goodbye2" });
//        });
        return true; // Добавляем эту строку для асинхронного ответа
    } else if (request.action === "stopRecording") {
        console.log("Received message:", request.action);
        stopRecording();
//        if (recorder && recorder.state === "recording") {
//            recorder.stop();
        sendResponse({ success: "goodbye3" });
//        } else {
//            sendResponse({ success: false });
        return true; // Добавляем эту строку для асинхронного ответа
    }
});

let recorder; // Глобальная переменная для доступа к MediaRecorder
let chunks = []; // Глобальная переменная для хранения фрагментов видео

function captureTabAudioVideo() {
    chrome.tabCapture.capture({ audio: true, video: true }, (stream) => {
        if (chrome.runtime.lastError || !stream) {
            console.error('Ошибка при захвате аудио и видео:', chrome.runtime.lastError);
            return;
        }

        const options = { mimeType: 'video/webm; codecs=vp8,opus' }; // Параметры MediaRecorder
        recorder = new MediaRecorder(stream, options);

        recorder.ondataavailable = (e) => {
            if (e.data.size > 0) {
                chunks.push(e.data); // Сохраняем фрагменты данных
            }
        };

        recorder.onstop = () => {
            saveToFile(new Blob(chunks, { type: 'video/webm' }), "recording.webm"); // Сохраняем запись при остановке
            chunks = []; // Очищаем массив фрагментов
        };

        recorder.start(); // Начинаем запись
    });
}

function stopRecording() {
    if (recorder && recorder.state === "recording") {
        recorder.stop(); // Останавливаем запись
    }
}


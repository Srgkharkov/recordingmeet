//home.js
function saveToFile(blob, name) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    document.body.appendChild(a);
    a.style = "display: none";
    a.href = url;
    a.download = name;
    a.click();
    URL.revokeObjectURL(url);
    a.remove();
}

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

let btnStart = document.createElement("button");
btnStart.id = "start-recording-button";

let btnStop = document.createElement("button");
btnStop.id = "stop-recording-button";

document.addEventListener('DOMContentLoaded', function () {
    document.getElementById("start-recording-button").addEventListener("click", function (){
        captureTabAudioVideo(); // Запуск записи
    });

    document.getElementById("stop-recording-button").addEventListener("click", function (){
        stopRecording(); // Остановка записи
    });
});

//document.querySelector('button[id="start-recording-button"]').click();

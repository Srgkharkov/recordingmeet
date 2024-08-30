document.getElementById('start-recording').addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: 'startRecording' }, (response) => {
        if (response.status === 'recording') {
            alert('Запись начата.');
        } else {
            alert('Ошибка при запуске записи: ' + (response.error || 'неизвестная ошибка'));
        }
    });
});

document.getElementById('stop-recording').addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: 'stopRecording' }, (response) => {
        if (response.status === 'stopped') {
            alert('Запись остановлена. Видео сохранено.');
        } else if (response.status === 'not_recording') {
            alert('Запись не велась.');
        } else {
            alert('Ошибка при остановке записи.');
        }
    });
});


// Добавление кнопок управления записью на текущую страницу
const startBtn = document.createElement('button');
startBtn.textContent = 'Start Recording';
startBtn.style.position = 'fixed';
startBtn.style.top = '10px';
startBtn.style.left = '10px';
startBtn.style.zIndex = 1000;
document.body.appendChild(startBtn);

const stopBtn = document.createElement('button');
stopBtn.textContent = 'Stop Recording';
stopBtn.style.position = 'fixed';
stopBtn.style.top = '40px';
stopBtn.style.left = '10px';
stopBtn.style.zIndex = 1000;
document.body.appendChild(stopBtn);

startBtn.addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: "startRecording" }, (response) => {
        if (response.success) {
            console.log('Recording started');
        } else {
            console.error('Failed to start recording');
        }
    });
});

stopBtn.addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: "stopRecording" }, (response) => {
        if (response.success) {
            console.log('Recording stopped');
        } else {
            console.error('Failed to stop recording');
        }
    });
});


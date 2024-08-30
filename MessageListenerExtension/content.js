// content.js

console.log('Content script loaded');

chrome.runtime.sendMessage({ action: "hello from content script" }, (response) => {
    if (response) {
        console.log('Response from background:', response.success);
    } else {
        console.error('No response from background script');
        if (chrome.runtime.lastError) {
            console.error('Runtime error:', chrome.runtime.lastError.message);
        }
    }
});


// Добавление кнопки на веб-страницу
const button = document.createElement('button');
button.textContent = 'Send Message to Background';
button.id = "button34232432432";
button.onclick = () => {
    chrome.runtime.sendMessage({ action: "hello from button click" }, (response) => {
        if (response) {
            console.log('Response from background:', response.success);
        } else {
            console.error('No response from background script');
        }
    });
};
document.body.appendChild(button);

const startBtn = document.createElement('button');
startBtn.textContent = 'Start Recording';
startBtn.style.position = 'fixed';
startBtn.style.top = '10px';
startBtn.style.left = '10px';
startBtn.style.zIndex = 1000;
startBtn.id = "button34232432432";
startBtn.onclick = () => {
    chrome.runtime.sendMessage({ action: "startRecording" }, (response) => {
        if (response) {
            console.log('Response from background:', response.success);
        } else {
            console.error('No response from background script');
        }
    });
};
document.body.appendChild(startBtn);

const stopBtn = document.createElement('button');
stopBtn.textContent = 'Stop Recording';
stopBtn.style.position = 'fixed';
stopBtn.style.top = '40px';
stopBtn.style.left = '10px';
stopBtn.style.zIndex = 1000;
stopBtn.id = "button34232432432stop";
stopBtn.onclick = () => {
    chrome.runtime.sendMessage({ action: "stopRecording" }, (response) => {
        if (response) {
            console.log('Response from background:', response.success);
        } else {
            console.error('No response from background script');
        }
    });
};
document.body.appendChild(stopBtn);

document.getElementById('startButton').addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: 'startCapture' }, (response) => {
        if (response.status === 'started') {
            alert('Recording started!');
        } else {
            alert('Error starting recording');
        }
    });
});

document.getElementById('stopButton').addEventListener('click', () => {
    chrome.runtime.sendMessage({ action: 'stopCapture' }, (response) => {
        if (response.status === 'stopped') {
            alert('Recording stopped and saved!');
        } else {
            alert('Error stopping recording: ' + response.error);
        }
    });
});


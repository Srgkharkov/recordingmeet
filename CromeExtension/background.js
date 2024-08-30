let recorder;
let recordedChunks = [];

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'startRecording') {
        chrome.tabCapture.capture({ video: true, audio: true }, (stream) => {
            if (chrome.runtime.lastError || !stream) {
                console.error('Ошибка при захвате потока:', chrome.runtime.lastError);
                sendResponse({ status: 'error', error: chrome.runtime.lastError.message });
                return;
            }

            recorder = new MediaRecorder(stream);
            recordedChunks = [];

            recorder.ondataavailable = (event) => {
                if (event.data.size > 0) {
                    recordedChunks.push(event.data);
                }
            };

            recorder.onstop = () => {
                const blob = new Blob(recordedChunks, { type: 'video/webm' });
                const url = URL.createObjectURL(blob);
                chrome.downloads.download({
                    url: url,
                    filename: 'recording.webm',
                    saveAs: true
                });
            };

            recorder.start();
            sendResponse({ status: 'recording' });
        });
        return true; // Required to indicate async response
    }

    if (message.action === 'stopRecording') {
        if (recorder) {
            recorder.stop();
            sendResponse({ status: 'stopped' });
        } else {
            sendResponse({ status: 'not_recording' });
        }
        return true;
    }
});


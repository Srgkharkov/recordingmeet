let mediaRecorder;
let recordedChunks = [];

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'startCapture') {
        chrome.tabCapture.capture({ audio: true, video: true }, (stream) => {
            if (!stream) {
                console.error('Ошибка захвата.');
                return;
            }

            recordedChunks = [];
            mediaRecorder = new MediaRecorder(stream, { mimeType: 'video/webm; codecs=vp9' });

            mediaRecorder.ondataavailable = (event) => {
                if (event.data.size > 0) {
                    recordedChunks.push(event.data);
                }
            };

            mediaRecorder.onstop = () => {
                const blob = new Blob(recordedChunks, { type: 'video/webm' });
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'recording.webm';
                a.click();
                URL.revokeObjectURL(url);
            };

            mediaRecorder.start();
            sendResponse({ status: 'recording started' });
        });
        return true; // Асинхронный ответ
    } else if (message.action === 'stopCapture') {
        if (mediaRecorder && mediaRecorder.state !== 'inactive') {
            mediaRecorder.stop();
            sendResponse({ status: 'recording stopped' });
        } else {
            sendResponse({ status: 'error', error: 'Recorder is not active' });
        }
        return true;
    }
});


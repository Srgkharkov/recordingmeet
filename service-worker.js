// service-worker.js

let recorder;
let chunks = [];

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    console.log(sender.tab ? "from a content script:" + sender.tab.url : "from the extension");
    if (request.action === "startRecording") {
//        chrome.tabCapture.capture({ audio: true, video: true }, (stream) => {
//            if (chrome.runtime.lastError || !stream) {
//                console.error('Ошибка при захвате аудио и видео:', chrome.runtime.lastError);
//                sendResponse({ success: false });
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
            sendResponse({ success: "true" });
//        });
    } else if (request.action === "stopRecording") {
//        if (recorder && recorder.state === "recording") {
//            recorder.stop();
            sendResponse({ success: "true" });
//        } else {
//            sendResponse({ success: false });
        }
    }

    return true; // Необходимо для асинхронного ответа
});

function saveToFile(blob, filename) {
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.style.display = "none";
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    URL.revokeObjectURL(url);
    a.remove();
}


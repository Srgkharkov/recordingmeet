chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'startCapture' || message.action === 'stopCapture') {
        // Проверяем, если sender.tab не определен (например, при вызове из popup)
        if (sender.tab && sender.tab.id) {
            chrome.tabs.sendMessage(sender.tab.id, { action: message.action });
            sendResponse({ status: message.action === 'startCapture' ? 'started' : 'stopped' });
        } else {
            // Получаем текущую активную вкладку, если сообщение пришло не из контентного скрипта
            chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
                if (tabs.length > 0) {
                    chrome.tabs.sendMessage(tabs[0].id, { action: message.action });
                    sendResponse({ status: message.action === 'startCapture' ? 'started' : 'stopped' });
                } else {
                    sendResponse({ status: 'error', error: 'No active tab found' });
                }
            });
        }
        // Указываем, что ответ будет отправлен асинхронно
        return true;
    }
});


// Слушатель сообщений
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    console.log('Received message:', message);

    // Отправляем ответ
    sendResponse({ status: 'received' });

    // Возвращаем true, чтобы оставить канал открытым для асинхронного ответа
    return true;
});


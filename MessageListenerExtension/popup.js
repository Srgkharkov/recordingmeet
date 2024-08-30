document.getElementById('sendMessageButton').addEventListener('click', () => {
    chrome.runtime.sendMessage({action: "startRecording"}, (response) => {
        if (response) {
            console.log('Received response:', response.success);
            alert('Response: ' + response.success);
        } else {
            console.error('No response or error occurred');
        }
    });
});


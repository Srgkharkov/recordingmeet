document.getElementById('sendMessageButton').addEventListener('click', () => {
    chrome.runtime.sendMessage("njipjaigpcnnamffjknniboapafkajnk", {greeting: "hello"}, (response) => {
        if (response) {
            console.log('Received response:', response.farewell);
            alert('Response: ' + response.farewell);
        } else {
            console.error('No response or error occurred');
        }
    });
});


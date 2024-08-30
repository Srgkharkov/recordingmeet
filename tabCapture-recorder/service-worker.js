// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

chrome.action.onClicked.addListener(async (tab) => {
  console.log('Click extension');
  const existingContexts = await chrome.runtime.getContexts({});
  let recording = false;

  const offscreenDocument = existingContexts.find(
    (c) => c.contextType === 'OFFSCREEN_DOCUMENT'
  );

  // If an offscreen document is not already open, create one.
  if (!offscreenDocument) {
    // Create an offscreen document.
    await chrome.offscreen.createDocument({
      url: 'offscreen.html',
      reasons: ['USER_MEDIA'],
      justification: 'Recording from chrome.tabCapture API'
    });
  } else {
    recording = offscreenDocument.documentUrl.endsWith('#recording');
  }

  if (recording) {
    chrome.runtime.sendMessage({
      type: 'stop-recording',
      target: 'offscreen'
    });
    chrome.action.setIcon({ path: 'icons/not-recording.png' });
    return;
  }

  console.log('tab:', tab.id);

  // Get a MediaStream for the active tab.
  const streamId = await chrome.tabCapture.getMediaStreamId({
    targetTabId: tab.id
  });

  // Send the stream ID to the offscreen document to start recording.
  chrome.runtime.sendMessage({
    type: 'start-recording',
    target: 'offscreen',
    data: streamId
  });

  chrome.action.setIcon({ path: '/icons/recording.png' });
});

// chrome.commands.onCommand.addListener((command) => {
//   if (command === 'toggle-recording') {
//     console.log('Toggle recording command triggered');
//     // Ваш код для обработки команды
//   }
// });
chrome.commands.onCommand.addListener(async (command, tab) => {
  console.log('HoyKeys');
  if (command === 'toggle-recording') {
    console.log('Toggle recording command triggered');
    // Ваш код для обработки команды
    const existingContexts = await chrome.runtime.getContexts({});
    let recording = false;
  
    const offscreenDocument = existingContexts.find(
      (c) => c.contextType === 'OFFSCREEN_DOCUMENT'
    );
  
    // If an offscreen document is not already open, create one.
    if (!offscreenDocument) {
      // Create an offscreen document.
      await chrome.offscreen.createDocument({
        url: 'offscreen.html',
        reasons: ['USER_MEDIA'],
        justification: 'Recording from chrome.tabCapture API'
      });
    } else {
      recording = offscreenDocument.documentUrl.endsWith('#recording');
    }
  
    if (recording) {
      chrome.runtime.sendMessage({
        type: 'stop-recording',
        target: 'offscreen'
      });
      chrome.action.setIcon({ path: 'icons/not-recording.png' });
      return;
    }
  
    console.log('tab:', tab.id);
  
    // Get a MediaStream for the active tab.
    const streamId = await chrome.tabCapture.getMediaStreamId({
      targetTabId: tab.id
    });
  
    // Send the stream ID to the offscreen document to start recording.
    chrome.runtime.sendMessage({
      type: 'start-recording',
      target: 'offscreen',
      data: streamId
    });
  
    chrome.action.setIcon({ path: '/icons/recording.png' });
    }
});

// service-worker.js или background.js

chrome.runtime.onMessageExternal.addListener(async (request, sender, sendResponse) => {
  console.log('Received request:', request);
  if (request.type === 'toggle-recording') {
      const existingContexts = await chrome.runtime.getContexts({});
      let recording = false;

      const offscreenDocument = existingContexts.find(
          (c) => c.contextType === 'OFFSCREEN_DOCUMENT'
      );

      // If an offscreen document is not already open, create one.
      if (!offscreenDocument) {
          // Create an offscreen document.
          await chrome.offscreen.createDocument({
              url: 'offscreen.html',
              reasons: ['USER_MEDIA'],
              justification: 'Recording from chrome.tabCapture API'
          });
      } else {
          recording = offscreenDocument.documentUrl.endsWith('#recording');
      }

      if (recording) {
          chrome.runtime.sendMessage({
              type: 'stop-recording',
              target: 'offscreen'
          });
          chrome.action.setIcon({ path: 'icons/not-recording.png' });
          sendResponse({ status: 'stopped' });
          return;
      }

      console.log('tab:', sender.tab.id);

      // Get a MediaStream for the active tab.
      // const streamId = await chrome.tabCapture.getMediaStreamId();
      const streamId = await chrome.tabCapture.getMediaStreamId({
        targetTabId: sender.tab.id
      });

      // Send the stream ID to the offscreen document to start recording.
      chrome.runtime.sendMessage({
          type: 'start-recording',
          target: 'offscreen',
          data: streamId
      });

      chrome.action.setIcon({ path: '/icons/recording.png' });
      sendResponse({ status: 'started' });
  }
  return true; // Keep the message channel open for asynchronous response
});

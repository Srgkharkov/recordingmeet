package utils

var Mediarecorderjs string = `
// Объект для хранения активных треков и их записей
const activeMediaRecorders = {};

// Массив для хранения информации о записях
const timeline = [];

// Функция для записи нового трека и добавления информации в таймлайн
function startRecordingTrack(track, index) {
    const stream = new MediaStream([track]);
    const mediaRecorder = new MediaRecorder(stream);
    const chunks = [];
    const startTime = new Date().toISOString();

    mediaRecorder.ondataavailable = event => {
        if (event.data.size > 0) {
            chunks.push(event.data);
        }
    };

    mediaRecorder.onstop = () => {
        const endTime = new Date().toISOString();
        const blob = new Blob(chunks, { type: 'video/webm' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        // a.download = 'stream-' + (index + 1) + '.webm';
        a.download = track.id + '.webm';
        a.click();

        // Добавление информации о записи в таймлайн
        timeline.push({
            type: track.kind,
            label: track.label,
            trackId: track.id,
            fileName: track.id + '.webm',//'stream-' + (index + 1) + '.webm',
            startTime: startTime,
            endTime: endTime
        });

        // Сохранение таймлайна в файл JSON
        saveTimelineToFile();
    };

    mediaRecorder.start();
    activeMediaRecorders[track.id] = mediaRecorder;
    console.log('Recording started for track', track.id);
}

// Функция для сохранения таймлайна в JSON файл
function saveTimelineToFile() {
    const blob = new Blob([JSON.stringify(timeline, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'timeline.json';
    a.click();
    console.log('Timeline saved to timeline.json');
}

// ... остальная часть кода ...
// // Функция для записи нового потока
// function startRecordingTrack(track, index) {
// 	const stream = new MediaStream([track]);
// 	const mediaRecorder = new MediaRecorder(stream);
// 	const chunks = [];

// 	mediaRecorder.ondataavailable = event => {
// 		if (event.data.size > 0) {
// 			chunks.push(event.data);
// 		}
// 	};

// 	mediaRecorder.onstop = () => {
// 		const blob = new Blob(chunks, { type: 'video/webm' });
// 		const url = URL.createObjectURL(blob);
// 		const a = document.createElement('a');
// 		a.href = url;
// 		a.download = 'stream-' + (index + 1) + '.webm';
// 		a.click();
// 		console.log('Recording stopped for track', track.id);
// 		delete activeMediaRecorders[track.id];
// 	};

// 	mediaRecorder.start();
// 	activeMediaRecorders[track.id] = mediaRecorder;
// 	console.log('Recording started for track', track.id);
// }

// Функция для обработки медиа-элементов
function handleMediaElements(mediaElements) {
	mediaElements.forEach((element, index) => {
		if (element.srcObject) {
			const tracks = element.srcObject.getTracks();
			tracks.forEach(track => {
				if (!activeMediaRecorders[track.id]) {
					startRecordingTrack(track, index);
				}
			});

			// Остановка записи, если поток прерывается
			element.srcObject.onremovetrack = (event) => {
				const track = event.track;
				if (activeMediaRecorders[track.id]) {
					activeMediaRecorders[track.id].stop();
					console.log('Track', track.id, 'removed and recording stopped');
				}
			};
		}
	});
}

// Функция для наблюдения за изменениями DOM
function observeMediaChanges() {
	const observer = new MutationObserver(mutations => {
		mutations.forEach(mutation => {
			const addedNodes = Array.from(mutation.addedNodes).filter(node => node.tagName === 'VIDEO' || node.tagName === 'AUDIO');
			if (addedNodes.length > 0) {
				handleMediaElements(addedNodes);
			}
		});
	});

	observer.observe(document.body, {
		childList: true,
		subtree: true
	});
}

// Начальная инициализация для уже существующих медиа-элементов
const initialMediaElements = Array.from(document.querySelectorAll('video, audio'));
handleMediaElements(initialMediaElements);

// Запускаем наблюдателя за изменениями DOM
observeMediaChanges();
`

package handlers

import (
	"fmt"
	"net/http"

	"github.com/Srgkharkov/recordingmeet/internal/meet"
	// "github.com/Srgkharkov/recordingmeet/internal/utils"
)

func HandleRecordRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	link := r.URL.Query().Get("link")
	if link == "" {
		http.Error(w, "Link parameter is required", http.StatusBadRequest)
		return
	}

	ms, err := meet.ParseLinkAndCreateDir(link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error processing link: %v", err)
		return
	}

	switch ms.ShortName {
	case "GM":
		ch := make(chan string)
		go meet.RecordGoogleMeet(ch, ms)
		for msg := range ch {
			fmt.Fprintf(w, "%s\n", msg)
			flusher.Flush()
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported platform: %s\n", ms.ShortName)
	}
}

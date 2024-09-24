package meet

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Record represents a meeting service.
type Record struct {
	Service      string `json:"service,omitempty"`
	ID           string `json:"id,omitempty"`
	Link         string `json:"link,omitempty"`
	Status       string `json:"status,omitempty"`
	LinkDownload string `json:"linkdownload,omitempty"`
	LinkLog      string `json:"linklog,omitempty"`
	StreamCount  int    `json:"streamcount,omitempty"`
	DirName      string
	log          *log.Logger
	file         *os.File
}

// Close закрывает файл, если он был открыт.
func (r *Record) CloseFile() error {
	if r.file != nil {
		err := r.file.Close()
		if err != nil {
			return err
		}
		r.file = nil // Обнуляем указатель на файл после закрытия
	}
	return nil
}

// NewRecordByLink parses the link and creates a directory for recordings.
func NewRecordByLink(link string) (*Record, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	host := parsedURL.Host
	pathURL := parsedURL.Path

	workDirName, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var record Record

	switch {
	case host == "meet.google.com":
		meetID := strings.TrimPrefix(pathURL, "/")
		record = Record{
			Service: "GM",
			ID:      fmt.Sprintf("%s_%s_%d", "GM", meetID, time.Now().Unix()),
			Link:    link,
		}

		// path.Split(pathURL)
	// case strings.Contains(host, "zoom.us"):
	// 	service = Record{
	// 		ShortName: "ZOOM",
	// 		ID:        strings.Split(pathURL, "/")[1],
	// 		Link:      link,
	// 	}
	default:
		return nil, fmt.Errorf("unknown service")
	}

	// service.DirName = fmt.Sprintf("%s_%s_%d", service.ShortName, service.ID, time.Now().Unix())
	record.DirName = path.Join(workDirName, "records", record.ID)
	record.LinkDownload = fmt.Sprintf("/download?recordsid=%s", record.ID)
	record.LinkLog = fmt.Sprintf("/log?recordsid=%s", record.ID)

	if err := os.Mkdir(record.DirName, 0755); err != nil {
		return nil, err
	}

	record.file, err = os.Create(path.Join(record.DirName, "log.log"))
	record.log = log.New(record.file, "INFO\t", log.Ldate|log.Ltime)

	return &record, nil
}

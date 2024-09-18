package meet

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// MeetService represents a meeting service.
type MeetService struct {
	ShortName     string
	ID            string
	Link          string
	ParentDirName string
}

// ParseLinkAndCreateDir parses the link and creates a directory for recordings.
func ParseLinkAndCreateDir(link string) (*MeetService, error) {
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

	var service MeetService

	switch {
	case host == "meet.google.com":
		meetID := strings.TrimPrefix(pathURL, "/")
		service = MeetService{
			ShortName:     "GM",
			ID:            fmt.Sprintf("%s_%s_%d", "GM", meetID, time.Now().Unix()),
			Link:          link,
			ParentDirName: path.Join(workDirName, "records"),
		}
		path.Split(pathURL)
	// case strings.Contains(host, "zoom.us"):
	// 	service = MeetService{
	// 		ShortName: "ZOOM",
	// 		ID:        strings.Split(pathURL, "/")[1],
	// 		Link:      link,
	// 	}
	default:
		return nil, fmt.Errorf("unknown service")
	}

	// service.DirName = fmt.Sprintf("%s_%s_%d", service.ShortName, service.ID, time.Now().Unix())
	recordDir := path.Join(service.ParentDirName, service.ID)

	if err := os.Mkdir(recordDir, 0755); err != nil {
		return nil, err
	}

	return &service, nil
}

package downloader

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

var Status map[string]*DownloadStatus

func ListSites() (time.Time, map[string]string, error) {
	data := map[string]string{"stackoverflow.com": ""}

	resp, err := soup.Get("https://archive.org/details/stackexchange")
	if err != nil {
		return time.Now(), data, err
	}

	doc := soup.HTMLParse(resp)

	date, err := time.Parse("2006-01-02", doc.FindAll("dd")[1].Find("span").Text())
	if err != nil {
		return time.Now(), data, err
	}

	links := doc.Find("div", "id", "quickdown1").FindAll("a", "class", "download-pill")
	for _, link := range links {
		name := strings.TrimSpace(link.Text())
		if strings.Index(name, "stackoverflow.com-") > -1 {
			data["stackoverflow.com"] += link.Attrs()["href"] + ";"
		} else {
			data[name[0:len(name)-3]] = "https://archive.org" + link.Attrs()["href"]
		}
	}
	return date, data, nil
}

type DownloadStatus struct {
	TotalFiles                int    `json:"totalFiles"`
	FinishedFiles             int    `json:"finishedFiles"`
	CurrentFile               string `json:"currentFile"`
	CurrentFileDownloadedSize int    `json:"currentFileDownloadedSize"`
	CurrentFileTotalSize      int    `json:"currentFileTotalSize"`
}

func GetDownloadStatus(name string) (*DownloadStatus, error) {
	st, ok := Status[name]
	if !ok {
		return nil, errors.New("site not found")
	}
	file, err := os.Stat(st.CurrentFile)
	if err != nil {
		return nil, err
	}
	st.CurrentFileDownloadedSize = int(file.Size())
	return st, nil
}

func DownloadSite(name string, onFinishCallback func(string)) {
	if Status == nil {
		Status = map[string]*DownloadStatus{}
	}
	st := DownloadStatus{}
	Status[name] = &st

	_, sites, err := ListSites()
	if err != nil {
		log.Println(err)
		return
	}

	data, ok := sites[name]
	if !ok {
		log.Printf("unknown site: %s\n", name)
		return
	}

	fmt.Printf("downloading files for %s\n", name)

	lst := strings.Split(data, ";")
	Status[name].TotalFiles = len(lst)
	Status[name].FinishedFiles = 0

	for _, url := range lst {
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
		defer resp.Body.Close()

		p := "data/" + path.Base(url)

		fmt.Printf("downloading %s ... \n", p)

		Status[name].CurrentFile = p
		Status[name].CurrentFileTotalSize = size

		// Create the file
		out, err := os.Create(p)
		if err != nil {
			log.Println(err)
			return
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)

		Status[name].FinishedFiles++
	}

	onFinishCallback(name)
}

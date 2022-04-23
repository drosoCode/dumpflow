package importer

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/drosocode/dumpflow/internal/config"
	"github.com/drosocode/dumpflow/internal/database"
)

var Status map[string]*ImportStatus

type ImportStatusItem struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type ImportStatus struct {
	Badges      ImportStatusItem `json:"badges"`
	Comments    ImportStatusItem `json:"comments"`
	PostHistory ImportStatusItem `json:"postHistory"`
	PostLinks   ImportStatusItem `json:"postLinks"`
	Posts       ImportStatusItem `json:"posts"`
	Tags        ImportStatusItem `json:"tags"`
	Users       ImportStatusItem `json:"users"`
	Votes       ImportStatusItem `json:"votes"`
}

func ListPaths() []string {
	ret := []string{}
	filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if (strings.Index(info.Name(), "stackoverflow.com") > -1 || strings.Index(info.Name(), "stackexchange.com") > -1) &&
			((!info.IsDir() && filepath.Ext(info.Name()) == ".7z") ||
				(info.IsDir() && verifXmlFiles(path) == nil)) {
			ret = append(ret, path)
		}
		return nil
	})

	return ret
}

func getSlug(path string) string {
	name := filepath.Base(path)
	if strings.Index(name, "stackoverflow.com-") > -1 {
		return "stackoverflow"
	} else {
		return strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(name, ".stackexchange.com.7z", ""),
				".stackoverflow.com.7z", ""),
			".", "")
	}
}

func verifXmlFiles(path string) error {
	files := []string{"Badges.xml", "Comments.xml", "PostHistory.xml", "PostLinks.xml", "Posts.xml", "Tags.xml", "Users.xml", "Votes.xml"}

	// check that all required files exists
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(path, f)); errors.Is(err, os.ErrNotExist) {
			return errors.New("file " + f + " is missing in " + path)
		}
	}
	return nil
}

func ImportFromPath(providedPath string, onFinishCallback func(string)) error {
	name := getSlug(providedPath)

	if Status == nil {
		Status = map[string]*ImportStatus{}
	}

	var wg sync.WaitGroup
	wg.Add(8)

	st := ImportStatus{}
	Status[providedPath] = &st

	path := providedPath
	if filepath.Ext(providedPath) == ".7z" {
		path = providedPath[0 : len(providedPath)-3]
		cmd := exec.Command("7z", "x", providedPath, "-o"+path, "-aoa")
		data, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New("7zip error: " + string(data))
		}
	}

	if err := verifXmlFiles(path); err != nil {
		return err
	}

	db, err := database.GetDB(name)
	if err != nil {
		return err
	}

	files := []string{"Badges.xml", "Comments.xml", "PostHistory.xml", "PostLinks.xml", "Posts.xml", "Tags.xml", "Users.xml", "Votes.xml"}

	// import all files
	for _, f := range files {
		fi, err := os.Stat(filepath.Join(path, f))
		if err != nil {
			return err
		}
		size := fi.Size()

		reader, err := os.Open(filepath.Join(path, f))
		if err != nil {
			return err
		}
		go importFile(db, reader, f, size, &st, &wg)
	}

	go onFinish(name, path, providedPath, &wg, onFinishCallback)

	return nil
}

func onFinish(name string, path string, providedPath string, wg *sync.WaitGroup, onFinishCallback func(string)) {
	wg.Wait()

	log.Printf("import finished for %s", name)

	// add the new site to db
	err := config.DB.AddSite(context.Background(), config.AddSiteParams{DbName: name, Link: getLinkFromPath(path), UpdateDate: time.Now(), AutoUpdate: false, Enabled: true})
	if err != nil {
		log.Println(err)
	}

	// remove files
	os.RemoveAll(path)
	if path != providedPath {
		os.Remove(providedPath)
	}

	onFinishCallback(providedPath)
}

func getLinkFromPath(path string) string {
	if strings.Index(path, "stackoverflow.com-") > -1 {
		return "stackoverflow.com"
	} else {
		return filepath.Base(path)
	}
}

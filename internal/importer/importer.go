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

	"github.com/drosocode/dumpflow/cmd"
	"github.com/drosocode/dumpflow/internal/config"
	"github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/utils"
)

var Status map[string]*ImportStatus

type ImportStatusItem struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type ImportStatus struct {
	Unzipping   ImportStatusItem `json:"unzipping"`
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
		if strings.Index(info.Name(), "stackoverflow.com-") > -1 {
			p := filepath.Join(filepath.Dir(path), "stackoverflow.com")

			if !info.IsDir() && filepath.Ext(info.Name()) == ".7z" {
				p += ".7z"
			}

			if !utils.Contains(ret, p) {
				ret = append(ret, p)
			}

		} else if (strings.Index(info.Name(), "stackoverflow.com") > -1 || strings.Index(info.Name(), "stackexchange.com") > -1) &&
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
	return strings.ReplaceAll(strings.Replace(name, ".7z", "", 1), ".", "_")
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
	base := filepath.Base(providedPath)

	if Status == nil {
		Status = map[string]*ImportStatus{}
	}

	var wg sync.WaitGroup
	wg.Add(8)

	st := ImportStatus{
		Badges:      ImportStatusItem{Current: 0, Total: 1},
		Comments:    ImportStatusItem{Current: 0, Total: 1},
		PostHistory: ImportStatusItem{Current: 0, Total: 1},
		PostLinks:   ImportStatusItem{Current: 0, Total: 1},
		Posts:       ImportStatusItem{Current: 0, Total: 1},
		Tags:        ImportStatusItem{Current: 0, Total: 1},
		Users:       ImportStatusItem{Current: 0, Total: 1},
		Votes:       ImportStatusItem{Current: 0, Total: 1},
		Unzipping:   ImportStatusItem{Current: 0, Total: 1},
	}
	Status[providedPath] = &st

	files := []string{"Badges.xml", "Comments.xml", "PostHistory.xml", "PostLinks.xml", "Posts.xml", "Tags.xml", "Users.xml", "Votes.xml"}
	path := providedPath

	if base == "stackoverflow.com.7z" {
		path = providedPath[0 : len(providedPath)-3]
		st.Unzipping.Total = 8
		for i := range files {
			p := path + "-" + files[i][0:len(files[i])-4] + ".7z"
			log.Println(p)

			cmd := exec.Command("7z", "x", p, "-o"+path, "-aoa")
			data, err := cmd.CombinedOutput()
			if err != nil {
				log.Println("7zip error: " + string(data))
			}
			st.Unzipping.Current++
		}
	} else if filepath.Ext(base) == ".7z" {
		if filepath.Ext(providedPath) == ".7z" {
			path = providedPath[0 : len(providedPath)-3]
			cmd := exec.Command("7z", "x", providedPath, "-o"+path, "-aoa")
			data, err := cmd.CombinedOutput()
			if err != nil {
				return errors.New("7zip error: " + string(data))
			}
		}
		st.Unzipping.Current++
	} else {
		st.Unzipping.Current = 1
	}

	if err := verifXmlFiles(path); err != nil {
		return err
	}

	db, err := database.GetDB(name)
	if err != nil {
		return err
	}

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
	if cmd.Config.DeleteOnFinish {
		os.RemoveAll(path)
		if path != providedPath {
			os.Remove(providedPath)
		}
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

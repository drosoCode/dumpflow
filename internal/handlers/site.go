package handlers

import (
	"context"
	"net/http"
	"time"

	config "github.com/drosocode/dumpflow/internal/config"
	database "github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/downloader"
	"github.com/drosocode/dumpflow/internal/importer"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
)

var status map[string]bool

var onFinish = func(name string) {
	delete(status, name)
}

// handle sites
func SetupSite(r *chi.Mux) {
	site := chi.NewRouter()
	r.Mount("/site", site)

	site.Get("/", ListSites())
	site.Get("/status", GetStatus())
	site.Get("/{name}", GetSite())
	site.Delete("/{name}", RemoveSite())

	site.Get("/download", ListDownloadSites())
	site.Post("/download", DownloadSite())
	site.Get("/download/status", DownloadSiteStatus())

	site.Get("/import", ListImportSites())
	site.Post("/import", ImportSite())
	site.Get("/import/status", ImportSiteStatus())

	status = map[string]bool{}
}

// GET site/{name}
func GetSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := config.DB.GetSite(context.Background(), chi.URLParam(r, "name"))
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// DELETE site/{name}
func RemoveSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		data, err := config.DB.GetSite(context.Background(), name)
		if srv.IfError(w, r, err) {
			return
		}

		err = database.DeleteDB(data.DbName)
		if srv.IfError(w, r, err) {
			return
		}

		err = config.DB.RemoveSite(context.Background(), name)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, "ok")
	}
}

// GET site
func ListSites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := config.DB.ListSites(context.Background())
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET site/download
func ListDownloadSites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type DownloadSitesList struct {
			Date  time.Time         `json:"date"`
			Sites map[string]string `json:"sites"`
		}

		date, lst, err := downloader.ListSites()
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, DownloadSitesList{Date: date, Sites: lst})
	}
}

// POST site/download
func DownloadSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		site := r.URL.Query().Get("site")
		status[site] = false
		go downloader.DownloadSite(site, onFinish)
		srv.JSON(w, r, 200, "ok")
	}
}

// GET site/download/status
func DownloadSiteStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := downloader.GetDownloadStatus(r.URL.Query().Get("site"))
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, status)
	}
}

// POST site/import
func ImportSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		status[path] = true
		err := importer.ImportFromPath(path, onFinish)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, "ok")
	}
}

// GET site/import/status
func ImportSiteStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, ok := importer.Status[r.URL.Query().Get("path")]
		if !ok {
			srv.Error(w, r, 400, "site not found")
			return
		}
		srv.JSON(w, r, 200, status)
	}
}

// GET site/import
func ListImportSites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.JSON(w, r, 200, importer.ListPaths())
	}
}

// GET site/status
func GetStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.JSON(w, r, 200, status)
	}
}

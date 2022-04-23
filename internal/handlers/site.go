package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	config "github.com/drosocode/dumpflow/internal/config"
	"github.com/drosocode/dumpflow/internal/downloader"
	"github.com/drosocode/dumpflow/internal/importer"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
)

// handle sites
func SetupSite(r *chi.Mux) {
	site := chi.NewRouter()
	r.Mount("/site", site)

	site.Get("/", ListSites())
	site.Get("/{id}", GetSite())
	site.Get("/db/{name}", GetSiteFromDB())

	site.Get("/download", ListDownloadSites())
	site.Post("/download", DownloadSite())
	site.Get("/download/status", DownloadSiteStatus())

	site.Get("/import", ListImportSites())
	site.Post("/import", ImportSite())
	site.Get("/import/status", ImportSiteStatus())
}

// GET site/{id}
func GetSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := config.DB.GetSite(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
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

// GET site/db/{name}
func GetSiteFromDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := config.DB.GetSiteFromDB(context.Background(), chi.URLParam(r, "name"))
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
		go downloader.DownloadSite(r.URL.Query().Get("site"))
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
		err := importer.ImportFromPath(path)
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

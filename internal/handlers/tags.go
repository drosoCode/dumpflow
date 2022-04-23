package handlers

import (
	"context"
	"net/http"
	"strconv"

	database "github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
)

// handle tags
func SetupTags(r *chi.Mux) {
	post := chi.NewRouter()
	r.Mount("/tag", post)

	post.Get("/", ListTags())
	post.Get("/{id}", GetTag())
	post.Get("/name/{name}", GetTagFromName())
}

// GET tags/{id}
func GetTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.GetTag(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET tags/name/{name}
func GetTagFromName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.GetTagFromName(context.Background(), chi.URLParam(r, "name"))
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET tag
func ListTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListTags(context.Background())
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

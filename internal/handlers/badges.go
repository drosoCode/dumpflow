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
func SetupBadges(r *chi.Mux) {
	post := chi.NewRouter()
	r.Mount("/badge", post)

	post.Get("/", ListBadges())
	post.Get("/{id}", GetBadge())
}

// GET badge/{id}
func GetBadge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.GetBadge(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET badge/
func ListBadges() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListBadges(context.Background())
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

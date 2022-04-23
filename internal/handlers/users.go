package handlers

import (
	"context"
	"net/http"
	"strconv"

	database "github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
)

// handle users
func SetupUsers(r *chi.Mux) {
	post := chi.NewRouter()
	r.Mount("/user", post)

	post.Get("/{id}", GetUser())
	post.Get("/{id}/badges", ListBadgesFromUser())
}

// GET user/{id}
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.GetUser(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET user/{id}/badges
func ListBadgesFromUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListBadgesFromUser(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

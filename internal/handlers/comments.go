package handlers

import (
	"context"
	"net/http"
	"strconv"

	database "github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
)

// handle comments
func SetupComments(r *chi.Mux) {
	post := chi.NewRouter()
	r.Mount("/comment", post)

	post.Get("/{id}", GetComment())
}

// GET comment/{id}
func GetComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.GetComment(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

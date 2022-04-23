package handlers

import (
	"context"
	"net/http"
	"strconv"

	database "github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
)

// handle posts
func SetupPosts(r *chi.Mux) {
	post := chi.NewRouter()
	r.Mount("/post", post)

	post.Get("/{id}", GetPost())
	post.Get("/{id}/answers", ListAnswers())
	post.Get("/{id}/history", ListHistoryFromPost())
	post.Get("/{id}/related", ListRelatedPosts())
	post.Get("/{id}/votes", ListVotesFromPost())
	post.Get("/{id}/comments", ListCommentsFromPost())
	post.Get("/{id}/users", ListUsersFromPost())
}

// GET post/{id}
func GetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.GetPost(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET post/{id}/answers
func ListAnswers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListAnswers(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET post/{id}/history
func ListHistoryFromPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListHistoryFromPost(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET post/{id}/related
func ListRelatedPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		l := r.URL.Query().Get("limit")
		limit := int32(10)
		if l != "" {
			lim, err := strconv.ParseInt(l, 10, 32)
			limit = int32(lim)
			if srv.IfError(w, r, err) {
				return
			}
		}

		data, err := db.ListRelatedPosts(context.Background(), database.ListRelatedPostsParams{PostID: id, Limit: limit})
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET post/{id}/votes
func ListVotesFromPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListVotesFromPost(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET post/{id}/comments
func ListCommentsFromPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListCommentsFromPost(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

// GET post/{id}/users
func ListUsersFromPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if srv.IfError(w, r, err) {
			return
		}

		data, err := db.ListUsersFromPost(context.Background(), id)
		if srv.IfError(w, r, err) {
			return
		}
		srv.JSON(w, r, 200, data)
	}
}

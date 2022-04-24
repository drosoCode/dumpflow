package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	database "github.com/drosocode/dumpflow/internal/database"
	"github.com/drosocode/dumpflow/internal/utils/srv"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

// handle posts
func SetupPosts(r *chi.Mux) {
	post := chi.NewRouter()
	r.Mount("/post", post)

	post.Post("/search", Search())
	post.Get("/{id}", GetPost())
	post.Get("/{id}/answers", ListAnswers())
	post.Get("/{id}/history", ListHistoryFromPost())
	post.Get("/{id}/related", ListRelatedPosts())
	post.Get("/{id}/votes", ListVotesFromPost())
	post.Get("/{id}/comments", ListCommentsFromPost())
	post.Get("/{id}/users", ListUsersFromPost())

	cachedRequestsSize = 10
	cachedRequests = map[string]CachedSearch{}
}

type SearchData struct {
	Search   string   `json:"search"`
	Username string   `json:"username"`
	Comments bool     `json:"comments"`
	History  bool     `json:"history"`
	Tags     []string `json:"tags"`
	Limit    int      `json:"limit"`
	Start    int      `json:"start"`
}

type SearchRow struct {
	ID   int
	Rank float64
}

type CachedSearch struct {
	Data []int
	Date int64
}

type SearchReturn struct {
	Posts   []int `json:"posts"`
	Results int   `json:"results"`
}

var cachedRequests map[string]CachedSearch
var cachedRequestsSize int

func getHash(hashSettings SearchData) string {
	hashSettings.Limit = 0
	hashSettings.Start = 0
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", hashSettings)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func addToCache(hash string, data []int) {
	oldest := ""
	oldestTime := time.Now().Unix()
	if len(cachedRequests) == cachedRequestsSize {
		for k, v := range cachedRequests {
			if v.Date < oldestTime {
				oldestTime = v.Date
				oldest = k
			}
		}
		delete(cachedRequests, oldest)
	}
	cachedRequests[hash] = CachedSearch{Data: data, Date: time.Now().Unix()}
}

func returnData(data []int, start int, limit int) SearchReturn {
	x := start * limit
	l := len(data)
	if x > l {
		return SearchReturn{Posts: []int{}, Results: l}
	}
	m := x + limit
	if m > l {
		m = l
	}
	return SearchReturn{Posts: data[x:m], Results: l}
}

// POST post/search
func Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := database.GetRawDB(chi.URLParam(r, "site"))
		if srv.IfError(w, r, err) {
			return
		}

		settings := SearchData{}
		err = json.NewDecoder(r.Body).Decode(&settings)
		if srv.IfError(w, r, err) {
			return
		}

		// if limit is not set, default to 10
		if settings.Limit == 0 {
			settings.Limit = 10
		}

		// error if search is empty
		if settings.Search == "" {
			srv.Error(w, r, 400, "You need to specify at least a search string")
			return
		}

		// get the hash of the request
		hash := getHash(settings)
		// if the request is already known, use the existing data and return
		if data, ok := cachedRequests[hash]; ok {
			srv.JSON(w, r, 200, returnData(data.Data, settings.Start, settings.Limit))
			return
		}

		// prepare the tags name with % for the like operator
		for i := range settings.Tags {
			settings.Tags[i] = "%" + settings.Tags[i] + "%"
		}

		// by default do not search in history (3)
		historyLimit := 3
		if settings.History {
			historyLimit = 30
		}

		// SQL Requests
		searchTags := "(SELECT  id, parent_id, body FROM posts WHERE tags LIKE ALL($3))"
		searchUsers := "(SELECT p.id, p.parent_id, p.body FROM posts p, post_history h, users u WHERE h.post_id = p.id AND h.user_id = u.id AND u.display_name LIKE $3)"

		query := `
			SELECT DISTINCT COALESCE(NULLIF(p.parent_id, 0), p.id) AS id, COALESCE(ts_rank_cd(vector1, query),ts_rank_cd(vector2, query)) AS rank 
			FROM posts p, post_history h, to_tsquery('english', $1) query, to_tsvector('english', p.body) vector1, to_tsvector('english', h.text) vector2 
			WHERE p.id = h.post_id AND h.post_history_type_id <= $2 AND (vector1 @@ query OR vector2 @@ query) 
			ORDER BY rank DESC;
		`
		if settings.Comments {
			query = `
				SELECT DISTINCT COALESCE(NULLIF(p.parent_id, 0), p.id) AS id, COALESCE(ts_rank_cd(vector1, query),ts_rank_cd(vector2, query),ts_rank_cd(vector3, query)) AS rank 
				FROM posts p, post_history h, comments c, to_tsquery('english', $1) query, to_tsvector('english', p.body) vector1, to_tsvector('english', h.text) vector2, to_tsvector('english', c.text) vector3 
				WHERE p.id = h.post_id AND COALESCE(NULLIF(p.parent_id, 0), p.id) = c.post_id AND h.post_history_type_id <= $2 AND (vector1 @@ query OR vector2 @@ query OR vector3 @@ query) 
				ORDER BY rank DESC;
			`
		}

		// execute query
		var rows *sql.Rows
		if len(settings.Tags) > 0 && settings.Username != "" {
			query = strings.Replace(query, "FROM posts p", "FROM ("+searchTags+" INTERSECT "+strings.Replace(searchUsers, "$3", "$4", 1)+") AS p", 1)
			rows, err = db.Query(query, settings.Search, historyLimit, pq.Array(settings.Tags), settings.Username)
		} else if len(settings.Tags) > 0 {
			query = strings.Replace(query, "FROM posts p", "FROM "+searchTags+" AS p", 1)
			rows, err = db.Query(query, settings.Search, historyLimit, pq.Array(settings.Tags))
		} else if settings.Username != "" {
			query = strings.Replace(query, "FROM posts p", "FROM "+searchUsers+" AS p", 1)
			rows, err = db.Query(query, settings.Search, historyLimit, settings.Username)
		} else {
			rows, err = db.Query(query, settings.Search, historyLimit)
		}

		if srv.IfError(w, r, err) {
			return
		}

		// gather data
		defer rows.Close()
		var results []int
		for rows.Next() {
			var res SearchRow
			if err := rows.Scan(&res.ID, &res.Rank); err != nil {
				break
			}
			results = append(results, res.ID)
		}

		// cache results
		addToCache(hash, results)

		srv.JSON(w, r, 200, returnData(results, settings.Start, settings.Limit))
	}
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

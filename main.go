package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/drosocode/dumpflow/internal/config"
	"github.com/drosocode/dumpflow/internal/database"
	handler "github.com/drosocode/dumpflow/internal/handlers"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed api static
var embedFS embed.FS

func main() {
	conn := database.DBMSConn{Host: "10.10.2.1", Port: 5432, User: "postgres", Password: "secret", Prefix: "so_"}
	database.ConfigDB(conn)
	config.ConfigDB(conn)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api := chi.NewRouter()
	r.Mount("/api", api)
	site := chi.NewRouter()
	api.Mount("/{site}/", site)
	handler.SetupBadges(site)
	handler.SetupComments(site)
	handler.SetupPosts(site)
	handler.SetupTags(site)
	handler.SetupUsers(site)
	handler.SetupSite(api)

	staticFS := fs.FS(embedFS)
	apiDir, _ := fs.Sub(staticFS, "api")
	handler.ServeStatic(r, "/swagger", http.FS(apiDir))
	staticDir, _ := fs.Sub(staticFS, "static")
	handler.ServeStatic(r, "/", http.FS(staticDir))

	http.ListenAndServe("0.0.0.0:3002", r)
}

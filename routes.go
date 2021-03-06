package chipmunk

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/sessions"
)

var prefix map[string]string

func addRoutes(sm *http.ServeMux, server *Server, staticFiles string) {
	prefix = map[string]string{
		"info":      "/info/",
		"static":    "/static/s/",
		"protected": "/static/",
		"login":     "/api/v0/login/",
		"logout":    "/api/v0/logout/",
		"oauth":     "/api/v0/oauth_cb/",
		"auth":      "/api/v0/auth/",
		"health":    "/healthz",
		"category":  "/api/v0/category/",
		"tranx":     "/api/v0/tranx/",
		"user":      "/api/v0/user/",

		"fake": "/fake/",
	}

	if staticFiles == "" {
		sm.Handle(
			prefix["static"],
			http.FileServer(
				&assetfs.AssetFS{
					Asset:     Asset,
					AssetDir:  AssetDir,
					AssetInfo: AssetInfo,
				},
			),
		)
		sm.HandleFunc(
			"/",
			func(w http.ResponseWriter, req *http.Request) {
				data, err := Asset("static/s/index.html")
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				r := bytes.NewReader(data)
				io.Copy(w, r)
			},
		)
	} else {
		sm.Handle(
			prefix["static"],
			http.StripPrefix(
				prefix["static"],
				http.FileServer(http.Dir(staticFiles)),
			),
		)
		sm.HandleFunc(
			"/",
			func(w http.ResponseWriter, req *http.Request) {
				http.ServeFile(w, req, filepath.Join(staticFiles, "index.html"))
			},
		)
	}

	store = sessions.NewCookieStore([]byte(server.cookieSecret))
	sm.HandleFunc(prefix["tranx"], server.tranx)
	sm.HandleFunc(prefix["fake"], server.fakeSetup)
	sm.HandleFunc(prefix["category"], server.category)
	sm.HandleFunc(prefix["user"], server.user)
	sm.HandleFunc(prefix["protected"], server.plist)
	sm.HandleFunc(prefix["info"], server.serverInfo)
	sm.HandleFunc(prefix["login"], server.login)
	sm.HandleFunc(prefix["logout"], server.logout)
	sm.HandleFunc(prefix["oauth"], server.oauthCallback)
	sm.HandleFunc(prefix["auth"], server.auth)
	sm.HandleFunc(prefix["health"], server.health)
}

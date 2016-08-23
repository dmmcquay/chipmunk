package chipmunk

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var store *sessions.CookieStore

var (
	oauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://127.0.0.1:8080/api/v0/oauth_cb/",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = strconv.Itoa(rand.Int())
)

var Version string = "dev"
var start time.Time
var users []user

type failure struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewFailure(msg string) *failure {
	return &failure{
		Success: false,
		Error:   msg,
	}
}

type Server struct {
	ClientID     string
	ClientSecret string
	CookieSecret string
}

func init() {
	log.SetFlags(log.Ltime)
	start = time.Now()
}

func NewServer(sm *http.ServeMux, clientId, clientSecret, cookieSecret, static string) *Server {
	server := &Server{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		CookieSecret: cookieSecret,
	}
	addRoutes(sm, server, static)
	return server
}

func (s *Server) fakeSetup(w http.ResponseWriter, r *http.Request) {
	u := userInfo{
		Email: "derekmcquay@gmail.com",
	}
	addUser(u)
}

func (s *Server) tranx(w http.ResponseWriter, r *http.Request) {
	//TODO add back in oauth
	//w.Header().Set("Content-Type", "application/json")
	//session, _ := store.Get(r, "creds")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if loggedIn := session.Values["authenticated"]; loggedIn != true {
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	switch r.Method {
	default:
		b, _ := json.Marshal(NewFailure("Allowed method: POST"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	case "GET":
		u, err := getUser("derekmcquay@gmail.com") //TODO will grab this from session
		if err != nil {
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users[u].txs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		u, err := getUser("derekmcquay@gmail.com") //TODO will grab this from session
		if err != nil {
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}

		t := tranx{}
		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		users[u].txs = append(users[u].txs,
			tranx{
				Cost:  t.Cost,
				Store: t.Store,
				Info:  t.Info,
				Month: t.Month,
			},
		)
	}
}

func (s *Server) costPerMonth(w http.ResponseWriter, r *http.Request) {
	//TODO add back in oauth
	//w.Header().Set("Content-Type", "application/json")
	//session, _ := store.Get(r, "creds")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if loggedIn := session.Values["authenticated"]; loggedIn != true {
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	switch r.Method {
	default:
		b, _ := json.Marshal(NewFailure("Allowed method: GET"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	case "GET":
		u, err := getUser("derekmcquay@gmail.com") //TODO will grab this from session
		if err != nil {
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
		monthCost := make(map[time.Month]float32)
		for _, t := range users[u].txs {
			c, ok := monthCost[t.Month]
			if !ok {
				monthCost[t.Month] = t.Cost
				continue
			}
			monthCost[t.Month] = t.Cost + c
		}
		err = json.NewEncoder(w).Encode(monthCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	//TODO add back in oauth
	//w.Header().Set("Content-Type", "application/json")
	//session, _ := store.Get(r, "creds")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if loggedIn := session.Values["authenticated"]; loggedIn != true {
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	switch r.Method {
	default:
		b, _ := json.Marshal(NewFailure("Allowed method: GET"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	case "GET":
		err := json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	oauthConf.ClientID = s.ClientID
	oauthConf.ClientSecret = s.ClientSecret
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) oauthCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)

	email, err := oauthClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Printf("failed with getting userinfo: '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	u := userInfo{}
	err = json.Unmarshal(data, &u)
	if err != nil {
		b, _ := json.Marshal(NewFailure(err.Error()))
		http.Error(w, string(b), http.StatusInternalServerError)
		return
	}

	if authorizedEmail(u.Email) {
		session, err := store.Get(r, "creds")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["authenticated"] = true
		session.Values["uname"] = u.Email
		if err := session.Save(r, w); err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		addUser(u)
		http.Redirect(w, r, "/static/", http.StatusTemporaryRedirect)
		return
	}
	b, _ := json.Marshal(NewFailure("Not a authorized user"))
	http.Error(w, string(b), http.StatusForbidden)
	return
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	output := struct {
		Auth bool `json:"auth"`
	}{
		Auth: false,
	}
	w.Header().Set("Content-Type", "application/json")
	session, err := store.Get(r, "creds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if loggedIn := session.Values["authenticated"]; loggedIn == true {
		output.Auth = true
		json.NewEncoder(w).Encode(output)
		return
	}
	b, _ := json.Marshal(output)
	http.Error(w, string(b), http.StatusUnauthorized)
}

func (s *Server) logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "creds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(session.Values, "authenticated")
	delete(session.Values, "uname")
	session.Save(req, w)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (s *Server) serverInfo(w http.ResponseWriter, req *http.Request) {
	output := struct {
		Version string `json:"version"`
		Start   string `json:"start"`
		Uptime  string `json:"uptime"`
	}{
		Version: Version,
		Start:   start.Format("2006-01-02 15:04:05"),
		Uptime:  time.Since(start).String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (s *Server) plist(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "creds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if loggedIn := session.Values["authenticated"]; loggedIn != true {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := Asset("static/list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	req := bytes.NewReader(data)
	io.Copy(w, req)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

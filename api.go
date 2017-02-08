package chipmunk

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *Server) category(w http.ResponseWriter, req *http.Request) {
	// TODO add back in
	//w.Header().Set("Content-Type", "application/json")
	//session, err := store.Get(r, "creds")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if loggedIn := session.Values["authenticated"]; loggedIn != true {
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	switch req.Method {
	default:
		b, _ := json.Marshal(NewFailure("Allowed methods: GET, POST, DELETE"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	case "GET":
		categories, err := s.db.getCategories()
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(categories)
	case "POST":
		cat := category{}
		err := json.NewDecoder(req.Body).Decode(&cat)
		req.Body.Close()
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		_, err = s.db.db.Exec(
			`INSERT INTO categories (name, budget) VALUES ($1, $2)`,
			cat.Name,
			cat.Budget,
		)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		cat := category{}
		err := json.NewDecoder(req.Body).Decode(&cat)
		req.Body.Close()
		if err != nil {
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		_, err = s.db.db.Exec("DELETE FROM categories WHERE name = $1", cat.Name)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) user(w http.ResponseWriter, req *http.Request) {
	// TODO add back in
	//w.Header().Set("Content-Type", "application/json")
	//session, err := store.Get(r, "creds")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if loggedIn := session.Values["authenticated"]; loggedIn != true {
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	switch req.Method {
	default:
		b, _ := json.Marshal(NewFailure("Allowed methods: GET, POST, DELETE"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	case "GET":
		users, err := s.db.getUsers()
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	case "POST":
		u := user{}
		err := json.NewDecoder(req.Body).Decode(&u)
		req.Body.Close()
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}

		// TODO add back in
		//// verify current user is an admin
		//session, err := store.Get(req, "creds")
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//email := ""
		//if session.Values["uname"] != nil {
		//	email = session.Values["uname"].(string)
		//}
		//if !s.db.adminUser(email) {
		//	log.Printf("user is not admin")
		//	b, _ := json.Marshal(NewFailure("not admin"))
		//	http.Error(w, string(b), http.StatusForbidden)
		//	return
		//}

		_, err = s.db.db.Exec(
			`INSERT INTO users (email, admin) VALUES ($1, $2)`,
			u.Email,
			u.Admin,
		)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
	case "DELETE":
		u := user{}
		err := json.NewDecoder(req.Body).Decode(&u)
		req.Body.Close()
		if err != nil {
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}

		// TODO add back in
		//// verify current user is an admin
		//session, err := store.Get(req, "creds")
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//email := ""
		//if session.Values["uname"] != nil {
		//	email = session.Values["uname"].(string)
		//}
		//if !s.db.adminUser(email) {
		//	log.Printf("user is not admin")
		//	b, _ := json.Marshal(NewFailure("not admin"))
		//	http.Error(w, string(b), http.StatusForbidden)
		//	return
		//}

		_, err = s.db.db.Exec("DELETE FROM users WHERE email = $1", u.Email)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) tranx(w http.ResponseWriter, req *http.Request) {
	// TODO add back in
	//w.Header().Set("Content-Type", "application/json")
	//session, err := store.Get(r, "creds")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if loggedIn := session.Values["authenticated"]; loggedIn != true {
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	switch req.Method {
	default:
		b, _ := json.Marshal(NewFailure("Allowed methods: GET, POST, DELETE"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	case "GET":
		searchreq := req.URL.Path[len(prefix["tranx"]):]
		if len(searchreq) == 0 {
			tranxs, err := s.db.getTranxs()
			if err != nil {
				log.Printf("%+v", err)
				b, _ := json.Marshal(NewFailure(err.Error()))
				http.Error(w, string(b), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(tranxs)
			return
		}
		if searchreq[len(searchreq)-1] != '/' {
			http.Redirect(w, req, prefix["tranx"]+searchreq+"/", http.StatusMovedPermanently)
			return
		}
		searchReqParsed := strings.Split(searchreq, "/")
		t, err := s.db.getTranx(searchReqParsed[0])
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(t)

	case "POST":
		t := tranx{}
		err := json.NewDecoder(req.Body).Decode(&t)
		req.Body.Close()
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		category_id, err := s.db.getCategoryID(t.Category)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		user_id, err := s.db.getUserID(t.User)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		_, err = s.db.db.Exec(
			`INSERT INTO tranx (cost, store, info, date, category_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)`,
			t.Cost,
			t.Store,
			t.Info,
			time.Now(),
			category_id,
			user_id,
		)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	case "DELETE":
		t := tranx{}
		err := json.NewDecoder(req.Body).Decode(&t)
		req.Body.Close()
		if err != nil {
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
		// TODO need to find better way to delete tranx
		_, err = s.db.db.Exec("DELETE FROM tranx WHERE store = $1 AND cost = $2", t.Store, t.Cost)
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
	}
}

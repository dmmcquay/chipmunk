package chipmunk

import (
	"encoding/json"
	"log"
	"net/http"
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
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
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
		if err != nil {
			log.Printf("%+v", err)
			b, _ := json.Marshal(NewFailure(err.Error()))
			http.Error(w, string(b), http.StatusBadRequest)
			return
		}
	}
}

package handler

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./website/pages/index.html")
	if err != nil {
		ErrorHandler(w, r, "500")
		return
	}
	tmpl.Execute(w, nil)
}

// func LogInHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Cache-Control", "no-store")

// 	is, _, err := forum.IsLoggedIn(r, "token")
// 	if err != nil {
// 		jsonResponse(w, http.StatusInternalServerError, "Something wrong")
// 		return
// 	}
// 	if is {
// 		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
// 		return
// 	}

// 	tmpl, err := template.ParseFiles("./website/pages/login.html")
// 	if err != nil {
// 		ErrorHandler(w, r, "500")
// 		return
// 	}
// 	tmpl.Execute(w, nil)
// }

// func Profil(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Cache-Control", "no-store")

// 	is, _, err := forum.IsLoggedIn(r, "token")
// 	if err != nil {
// 		jsonResponse(w, http.StatusInternalServerError, "Something wrong")
// 		return
// 	}
// 	if !is {
// 		http.Redirect(w, r, "/log-in", http.StatusPermanentRedirect)
// 		return
// 	}

// 	user := r.FormValue("username")
// 	if user == "" {
// 		ErrorHandler(w, r, "404")
// 		return
// 	}
// 	tmpl, err := template.ParseFiles("./website/pages/profil.html")
// 	if err != nil {
// 		ErrorHandler(w, r, "500")
// 		return
// 	}
// 	tmpl.Execute(w, nil)
// }

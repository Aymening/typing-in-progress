package handler

import (
	"box/server/forum"
	"database/sql"
	"net/http"
)

type user struct {
	Username   string `json:"username"`
	Created_at string `json:"created_at"`
}

// get name and create date account for desplay in profile
func UserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		username := r.FormValue("username")
		if username == "" {
			jsonResponse(w, http.StatusBadRequest, "username is required")
			return
		}
		row, err := forum.SelectOneRow("SELECT username, created_at FROM user WHERE username = ?", username)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		var userLogged user
		err = row.Scan(&userLogged.Username, &userLogged.Created_at)
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "User not found")
			return
		}
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		jsonResponse(w, http.StatusOK, userLogged)
	} else {
		jsonResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

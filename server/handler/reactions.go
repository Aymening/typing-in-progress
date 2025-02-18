package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"box/server/forum"
)

func GetUserReaction(res http.ResponseWriter, req *http.Request, itemType string) {
	if req.Method != http.MethodGet {
		jsonResponse(res, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	ok, userID, errlogin := forum.IsLoggedIn(req, "token")
	if errlogin != nil{
		jsonResponse(res, http.StatusInternalServerError, "Something wrong")
	}
	if !ok {
		jsonResponse(res, http.StatusUnauthorized, "401 Unauthorized")
		return
	}
	itemID, err1 := strconv.Atoi(req.URL.Query().Get("item_id"))
	statusLike, err2 := strconv.Atoi(req.URL.Query().Get("status_like"))

	if err1 != nil || err2 != nil {
		jsonResponse(res, http.StatusBadRequest, "Missing itemID or action")
		return
	}

	var err error
	if itemType == "post" {
		err = forum.LikePost(itemID, userID, statusLike)
	} else if itemType == "comment" {
		err = forum.LikeComment(itemID, userID, statusLike)
	} else {
		jsonResponse(res, http.StatusBadRequest, "Invalid item type")
		return
	}

	if err != nil {
		jsonResponse(res, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(res, http.StatusOK, "Success")
}

type Like struct {
	ID        int `json:"id"`
	ItemID    int `json:"item_id"` // post_id or comment_id
	LikeValue int `json:"like"`
}

// for like post or comment
func CheckIfUserLike(w http.ResponseWriter, r *http.Request, itemType string) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Check if user is logged in
	ok, userID, errlogin := forum.IsLoggedIn(r, "token")
	if errlogin != nil{
		jsonResponse(w, http.StatusInternalServerError, "Something wrong")
	}
	if !ok {
		jsonResponse(w, http.StatusUnauthorized, "401 Unauthorized")
		return
	}

	var itemID int
	var err error

	if itemType == "post" {
		itemID, err = strconv.Atoi(r.FormValue("post_id"))
	} else if itemType == "comment" {
		itemID, err = strconv.Atoi(r.FormValue("comment_id"))
	} else {
		jsonResponse(w, http.StatusBadRequest, "Invalid item type")
		return
	}

	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Invalid Item ID")
		return
	}

	tableName := ""
	itemIDKey := ""

	if itemType == "post" {
		tableName = "like"
		itemIDKey = "post_id"
	} else if itemType == "comment" {
		tableName = "commentlike"
		itemIDKey = "comment_id"
	}

	query := fmt.Sprintf("SELECT id, %s, like FROM %s WHERE %s = ? AND user_id = ?", itemIDKey, tableName, itemIDKey)
	row, err := forum.SelectOneRow(query, itemID, userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var like Like
	err = row.Scan(&like.ID, &like.ItemID, &like.LikeValue)
	if err != nil {
		jsonResponse(w, http.StatusOK, nil)
		return
	}

	jsonResponse(w, http.StatusOK, like)
}

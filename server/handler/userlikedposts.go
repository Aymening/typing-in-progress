package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"box/server/forum"
)

func Likedposts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	ok, userID, err := forum.IsLoggedIn(r, "token")
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "some thing wrong")
		return
	}
	if !ok {
		jsonResponse(w, http.StatusUnauthorized, "not logged")
		return
	}

	page, pageErr := strconv.Atoi(r.FormValue("page"))

	if pageErr != nil {
		jsonResponse(w, http.StatusBadRequest, "Invalid page number or username")
		return
	}
	query := fmt.Sprintf(`
	SELECT
		p.id,
		(SELECT GROUP_CONCAT(cat.category, ', ') FROM category cat JOIN postcategory AS pc ON pc.category_id = cat.id WHERE pc.post_id = p.id) AS categories, 
		(SELECT COUNT(li.id) FROM like AS li WHERE li.post_id = p.id AND li.like = "1") AS likes, 
		(SELECT COUNT(li.id) FROM like AS li WHERE li.post_id = p.id AND li.like = "-1") AS dislikes, 
		(SELECT COUNT(com.id) FROM comment AS com WHERE com.post_id = p.id) AS comments,
		u.username,
		p.title,
		p.content,
		p.created_at,
		p.image
	FROM
        like AS l 
		JOIN post AS p ON p.id = l.post_id
		JOIN user AS u ON p.user_id = u.id
		LEFT JOIN postcategory AS pc ON pc.post_id = p.id
		LEFT JOIN category AS cat ON pc.category_id = cat.id
	WHERE
		l.user_id = %d AND l.like = 1
	GROUP BY
		p.id
	LIMIT
		7
	OFFSET 
		%d
	`, userID, (page*7)-7)
	posts, err := forum.SelectQuery(query)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, "Not Found")
		return
	}
	var Post FilterPostData
	var likedPosts []FilterPostData
	for posts.Next() {
		posts.Scan(&Post.PostId, &Post.Categories, &Post.Likes, &Post.Dislikes, &Post.Comments, &Post.Username, &Post.Title, &Post.Content, &Post.CreatedAt, &Post.Image)
		likedPosts = append(likedPosts, Post)
	}
	jsonResponse(w, http.StatusOK, likedPosts)
}

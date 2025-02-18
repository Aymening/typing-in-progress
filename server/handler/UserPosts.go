package handler

import (
	"box/server/forum"
	"fmt"
	"net/http"
	"strconv"
)

func UserPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	page, pageErr := strconv.Atoi(r.FormValue("page"))
	user := r.FormValue("username")

	if pageErr != nil || user == "" {
		jsonResponse(w, http.StatusBadRequest, "Invalid page number or username")
		return
	}

	myQuery := fmt.Sprintf(`
	SELECT
		p.id,
		(SELECT GROUP_CONCAT(cat.category, ', ') FROM category cat JOIN postcategory AS pc ON pc.category_id = cat.id WHERE	pc.post_id = p.id) AS categories,
		(SELECT	COUNT(li.id) FROM like AS li WHERE li.post_id = p.id AND li.like = "1") AS likes,
		(SELECT	COUNT(li.id) FROM like AS li WHERE li.post_id = p.id AND li.like = "-1") AS dislikes,
		(SELECT	COUNT(com.id) FROM comment AS com WHERE com.post_id = p.id) AS comments,
		u.username,
		p.title,
		p.content,
		p.created_at,
		p.image
	FROM
		post AS p
		JOIN user AS u ON p.user_id = u.id
		LEFT JOIN postcategory AS pc ON pc.post_id = p.id
		LEFT JOIN category AS cat ON pc.category_id = cat.id
	WHERE
		u.username = %s
	GROUP BY
		p.id
	ORDER BY 
		p.id DESC
	LIMIT
		7
	OFFSET
		%d
	`, user, (page*7)-7)
	posts, err := forum.SelectQuery(myQuery)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, "Not Found")
		return
	}
	var Post FilterPostData
	var allPosts []FilterPostData
	for posts.Next() {
		posts.Scan(&Post.PostId, &Post.Categories, &Post.Likes, &Post.Dislikes, &Post.Comments, &Post.Username, &Post.Title, &Post.Content, &Post.CreatedAt, &Post.Image)
		allPosts = append(allPosts, Post)
	}
	jsonResponse(w, http.StatusOK, allPosts)
}

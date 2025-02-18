package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"box/server/forum"
)

type FilterPostData struct {
	PostId     int    `json:"id"`
	Categories any    `json:"categories"`
	Likes      int    `json:"likes"`
	Dislikes   int    `json:"dislikes"`
	Comments   int    `json:"comments"`
	Username   string `json:"username"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	Image      string `json:"image"`
}

func FilterPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	page, pageErr := strconv.Atoi(r.FormValue("page"))
	categ := r.FormValue("categ")
	date := r.FormValue("date")
	like := r.FormValue("like")

	// check general bad request
	if pageErr != nil || (date != "" && date != "ASC" && date != "DESC") || (like != "" && like != "ASC" && like != "DESC") {
		jsonResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// if user not logged in cant filter with date and like
	ok, _, err := forum.IsLoggedIn(r, "token")
	if err != nil {
		forum.ErrorLog.Println(err)
		jsonResponse(w, http.StatusInternalServerError, "Something wrong")
	}
	if !ok{
		jsonResponse(w, http.StatusUnauthorized, "You must be logged in to filter posts")
		return
	}

	if date == "" && like == "" {
		date = "DESC"
	}

	myQuery := manageFilterQuery(page, categ, date, like)
	posts, err := forum.SelectQuery(myQuery)
	if err != nil {
		forum.ErrorLog.Println(err)
		jsonResponse(w, http.StatusInternalServerError, "Internal Server Error")
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

func manageFilterQuery(page int, cat string, date string, like string) string {
	numPage := strconv.Itoa((page * 7) - 7)
	whereCat := ""
	orderBy := ""

	if cat != "" {
		whereCat = "WHERE cat.category = '" + cat + "'"
	}

	if date == "ASC" || date == "DESC" {
		orderBy = "ORDER BY p.id " + date
	} else if like == "ASC" || like == "DESC" {
		orderBy = "ORDER BY likes " + like
	}

	return fmt.Sprintf(`
	SELECT
		p.id,
		(SELECT GROUP_CONCAT(cat.category, ', ') FROM category cat JOIN postcategory AS pc ON pc.category_id = cat.id WHERE pc.post_id = p.id) AS categories, 
		(SELECT COUNT(li.id) FROM like AS li WHERE li.post_id = p.id AND li.like = "1") AS likes, 
		(SELECT COUNT(li.id) FROM like AS li WHERE li.post_id = p.id AND li.like = "-1") AS dislikes, 
		(SELECT COUNT(com.id) FROM comment AS com WHERE com.post_id = p.id) AS comments,
		u.nickname,
		p.title,
		p.content,
		p.created_at,
		p.image
	FROM
		post AS p
		JOIN user AS u ON p.user_id = u.id
		LEFT JOIN postcategory AS pc ON pc.post_id = p.id
		LEFT JOIN category AS cat ON pc.category_id = cat.id
	%s
	GROUP BY
		p.id
	%s
	LIMIT
		7
	OFFSET
		%s
	`, whereCat, orderBy, numPage)
}

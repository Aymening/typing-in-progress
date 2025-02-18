package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"box/server/forum"
	"box/server/handler"
)

func main() {
	forum.CreateDataBase()
	addr := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/public/", handleStaticFile)

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		HandlerFunc(res, req)
	})

	fmt.Println("Server started...")
	err := http.ListenAndServe(addr, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func handleStaticFile(res http.ResponseWriter, req *http.Request) {
	directories := []string{"img", "js", "pages", "css"}
	for _, dir := range directories {
		if strings.HasPrefix(req.URL.Path, "/public/"+dir) {

			if req.URL.Path == "/public/"+dir || strings.HasSuffix(req.URL.Path, "/") {
				handler.ErrorHandler(res, req, "404")
				return
			}

			filePath := strings.TrimPrefix(req.URL.Path, "/public/"+dir+"/")

			if strings.Contains(filePath, "/") {
				handler.ErrorHandler(res, req, "404")
				return
			}

			file := "./website/" + dir + "/" + filePath

			_, err := os.Stat(file)
			if err != nil {
				handler.ErrorHandler(res, req, "404")
				return
			}

			http.ServeFile(res, req, file)
			return
		}
	}
	handler.ErrorHandler(res, req, "404")
}

func HandlerFunc(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	// handler
	case "/":
		handler.HomeHandler(res, req)

	// api EndPoints
	case "/api/sign-up":
		handler.SignUpHandlerApi(res, req)

	case "/api/log-in":
		handler.LogInHandlerApi(res, req)

	case "/api/logout":
		handler.LogoutHandlerApi(res, req)

	case "/api/is-logged-in":
		handler.IsLoggedInHandler(res, req)

	case "/api/categ":
		handler.Categorie(res, req)

	case "/api/likedPost":
		handler.Likedposts(res, req)

	case "/api/user/profil":
		handler.UserProfile(res, req)

	case "/api/filterPost":
		handler.FilterPost(res, req)

	case "/api/userPosts":
		handler.UserPosts(res, req)

	case "/api/getComments":
		handler.Get_comments(res, req)

	case "/api/addComment":
		handler.Add_comment(res, req)

	case "/api/create-post":
		handler.CreatePostHandlerApi(res, req)

	case "/api/reactions":
		handler.GetUserReaction(res, req, "post")

	case "/api/commentReactions":
		handler.GetUserReaction(res, req, "comment")

	case "/api/checklike":
		handler.CheckIfUserLike(res, req, "post")

	case "/api/check-comment-like":
		handler.CheckIfUserLike(res, req, "comment")

	// web socket
	case "/ws":
		handler.WsEndpoint(res, req)

	case "/api/getUsers":
		handler.GetUsersApi(res, req)

	case "/api/lastUsersChat":
		handler.GetUserOrganizedlastMsgApi(res, req)

	case "/api/getMessages":
		handler.GetMessages(res, req)

	default:
		res.WriteHeader(http.StatusNotFound)
		handler.HomeHandler(res, req)
		return
	}
}

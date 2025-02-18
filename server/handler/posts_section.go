package handler

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"box/server/forum"

	"github.com/gofrs/uuid"
)

// api for create post
func CreatePostHandlerApi(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		jsonResponse(res, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	flag, userID, errlogin := forum.IsLoggedIn(req, "token")
	if errlogin != nil {
		jsonResponse(res, http.StatusInternalServerError, "Something wrong")
		return
	}
	if !flag {
		jsonResponse(res, http.StatusUnauthorized, "you need to log in")
		return
	}

	// start handling image
	imagePublic := ""

	img, _, err := req.FormFile("image")

	if err != nil && err.Error() != "http: no such file" {
		jsonResponse(res, http.StatusInternalServerError, "Error copying img")
		return
	}
	// only if add image
	if err == nil {
		imgBytes := make([]byte, 512)
		_, err = img.Read(imgBytes)
		if err != nil {
			jsonResponse(res, http.StatusBadRequest, "Error reading img")
			return
		}

		imgType := http.DetectContentType(imgBytes)

		// check if avif type
		if len(imgBytes) >= 12 && string(imgBytes[4:12]) == "ftypavif" {
			imgType = "image/avif"
		}

		if imgType != "image/jpeg" && imgType != "image/gif" && imgType != "image/png" && imgType != "image/avif" {
			jsonResponse(res, http.StatusConflict, "Type not supported only : (jpeg, gif, png, avif)")
			return
		}

		fileSize := req.ContentLength

		// 3000 for content title , content , categories
		if fileSize > (20*1024*1024)+3000 {
			jsonResponse(res, http.StatusConflict, "The image is too big, max size is 20 MB")
			return
		}

		// Standard packages have limited image format support only (gif/jpg/png)
		// we skip avif in dimension size check
		if imgType != "image/avif" {
			img.Seek(0, 0)

			config, _, errConf := image.DecodeConfig(img)
			if errConf != nil {
				jsonResponse(res, http.StatusInternalServerError, "Error decoding image")
				return
			}

			width := config.Width
			height := config.Height

			if width < 200 || height < 200 || width*5 < height || height*10 < width {
				jsonResponse(res, http.StatusConflict, "Image dimensions must be at least 200px by 200px or max WIDTH[5:1] or HEIGHT[1:10]")
				return
			}
		}

		uuid, _ := uuid.NewV4()

		imgName := fmt.Sprintf("%s.%s", uuid, strings.Split(imgType, "/")[1]) // example uuid.jpg
		imagePath := "./website/img/" + imgName
		imagePublic = "/public/img/" + imgName

		dest, errCreate := os.Create(imagePath)
		if errCreate != nil {
			jsonResponse(res, http.StatusInternalServerError, "Something wrong")
			return
		}

		defer dest.Close()

		img.Seek(0, 0)

		_, err = io.Copy(dest, img)
		if err != nil {
			jsonResponse(res, http.StatusInternalServerError, "Error copying img")
			return
		}

		defer img.Close()
	}

	title := req.FormValue("title")
	content := req.FormValue("content")

	category := req.FormValue("category")
	var categories []string
	err = json.Unmarshal([]byte(category), &categories)
	if err != nil {
		jsonResponse(res, http.StatusInternalServerError, "JSON Unmarshal error")
		return
	}

	if len(title) >= 200 || len(content) >= 2500 {
		jsonResponse(res, http.StatusBadRequest, "Size of title or content is too large")
		return
	}

	// Log final values before inserting into the database
	err = forum.InsertPost(strconv.Itoa(userID), title, content, imagePublic, time.Now(), categories)
	if err != nil {
		if err.Error() == "categoryNotFound" {
			jsonResponse(res, http.StatusBadRequest, "Invalid request body")
			return
		}
		jsonResponse(res, http.StatusInternalServerError, "Error add post")
		return
	}

	jsonResponse(res, http.StatusOK, "Post created successfully")
}

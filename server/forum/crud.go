package forum

import (
	"database/sql"
	"fmt"
	"html"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LastUsers(user string) (*sql.Rows, error) {
	var query = `
	SELECT DISTINCT 
    	CASE 
        	WHEN user_one = ? THEN user_tow
        	ELSE user_one 
    	END AS user
	FROM chat
	WHERE 
    	user_one = ? OR user_tow = ?
	ORDER BY 
    	last_send DESC;`

	rows, err := SelectQuery(query, user, user, user)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func GetChatMessages(userone, usertow, lastMessage string) (*sql.Rows, error) {
	row, err := SelectOneRow("SELECT id FROM chat WHERE user_one = ? AND user_tow = ? OR user_tow = ? AND user_one = ?",
		userone, usertow, userone, usertow)
	if err != nil {
		return nil, err
	}
	var idchat int
	err = row.Scan(&idchat)
	if err != nil {
		return nil, err
	}
	where := ""
	if lastMessage != "" {
		where = "AND send_at < ?"
	}
	query := fmt.Sprintf("SELECT sender_user,message_text,send_at FROM message WHERE chat_id = ? %v ORDER BY send_at DESC, id DESC limit 10", where)
	rows, err := SelectQuery(query, idchat, lastMessage)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func GetUsers(lastuser string) (*sql.Rows, error) {
	where := ""
	if lastuser != "" {
		where = "WHERE nickname > ?"
	}
	qury := fmt.Sprintf("SELECT nickname,first_name,last_name FROM user %v ORDER BY nickname LIMIT 10", where)
	rows, err := SelectQuery(qury, lastuser)
	return rows, err
}

func Insertmessage(chatID int, sender, mesage string, sendat time.Time) error {
	mesage = html.EscapeString(mesage)
	_, err := ExecQuery("INSERT INTO message(chat_id,sender_user,message_text,send_at) VALUES(?,?,?,?)",
		chatID, sender, mesage, sendat)
	return err
}

func InsertChat(userOne, userTow string, lastSend time.Time) error {
	_, err := ExecQuery("INSERT INTO chat(user_one, user_tow, last_send) VALUES(?,?,?)",
		userOne, userTow, lastSend)
	return err
}

func GetChatID(userOne, userTow string, lastSend time.Time) (int, error) {
	var id int
	var row *sql.Row

	row, errSelct := SelectOneRow("SELECT id FROM chat WHERE user_one = ? AND user_tow = ? OR user_tow = ? AND user_one = ?",
		userOne, userTow, userOne, userTow)
	if errSelct != nil {
		return -1, errSelct
	}
	errScan := row.Scan(&id)
	if errScan != nil {
		if errScan == sql.ErrNoRows {
			errChat := InsertChat(userOne, userTow, lastSend)
			if errChat != nil {
				return -1, fmt.Errorf("inser chat: %v", errChat)
			}
		} else {
			return -1, errScan
		}
	} else {
		return id, nil
	}

	row2, errSelct2 := SelectOneRow("SELECT id FROM chat WHERE user_one = ? AND user_tow = ? OR user_tow = ? AND user_one = ?",
		userOne, userTow, userOne, userTow)
	if errSelct2 != nil {
		return -1, fmt.Errorf("select id chat after insert: %v", errSelct2)
	}

	errScan = row2.Scan(&id)
	if errScan != nil {
		return -1, errScan
	}

	return id, errScan
}

func UpdateLastTimeChat(chatID int, lastmessage time.Time) error {
	_, err := ExecQuery("UPDATE chat SET last_send = ? WHERE id = ?", lastmessage, chatID)
	return err
}

func GetUserById(id int) (string, error) {
	var username string
	row, err := SelectOneRow("SELECT nickname FROM user WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	err = row.Scan(&username)
	return username, err
}

func GetUserByUuid(uuid string) (string, error) {
	var username string
	row, err := SelectOneRow("SELECT nickname FROM user WHERE uuid = ?", uuid)
	if err != nil {
		return "", err
	}
	err = row.Scan(&username)
	return username, err
}

//-------------------------------------------------------------------------------------

// add user to database
func InsertUser(username, firstname, lastname, email, passowrd, gender string, age int) error {
	_, err := ExecQuery("INSERT INTO user(nickname,age,gender,first_name,last_name,email,password,created_at) VALUES(?,?,?,?,?,?,?,?)",
		username, age, gender, firstname, lastname, email, passowrd, time.Now().Format(time.ANSIC))
	return err
}

// update user uuid
func InsertUuid(uid, nameEmail string, exp time.Time) error {
	_, err := ExecQuery("UPDATE user SET uuid = ?,exp = ? WHERE nickname = ? OR email = ?",
		uid, exp, nameEmail, nameEmail)
	return err
}

// give user uuid value null
func DeletUuid(uid string) error {
	_, err := ExecQuery("UPDATE user SET uuid = NULL WHERE uuid = ?", uid)
	if err != nil {
		return fmt.Errorf("error find user %v", err)
	}
	return nil
}

// check if password corect in log in
func CheckHashPasword(password, nameEmail string) error {
	row, err := SelectOneRow("SELECT password FROM user WHERE nickname = ? OR email = ?", nameEmail, nameEmail)
	if err != nil {
		return fmt.Errorf("ERROR SELECT PASSWORD %v", err)
	}
	var hashPass string
	err = row.Scan(&hashPass)
	if err != nil {
		return fmt.Errorf("userNotFound")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	if err != nil {
		return fmt.Errorf("ERROR COMPARE PASSWORD %v", err)
	}
	return nil
}

// get ids categorys choose user in create post and check is valid
func getIdCategoryByName(categorys []string) ([]int, error) {
	for i := 0; i < len(categorys); i++ {
		for j := i + 1; j < len(categorys); j++ {
			if categorys[i] == categorys[j] {
				return nil, fmt.Errorf("duplicate categories")
			}
		}
	}
	var res []int
	for _, catgo := range categorys {
		row, err := SelectOneRow("SELECT id FROM category WHERE category = ?", catgo)
		if err != nil {
			return nil, fmt.Errorf("erorr not found catgo %v", err)
		}
		var id int
		if err := row.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scan row'%s': %v", catgo, err)
		}
		res = append(res, id)
	}
	return res, nil
}

// add post to database
func InsertPost(userID string, title string, content string, image string, createdAt time.Time, categories []string) error {
	categoryIDs, err := getIdCategoryByName(categories)
	if err != nil {
		return fmt.Errorf("categoryNotFound")
	}
	title = html.EscapeString(title)
	content = html.EscapeString(content)

	result, err := ExecQuery("INSERT INTO post (user_id, title, content, image, created_at) VALUES (?, ?, ? ,?, ?)", userID, title, content, image, time.Now())
	if err != nil {
		return err
	}

	// Get the last inserted post ID
	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert categories for the post
	for _, categoryID := range categoryIDs {
		_, err = ExecQuery("INSERT INTO postcategory (post_id, category_id) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			return fmt.Errorf("failed to insert category ID %d: %v", categoryID, err)
		}
	}

	return nil
}

func LikePost(postID int, userID int, statuslike int) error {
	if statuslike < -1 || statuslike > 1 {
		return fmt.Errorf("unsupported statuslike")
	}
	var like int
	row, err := SelectOneRow("SELECT like FROM like WHERE post_id = ? AND user_id = ?", postID, userID)
	if err != nil {
		return err
	}
	errlike := row.Scan(&like)

	var post int
	row, err = SelectOneRow("SELECT id FROM post WHERE id = ?", postID)
	if err != nil {
		return err
	}
	errpost := row.Scan(&post)

	if errpost != nil {
		return fmt.Errorf("post id makanch")
	}

	if errlike == sql.ErrNoRows {
		if statuslike == 0 {
			return fmt.Errorf("status like = 0 ")
		}
		_, err = ExecQuery("INSERT INTO like (post_id, user_id, like) VALUES (?, ?, ?)", postID, userID, statuslike)
		return err
	} else if errlike != nil {
		return fmt.Errorf("err scan like")
	}

	if statuslike == 0 {
		_, err = ExecQuery("DELETE FROM like WHERE user_id = ? AND post_id = ?", userID, postID)
	} else {
		_, err = ExecQuery("UPDATE like SET like = ? WHERE post_id = ? AND user_id = ?", statuslike, postID, userID)
	}
	if err != nil {
		fmt.Println("Error during update/delete:", err)
	}
	return err
}

func LikeComment(commentID int, userID int, statuslike int) error {
	if statuslike < -1 || statuslike > 1 {
		return fmt.Errorf("unsupported statuslike")
	}
	var like int
	row, err := SelectOneRow("SELECT like FROM commentlike WHERE comment_id = ? AND user_id = ?", commentID, userID)
	if err != nil {
		return err
	}
	errlike := row.Scan(&like)

	var comment int
	row, err = SelectOneRow("SELECT id FROM comment WHERE id = ?", commentID)
	if err != nil {
		return err
	}
	errcomment := row.Scan(&comment)

	if errcomment != nil {
		return fmt.Errorf("comment id makanch")
	}

	if errlike == sql.ErrNoRows {
		if statuslike == 0 {
			return fmt.Errorf("status like = 0 ")
		}
		_, err = ExecQuery("INSERT INTO commentlike (comment_id, user_id, like) VALUES (?, ?, ?)", commentID, userID, statuslike)
		return err
	}

	if statuslike == 0 {
		_, err = ExecQuery("DELETE FROM commentlike WHERE user_id = ? AND comment_id = ?", userID, commentID)
	} else {
		_, err = ExecQuery("UPDATE commentlike SET like = ? WHERE comment_id = ? AND user_id = ?", statuslike, commentID, userID)
	}
	if err != nil {
		fmt.Println("Error during update/delete:", err)
	}
	return err
}

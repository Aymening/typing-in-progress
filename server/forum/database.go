package forum

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDataBase() {
	db, err := sql.Open("sqlite3", "./server/database/database.db")

	CheckError("open file database", err)

	// Enable foreign key support
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	CheckError("pragma on error :", err)

	// User table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL,
		gender TEXT,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
		created_at TEXT NOT NULL,
		uuid TEXT UNIQUE,
		exp TEXT 
	);`)
	CheckError("table user error :", err)

	// message table
	_, err = db.Exec(`CREATE TABLE IF not EXISTS message (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
        sender_user TEXT NOT NULL,
        message_text TEXT NOT NULL,
        send_at TEXT NOT NULL,
		FOREIGN KEY (chat_id) REFERENCES chat (id) ON DELETE CASCADE,
        FOREIGN KEY (sender_user) REFERENCES user (nickname) ON DELETE CASCADE
    );`)
	CheckError("table message error :", err)

	// chat table
	_, err = db.Exec(`CREATE TABLE IF not EXISTS chat (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_one TEXT NOT NULL,
			user_tow TEXT NOT NULL,
			last_send TEXT NOT NULL,
			FOREIGN KEY (user_one) REFERENCES user (nickname) ON DELETE CASCADE,
			FOREIGN KEY (user_tow) REFERENCES user (nickname) ON DELETE CASCADE
		);`)
	CheckError("table chat error :", err)

	// Post table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS post (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT,
        content TEXT,
		image TEXT,
        created_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
    );`)
	CheckError("table post error :", err)

	// Like table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS like (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        like INTEGER NOT NULL,
        FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
    );`)
	CheckError("table like error :", err)

	// Comment table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        comment_content TEXT NOT NULL,
        created_at TEXT NOT NULL,
        hasChange INTEGER,
        FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
    );`)
	CheckError("table comment error :", err)

	// comment Like table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS commentlike (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        comment_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        like INTEGER NOT NULL,
        FOREIGN KEY (comment_id) REFERENCES comment(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
    );`)
	CheckError("table like error :", err)

	// Category table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category TEXT NOT NULL,
        description TEXT
    );`)
	CheckError("table categort error :", err)

	// Insert categories
	// _, err = db.Exec(`INSERT INTO category (category, description) VALUES
	// 	('Technology', 'All things related to technology and innovation'),
	// 	('Science', 'Posts about scientific discoveries, research, and theories'),
	// 	('Health', 'Health tips, fitness, and wellness-related posts'),
	// 	('Lifestyle', 'Posts related to personal style, habits, and daily living'),
	// 	('Education', 'Content on learning, development, and educational resources'),
	// 	('Business', 'Business news, trends, and entrepreneurial content'),
	// 	('Entertainment', 'Movies, music, gaming, and other forms of entertainment'),
	// 	('Funny', 'Ha Ha Ha Ha');`)
	// CheckError("insert categories error :", err)

	// PostCategory table to link posts and categories
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS postcategory (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
        FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
    );`)
	CheckError("table postcategory error :", err)
	db.Close()
}

func CheckError(str string, err error) {
	if err != nil {
		log.Fatal(str, err.Error())
	}
}

// for to many rows
func SelectQuery(query string, args ...any) (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", "./server/database/database.db")
	if err != nil {
		return nil, fmt.Errorf("OPEN ERROR: %v", err)
	}
	rows, er := db.Query(query, args...)
	if er != nil {
		return nil, fmt.Errorf("QUERY ERROR: %v", er)
	}
	db.Close()
	return rows, nil
}

// for single row
func SelectOneRow(query string, args ...any) (*sql.Row, error) {
	db, err := sql.Open("sqlite3", "./server/database/database.db")
	if err != nil {
		return nil, fmt.Errorf("OPEN ERROR: %v", err)
	}
	row := db.QueryRow(query, args...)
	db.Close()
	return row, nil
}

// for run query
func ExecQuery(query string, args ...any) (sql.Result, error) {
	db, err := sql.Open("sqlite3", "./server/database/database.db")
	if err != nil {
		return nil, fmt.Errorf("OPEN ERROR: %v", err)
	}
	defer db.Close()

	rs, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("EXEC ERROR: %v", err)
	}
	return rs, nil
}

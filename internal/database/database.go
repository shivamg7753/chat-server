package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(filepath string) (*DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &DB{conn: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return database, nil
}

func (db *DB) createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		text TEXT NOT NULL,
		room TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := db.conn.Exec(usersTable); err != nil {
		return err
	}

	if _, err := db.conn.Exec(messagesTable); err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateUser(username, hashedPassword string) (int64, error) {
	result, err := db.conn.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		username, hashedPassword,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) GetUserByUsername(username string) (int64, string, error) {
	var id int64
	var password string
	err := db.conn.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		username,
	).Scan(&id, &password)
	return id, password, err
}

func (db *DB) SaveMessage(userID int64, username, text, room string) error {
	_, err := db.conn.Exec(
		"INSERT INTO messages (user_id, username, text, room) VALUES (?, ?, ?, ?)",
		userID, username, text, room,
	)
	return err
}

type Message struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"user"`
	Text      string    `json:"text"`
	Room      string    `json:"room"`
	CreatedAt time.Time `json:"created_at"`
}

func (db *DB) GetRecentMessages(limit int) ([]Message, error) {
	rows, err := db.conn.Query(
		"SELECT id, user_id, username, text, room, created_at FROM messages ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Username, &msg.Text, &msg.Room, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (db *DB) GetMessagesByRoom(room string, limit int) ([]Message, error) {
	rows, err := db.conn.Query(
		"SELECT id, user_id, username, text, room, created_at FROM messages WHERE room = ? ORDER BY created_at DESC LIMIT ?",
		room, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Username, &msg.Text, &msg.Room, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

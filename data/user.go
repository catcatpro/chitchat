package data

import (
	"fmt"
	"time"
)

type User struct {
	Id       int
	Uuid     string
	Name     string
	Email    string
	Password string
	CreateAt time.Time
}

type Session struct {
	Id       int
	Uuid     string
	Email    string
	UserId   int
	CreateAt time.Time
}

func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreateAt)
	return
}

func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreateAt)
	return
}

func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreateAt)

	if err != nil {
		valid = false
		return
	}

	if session.Id != 0 {
		valid = true
	}

	return
}

func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreateAt)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	fmt.Println(err)
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now()).Scan(&user.Id, &user.Uuid, &user.CreateAt)
	return
}

// Get the user from the session
func (session *Session) User() (user User, err error) {
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreateAt)
	return
}

// Create a new thread
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	statement := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2,$3,$4) returning id, uuid, topic, user_id, created_at"
	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()
	//user QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), topic, user.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreateAt)
	return
}

// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	sttement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(sttement)
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

package models

import (
	"log"
	"time"
)

// User はユーザー情報を表します
type User struct {
	ID         int
	UUID       string
	Name       string
	Email      string
	PassWord   string
	Created_At time.Time
	Todos      []Todo
}

// Session はセッション情報を表します
type Session struct {
	ID         int
	UUID       string
	Email      string
	UserID     int
	Created_At time.Time
}

// CreateUser は新しいユーザーをデータベースに作成します
func (u *User) CreateUser() error {
	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values (?, ?, ?, ?, ?)`

	_, err := Db.Exec(cmd,
		createUUID(),
		u.Name,
		u.Email,
		Encrypt(u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// GetUser は特定のIDを持つユーザーを取得します
func GetUser(id int) (User, error) {
	user := User{}
	cmd := `select id, uuid, name, email, password, created_at from users where id = ?`

	err := Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.Created_At,
	)
	return user, err
}

// UpdateUser はユーザー情報を更新します
func (u *User) UpdateUser() error {
	cmd := `update users set name = ?, email = ? where id = ?`

	_, err := Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// DeleteUser はユーザーを削除します
func (u *User) DeleteUser() error {
	cmd := `delete from users where id = ?`

	_, err := Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// GetUserByEmail は特定のEmailを持つユーザーを取得します
func GetUserByEmail(email string) (User, error) {
	user := User{}
	cmd := `select id, uuid, name, email, password, created_at from users where email = ?`

	err := Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.Created_At)

	return user, err
}

// CreateSession は新しいセッションを作成します
func (u *User) CreateSession() (Session, error) {
	session := Session{}

	cmd1 := `insert into sessions(
	uuid,
	email,
	user_id,
	created_at) values(?,?,?,?)`

	_, err := Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	cmd2 := `select id, uuid, email, user_id, created_at from sessions where user_id = ? and email = ?`

	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.Created_At,
	)
	if err != nil {
		log.Fatalln(err)
	}

	return session, err
}

// CheckSession はセッションが有効かどうかを確認します
func (sess *Session) CheckSession() (bool, error) {
	cmd := `select id, uuid, email, user_id, created_at from sessions where uuid = ?`

	err := Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.Created_At,
	)

	valid := sess.ID != 0

	return valid, err
}

// DeleteSessionByUUID はUUIDに基づいてセッションを削除します
func (sess *Session) DeleteSessionByUUID() error {
	cmd := `delete from sessions where uuid = ?`

	_, err := Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// GetUserBySession はセッションに関連付けられたユーザーを取得します
func (sess Session) GetUserBySession() (User, error) {
	user := User{}
	cmd := `select id, uuid, name, email, created_at from users where id = ?`

	err := Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.Created_At,
	)
	return user, err
}

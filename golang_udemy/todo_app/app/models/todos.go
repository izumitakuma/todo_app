package models

import (
	"log"
	"time"
)

// Todo は各Todo項目を表します
type Todo struct {
	ID         int
	Content    string
	User_id    int
	Created_At time.Time
}

// CreateTodo は新しいTodoをデータベースに挿入します
func (u *User) CreateTodo(content string) error {
	cmd := `insert into todos (
		content, 
		user_id, 
		created_at) values (?, ?, ?)`

	_, err := Db.Exec(cmd, content, u.ID, time.Now())
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// UpdateTodo は既存のTodoを更新します
func (t *Todo) UpdateTodo() error {
	cmd := `update todos set content=?,user_id=? where id=?`

	_, err := Db.Exec(cmd, t.Content, t.User_id, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// GetTodo は特定のIDを持つTodoを取得します
func GetTodo(id int) (Todo, error) {
	cmd := `select id, content, user_id, created_at from todos where id = ?`
	todo := Todo{}

	err := Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.User_id,
		&todo.Created_At)

	return todo, err
}

// GetTodos はすべてのTodoを取得します
func GetTodos() ([]Todo, error) {
	cmd := `select id, content, user_id, created_at from todos`

	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.User_id,
			&todo.Created_At)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	return todos, err
}

// GetTodosByUser は特定のユーザーに関連するTodoをすべて取得します
func (u User) GetTodosByUser() ([]Todo, error) {
	cmd := `select * from todos where user_id=?`

	rows, err := Db.Query(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Content,
			&todo.User_id,
			&todo.Created_At)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	return todos, err
}

// DeleteTodo は特定のTodoを削除します
func (t *Todo) DeleteTodo() error {
	cmd := `delete from todos where id = ?`
	_, err := Db.Exec(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

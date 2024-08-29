package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func top(w http.ResponseWriter, r *http.Request) {
	if _, err := session(w, r); err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", http.StatusFound) // 302
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println("Error getting user by session:", err)
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	log.Println(user.Name)

	todos, err := user.GetTodosByUser()
	if err != nil {
		log.Println("Error getting todos:", err)
	}

	user.Todos = todos
	generateHTML(w, user, "layout", "index", "private_navbar")
}

func todoNew(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	if r.Method == "GET" {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	} else if r.Method == "POST" {
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println("Error getting user by session:", err)
			http.Redirect(w, r, "/login", http.StatusFound) // 302
			return
		}

		if err := user.CreateTodo(r.PostFormValue("content")); err != nil {
			log.Println("Error creating todo:", err)
		}

		http.Redirect(w, r, "/todos", http.StatusFound) // 302
	}
}

func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	if _, err := sess.GetUserBySession(); err != nil {
		log.Fatalln("Error getting user by session:", err)
	}

	todo, err := models.GetTodo(id)
	if err != nil {
		log.Fatalln("Error getting todo:", err)
	}

	generateHTML(w, todo, "layout", "private_navbar", "todo_edit")
}

func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Fatalln("Error parsing form:", err)
	}

	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println("Error getting user by session:", err)
	}

	todo := &models.Todo{ID: id, Content: r.PostFormValue("content"), User_id: user.ID}

	if err := todo.UpdateTodo(); err != nil {
		log.Fatalln("Error updating todo:", err)
	}

	http.Redirect(w, r, "/todos", http.StatusFound) // 302
}

func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	if _, err := sess.GetUserBySession(); err != nil {
		log.Println("Error getting user by session:", err)
	}

	todo, err := models.GetTodo(id)
	if err != nil {
		log.Println("Error getting todo:", err)
	}

	if err := todo.DeleteTodo(); err != nil {
		log.Fatalln("Error deleting todo:", err)
	}

	http.Redirect(w, r, "/todos", http.StatusFound) // 302
}

package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if _, err := session(w, r); err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", http.StatusFound) // 302
		}
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println("Error parsing form:", err)
		}

		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}

		if err := user.CreateUser(); err != nil {
			log.Println("Error creating user:", err)
		}

		http.Redirect(w, r, "/", http.StatusFound) // 302
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if _, err := session(w, r); err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "/todos", http.StatusFound) // 302
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println("Error getting user by email:", err)
		http.Redirect(w, r, "/login", http.StatusFound) // 302
		return
	}

	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println("Error creating session:", err)
			http.Redirect(w, r, "/login", http.StatusFound) // 302
			return
		}

		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/todos", http.StatusFound) // 302
	} else {
		http.Redirect(w, r, "/login", http.StatusFound) // 302
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil && err != http.ErrNoCookie {
		log.Println("Error retrieving cookie:", err)
	}

	if err == nil {
		session := models.Session{UUID: cookie.Value}
		if err := session.DeleteSessionByUUID(); err != nil {
			log.Println("Error deleting session:", err)
		}
	}

	http.Redirect(w, r, "/login", http.StatusFound) // 302
}

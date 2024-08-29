package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"todo_app/app/models"
	"todo_app/config"
)

// generateHTML generates and renders HTML templates
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// session retrieves the session from the request cookie
func session(w http.ResponseWriter, r *http.Request) (models.Session, error) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		return models.Session{}, err
	}

	sess := models.Session{UUID: cookie.Value}
	if ok, _ := sess.CheckSession(); !ok {
		return sess, fmt.Errorf("Invalid session")
	}

	return sess, nil
}

// validPath is a regular expression to validate todo paths
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// ParseURL extracts the ID from the URL and passes it to the handler function
func ParseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := validPath.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(matches[2])
		if err != nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, id)
	}
}

// StartMainServer sets up the routes and starts the HTTP server
func StartMainServer() error {
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	// Route handlers
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/edit/", ParseURL(todoEdit))
	http.HandleFunc("/todos/update/", ParseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", ParseURL(todoDelete))

	// Start the server
	return http.ListenAndServe(":"+config.Config.Port, nil)
}

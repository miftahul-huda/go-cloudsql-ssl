package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"go-cloud-ssl/db"
	"go-cloud-ssl/models"
)

func Index(w http.ResponseWriter, r *http.Request, config db.Config) {
	rows, err := db.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			http.Error(w, "Failed to scan row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	renderTemplate(w, "index", users)

}

func Create(w http.ResponseWriter, r *http.Request, config db.Config) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		db.Execute(config, "INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Edit(w http.ResponseWriter, r *http.Request, config db.Config) {
	id := r.URL.Query().Get("id")
	row := db.ExecuteRow(config, "SELECT id, name, email FROM users WHERE id = ?", id)
	var u models.User
	row.Scan(&u.ID, &u.Name, &u.Email)

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/edit.html"))
	tmpl.ExecuteTemplate(w, "layout", u)
}

func Update(w http.ResponseWriter, r *http.Request, config db.Config) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	db.Execute(config, "UPDATE users SET name = ?, email = ? WHERE id = ?", name, email, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request, config db.Config) {
	id := r.URL.Query().Get("id")
	db.Execute(config, "DELETE FROM users WHERE id = ?", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := fmt.Sprintf("templates/%s.html", tmpl)
	t, err := template.ParseFiles("templates/layout.html", tmplPath)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}
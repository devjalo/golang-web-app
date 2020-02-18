package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"text/template"
)

type Contact struct {
	Id    int
	Name  string
	Email string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "goapp"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Contact ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Contact{}
	var res []Contact
	for selDB.Next() {
		var id int
		var name, email string
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Email = email
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Contact WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Contact{}
	for selDB.Next() {
		var id int
		var name, email string
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Email = email
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Contact WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Contact{}
	for selDB.Next() {
		var id int
		var name, email string
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Email = email
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		insForm, err := db.Prepare("INSERT INTO Contact(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Contact SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, id)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Contact WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}

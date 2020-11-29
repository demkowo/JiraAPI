package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/callback", callback)
	http.ListenAndServe(":5000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func callback(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")

	d := struct {
		Fname string
		Lname string
	}{
		Fname: fname,
		Lname: lname,
	}

	tpl.ExecuteTemplate(w, "callback.gohtml", d)
	fmt.Println(fname, lname)
}

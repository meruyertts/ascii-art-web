package main

import (
	"ascii-art-web/printascii"
	"fmt"
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ascii-art", processor)
	fmt.Println("Server launched ...")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func processor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/ascii-art" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	fname := r.FormValue("string")
	f := r.FormValue("font")

	ascii, err := printascii.AsciiWeb(fname, f)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	d := struct {
		AsciiPrint string
	}{
		AsciiPrint: ascii,
	}
	tpl.ExecuteTemplate(w, "index.html", d)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		tpl.ExecuteTemplate(w, "notfound.html", nil)
	} else if status == http.StatusBadRequest {
		tpl.ExecuteTemplate(w, "badreq.html", nil)
	} else if status == http.StatusInternalServerError {
		tpl.ExecuteTemplate(w, "internal.html", nil)
	}
}

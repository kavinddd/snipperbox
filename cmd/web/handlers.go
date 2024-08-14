package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.NotFound(w, r)
		return
	}

	templates, err := template.ParseFiles(
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	fmt.Println(templates)

	err = templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("Display a specific snippet %s", idString)))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, fmt.Sprintf("Method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet"))
}

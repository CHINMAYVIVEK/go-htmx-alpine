package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Todo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"is_completed"`
}

var todos = []Todo{
	{ID: 1, Name: "Learn Go", IsCompleted: false},
	{ID: 2, Name: "Learn Alpine", IsCompleted: false},
	{ID: 3, Name: "Learn Htmx", IsCompleted: false},
}
var templates map[string]*template.Template

var indexhtml = "index.html"
var todohtml = "todo.html"

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates[indexhtml] = template.Must(template.ParseFiles(indexhtml))
	templates[todohtml] = template.Must(template.ParseFiles(todohtml))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates[indexhtml].ExecuteTemplate(w, indexhtml, map[string]template.JS{"Todos": template.JS(json)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func submitTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	completed := r.PostFormValue("is_completed") == "true"
	todo := Todo{ID: len(todos) + 1, Name: name, IsCompleted: completed}
	todos = append(todos, todo)

	tmpl := templates[todohtml]
	err := tmpl.ExecuteTemplate(w, todohtml, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit-todo", submitTodoHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

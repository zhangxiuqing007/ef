package controllers

import (
	"html/template"
	"net/http"
)

var errorTemplate = template.Must(template.ParseFiles("views/error.html"))

func sendErrorPage(w http.ResponseWriter, err string) {
	errorTemplate.Execute(w, err)
}

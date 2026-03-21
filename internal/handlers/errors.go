package handlers

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func InitTemplates() error {
	result, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return err
	}
	tmpl = result
	return nil
}
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := tmpl.ExecuteTemplate(w, "404.html", nil)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}
func StatusInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	err := tmpl.ExecuteTemplate(w, "500.html", nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered", err)
				StatusInternalServerError(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

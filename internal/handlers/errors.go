package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func NotFoundHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		err := tmpl.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			log.Print(err)
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}
}
func StatusInternalServerError(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		err := tmpl.ExecuteTemplate(w, "500.html", nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
func RecoveryMiddleware(tmpl *template.Template, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered", err)
				StatusInternalServerError(tmpl)(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

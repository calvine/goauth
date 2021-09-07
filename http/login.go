package http

import (
	"html/template"
	"net/http"
	"sync"
)

func (s *server) handleLoginGet() http.HandlerFunc {
	var (
		once          sync.Once
		loginTemplate *template.Template
		templateErr   error
	)
	type requestData struct {
		CSRFToken string
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			templateFileData, err := s.templateFS.ReadFile("http/templates/login.tmpl")
			templateErr = err
			if templateErr == nil {
				loginTemplate, templateErr = template.New("loginPage").Parse(string(templateFileData))
			}
		})
		if templateErr != nil {
			http.Error(rw, templateErr.Error(), http.StatusInternalServerError)
		}
		data := requestData{"Test CSRF TOKEN"}
		err := loginTemplate.Execute(rw, data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *server) handleLoginPost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {}
}

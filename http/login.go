package http

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities"
)

func (s *server) handleLoginGet() http.HandlerFunc {
	var (
		once          sync.Once
		loginTemplate *template.Template
		templateErr   error
		templatePath  string = "http/templates/login.tmpl"
	)
	type requestData struct {
		CSRFToken string
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			templateFileData, err := s.templateFS.ReadFile(templatePath)
			templateErr = err
			if templateErr == nil {
				loginTemplate, templateErr = template.New("loginPage").Parse(string(templateFileData))
			}
		})
		if templateErr != nil {
			http.Error(rw, templateErr.Error(), http.StatusInternalServerError)
		}
		tokenString, err := utilities.NewTokenString()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		// TODO: make CSRF token life span configurable
		token := models.NewToken(tokenString, "", models.TokenTypeCSRF, time.Minute*10)
		s.tokenService.PutToken(token)
		templateRenderError := loginTemplate.Execute(rw, requestData{token.Value})
		if templateRenderError != nil {
			err = coreerrors.NewTemplateRenderErrorError(templatePath, templateRenderError, true)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *server) handleLoginPost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {}
}

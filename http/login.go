package http

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
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
			return
		}
		// TODO: make CSRF token life span configurable
		ctx := r.Context()
		token, err := models.NewToken("", models.TokenTypeCSRF, time.Minute*10)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.tokenService.PutToken(ctx, token)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		templateRenderError := loginTemplate.Execute(rw, requestData{token.Value})
		if templateRenderError != nil {
			err = coreerrors.NewTemplateRenderErrorError(templatePath, templateRenderError, true)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleLoginPost() http.HandlerFunc {
	// var (
	// 	once          sync.Once
	// 	loginTemplate *template.Template
	// 	templateErr   error
	// 	templatePath  string = "http/templates/login.tmpl"
	// )
	type requestData struct {
		CSRFToken string
		Email     string
		Password  string
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		data := requestData{}
		ctx := r.Context()
		data.CSRFToken = r.FormValue("csrf_token")
		data.Email = r.FormValue("email")
		data.Password = r.FormValue("password")

		_, err := s.tokenService.GetToken(ctx, data.CSRFToken, models.TokenTypeCSRF)
		if err != nil {
			http.Error(rw, err.GetErrorMessage(), http.StatusBadRequest)
			return
		}
		err = s.tokenService.DeleteToken(ctx, data.CSRFToken)
		if err != nil {
			// uh of the token was not deleted! need to log this...
		}
		_, err = s.loginService.LoginWithPrimaryContact(ctx, data.Email, core.CONTACT_TYPE_EMAIL, data.Password, "login post handler")
		if err != nil {
			http.Error(rw, err.GetErrorMessage(), http.StatusUnauthorized)
			return
		}
		rw.Header().Add("test-header", "JWT")
		http.Redirect(rw, r, "/static/hooray.html", http.StatusFound)
	}
}

// func (s *server) handleAuthGet() http.HandlerFunc {
// 	var (
// 		once          sync.Once
// 		loginTemplate *template.Template
// 		templateErr   error
// 		templatePath  string = "http/templates/login.tmpl"
// 	)
// 	type requestData struct {
// 		clientID      string
// 		codeChallenge string // PKCE
// 		redirectURI   string
// 		responseType  string
// 		scope         string
// 		state         string
// 	}
// 	return func(rw http.ResponseWriter, r *http.Request) {
// 		once.Do(func() {
// 			templateFileData, err := s.templateFS.ReadFile(templatePath)
// 			templateErr = err
// 			if templateErr == nil {
// 				loginTemplate, templateErr = template.New("loginPage").Parse(string(templateFileData))
// 			}
// 		})
// 	}
// }

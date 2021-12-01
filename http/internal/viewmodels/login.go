package viewmodels

type LoginTemplateData struct {
	CallbackURL     string
	CSRFToken       string
	Email           string
	HasErrorMessage bool
	ErrorMsg        string
	Password        string
}

package viewmodels

type LoginTemplateData struct {
	CallbackURL string
	CSRFToken   string
	Email       string
	ErrorMsg    string
	Password    string
	RememberMe  bool
}

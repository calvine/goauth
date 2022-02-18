package service

import (
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/calvine/goauth/core/constants/template"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

type staticTemplateService struct {
	textTemplates map[string]*texttemplate.Template
	htmlTemplates map[string]*htmltemplate.Template
}

func NewStaticTemplateService() (services.TemplateService, errors.RichError) {
	textTemplates := make(map[string]*texttemplate.Template)
	htmlTemplates := make(map[string]*htmltemplate.Template)
	return staticTemplateService{
		textTemplates: textTemplates,
		htmlTemplates: htmlTemplates,
	}, nil
}

func (staticTemplateService) GetName() string {
	return "staticTemplateService"
}

func (sts staticTemplateService) GetTextTemplate(name string) (*texttemplate.Template, bool) {
	t, ok := sts.textTemplates[name]
	if !ok {
		switch name {
		case template.ConfirmContactTextEmailTemplateName:

		case template.PasswordResetTextEmailTemplateName:

		case template.MagicLoginTextEmailTemplateName:

		default:
			return nil, false
		}
	}
	return t, true
}

func (sts staticTemplateService) GetHTMLTemplate(name string) (*htmltemplate.Template, bool) {
	t, ok := sts.htmlTemplates[name]
	if !ok {
		switch name {
		default:
			return nil, false
		}
	}
	return t, true
}

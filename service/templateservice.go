package service

import (
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

type staticTemplateService struct {
}

func NewStaticTemplateService() (services.TemplateService, errors.RichError) {
	return staticTemplateService{}, nil
}

func (staticTemplateService) GetName() string {
	return "staticTemplateService"
}

func (sts staticTemplateService) GetTextTemplate(name string) (*texttemplate.Template, bool) {
	switch name {
	default:
		return nil, false
	}
}

func (sts staticTemplateService) GetHTMLTemplate(name string) (*htmltemplate.Template, bool) {
	switch name {
	default:
		return nil, false
	}
}

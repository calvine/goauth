package service

import (
	"context"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/constants/template"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

// TODO: implement a repo backed version of this

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

func (sts staticTemplateService) GetTextTemplate(ctx context.Context, logger *zap.Logger, name string) (*texttemplate.Template, bool) {
	span := apptelemetry.CreateFunctionSpan(ctx, sts.GetName(), "GetTextTemplate")
	defer span.End()
	cachedTemplate, ok := sts.textTemplates[name]
	if !ok {
		logger.Info("template not found in memory", zap.String("templateName", name))
		var templateString string
		switch name {
		case template.ConfirmContactTextEmailTemplateName:
			templateString = template.ConfirmContactTextEmailTemplate
		// case template.PasswordResetTextEmailTemplateName:
		// TODO: write these
		// case template.MagicLoginTextEmailTemplateName:
		// TODO: write these
		default:
			span.AddEvent("template not found")
			return nil, false
		}
		parsedTemplate, err := texttemplate.New(name).Parse(templateString)
		if err != nil {
			logger.Error("failed to parse template", zap.String("templateString", templateString), zap.Error(err))
		}
		sts.textTemplates[name] = parsedTemplate
	}
	span.AddEvent("template found")
	return cachedTemplate, true
}

func (sts staticTemplateService) GetHTMLTemplate(ctx context.Context, logger *zap.Logger, name string) (*htmltemplate.Template, bool) {
	span := apptelemetry.CreateFunctionSpan(ctx, sts.GetName(), "GetHTMLTemplate")
	defer span.End()
	cachedTemplate, ok := sts.htmlTemplates[name]
	if !ok {
		logger.Info("template not found in memory", zap.String("templateName", name))
		// var templateString string
		switch name {
		// case template.ConfirmContactHtmlEmailTemplateName:
		// 	// TODO: write these
		// case template.PasswordResetHtmlEmailTemplateName:
		// 	// TODO: write these
		// case template.MagicLoginHtmlEmailTemplateName:
		// 	// TODO: write these
		default:
			span.AddEvent("template not found")
			return nil, false
		}
		// parsedTemplate, err := htmltemplate.New(name).Parse(templateString)
		// if err != nil {
		// 	logger.Error("failed to parse template", zap.String("templateString", templateString), zap.Error(err))
		// }
		// sts.htmlTemplates[name] = parsedTemplate
	}
	span.AddEvent("template found")
	return cachedTemplate, true
}

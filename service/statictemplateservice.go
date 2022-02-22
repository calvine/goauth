package service

import (
	"context"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"

	"github.com/calvine/goauth/core/constants/templates"
	coreerrors "github.com/calvine/goauth/core/errors"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

// TODO: implement a repo backed version of this

type staticTemplateService struct {
	textTemplates map[templates.TemplateName]*texttemplate.Template
	htmlTemplates map[templates.TemplateName]*htmltemplate.Template
}

func NewStaticTemplateService() (services.TemplateService, errors.RichError) {
	textTemplates := make(map[templates.TemplateName]*texttemplate.Template)
	htmlTemplates := make(map[templates.TemplateName]*htmltemplate.Template)
	return staticTemplateService{
		textTemplates: textTemplates,
		htmlTemplates: htmlTemplates,
	}, nil
}

func (staticTemplateService) GetName() string {
	return "staticTemplateService"
}

func (sts staticTemplateService) GetTextTemplate(ctx context.Context, logger *zap.Logger, templateName templates.TemplateName) (*texttemplate.Template, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, sts.GetName(), "GetTextTemplate")
	defer span.End()
	cachedTemplate, ok := sts.textTemplates[templateName]
	if !ok {
		logger.Info("template not found in memory", zap.String("templateName", string(templateName)))
		var templateString string
		switch templateName {
		case templates.ConfirmContactTextEmail:
			templateString = templates.ConfirmContactTextEmailTemplate
		default:
			err := coreerrors.NewTemplateNotFoundError(string(templateName), true)
			evtString := "template not found"
			logger.Error(evtString, zap.Reflect("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, evtString)
			return nil, err
		}
		parsedTemplate, err := texttemplate.New(string(templateName)).Parse(templateString)
		if err != nil {
			logger.Error("failed to parse template", zap.String("templateString", templateString), zap.Error(err))
		}
		sts.textTemplates[templateName] = parsedTemplate
		cachedTemplate = parsedTemplate
	}
	span.AddEvent("template found")
	return cachedTemplate, nil
}

func (sts staticTemplateService) GetHTMLTemplate(ctx context.Context, logger *zap.Logger, templateName templates.TemplateName) (*htmltemplate.Template, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, sts.GetName(), "GetHTMLTemplate")
	defer span.End()
	cachedTemplate, ok := sts.htmlTemplates[templateName]
	if !ok {
		logger.Info("template not found in memory", zap.String("templateName", string(templateName)))
		// var templateString string
		switch templateName {
		// case templates.ConfirmContactHTMLEmail:
		// 	templateString = templates.ConfirmContactHTMLEmailTemplate
		default:
			err := coreerrors.NewTemplateNotFoundError(string(templateName), true)
			evtString := "template not found"
			logger.Error(evtString, zap.Reflect("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, evtString)
			return nil, err
		}
		// parsedTemplate, err := htmltemplate.New(name).Parse(templateString)
		// if err != nil {
		// 	logger.Error("failed to parse template", zap.String("templateString", templateString), zap.Error(err))
		// }
		// sts.htmlTemplates[name] = parsedTemplate
		// cachedTemplate = parsedTemplate
	}
	span.AddEvent("template found")
	return cachedTemplate, nil
}

func (sts staticTemplateService) ExecuteTextTemplate(ctx context.Context, logger *zap.Logger, templateName templates.TemplateName, templateData interface{}) (string, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, sts.GetName(), "ExecuteTextTemplate")
	defer span.End()
	template, err := sts.GetTextTemplate(ctx, logger, templateName)
	if err != nil {
		evtString := "template not found"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return "", err
	}
	buf := strings.Builder{}
	tErr := template.Execute(&buf, templateData)
	if tErr != nil {
		err := coreerrors.NewFailedTemplateRenderError(string(templateName), tErr, true)
		evtString := "failed to render template"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return "", err
	}
	return buf.String(), nil
}

func (sts staticTemplateService) ExecuteHTMLTemplate(ctx context.Context, logger *zap.Logger, templateName templates.TemplateName, templateData interface{}) (string, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, sts.GetName(), "ExecuteHTMLTemplate")
	defer span.End()
	template, err := sts.GetHTMLTemplate(ctx, logger, templateName)
	if err != nil {
		evtString := "template not found"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return "", err
	}
	buf := strings.Builder{}
	tErr := template.Execute(&buf, templateData)
	if tErr != nil {
		err := coreerrors.NewFailedTemplateRenderError(string(templateName), tErr, true)
		evtString := "failed to render template"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return "", err
	}
	return buf.String(), nil
}

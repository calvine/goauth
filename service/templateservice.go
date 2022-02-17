package service

import "github.com/calvine/richerror/errors"

type staticTemplateService struct {
}

func NewStaticTemplateService() (staticTemplateService, errors.RichError) {
	return staticTemplateService{}, nil
}

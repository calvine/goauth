package factory

import (
	"net/url"
	"path"

	coreerrors "github.com/calvine/goauth/core/errors"
	corefactory "github.com/calvine/goauth/core/factory"
	"github.com/calvine/richerror/errors"
)

type serviceLinkFactory struct {
	servicePublicURL string
}

func NewServiceLinkFactory(servicePublicURL string) (corefactory.ServiceLinkFactory, errors.RichError) {
	if len(servicePublicURL) == 0 {
		return nil, errors.NewRichError("ServiceLinkFactoryMissingServicePublicURL", "service link factory is missing service public url").WithStack(0)
	}
	return serviceLinkFactory{
		servicePublicURL: servicePublicURL,
	}, nil
}

func (slf serviceLinkFactory) CreateLink(linkPath string) (string, errors.RichError) {
	// This seems wasteful to reparse this each time...
	// TODO: look into alternate solutions
	u, err := url.Parse(slf.servicePublicURL)
	if err != nil {
		return "", coreerrors.NewInvalidValueError(slf.servicePublicURL, true)
	}
	u.Path = path.Join(u.Path, linkPath)
	return u.String(), nil
}

func (slf serviceLinkFactory) CreatePasswordResetLink(passwordResetToken string) (string, errors.RichError) {
	linkPath := ""
	return slf.CreateLink(linkPath)
}

func (slf serviceLinkFactory) CreateConfirmContactLink(confirmContactToken string) (string, errors.RichError) {
	linkPath := ""
	return slf.CreateLink(linkPath)
}

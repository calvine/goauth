package http

import (
	"net/http"
	"strings"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *server) handleConfirmContactGet() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		defer span.End()
		pathParts := strings.Split(r.URL.Path, "/")
		partsLen := len(pathParts)
		confirmationCode := pathParts[partsLen-1]
		err := s.userService.ConfirmContact(ctx, logger, confirmationCode, "contact confirmation page")
		if err != nil {
			var errorMsg string
			var errorCode int
			switch err.GetErrorCode() {
			case coreerrors.ErrCodeInvalidToken:
				errorMsg = "confirmation code invalid"
				errorCode = http.StatusBadRequest
			case coreerrors.ErrCodeExpiredToken:
				errorMsg = "confirmation code has expired"
				errorCode = http.StatusBadRequest
			case coreerrors.ErrCodeWrongTokenType:
				errorMsg = "confirmation code invalid"
				errorCode = http.StatusBadRequest
			case coreerrors.ErrCodeContactAlreadyConfirmed:
				errorMsg = "contact is already confirmed"
				errorCode = http.StatusBadRequest
			default:
				errorMsg = "An unexpected error occurred"
				errorCode = http.StatusInternalServerError
			}
			logger.Error(errorMsg,
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, errorMsg)
			// FIXME: redirect to specific error page... or alternatively just render an error page template here.
			http.Error(rw, errorMsg, errorCode)
			return
		}
		// TODO: success code here
		// The idea here is to redirect to a set password page, so we woud need to make a password reset token and all that jazz?
	}
}

package constants

import "time"

const (
	LoginCookieName             = "x-goauth-session"
	Default_CSRF_Token_Duration = time.Minute * 10
)

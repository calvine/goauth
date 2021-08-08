package codes

const (
	// General Errors

	// ErrCodeWrongType indicates that a type was not what was expected.
	ErrCodeWrongType = "WrongType"
	// ErrCodeNilNotAllowed indicates that a nil value was encountered but not allowed
	ErrCodeNilNotAllowed = "NilNotAllowed"
	// ErrTypeNotAllowed indicates that a type is not allowed
	ErrCodeInvalidType = "InvalidType"
	// ErrCodeInvalidValue
	ErrCodeInvalidValue = "InvalidValue"
	// ErrCodeRepoQueryFailed indicates that an error occurred while performing a db query.
	ErrCodeRepoQueryFailed = "RepoQueryFailed"
	// ErrCodeContactNotPrimary
	ErrCodeContactNotPrimary = "ContactNotPrimary"

	//No Data Found

	// ErrCodeNoUserFound means no users was found for a given request.
	ErrCodeNoUserFound = "NoUserFound"
	// ErrCodeNoContactFound means no contact was found for a given request.
	ErrCodeNoContactFound = "NoContactFound"
	// ErrCodeNoAddressFound means no address was found for a given request.
	ErrCodeNoAddressFound = "NoAddressFound"

	// Login Errors

	// ErrCodeUserLockedOut means that a user is currently locked out.
	ErrCodeUserLockedOut = "UserLockedOut"
)

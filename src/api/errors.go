package api

var (
	// ReasonInternalError is used for an internal server error. This should only
	// be used for CRITICAL and FATAL level errors.
	ReasonInternalError ReasonCode = "INTERNAL_ERROR"

	// ReasonBadRequestData is used if the response body is invalid or fails to meet
	// requirements of a method when its consumed.
	ReasonBadRequestData ReasonCode = "BAD_DATA"

	// ReasonUserAlreadyExists is used for when the SQL database rejects the user
	// due to a duplicate entry.
	ReasonUserAlreadyExists ReasonCode = "USER_ALREADY_EXISTS"

	// ReasonBadUsername is used to indicate the given username failed validation.
	ReasonBadUsername ReasonCode = "BAD_USERNAME"

	// ReasonUnauthorized is used for unauthenticated requests.
	ReasonUnauthorized ReasonCode = "UNAUTHORIZED"
)

var (
	ErrorInternalErrorMsg = "An unexpected internal error has occurred"
	ErrorUnauthorizedMsg  = "Unauthorized access"
	ErrorBadDataMsg       = "Bad request data"
)

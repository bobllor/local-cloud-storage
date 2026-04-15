package api

// StatusType is the status of the response.
type StatusType string

type Error struct {
	// Code is the status code of the error.
	Code int `json:"code"`

	// Message is the reason why the error had occurred.
	Message string `json:"message"`
}

const (
	StatusSuccess StatusType = "success"
	StatusError              = "error"
)

// ApiResponse is the response streamed to the client from the server.
// This is the standardized response for the backend system for all handlers.
type ApiResponse struct {
	// Status is the status of the response.
	Status StatusType `json:"status"`

	// Output is the output to the client of any type. It can be omitted.
	Output any `json:"output,omitempty"`

	// Error is an Error type that occurred due to an issue with the request or backend.
	// It can be omitted.
	Error *Error `json:"error,omitempty"`
}

// NewApiResponse creates a new "success" API response by default.
//
// output is any output type that is added to the output of the response.
func NewApiResponse(output any) *ApiResponse {
	res := &ApiResponse{
		Status: StatusSuccess,
		Output: output,
	}

	return res
}

// NewApiResponseError creates a new error API response.
//
// statusCode is the type of status code int.
//
// errMessage is the message string to send back to the client.
func NewApiResponseError(statusCode int, errMessage string) *ApiResponse {
	res := &ApiResponse{
		Status: StatusError,
		Error: &Error{
			Code:    statusCode,
			Message: errMessage,
		},
	}

	return res
}

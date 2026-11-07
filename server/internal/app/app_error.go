package app

type ErrorCode string

const (
	BadRequest ErrorCode = "bad_request"
	NotFound   ErrorCode = "not_found"
	Invalid    ErrorCode = "invalid"
	Conflict   ErrorCode = "conflict"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func New(
	code ErrorCode,
	message string,
	err error,
) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

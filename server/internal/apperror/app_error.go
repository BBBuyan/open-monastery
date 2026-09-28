package apperror

type Code string

const (
	BadRequest Code = "bad_request"
	NotFound   Code = "not_found"
	Invalid    Code = "invalid"
	Conflict   Code = "conflict"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func New(code Code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

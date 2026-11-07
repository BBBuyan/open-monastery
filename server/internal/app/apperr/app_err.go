package apperr

type AppErrCode string

const (
	BadRequest AppErrCode = "bad_request"
	NotFound   AppErrCode = "not_found"
	Invalid    AppErrCode = "invalid"
	Conflict   AppErrCode = "conflict"
)

type AppErr struct {
	Code    AppErrCode
	Message string
	Err     error
}

// Implements error
var _ error = (*AppErr)(nil)

func (e *AppErr) Error() string {
	return e.Message
}

func (e *AppErr) Unwrap() error {
	return e.Err
}

func NewConflict(msg string, err error) *AppErr {
	return &AppErr{
		Code:    Conflict,
		Message: msg,
		Err:     err,
	}
}

func NewNotFound(msg string, err error) *AppErr {
	return &AppErr{
		Code:    NotFound,
		Message: msg,
		Err:     err,
	}
}

func NewInvalid(msg string, err error) *AppErr {
	return &AppErr{
		Code:    Invalid,
		Message: msg,
		Err:     err,
	}
}

func NewBadRequest(msg string, err error) *AppErr {
	return &AppErr{
		Code:    BadRequest,
		Message: msg,
		Err:     err,
	}
}

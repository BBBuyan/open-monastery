package apperr

type AppErrType string

const (
	BadRequest AppErrType = "bad_request"
	NotFound   AppErrType = "not_found"
	Invalid    AppErrType = "invalid"
	Conflict   AppErrType = "conflict"
)

type AppErr struct {
	Type    AppErrType
	Code    string
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

func NewConflict(
	code string,
	msg string,
	err error,
) *AppErr {
	return &AppErr{
		Type:    Conflict,
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func NewNotFound(
	code string,
	msg string,
	err error,
) *AppErr {
	return &AppErr{
		Type:    NotFound,
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func NewInvalid(
	code string,
	msg string,
	err error,
) *AppErr {
	return &AppErr{
		Type:    Invalid,
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func NewBadRequest(
	code string,
	msg string,
	err error,
) *AppErr {
	return &AppErr{
		Type:    BadRequest,
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

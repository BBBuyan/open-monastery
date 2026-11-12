package apperr

type AppErrType string

const (
	BadRequest   AppErrType = "bad_request"
	Unauthorized AppErrType = "unauthorized"
	NotFound     AppErrType = "not_found"
	Invalid      AppErrType = "invalid"
	Conflict     AppErrType = "conflict"
)

type AppErr struct {
	Type    AppErrType
	Code    string
	Message *string
	Err     error
}

// Implements error
var _ error = (*AppErr)(nil)

func (e *AppErr) Error() string {
	return e.Code
}

func (e *AppErr) Unwrap() error {
	return e.Err
}

func NewConflict(code string, err error) *AppErr {
	return &AppErr{
		Type: Conflict,
		Code: code,
		Err:  err,
	}
}

func NewConflictWithMsg(code string, msg string, err error) *AppErr {
	appErr := NewConflict(code, err)
	appErr.Message = &msg
	return appErr
}

func NewNotFound(code string, err error) *AppErr {
	return &AppErr{
		Type: NotFound,
		Code: code,
		Err:  err,
	}
}

func NewNotFoundWithMsg(code string, msg string, err error) *AppErr {
	appErr := NewNotFound(code, err)
	appErr.Message = &msg
	return appErr
}

func NewInvalid(code string, err error) *AppErr {
	return &AppErr{
		Type: Invalid,
		Code: code,
		Err:  err,
	}
}

func NewInvalidWithMsg(code string, msg string, err error) *AppErr {
	appErr := NewInvalid(code, err)
	appErr.Message = &msg
	return appErr
}

func NewBadRequest(code string, err error) *AppErr {
	return &AppErr{
		Type: BadRequest,
		Code: code,
		Err:  err,
	}
}

func NewBadRequestWithMsg(code string, msg string, err error) *AppErr {
	appErr := NewBadRequest(code, err)
	appErr.Message = &msg
	return appErr
}

func NewUnauthorized(code string, err error) *AppErr {
	return &AppErr{
		Type: Unauthorized,
		Code: code,
		Err:  err,
	}
}

package apimodels

type ProblemDetails struct {
	Type     *string `json:"type,omitempty"`
	Title    *string `json:"title,omitempty"`
	Status   int     `json:"status,omitempty"`
	Detail   *string `json:"detail,omitempty"`
	Instance *string `json:"instance,omitempty"`
}

type ValidationProblemDetails struct {
	ProblemDetails
	Errors []string `json:"errors"`
}

package model

type ErrorResponse struct {
	Error            bool                     `json:"error"`
	Message          string                   `json:"message"`
	ValidationErrors []*ValidateErrorResponse `json:"validation,omitempty"`
}

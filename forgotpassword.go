package app

import (
	"context"

	lferr "github.com/lofyd/common-lib-errors"
)

// ForgotpasswordRequest - http request querystring data
type ForgotPasswordRequest struct {
	Email string `json:"email,omitempty"`
}

// SignupResponse - http response body for signup
type ForgotPasswordResponse struct {
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CognitoData struct {
	UserPoolID      string `json:"userPoolId,omitempty"`
	AppClientID     string `json:"appClientId,omitempty"`
	AppClientSecret string `json:"appClientSecret,omitempty"`
}

type LambdaHandler interface {
	Start(handler interface{})
	Respond(statusCode, body string) (string, error)
}

type ForgotPasswordService interface {
	ForgotPassword(ctx context.Context, request ForgotPasswordRequest) (*ForgotPasswordResponse, lferr.Error)
}

package mock

import (
	"context"

	app "github.com/lofyd/auth/api/forgotpassword-func"
	lferr "github.com/lofyd/common-lib-errors"
)

// This ensures the ForgotPasswordService type implements the ForgotPasswordService interface via a compiler check,
// even if it is not used elsewhere.  You can read more on this pattern on the effective go site
// https://golang.org/doc/effective_go#blank_implements
var _ app.ForgotPasswordService = (*ForgotPasswordService)(nil)

type ForgotPasswordService struct {
	SignupFn func(ctx context.Context, req app.ForgotPasswordRequest) (*app.ForgotPasswordResponse, lferr.Error)
}

func (s *ForgotPasswordService) ForgotPassword(ctx context.Context, request app.ForgotPasswordRequest) (*app.ForgotPasswordResponse, lferr.Error) {
	return s.SignupFn(ctx, request)
}

package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"go.uber.org/zap"

	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	app "github.com/lofyd/auth/api/forgotpassword-func"
	lferr "github.com/lofyd/common-lib-errors"
	rezap "github.com/lofyd/common-lib-rezap"
)

func NewForgotPasswordService(client *cognito.Client, data *app.CognitoData) *ForgtoPasswordService {
	return &ForgtoPasswordService{
		CognitoClient: client,
		CognitoData:   data,
	}
}

type ForgtoPasswordService struct {
	CognitoClient *cognito.Client
	CognitoData   *app.CognitoData
}

func (s *ForgtoPasswordService) ForgotPassword(ctx context.Context, req app.ForgotPasswordRequest) (*app.ForgotPasswordResponse, lferr.Error) {
	log := rezap.WithContext(ctx)
	log = log.With(zap.String("email", req.Email), zap.String(("appClientId"), s.CognitoData.AppClientID))

	secretHash := generateSecretHash(req.Email, s.CognitoData.AppClientID, s.CognitoData.AppClientSecret)

	_, err := s.CognitoClient.ForgotPassword(ctx, &cognito.ForgotPasswordInput{
		ClientId:   aws.String(s.CognitoData.AppClientID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(req.Email),
	})
	if err != nil {
		var usernameExistsErr *types.UsernameExistsException
		if errors.As(err, &usernameExistsErr) {
			return nil, lferr.New(lferr.AlreadyExists, *usernameExistsErr.Message)
		}

		log.Errorw("error occurred during forgot password", zap.Error(err))
		return nil, lferr.New(lferr.Internal, "An error occurred during forgot password request. Please try again.")
	}

	return &app.ForgotPasswordResponse{}, nil
}

func generateSecretHash(username, clientId, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	secretHash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return secretHash
}

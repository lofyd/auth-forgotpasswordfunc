package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"go.uber.org/zap"

	app "github.com/lofyd/auth/api/forgotpassword-func"
	"github.com/lofyd/auth/api/forgotpassword-func/apigw"
	"github.com/lofyd/auth/api/forgotpassword-func/service"
	rezap "github.com/lofyd/common-lib-rezap"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awsxray "github.com/aws/aws-xray-sdk-go/xray"
)

const STAGE_ENV_VAR_KEY = "STAGE"

func run() error {
	clf := rezap.ContextLoggerFunc(func(ctx context.Context, l *zap.SugaredLogger) *zap.SugaredLogger {
		slog := l
		if traceID := awsxray.TraceID(ctx); traceID != "" {
			slog = slog.With(zap.String("traceId", traceID))
		}
		return slog
	})
	ctx := context.Background()
	slog, err := rezap.InitLoggingFromEnv(clf)
	if err != nil {
		return err
	}
	defer slog.Sync() // nolint

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	cd, err := getCognitoSecrets(ctx, cfg)
	if err != nil {
		return err
	}

	svc := service.NewForgotPasswordService(client, cd)

	rh := apigw.NewRequestHandler(svc)

	slog.Infow("lambda started", "runtime", runtime.Version())
	lambda.Start(rh.HandleRequest)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func getCognitoSecrets(ctx context.Context, cfg aws.Config) (*app.CognitoData, error) {
	log := rezap.WithContext(ctx)

	stage := os.Getenv(STAGE_ENV_VAR_KEY)
	if stage == "" {
		msg := fmt.Sprintf("missing environment variable: %s", STAGE_ENV_VAR_KEY)
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	cognitoSecretKey := fmt.Sprintf("/lofyd/%s/mobileapp/cognito", stage)

	sm := secretsmanager.NewFromConfig(cfg)
	result, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: aws.String(cognitoSecretKey)})
	if err != nil {
		log.Errorw("error getting secret value", zap.String("secretKey", cognitoSecretKey), zap.Error(err))
		return nil, err
	}

	if result.SecretString == nil {
		log.Errorw("secret not found", zap.String("secretKey", cognitoSecretKey))
		return nil, fmt.Errorf("secret not found: %s", cognitoSecretKey)
	}

	var cd app.CognitoData

	err = json.Unmarshal([]byte(*result.SecretString), &cd)
	if err != nil {
		log.Errorw("error unmarshalling cognito secrets", zap.Error(err))
		return nil, err
	}

	return &cd, nil
}

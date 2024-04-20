package apigw

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	app "github.com/lofyd/auth/api/forgotpassword-func"
	"github.com/lofyd/auth/api/forgotpassword-func/service"
	rezap "github.com/lofyd/common-lib-rezap"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/events"
	lferr "github.com/lofyd/common-lib-errors"
)

type RequestHandler struct {
	svc *service.ForgtoPasswordService
}

func NewRequestHandler(svc *service.ForgtoPasswordService) *RequestHandler {
	return &RequestHandler{
		svc: svc,
	}
}

func (rh RequestHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log := rezap.WithContext(ctx)

	var forgotPasswordReq app.ForgotPasswordRequest
	err := json.Unmarshal([]byte(request.Body), &forgotPasswordReq)
	if err != nil {
		log.Errorw("error unmarshalling forgot password request body", zap.Error(err))
		return rh.respondError(lferr.New(lferr.UnprocessableResource, "invalid request body")), nil
	}

	e := validateRequest(forgotPasswordReq)
	if e != nil {
		log.Errorf("validation error: %v", e)
		return rh.respondError(e), nil
	}

	results, suErr := rh.svc.ForgotPassword(ctx, forgotPasswordReq)
	if suErr != nil {
		return rh.respondError(suErr), nil
	}

	marshalRsp, err := json.Marshal(results)
	if err != nil {
		errorId := uuid.New()
		log.Errorw("error marshalling response", zap.Error(err), zap.String("errorId", errorId.String()))
		newErr := lferr.New(lferr.Internal, "error processing request: "+errorId.String())
		rh.respondError(newErr)
	}

	return rh.respond(http.StatusOK, string(marshalRsp), nil)
}

func (rh RequestHandler) respond(statusCode int, body string, headers map[string]string) (*events.APIGatewayProxyResponse, error) {
	isBase64Encode := true
	_, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		// error means body is not base64
		isBase64Encode = false
	}
	rsp := &events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		Body:            body,
		Headers:         headers,
		IsBase64Encoded: isBase64Encode,
	}
	return rsp, nil
}

func (rh RequestHandler) respondError(e lferr.Error) *events.APIGatewayProxyResponse {
	statusCode := lferr.AsStatusCodeOrDefault(e, http.StatusInternalServerError)
	msg := e.Error()

	ed, ok := e.(*lferr.ErrorDetail)
	if ok && ed.Detail != "" {
		msg = ed.Detail
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		Body:            fmt.Sprintf("{\"error\":\"%s\"}", msg),
		IsBase64Encoded: false,
	}
}

func validateRequest(req app.ForgotPasswordRequest) lferr.Error {
	if req.Email == "" {
		return lferr.New(lferr.FailedPrecondition, "Email is required")
	}

	return nil
}

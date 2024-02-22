package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/common"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
	ot "github.com/lengocson131002/go-clean/pkg/trace/opentelemetry"
	"github.com/lengocson131002/go-clean/pkg/util"
)

// ERROR HANDLING FOR RECOVERING FROM PANIC
type ErrorHandlingBehavior struct {
	logger logger.Logger
}

func NewErrorHandlingBehavior(logger logger.Logger) *ErrorHandlingBehavior {
	return &ErrorHandlingBehavior{
		logger: logger,
	}
}

func (b *ErrorHandlingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (res interface{}, err error) {
	// TODO: recover from error panic to prevent stop application
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("Recovered from panic: %v", r)
			err = common.ErrInternalServer
		}
	}()
	response, err := next(ctx)
	return response, err
}

// TRACING
type RequestTracingBehavior struct {
	logger logger.Logger
	tracer *ot.OpenTelemetryTracer
}

func NewTracingBehavior(logger logger.Logger, tracer *ot.OpenTelemetryTracer) *RequestTracingBehavior {
	return &RequestTracingBehavior{
		logger: logger,
		tracer: tracer,
	}
}

func (b *RequestTracingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (interface{}, error) {
	reqType := util.GetType(request)
	opName := fmt.Sprintf("request pipeline - %s", reqType)
	traceCtx, span := b.tracer.StartSpanFromContext(ctx, opName)
	defer span.End()

	response, err := next(traceCtx)
	return response, err
}

// METRICS
type RequestMetricBehavior struct {
	logger   logger.Logger
	metricer *Metricer
}

func NewMetricBehavior(logger logger.Logger, metricer *Metricer) *RequestMetricBehavior {
	return &RequestMetricBehavior{
		logger:   logger,
		metricer: metricer,
	}
}

func (b *RequestMetricBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (response interface{}, err error) {
	defer func() {
		reqType := util.GetType(request)
		reqTypeCounter := b.metricer.requestCountMetrics.With(METRIC_LABEL_REQUEST_TYPE, reqType)
		if err == nil {
			reqTypeCounter.With(METRIC_LABEL_REQUEST_STATUS, METRIC_LABEL_VALUE_REQUEST_STATUS_SUCCESS).Add(1)
		} else {
			reqTypeCounter.With(METRIC_LABEL_REQUEST_STATUS, METRIC_LABEL_VALUE_REQUEST_STATUS_ERROR).Add(1)
		}
	}()

	response, err = next(ctx)
	return response, err
}

// LOGGING
type RequestLoggingBehavior struct {
	logger logger.Logger
}

func NewRequestLoggingBehavior(logger logger.Logger) *RequestLoggingBehavior {
	return &RequestLoggingBehavior{
		logger: logger,
	}
}

func (b *RequestLoggingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (response interface{}, err error) {
	isSuccess := true
	start := time.Now()

	defer func() {
		if err != nil {
			isSuccess = false
		}

		requestJson, _ := json.Marshal(request)
		responseJson, _ := json.Marshal(response)
		errJson, _ := json.Marshal(err)

		b.logger.Info("Success: %t, Request: %s, Response: %s. Error: %s, Duration: %dms",
			isSuccess,
			string(requestJson),
			string(responseJson),
			string(errJson),
			time.Since(start).Milliseconds())
	}()

	response, err = next(ctx)

	return response, err
}

func RegisterPipeline(
	// request handlers
	verifyUserHandler domain.VerifyUserHandler,
	loginUserHandler domain.LoginUserHandler,
	createUserHandler domain.CreateUserHandler,
	getUserHandler domain.GetUserHandler,
	logoutUserHandler domain.LogoutUserHandler,
	updateUserHandler domain.UpdateUserHandler,
	openAccountHandler domain.OpenAccountHandler,

	// request behaviors
	requestLoggingBehavior *RequestLoggingBehavior,
	requestTracingBehavior *RequestTracingBehavior,
	requestMetricBehavior *RequestMetricBehavior,
	errorHandlingBehavior *ErrorHandlingBehavior,

) {
	// Register request handlers
	pipeline.RegisterRequestHandler[*domain.VerifyUserRequest, *domain.VerifyUserResponse](verifyUserHandler)
	pipeline.RegisterRequestHandler[*domain.CreateUserRequest, *domain.CreateUserResponse](createUserHandler)
	pipeline.RegisterRequestHandler[*domain.GetUserRequest, *domain.GetUserResponse](getUserHandler)
	pipeline.RegisterRequestHandler[*domain.LogoutUserRequest, bool](logoutUserHandler)
	pipeline.RegisterRequestHandler[*domain.UpdateUserRequest, *domain.UpdateUserResponse](updateUserHandler)
	pipeline.RegisterRequestHandler[*domain.LoginUserRequest, *domain.LoginUserResponse](loginUserHandler)
	pipeline.RegisterRequestHandler[*domain.OpenAccountRequest, *domain.OpenAccountResponse](openAccountHandler)

	// Register request behaviors
	pipeline.RegisterRequestPipelineBehaviors(requestTracingBehavior, requestLoggingBehavior, requestMetricBehavior, errorHandlingBehavior)

}

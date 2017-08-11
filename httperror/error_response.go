package httperror

import (
	"fmt"
	"net/http"
	"time"

	"code.cloudfoundry.org/lager"
)

const HTTP_ERROR_METRIC_NAME = "http_error"

//go:generate counterfeiter -o ../fakes/metrics_sender.go --fake-name MetricsSender . metricsSender
type metricsSender interface {
	SendDuration(string, time.Duration)
	IncrementCounter(string)
}

type ErrorResponse struct {
	MetricsSender metricsSender
}

func (e *ErrorResponse) InternalServerError(logger lager.Logger, w http.ResponseWriter, err error, description string) {
	e.respondWithCode(http.StatusInternalServerError, logger, w, err, description)
}

func (e *ErrorResponse) BadRequest(logger lager.Logger, w http.ResponseWriter, err error, description string) {
	e.respondWithCode(http.StatusBadRequest, logger, w, err, description)
}

func (e *ErrorResponse) Forbidden(logger lager.Logger, w http.ResponseWriter, err error, description string) {
	e.respondWithCode(http.StatusForbidden, logger, w, err, description)
}

func (e *ErrorResponse) Unauthorized(logger lager.Logger, w http.ResponseWriter, err error, description string) {
	e.respondWithCode(http.StatusUnauthorized, logger, w, err, description)
}

func (e *ErrorResponse) Conflict(logger lager.Logger, w http.ResponseWriter, err error, description string) {
	e.respondWithCode(http.StatusConflict, logger, w, err, description)
}

func (e *ErrorResponse) NotAcceptable(logger lager.Logger, w http.ResponseWriter, err error, description string) {
	e.respondWithCode(http.StatusNotAcceptable, logger, w, err, description)
}

func (e *ErrorResponse) respondWithCode(statusCode int, logger lager.Logger, w http.ResponseWriter, err error, description string) {
	logger.Error(fmt.Sprintf("%s", description), err)
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

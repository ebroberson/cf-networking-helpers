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
	Logger        lager.Logger
	MetricsSender metricsSender
}

func (e *ErrorResponse) InternalServerError(w http.ResponseWriter, err error, message, description string) {
	e.Logger.Error(fmt.Sprintf("%s: %s", message, description), err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s: %s"}`, message, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

func (e *ErrorResponse) BadRequest(w http.ResponseWriter, err error, message, description string) {
	e.Logger.Error(fmt.Sprintf("%s: %s", message, description), err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s: %s"}`, message, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

func (e *ErrorResponse) Forbidden(w http.ResponseWriter, err error, message, description string) {
	e.Logger.Error(fmt.Sprintf("%s: %s", message, description), err)
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s: %s"}`, message, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

func (e *ErrorResponse) Unauthorized(w http.ResponseWriter, err error, message, description string) {
	e.Logger.Error(fmt.Sprintf("%s: %s", message, description), err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s: %s"}`, message, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

func (e *ErrorResponse) Conflict(w http.ResponseWriter, err error, message, description string) {
	e.Logger.Error(fmt.Sprintf("%s: %s", message, description), err)
	w.WriteHeader(http.StatusConflict)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s: %s"}`, message, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

func (e *ErrorResponse) NotAcceptable(w http.ResponseWriter, err error, message, description string) {
	e.Logger.Error(fmt.Sprintf("%s: %s", message, description), err)
	w.WriteHeader(http.StatusNotAcceptable)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s: %s"}`, message, description)))
	e.MetricsSender.IncrementCounter(HTTP_ERROR_METRIC_NAME)
}

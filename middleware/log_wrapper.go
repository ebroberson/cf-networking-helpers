package middleware

import (
	"context"
	"net/http"

	"code.cloudfoundry.org/lager"
)

type Key string

const LoggerKey = Key("logger")

func LogWrap(logger lager.Logger, wrappedHandler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		requestLogger := logger.Session("request", lager.Data{
			"method":  r.Method,
			"request": r.URL.String(),
		})
		contextWithLogger := context.WithValue(r.Context(), LoggerKey, requestLogger)
		r = r.WithContext(contextWithLogger)

		requestLogger.Debug("serving")
		defer requestLogger.Debug("done")

		wrappedHandler(w, r)
	}
}

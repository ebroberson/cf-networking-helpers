package middleware

import (
	"net/http"
	"code.cloudfoundry.org/lager"

	semver "github.com/hashicorp/go-version"
	"fmt"
)

//go:generate counterfeiter -o fakes/error_response.go --fake-name ErrorResponse . errorResponse
type errorResponse interface {
	NotAcceptable(http.ResponseWriter, error, string, string)
}

type CheckVersionWrapper struct {
	ErrorResponse errorResponse
}

func (c *CheckVersionWrapper) CheckVersion(handler LoggableHandlerFunc) LoggableHandlerFunc {
	return LoggableHandlerFunc(func(logger lager.Logger, rw http.ResponseWriter, r *http.Request) {
		version := r.Header["Accept"][0]

		v1, err := semver.NewVersion(version)
		if err != nil {
			c.ErrorResponse.NotAcceptable(rw, nil, "check api version", fmt.Sprintf("api version '%s' not supported", version))
			return
		}

		// Constraints example.
		constraints, err := semver.NewConstraint(">= 1.0.0, < 2.0")
		if constraints.Check(v1) {
			handler(logger, rw, r)
		} else {
			c.ErrorResponse.NotAcceptable(rw, nil, "check api version", fmt.Sprintf("api version '%s' not supported", version))
		}
	})
}

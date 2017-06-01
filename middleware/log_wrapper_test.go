package middleware_test

import (
	"net/http"

	"code.cloudfoundry.org/cf-networking-helpers/middleware"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var LogsWith = func(level lager.LogLevel, msg string) types.GomegaMatcher {
	return And(
		WithTransform(func(log lager.LogFormat) string {
			return log.Message
		}, Equal(msg)),
		WithTransform(func(log lager.LogFormat) lager.LogLevel {
			return log.LogLevel
		}, Equal(level)),
	)
}

var HaveLogData = func(nextMatcher types.GomegaMatcher) types.GomegaMatcher {
	return WithTransform(func(log lager.LogFormat) lager.Data {
		return log.Data
	}, nextMatcher)
}

var _ = Describe("LogWrap", func() {
	var (
		logger              *lagertest.TestLogger
		loggableHandlerFunc middleware.LoggableHandlerFunc
	)
	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test-session")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))
		loggableHandlerFunc = func(logger lager.Logger, w http.ResponseWriter, r *http.Request) {
			logger = logger.Session("logger-group")
			logger.Info("written-in-loggable-handler")
		}
	})
	It("creates \"request\" session and passes it to LoggableHandlerFunc", func() {
		handler := middleware.LogWrap(logger, loggableHandlerFunc)
		req, err := http.NewRequest("GET", "http://example.com", nil)
		Expect(err).NotTo(HaveOccurred())
		handler.ServeHTTP(nil, req)

		Expect(logger.Logs()).To(HaveLen(3))
		Expect(logger.Logs()[0]).To(SatisfyAll(
			LogsWith(lager.DEBUG, "test-session.request.serving"),
			HaveLogData(SatisfyAll(
				HaveKeyWithValue("session", Equal("1")),
				HaveKeyWithValue("method", Equal("GET")),
				HaveKeyWithValue("request", Equal("http://example.com")),
			)),
		))

		Expect(logger.Logs()[1]).To(SatisfyAll(
			LogsWith(lager.INFO, "test-session.request.logger-group.written-in-loggable-handler"),
			HaveLogData(HaveKeyWithValue("session", Equal("1.1"))),
		))

		Expect(logger.Logs()[2]).To(SatisfyAll(
			LogsWith(lager.DEBUG, "test-session.request.done"),
			HaveLogData(SatisfyAll(
				HaveKeyWithValue("session", Equal("1")),
				HaveKeyWithValue("method", Equal("GET")),
				HaveKeyWithValue("request", Equal("http://example.com")),
			)),
		))
	})
})

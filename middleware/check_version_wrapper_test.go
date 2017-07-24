package middleware_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"bytes"
	"code.cloudfoundry.org/lager"
	"net/http/httptest"
	"code.cloudfoundry.org/lager/lagertest"

	"code.cloudfoundry.org/cf-networking-helpers/middleware/fakes"
	"code.cloudfoundry.org/cf-networking-helpers/middleware"
)

var _ = FDescribe("CheckVersionWrapper", func() {
	var checkVersionHandler middleware.LoggableHandlerFunc
	var request *http.Request
	var fakeHandler *fakeLoggableHandler
	var resp *httptest.ResponseRecorder
	var logger *lagertest.TestLogger
	var fakeErrorResponse *fakes.ErrorResponse

	BeforeEach(func() {
		var err error
		request, err = http.NewRequest("GET", "/some/resource", bytes.NewBuffer([]byte(`{}`)))
		Expect(err).NotTo(HaveOccurred())
		request.Header["Accept"] = []string{"1.0.0+policy-server-json"}

		fakeHandler = &fakeLoggableHandler{}
		fakeErrorResponse = &fakes.ErrorResponse{}

		logger = lagertest.NewTestLogger("test")

		resp = httptest.NewRecorder()
		checkVersionWrapper := middleware.CheckVersionWrapper{ErrorResponse: fakeErrorResponse}
		checkVersionHandler = checkVersionWrapper.CheckVersion(fakeHandler.LoggableHandler)

	})

	It("should delegate to handler if version is supported", func() {
		checkVersionHandler(logger, resp, request)

		Expect(fakeErrorResponse.NotAcceptableCallCount()).To(Equal(0))
		Expect(fakeHandler.invocationCount).To(Equal(1))
		Expect(fakeHandler.actualLogger).To(Equal(logger))
		Expect(fakeHandler.actualWriter).To(Equal(resp))
		Expect(fakeHandler.actualRequest).To(Equal(request))
	})

	Context("when the requested version has a different major version", func() {
		BeforeEach(func() {
			request.Header["Accept"] = []string{"6.2.3+policy-server-json"}
		})
		It("Rejects the request with a 406 status code", func() {
			checkVersionHandler(logger, resp, request)

			Expect(fakeErrorResponse.NotAcceptableCallCount()).To(Equal(1))
		})
	})

	Context("when the requested version has the same major but different minor version", func() {
		BeforeEach(func() {
			request.Header["Accept"] = []string{"1.2.3+policy-server-json"}
		})
		It("allow the call to the wrapped handler", func() {
			checkVersionHandler(logger, resp, request)

			Expect(fakeErrorResponse.NotAcceptableCallCount()).To(Equal(0))
			Expect(fakeHandler.invocationCount).To(Equal(1))
		})
	})

	Context("when the api header version is not compatible", func() {
		BeforeEach(func() {
			request.Header["Accept"] = []string{"0.0.0+policy-server-json"}
		})

		It("calls the 406 Not Acceptable handler", func() {
			checkVersionHandler(logger, resp, request)

			Expect(fakeErrorResponse.NotAcceptableCallCount()).To(Equal(1))
			rw, err, message, desc := fakeErrorResponse.NotAcceptableArgsForCall(0)
			Expect(rw).To(Equal(resp))
			Expect(err).To(BeNil())
			Expect(message).To(Equal("check api version"))
			Expect(desc).To(Equal("api version '0.0.0+policy-server-json' not supported"))

		})

		XContext("when multiple accept values are provided", func() {
			BeforeEach(func() {
				request.Header["Accept"] = []string{"1.0.0", "2.0.0"}
			})

			It("should return a sensible error", func() {
				checkVersionHandler(logger, resp, request)


			})
		})

		Context("when no accept header is provided", func() {
			BeforeEach(func() {
				delete(request.Header, "Accept")
			})

			It("should return a sensible error", func() {
				checkVersionHandler(logger, resp, request)

				Expect(fakeErrorResponse.NotAcceptableCallCount()).To(Equal(1))
				rw, err, message, desc := fakeErrorResponse.NotAcceptableArgsForCall(0)
				Expect(rw).To(Equal(resp))
				Expect(err).To(BeNil())
				Expect(message).To(Equal("check api version"))
				Expect(desc).To(Equal("api version not provided. please update your client."))
			})
		})
	})

	Context("when the version is not valid", func(){
		BeforeEach(func(){
			request.Header["Accept"] = []string{"banana"}
		})
		It("returns a 406 error", func() {
			checkVersionHandler(logger, resp, request)

			Expect(fakeErrorResponse.NotAcceptableCallCount()).To(Equal(1))
		})
	})
})

type fakeLoggableHandler struct {
	invocationCount int
	actualLogger    lager.Logger
	actualWriter    http.ResponseWriter
	actualRequest   *http.Request
}

func (f *fakeLoggableHandler) LoggableHandler(logger lager.Logger, w http.ResponseWriter, r *http.Request) {
	f.invocationCount++
	f.actualLogger = logger
	f.actualWriter = w
	f.actualRequest = r
}

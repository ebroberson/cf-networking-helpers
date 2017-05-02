package testsupport_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"code.cloudfoundry.org/go-db-helpers/testsupport"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BadTcpProxy", func() {
	var (
		testHandler       *testsupport.TestHandler
		destinationServer *httptest.Server
		badProxy          *testsupport.BadTCPProxy
		proxyURL          string
	)

	BeforeEach(func() {
		testHandler = &testsupport.TestHandler{}

		destinationServer = httptest.NewServer(testHandler)

		var err error
		badProxy, err = testsupport.NewBadTCPProxy(destinationServer.Listener.Addr().String())
		Expect(err).NotTo(HaveOccurred())

		proxyURL = fmt.Sprintf("http://%s", badProxy.ListenAddress())
	})

	AfterEach(func() {
		destinationServer.Close()
		badProxy.Close()

		Expect(testHandler.NumRequestsInFlight()).To(Equal(0))
	})

	It("proxies its inputs to its outputs", func(done Done) {
		resp, err := http.Post(proxyURL, "text/plain", strings.NewReader("hello world"))
		Expect(err).NotTo(HaveOccurred())
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(string(bodyBytes)).To(Equal("HELLO WORLD"))

		close(done)
	}, 2 /* timeout in seconds */)

	FDescribe("partition", func() {
		It("halts all traffic", func(done Done) {
			requestCompleted := make(chan error)
			go func() {
				defer GinkgoRecover()
				resp, err := http.Post(proxyURL, "text/plain", strings.NewReader("hello! world"))
				Expect(err).NotTo(HaveOccurred())
				_, err = ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				requestCompleted <- err
			}()
			Consistently(requestCompleted, "1s").ShouldNot(Receive())

			Expect(testHandler.NumRequestsInFlight()).To(Equal(1))
			close(done)
		}, 2 /* timeout in seconds */)
	})
})

package internal

import (
	"fmt"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

//tutorial
//https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318
//https://www.youtube.com/watch?v=sOeUf1YICSA&t=49s
//https://betterprogramming.pub/easy-guide-to-unit-testing-in-golang-4fc1e9d96679

// Auxiliary function to factor code for checking that the function returns the right value
func auxTestParseQueryParametersValid(query string, expectedQueryParameters queryParameters) {
	u, err := url.Parse(query)
	Ω(err).Should(BeNil())
	q := u.Query()
	Ω(parseQueryParameters(q)).Should(BeIdenticalTo(expectedQueryParameters))
}

func TestParseQueryParametersThreeParam(t *testing.T) {
	RegisterTestingT(t)
	expectedParam := queryParameters{
		namespace: "ns1",
		start:     "2",
		end:       "1",
	}
	auxTestParseQueryParametersValid("https://example.org/?namespace=ns1&start=2&end=1", expectedParam)
}

func TestParseQueryParametersTwoParam(t *testing.T) {
	RegisterTestingT(t)
	expectedParam2 := queryParameters{
		start: "2",
		end:   "1",
	}
	auxTestParseQueryParametersValid("https://example.org/?start=2&end=1", expectedParam2)
}

func TestServerFailsWhenAddressEnvNotPresent(t *testing.T) {
	RegisterTestingT(t)
	Ω(os.Setenv("PROMETHEUS_URL", "test")).Should(Succeed(), "Setup failed")
	Ω(Server).Should(PanicWith(ContainSubstring("SERVE_ADDRESS")))
}

func TestServerSendsCorrectHeaders(t *testing.T) {
	RegisterTestingT(t)
	api := doNothing(t)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/test?start=%s&end=%s", validParameters.start, validParameters.end), nil)
	handlerFactory("test_query", api)(writer, req)
	Ω(writer).Should(HaveHTTPStatus(http.StatusOK))
	// CORS
	Ω(writer).Should(HaveHTTPHeaderWithValue("Access-Control-Allow-Origin", "*"))
	// JSON response
	Ω(writer).Should(HaveHTTPHeaderWithValue("Content-Type", "application/json; charset=utf-8"))
}

func TestParsesNamespace(t *testing.T) {
	RegisterTestingT(t)
	ns := "a_namespace"
	writer := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/test?start=%s&end=%s&namespace=%s", validParameters.start, validParameters.end, ns), nil)

	api := doQuery(t, func(query string) {
		Ω(query).Should(ContainSubstring(fmt.Sprintf("namespace=~\"%s\"", ns)))
	})
	handlerFactory("test_query{namespace=~\"%s\"}", api)(writer, req)
	Ω(writer).Should(HaveHTTPStatus(http.StatusOK))
}

func TestHasEmptyNamespaceWhenNotGiven(t *testing.T) {
	RegisterTestingT(t)
	api := doQuery(t, func(query string) {
		Ω(query).Should(ContainSubstring("namespace!=\"\""))
	})
	writer := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/test?start=%s&end=%s", validParameters.start, validParameters.end), nil)
	handlerFactory("test_query{namespace=~\"%s\"}", api)(writer, req)
}

func TestMissingTimestampSendsError(t *testing.T) {
	RegisterTestingT(t)
	api := notCalled(t)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	handlerFactory("test_query{namespace=~\"%s\"}", api)(writer, req)
	Ω(writer).Should(HaveHTTPStatus(http.StatusBadRequest))
}

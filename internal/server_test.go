package internal

import (
	"net/url"
	"testing"
)

//tutorial
//https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318
//https://www.youtube.com/watch?v=sOeUf1YICSA&t=49s
//https://betterprogramming.pub/easy-guide-to-unit-testing-in-golang-4fc1e9d96679

/// Auxiliary function to factor code for checking that the function throws an error
func auxTestParseQueryParametersFailing(query string, t *testing.T) {
	u, err := url.Parse(query)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	body, err := parseQueryParameters(q)
	if err == nil {
		t.Error("Error expected")
	}
	_ = body
}

/// Auxiliary function to factor code for checking that the function returns the right value
func auxTestParseQueryParametersValid(query string, expectedQueryParameters queryParameters, t *testing.T) {
	u, err := url.Parse(query)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	body, err := parseQueryParameters(q)
	if err != nil {
		t.Error("Error unexpected")
	}
	if body != expectedQueryParameters {
		t.Error("Unexpected results")
	}
}

func TestParseQueryParametersWrongQuery(t *testing.T) {
	auxTestParseQueryParametersFailing("https://example.org/?a=1&a=2&b=&=3&&&&", t)
}

func TestParseQueryParametersEmptyQuery(t *testing.T) {
	auxTestParseQueryParametersFailing("", t)
}

func TestParseQueryParametersThreeParam(t *testing.T) {

	expectedParam := queryParameters{
		namespace:    "ns1",
		timeInterval: "2",
		step:         "1",
	}
	auxTestParseQueryParametersValid("https://example.org/?namespace=ns1&interval=2&step=1", expectedParam, t)
}

func TestParseQueryParametersTwoParam(t *testing.T) {
	expectedParam2 := queryParameters{
		timeInterval: "2",
		step:         "1",
	}
	auxTestParseQueryParametersValid("https://example.org/?interval=2&step=1", expectedParam2, t)
}
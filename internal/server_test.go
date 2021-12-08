package internal

import (
	"net/url"
	"testing"
)

//tutorial
//https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318
//https://www.youtube.com/watch?v=sOeUf1YICSA&t=49s
//https://betterprogramming.pub/easy-guide-to-unit-testing-in-golang-4fc1e9d96679

//test template

func myTestCase(urlQuery url.Values, expectedValue url.Values, t *testing.T) {
	output, err := parseQueryParameters(urlQuery)
	if err != nil {
		t.Errorf("got %v, want %v", output, expectedValue)
	}
	if output != expectedValue {
		t.Errorf("got %v, want %v", output, expectedValue)
	}
}

func TestParameter(t *testing.T) {
	myString := [2]string{"interval", "step"}
	got := parseQueryParameters("interval")
	want := "interval"
	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, myString)
	}
}

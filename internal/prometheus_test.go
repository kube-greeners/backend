package internal

import (
	mocks "backend/internal/generated_mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var validParameters = queryParameters{
	"",
	"1638445155310",
	"1639654755310",
}

func notCalled(t *testing.T) prometheus {
	cntr := gomock.NewController(t)
	promApi := mocks.NewMockAPI(cntr)
	return prometheus{promApi}
}
func auxApi(t *testing.T, f func(context.Context, string, v1.Range)) prometheus {

	cntr := gomock.NewController(t)
	promApi := mocks.NewMockAPI(cntr)

	promApi.EXPECT().
		QueryRange(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.Matrix{}, nil, nil).
		Do(f)
	return prometheus{api: promApi}
}

func doRange(t *testing.T, f func(range_ v1.Range)) prometheus {
	return auxApi(t, func(ctx context.Context, query string, range_ v1.Range) {
		f(range_)
	})
}

func doQuery(t *testing.T, f func(query string)) prometheus {
	return auxApi(t, func(ctx context.Context, query string, range_ v1.Range) {
		f(query)
	})
}

func doNothing(t *testing.T) prometheus {
	return auxApi(t, func(ctx context.Context, query string, range_ v1.Range) {})
}

func TestNamespaceIsReplaced(t *testing.T) {
	RegisterTestingT(t)
	api := doQuery(t, func(query string) {
		Ω(query).Should(ContainSubstring("namespace=~\"test\""))
	})

	parameters := validParameters
	parameters.namespace = "test"
	_, _ = api.executeQuery("some_query{namespace=~\"%s\"}", parameters)
}

func TestEmptyNamespaceIsNegated(t *testing.T) {
	RegisterTestingT(t)
	api := doQuery(t, func(query string) {
		Ω(query).Should(ContainSubstring("namespace!=\"\""))
	})
	var query = "some_query{namespace=~\"%s\"}"
	_, _ = api.executeQuery(query, validParameters)
}

func helperInvalidParameter(t *testing.T, parameters queryParameters, errIncludes string) {
	RegisterTestingT(t)
	api := notCalled(t)
	res, err := api.executeQuery("some_query", parameters)
	Ω(res).Should(BeEmpty())
	Ω(err).Should(Not(BeNil()), "error needed since timing is wrongly formatted")
	Ω(strings.ToLower(err.Error())).Should(ContainSubstring(errIncludes), "wrong error is thrown")
}
func TestInvalidTimesFailQueries(t *testing.T) {
	RegisterTestingT(t)

	invalidStart := validParameters
	invalidStart.start = "asd"
	helperInvalidParameter(t, invalidStart, "invalid")

	invalidEnd := validParameters
	invalidEnd.end = "asd"
	helperInvalidParameter(t, invalidEnd, "invalid")

	negativeStart := validParameters
	negativeStart.start = "-1"
	helperInvalidParameter(t, negativeStart, "invalid")

	negativeEnd := validParameters
	negativeEnd.end = "-1"
	helperInvalidParameter(t, negativeEnd, "invalid")

}

func TestTooCloseTimes(t *testing.T) {
	RegisterTestingT(t)
	tooCloseStartEnd := validParameters
	tooCloseStartEnd.start = "1640139779221"
	tooCloseStartEnd.end = "1640139779223"
	helperInvalidParameter(t, tooCloseStartEnd, "too close")
}

func TestEndDateSmallerThanStart(t *testing.T) {

	RegisterTestingT(t)
	tooCloseStartEnd := validParameters
	tooCloseStartEnd.start = "1640139779221"
	tooCloseStartEnd.end = "164013"
	helperInvalidParameter(t, tooCloseStartEnd, "is smaller than")
}

func TestConstructorFinishesWithEnv(t *testing.T) {
	RegisterTestingT(t)
	Ω(os.Setenv("PROMETHEUS_URL", "test")).Should(Succeed(), "Setup failed")
	Ω(prometheusClient()).Error().Should(Succeed())
}

func TestConstructorFailsWithoutEnv(t *testing.T) {
	RegisterTestingT(t)
	Ω(os.Unsetenv("PROMETHEUS_URL")).Should(Succeed(), "Setup failed")
	Ω(prometheusClient()).Error().Should(MatchError(ContainSubstring("PROMETHEUS_URL")))
}

func TestIntervalComputedCorrectly(t *testing.T) {
	RegisterTestingT(t)
	start := time.Now()
	end := start.Add(3 * time.Hour)
	api := doRange(t, func(range_ v1.Range) {
		expectedStep := 3 * time.Hour / 20
		Ω(range_.Step).Should(BeIdenticalTo(expectedStep))
	})
	parameters := validParameters
	parameters.start = strconv.FormatInt(start.UnixMilli(), 10)
	parameters.end = strconv.FormatInt(end.UnixMilli(), 10)
	_, _ = api.executeQuery("some_query{namespace=~\"%s\"}", parameters)
}

func TestTimestampsAreCorrect(t *testing.T) {
	RegisterTestingT(t)
	start := time.Now().Round(time.Second)
	end := start.Add(3 * time.Hour)
	api := doRange(t, func(range_ v1.Range) {
		if range_.Start != start {
			t.Error("Start timestamps don't match expected: " + start.String() + " got: " + range_.Start.String())
		}
		if range_.End != end {
			t.Error("End timestamps don't match expected: " + end.String() + " got: " + range_.End.String())
		}
	})
	parameters := validParameters
	parameters.start = strconv.FormatInt(start.UnixMilli(), 10)
	parameters.end = strconv.FormatInt(end.UnixMilli(), 10)
	_, _ = api.executeQuery("some_query{namespace=~\"%s\"}", parameters)
}

package internal

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	s "strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"golang.org/x/net/context"
)

type queryParameters struct {
	namespace     string
	start         string
	end           string
	measureTiming bool
}

type prometheus struct {
	api v1.API
}

func prometheusClient() (prometheus, error) {

	address := os.Getenv("PROMETHEUS_URL")
	if len(address) == 0 {
		return prometheus{}, errors.New("can't read PROMETHEUS_URL env variable, please set it")
	}
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return prometheus{}, err
	}
	return prometheus{api: v1.NewAPI(client)}, nil
}

func (client prometheus) rawQuery(query string, start time.Time, end time.Time, step time.Duration, measure bool) (string, error) {
	println(query)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	timer := time.Now()
	value, warning, err := client.api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})
	if err != nil {
		return "", err
	}
	if warning != nil {
		for warn := range warning {
			println(warn)
		}
	}
	diff := time.Now().Sub(timer)
	marshaledJson, err := json.Marshal(value)
	if measure {
		var resultDecoded = []map[string]json.RawMessage{}
		err := json.Unmarshal(marshaledJson, &resultDecoded)
		if err != nil {
			return "", err
		}
		resultDecoded[0]["time"], err = json.Marshal(diff.Milliseconds())
		if err != nil {
			return "", err
		}
		marshaledJson, err = json.Marshal(resultDecoded)
		if err != nil {
			return "", err
		}
	}
	if err != nil {
		return "", err
	}
	return string(marshaledJson), nil
}

func (client prometheus) executeQuery(query string, parameters queryParameters) (string, error) {

	for s.Contains(query, "\"%s\"") {
		if parameters.namespace != "" {
			query = s.Replace(query, "%s", parameters.namespace, 1)
		} else {
			query = s.Replace(query, "namespace=~\"%s\"", "namespace!=\"\"", 1)
		}
	}
	intStart, err := strconv.ParseInt(parameters.start, 0, 0)
	if intStart < 1 {
		return "", errors.New("Invalid start date: " + parameters.start + " because: non-positive timestamp")
	}
	timestampStart := time.Unix(intStart/1000, 0)
	if err != nil {
		return "", errors.New("Invalid start date: " + parameters.start + " because: " + err.Error())
	}

	intEnd, err := strconv.ParseInt(parameters.end, 0, 0)
	if intEnd < 1 {
		return "", errors.New("Invalid end date: " + parameters.end + " because: non-positive timestamp")
	}
	timestampEnd := time.Unix(intEnd/1000, 0)
	if err != nil {
		return "", errors.New("Invalid end date: " + parameters.end + " because: " + err.Error())
	}
	const datapointCount = 20
	step := (intEnd - intStart) / 1000 / datapointCount
	if step == 0 {
		return "", errors.New("start and end time are too close to one another, can't compute step")
	}
	if step < 0 {
		return "", errors.New("end time " + timestampEnd.String() + " is smaller than start time " + timestampStart.String())
	}

	return client.rawQuery(query, timestampStart, timestampEnd, time.Duration(step*time.Second.Nanoseconds()), parameters.measureTiming)

}

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	s "strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"golang.org/x/net/context"
)

type queryParameters struct {
	namespace    string
	// timeInterval string
	start string
	end string
	step         string
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

func (client prometheus) rawQuery(query string, start time.Time, end time.Time, step time.Duration) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
	marshaledJson, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(marshaledJson), nil
}

func parseInterval(input string) (time.Duration, error) {
	var amount int
	var interval byte
	nTokens, err := fmt.Sscanf(input, "%d%c", &amount, &interval)
	if err != nil {
		return 0, err
	}
	if nTokens != 2 {
		return 0, errors.New("The format of interval is not valid" + input)
	}
	var factor time.Duration
	switch interval {
	case 'd':
		factor = time.Hour * 24
	case 'h':
		factor = time.Hour
	case 'm':
		factor = time.Minute
	case 'w':
		factor = time.Hour * 24 * 7
	case 's':
		factor = time.Second
	default:
		return 0, errors.New("The format of date is not valid, valid formats are w, d, h, m, s: " + string(interval))
	}
	return factor * time.Duration(amount), nil

}

func (client prometheus) executeQuery(query string, parameters queryParameters) (string, error) {
	for s.Contains(query, "\"%s\"") {
		query = s.Replace(query, "%s", parameters.namespace, 1)
	}
	parsedInterval, err := parseInterval(parameters.timeInterval)
	if err != nil {
		return "", errors.New("Invalid interval parameter: " + parameters.timeInterval + " because: " + err.Error())
	}
	parsedStep, err := parseInterval(parameters.step)
	if err != nil {
		return "", errors.New("Invalid step parameter: " + parameters.step + " because: " + err.Error())
	}
	intStart, err := strconv.ParseInt(parameters.start, 0, 0)
	timestampStart := time.Unix(intStart/1000, 0)

	intEnd, err := strconv.ParseInt(parameters.end, 0, 0)
	timestampEnd := time.Unix(intEnd/1000, 0)
	return := client.rawQuery(query,timestampStart, timestampEnd, parsedStep)
	
}

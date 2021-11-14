package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"golang.org/x/net/context"
	"os"
	s "strings"
	"time"
)

type queryParameters struct {
	namespace    string
	timeInterval string
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

func (client prometheus) rawQuery(query string, interval time.Duration, step time.Duration) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	value, warning, err := client.api.QueryRange(ctx, query, v1.Range{
		Start: time.Now().Add(-interval),
		End:   time.Now(),
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
		return 0, errors.New("Invalid format for the interval string" + input)
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
		return 0, errors.New("Not a valid letter: " + string(interval))
	}
	return factor * time.Duration(amount), nil

}

func (client prometheus) executeQuery(query string, parameters queryParameters) (string, error) {
	if s.Contains(query, "\"%s\"") {
		query = fmt.Sprintf(query, parameters.namespace)
	}
	parsedInterval, err := parseInterval(parameters.timeInterval)
	if err != nil {
		return "", errors.New("Wrong interval: " + parameters.timeInterval + " because: " + err.Error())
	}
	parsedStep, err := parseInterval(parameters.step)
	if err != nil {
		return "", errors.New("Wrong step: " + parameters.step + " because: " + err.Error())
	}
	return client.rawQuery(query, parsedInterval, parsedStep)
}

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/cors"
	"net/http"
	"net/url"
	"os"
)

func logJSONError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if err != nil {
		panic(err)
	}
}

func parseQueryParameters(urlQuery url.Values) (queryParameters, error) {
	if !urlQuery.Has("interval") {
		return queryParameters{}, errors.New("Interval parameter is missing or is not valid.")
	}
	if !urlQuery.Has("step") {
		return queryParameters{}, errors.New("Step parameter is missing or is not valid.")
	}
	return queryParameters{
		namespace:    urlQuery.Get("namespace"),
		timeInterval: urlQuery.Get("interval"),
		step:         urlQuery.Get("step"),
	}, nil
}
func handlerFactory(query string, prometheusClient prometheus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		body, err := parseQueryParameters(r.URL.Query())
		if err != nil {
			logJSONError(w, err, http.StatusBadRequest)
			return
		}
		result, err := prometheusClient.executeQuery(query, body)
		if err != nil {
			logJSONError(w, err, http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		_, err = fmt.Fprint(w, result)
		if err != nil {
			panic(err)
		}
	}
}

func Server() {
	prometheusClient, err := prometheusClient()
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	for path := range queryDict {
		mux.HandleFunc("/"+path, handlerFactory(queryDict[path], prometheusClient))
	}
	address := os.Getenv("SERVE_ADDRESS")
	if len(address) == 0 {
		panic("define env variable SERVE_ADDRESS")
	}
	err = http.ListenAndServe(address, cors.Default().Handler(mux))
	if err != nil {
		panic(err)
	}
}

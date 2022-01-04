package internal

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/rs/cors"
)

func logJSONError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if err != nil {
		panic(err)
	}
}

func parseQueryParameters(urlQuery url.Values) (queryParameters, error) {
	return queryParameters{
		namespace: urlQuery.Get("namespace"),
		start:     urlQuery.Get("start"),
		end:       urlQuery.Get("end"),
	}, nil
}
func handlerFactory(query string, prometheusClient prometheus) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

//go:embed swaggerui
var swaggerFs embed.FS

func Server() {
	prometheusClient, err := prometheusClient()
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	for path := range queryDict {
		mux.HandleFunc("/"+path, handlerFactory(queryDict[path], prometheusClient))
	}
	mux.Handle("/swaggerui/", http.FileServer(http.FS(swaggerFs)))

	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/", fs)
	address := os.Getenv("SERVE_ADDRESS")
	if len(address) == 0 {
		panic("define env variable SERVE_ADDRESS")
	}
	err = http.ListenAndServe(address, cors.Default().Handler(mux))
	if err != nil {
		panic(err)
	}

}

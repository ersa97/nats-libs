package main

import (
	"context"
	"encoding/json"
	"net/http"

	natslibs "github.com/ersa97/nats-libs"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {
	r := mux.NewRouter()

	n := natslibs.New(context.Background(), nats.DefaultURL, "test")

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(rw, r)
		})
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		err := n.Publish("simple", body)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]any{
				"message": "NOT OK",
			})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"message": "OK",
		})
	}).Methods(http.MethodPost)

	http.ListenAndServe(":8000", r)
}

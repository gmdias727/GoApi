package main

import (
	"encoding/json"
	"net/http"
	mongo "src/database/mongo"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	client, ctx, cancel, err := mongo.Connect("mongodb+srv://Potato:PotatoChips@pudim.riqj4an.mongodb.net/?retryWrites=true&w=majority")

	if err != nil {
		panic(err)
	}

	defer mongo.Close(client, ctx, cancel)

	mongo.Ping(client, ctx)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Get("/test", makeAPIFunc(handleUser))

	http.ListenAndServe(":3000", r)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, map[string]string{"message": "Hello World!"})
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

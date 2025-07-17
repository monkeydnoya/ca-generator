package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type RequestHandler func(w http.ResponseWriter, r *http.Request) error

func (fn RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ManualHandler(w http.ResponseWriter, r *http.Request) error {
	var applications []CreditApplication
	err := json.NewDecoder(r.Body).Decode(&applications)
	if err != nil {
		return fmt.Errorf("failed to decode request body: %s", err.Error())
	}
	defer r.Body.Close()

	if err := GenerateManualCAs(r.Context(), applications); err != nil {
		return err
	}

	return nil
}

func LoadHandler(w http.ResponseWriter, r *http.Request) error {
	countString := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countString)
	if err != nil {
		return err
	}

	durationString := r.URL.Query().Get("duration")
	duraiton, err := strconv.Atoi(durationString)
	if err != nil {
		return err
	}

	duration, err := GenerateLoadCAs(r.Context(), count, duraiton)
	if err != nil {
		return err
	}

	data, err := json.Marshal(&duration)
	if err != nil {
		return err
	}

	w.Write(data)
	return nil
}

func RegisterHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/api/manual", RequestHandler(ManualHandler))
	mux.Handle("/api/load", RequestHandler(LoadHandler))

	return mux
}

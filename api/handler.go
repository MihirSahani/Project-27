package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) writeJSON(writer http.ResponseWriter, status int, data any) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if data == nil {
		writer.WriteHeader(http.StatusNoContent)
		return nil
	}
	return json.NewEncoder(writer).Encode(&data)
}

func (app *Application) readJSON(request *http.Request, dst any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(&dst)
}

func (app *Application) errorJSON(w http.ResponseWriter, status int) {
	app.writeJSON(w, status, nil)
}

package main

import "net/http"

func (app *Application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {

	data := map[string]string{
		"status":      "Up",
		"environment": "development",
		"version":     "0.0.1",
	}
	app.writeJSON(writer, http.StatusOK, &data)

}

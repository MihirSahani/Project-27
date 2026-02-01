package main

import (
	"net/http"

	"github.com/MihirSahani/Project-27/internal"
)

func (app *Application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":      "UP",
		"environment": internal.GetEnvAsString("ENVIRONMENT", "DEVELOPMENT"),
		"version":     "0.0.1",
	}

	app.writeJSON(writer, http.StatusOK, &data)
}

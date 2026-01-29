package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/MihirSahani/Project-27/internal"
	"github.com/MihirSahani/Project-27/storage/entity"
)


type AuthenticationPayload struct {
	Email string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (app *Application) authenticationHandler(writer http.ResponseWriter, request *http.Request)  {
	// get the payload
	var payload AuthenticationPayload
	err := app.readJSON(request, &payload)
	if err != nil {
		app.errorJSON(writer, http.StatusBadRequest)
		return
	}
	// validate the payload
	if err = internal.Validate.Struct(payload); err != nil {
		app.errorJSON(writer, http.StatusBadRequest)
		return
	}
	data, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.UserStorageManager.GetUserByEmail(ctx, tx, payload.Email)
	})
	if err != nil {
		app.ErrorLogger("Failed to get user from database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}
	user, ok := data.(*entity.User)
	if !ok {
		app.errorJSON(writer, http.StatusInternalServerError)
		return
	}
	if user == nil {
		app.errorJSON(writer, http.StatusUnauthorized)
		return
	}
	// generate token
	token, err := app.authenticator.GenerateToken(int64(user.Id))
	if err != nil {
		app.ErrorLogger("Failed to generate token", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}

	app.writeJSON(writer, http.StatusCreated, token)
}
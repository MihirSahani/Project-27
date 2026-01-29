package main

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"github.com/MihirSahani/Project-27/internal"
	"github.com/MihirSahani/Project-27/storage/entity"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	LOGGED_IN_USER_ID = "logged_in_user_id"
	ARGUMENTED_USER   = "arg_user"
)


func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		Email     string `json:"email" validate:"required,max=100"`
		Password  string `json:"password" validate:"required,min=8,max=72"`
		FirstName string `json:"first_name" validate:"required,max=50"`
		LastName  string `json:"last_name" validate:"max=50"`
	}
	// Reading JSON
	err := app.readJSON(r, &payload)
	if err != nil {
		app.ErrorLogger("Failed to read JSON payload", err, http.StatusBadRequest, w, WarnLog)
		return
	}
	// Validating JSON
	err = internal.Validate.Struct(payload)
	if err != nil {
		app.ErrorLogger("Payload validation failed", err, http.StatusBadRequest, w, WarnLog)
		return
	}
	// Hashing password
	hashedPassword, err := hashPassword([]byte(payload.Password))
	if err != nil {
		app.ErrorLogger("Failed to hash password", err, http.StatusInternalServerError, w, ErrorLog)
		return
	}
	// Writing to database
	user, err := app.storageManager.WithTx(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.UserStorageManager.CreateUser(ctx, tx, &entity.User{
			Email:     payload.Email,
			Password:  hashedPassword,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
		})
	})

	if err != nil {
		app.ErrorLogger("Failed to write to database", err, http.StatusInternalServerError, w, ErrorLog)
		return
	}
	app.writeJSON(w, http.StatusCreated, &user)
}

func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := r.Context().Value(ARGUMENTED_USER).(*entity.User)
	app.writeJSON(w, http.StatusOK, &user)
}

func (app *Application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get userId from URL
	userId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.ErrorLogger("Error in string argument", err, http.StatusBadRequest, w, WarnLog)
		return
	}
	// Delete user from database
	_, err = app.storageManager.WithTx(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.UserStorageManager.DeleteUser(ctx, tx, userId)
	})
	if err != nil {
		app.ErrorLogger("Failed to delete user from database", err, http.StatusInternalServerError, w, ErrorLog)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email     *string `json:"email" validate:"omitempty,max=100"`
		Password  *string `json:"password" validate:"omitempty,min=8,max=72"`
		FirstName *string `json:"first_name" validate:"omitempty,max=50"`
		LastName  *string `json:"last_name" validate:"omitempty,max=50"`
	}
	// Reading JSON
	err := app.readJSON(r, &payload)
	if err != nil {
		app.ErrorLogger("Failed to read JSON payload", err, http.StatusBadRequest, w, WarnLog)
		return
	}
	// Validating JSON
	err = internal.Validate.Struct(payload)
	if err != nil {
		app.ErrorLogger("Payload validation failed", err, http.StatusBadRequest, w, WarnLog)
		return
	}
	// Get user from context
	originalUser := r.Context().Value(ARGUMENTED_USER).(*entity.User)
	// Update the required fields from payload
	if payload.Email != nil {
		originalUser.Email = *payload.Email
	}
	if payload.FirstName != nil {
		originalUser.FirstName = *payload.FirstName
	}
	if payload.LastName != nil {
		originalUser.LastName = *payload.LastName
	}
	
	if payload.Password != nil {
		hashedPassword, err := hashPassword([]byte(*payload.Password))
		if err != nil {
			app.ErrorLogger("Error Hashing the password", err, http.StatusInternalServerError, w, ErrorLog)
			return
		}
		originalUser.Password = hashedPassword
	}
	// Update the user in database
	_, err = app.storageManager.WithTx(r.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.UserStorageManager.UpdateUser(ctx, tx, originalUser)
	})

	if err != nil {
		app.ErrorLogger("Error updating user", err, http.StatusInternalServerError, w, ErrorLog)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
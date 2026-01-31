package main

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/MihirSahani/Project-27/internal"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (app *Application) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromHeader(r)
		if err != nil {
			app.errorJSON(w, http.StatusUnauthorized)
			return
		}
		userId, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.errorJSON(w, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(),LOGGED_IN_USER_ID, userId)
		app.logger.Info("Authenticated User", zap.Int64("user_id", userId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) addUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
		if err != nil {
			app.ErrorLogger("Error in string argument", err, http.StatusBadRequest, writer, WarnLog)
			return
		}
		
		user, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
			return app.storageManager.UserStorageManager.GetUserByID(ctx, tx, userId)
		})

		if err != nil {
			app.ErrorLogger("Failed to get user from database", err, http.StatusInternalServerError, writer, ErrorLog)
			return
		}
		
		ctx := context.WithValue(request.Context(), ARGUMENTED_USER, user)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func getTokenFromHeader(r *http.Request) (string, error) {
	// Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", internal.MissingAuthenticationError 
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", internal.InvalidAuthenticationError
	}
	return authHeader[len(prefix):], nil
}

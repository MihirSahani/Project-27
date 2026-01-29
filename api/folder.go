package main

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/MihirSahani/Project-27/internal"
	"github.com/MihirSahani/Project-27/storage/entity"
	"github.com/go-chi/chi/v5"
)

func (app *Application) createFolderHandler(writer http.ResponseWriter, request *http.Request) {
	var payload struct {
		Name string `json:"name" validate:"required,min=1,max=100"`
	}
	// Reading JSON
	err := app.readJSON(request, &payload)
	if err != nil {
		app.ErrorLogger("Failed to read JSON payload", err, http.StatusBadRequest, writer, WarnLog)
		return
	}
	// Validating JSON
	err = internal.Validate.Struct(payload)
	if err != nil {
		app.ErrorLogger("Payload validation failed", err, http.StatusBadRequest, writer, WarnLog)
		return
	}
	// Write to database
	folder, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.FolderStorageManager.CreateFolder(ctx, tx, &entity.Folder{
			Name: payload.Name,
			UserId: request.Context().Value(LOGGED_IN_USER_ID).(int64),
		})
	})
	if err != nil {
		app.ErrorLogger("Failed to write to database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}
	app.writeJSON(writer, http.StatusCreated, &folder)
}

func (app *Application) getAllFoldersHandler(writer http.ResponseWriter, request *http.Request) {
	// Writing to database
	folders, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.FolderStorageManager.GetAllFolders(ctx, tx, ctx.Value(LOGGED_IN_USER_ID).(int64))
	})
	if err != nil {
		app.ErrorLogger("Error fetching folders from database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}

	app.writeJSON(writer, http.StatusOK, &folders)
}

func (app *Application) deleteFolderHandler(writer http.ResponseWriter, request *http.Request) {
	// Get folder Id from URL
	folderId, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
	// Write to database
	folder, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.FolderStorageManager.DeleteFolder(ctx, tx, folderId)
	})
	if err != nil {
		app.ErrorLogger("Error deleting folder from database", err, http.StatusInternalServerError, writer, ErrorLog)
	}
	app.writeJSON(writer, http.StatusNoContent, &folder)
}

func (app *Application) GetNotesInFolder(writer http.ResponseWriter, request *http.Request) {
	// TODO: Implement
	app.writeJSON(writer, http.StatusNotImplemented, nil)
}
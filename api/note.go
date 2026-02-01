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

func (app *Application) deleteNoteHandler(writer http.ResponseWriter, request *http.Request) {
	// Get note Id from URL
	noteId, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
	// Write to database
	_, err = app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.NoteStorageManager.DeleteNote(ctx, tx, noteId)
	})
	if err != nil {
		app.ErrorLogger("Error deleting note from database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}
	// Remove from cache
	app.cacheManager.DeleteNote(noteId)

	app.writeJSON(writer, http.StatusOK, nil)
}

func (app *Application) getNoteByIDHandler(writer http.ResponseWriter, request *http.Request) {
	// Get note Id from URL
	noteId, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
	if err != nil {
		app.ErrorLogger("Invalid Notes Id", err, http.StatusBadRequest, writer, WarnLog)
		return
	}
	// Read from cache
	cachedNote, err := app.cacheManager.GetNote(noteId)
	if err == nil {
		app.writeJSON(writer, http.StatusOK, &cachedNote)
		return
	}
	// Writing to database
	note, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.NoteStorageManager.GetNoteByID(ctx, tx, noteId)
	})
	if err != nil {
		app.ErrorLogger("Error fetching note from database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}
	// Writing to cache
	app.cacheManager.SetNote(note.(*entity.Note))

	app.writeJSON(writer, http.StatusOK, &note)
}

func (app *Application) updateNoteHandler(writer http.ResponseWriter, request *http.Request) {
	// Get note Id from URL
	noteId, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
	if err != nil {
		app.ErrorLogger("Invalid note ID", err, http.StatusBadRequest, writer, WarnLog)
		return
	}

	var payload struct {
		Title   *string `json:"title" validate:"required,min=1,max=200"`
		Content *string `json:"content" validate:"required"`
	}
	// Reading JSON
	err = app.readJSON(request, &payload)
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
	note, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.NoteStorageManager.UpdateNote(ctx, tx, &entity.Note{
			Id:      noteId,
			Title:   *payload.Title,
			Content: *payload.Content,
		})
	})
	if err != nil {
		app.ErrorLogger("Failed to write to database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}
	// Write to Cache
	app.cacheManager.SetNote(note.(*entity.Note))

	app.writeJSON(writer, http.StatusOK, &note)
}

func (app *Application) createNoteHandler(writer http.ResponseWriter, request *http.Request) {
	var payload struct {
		Title    string `json:"title" validate:"required,min=1,max=200"`
		Content  string `json:"content" validate:"required"`
		FolderId int64  `json:"folder_id" validate:"required"`
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
	note, err := app.storageManager.WithTx(request.Context(), func(ctx context.Context, tx *sql.Tx) (any, error) {
		return app.storageManager.NoteStorageManager.CreateNote(ctx, tx, &entity.Note{
			Title:    payload.Title,
			Content:  payload.Content,
			UserId:   request.Context().Value(LOGGED_IN_USER_ID).(int64),
			FolderId: payload.FolderId,
		})
	})
	if err != nil {
		app.ErrorLogger("Failed to write to database", err, http.StatusInternalServerError, writer, ErrorLog)
		return
	}
	// Write to Cache
	app.cacheManager.SetNote(note.(*entity.Note))

	app.writeJSON(writer, http.StatusCreated, &note)
}

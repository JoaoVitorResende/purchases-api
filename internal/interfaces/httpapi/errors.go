// internal/interfaces/httpapi/errors.go
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	application "purchases-api/internal/application"
	"purchases-api/internal/domain/purchase"
)

// Equivale ao ExceptionFilter do .NET
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrItemNotFound),
		errors.Is(err, application.ErrPurchaseNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})

	case errors.Is(err, application.ErrInvalidID),
		errors.Is(err, purchase.ErrEmptyPurchase),
		errors.Is(err, purchase.ErrEmptyItemID),
		errors.Is(err, purchase.ErrInvalidQuantity),
		errors.Is(err, purchase.ErrInvalidPrice):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})

	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "erro desconhecido"})
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

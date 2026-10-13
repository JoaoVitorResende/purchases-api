// internal/interfaces/httpapi/purchase_handler.go
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	application "purchases-api/internal/application"
	command "purchases-api/internal/application/command"
	"purchases-api/internal/application/query"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Interfaces definidas onde são usadas
type createPurchaseHandler interface {
	Handle(ctx context.Context, cmd command.CreatePurchase) (command.CreatePurchaseResult, error)
}

type getPurchaseHandler interface {
	Handle(ctx context.Context, q query.GetPurchase) (*application.PurchaseView, error)
}

type PurchaseHandler struct {
	create createPurchaseHandler
	get    getPurchaseHandler
}

func NewPurchaseHandler(create createPurchaseHandler, get getPurchaseHandler) *PurchaseHandler {
	return &PurchaseHandler{create: create, get: get}
}

func (h *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limita o corpo a 1 MB

	var req createPurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "corpo da requisicao invalido"})
		return
	}

	cmd := command.CreatePurchase{Items: make([]command.CreatePurchaseItem, 0, len(req.Items))}
	for _, i := range req.Items {
		cmd.Items = append(cmd.Items, command.CreatePurchaseItem{ItemID: i.ItemID, Quantity: i.Quantity})
	}

	out, err := h.create.Handle(r.Context(), cmd)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createPurchaseResponse{
		ID: out.ID, TotalCents: out.TotalCents, CreatedAt: out.CreatedAt,
	})
}

func (h *PurchaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		handleError(w, application.ErrInvalidID)
		return
	}

	view, err := h.get.Handle(r.Context(), query.GetPurchase{ID: id})
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPurchaseResponse(view))
}

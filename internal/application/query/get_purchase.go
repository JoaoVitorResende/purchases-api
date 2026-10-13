package query

import (
	"context"
	application "purchases-api/internal/application"
)

type GetPurchase struct {
	ID int64
}

type GetPurchaseHandler struct {
	reader application.PurchaseReader
}

func NewGetPurchaseHandler(r application.PurchaseReader) *GetPurchaseHandler {
	return &GetPurchaseHandler{reader: r}
}

func (h *GetPurchaseHandler) Handle(ctx context.Context, q GetPurchase) (*application.PurchaseView, error) {
	if q.ID <= 0 {
		return nil, application.ErrInvalidID
	}
	return h.reader.FindByID(ctx, q.ID)
}

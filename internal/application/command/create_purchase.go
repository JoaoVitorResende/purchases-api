package command

import (
	"context"
	application "purchases-api/internal/application"
	"purchases-api/internal/domain/purchase"
	"time"
)

type CreatePurchaseItem struct {
	ItemID   string
	Quantity int
}

type CreatePurchase struct {
	Items []CreatePurchaseItem
}

type CreatePurchaseResult struct {
	ID         int64
	TotalCents int64
	CreatedAt  time.Time
}

type CreatePurchaseHandler struct {
	purchases purchase.Repository
	catalog   application.CatalogReader
	now       func() time.Time
}

func NewCreatePurchaseHandler(p purchase.Repository, c application.CatalogReader, now func() time.Time) *CreatePurchaseHandler {
	return &CreatePurchaseHandler{purchases: p, catalog: c, now: now}
}

func (h *CreatePurchaseHandler) Handle(ctx context.Context, cmd CreatePurchase) (CreatePurchaseResult, error) {
	lines := make([]purchase.LineItem, 0, len(cmd.Items))

	for _, i := range cmd.Items {
		ci, err := h.catalog.FindByID(ctx, i.ItemID)
		if err != nil {
			return CreatePurchaseResult{}, err
		}

		line, err := purchase.NewLineItem(ci.ItemID, ci.Title, ci.PriceCents, i.Quantity)
		if err != nil {
			return CreatePurchaseResult{}, err
		}
		lines = append(lines, line)
	}

	p, err := purchase.New(lines, h.now())
	if err != nil {
		return CreatePurchaseResult{}, err
	}

	id, err := h.purchases.Save(ctx, p)
	if err != nil {
		return CreatePurchaseResult{}, err
	}

	return CreatePurchaseResult{ID: id, TotalCents: p.TotalCents(), CreatedAt: p.CreatedAt()}, nil
}

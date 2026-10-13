// internal/application/command/sync_catalog_item.go
package command

import (
	"context"
	application "purchases-api/internal/application"
	"strings"
)

type SyncCatalogItem struct {
	ItemID     string
	Title      string
	PriceCents int64
}

type SyncCatalogItemHandler struct {
	catalog application.CatalogWriter
}

func NewSyncCatalogItemHandler(c application.CatalogWriter) *SyncCatalogItemHandler {
	return &SyncCatalogItemHandler{catalog: c}
}

func (h *SyncCatalogItemHandler) Handle(ctx context.Context, cmd SyncCatalogItem) error {
	itemID := strings.TrimSpace(cmd.ItemID)
	title := strings.TrimSpace(cmd.Title)

	if itemID == "" || title == "" || cmd.PriceCents <= 0 {
		return application.ErrInvalidCatalogItem
	}

	return h.catalog.Upsert(ctx, application.CatalogItem{
		ItemID:     itemID,
		Title:      title,
		PriceCents: cmd.PriceCents,
	})
}

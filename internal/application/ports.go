// internal/application/ports.go
package application

import (
	"context"
	"time"
)

// Cópia local do catálogo (alimentada pelos eventos do StoreApi)
type CatalogItem struct {
	ItemID     string
	Title      string
	PriceCents int64
}

type CatalogReader interface {
	FindByID(ctx context.Context, itemID string) (*CatalogItem, error) // ErrItemNotFound se não existir
}

type CatalogWriter interface {
	Upsert(ctx context.Context, item CatalogItem) error
}

// Modelos de leitura
type PurchaseItemView struct {
	ItemID         string
	Title          string
	UnitPriceCents int64
	Quantity       int
	SubtotalCents  int64
}

type PurchaseView struct {
	ID         int64
	TotalCents int64
	CreatedAt  time.Time
	Items      []PurchaseItemView
}

type PurchaseReader interface {
	FindByID(ctx context.Context, id int64) (*PurchaseView, error) // ErrPurchaseNotFound se não existir
}
// internal/infrastructure/postgres/catalog_store.go
package postgres

import (
	"context"
	"errors"
	application "purchases-api/internal/application"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CatalogStore implementa CatalogReader e CatalogWriter
type CatalogStore struct {
	pool *pgxpool.Pool
}

func NewCatalogStore(pool *pgxpool.Pool) *CatalogStore {
	return &CatalogStore{pool: pool}
}

func (s *CatalogStore) FindByID(ctx context.Context, itemID string) (*application.CatalogItem, error) {
	var c application.CatalogItem
	err := s.pool.QueryRow(ctx,
		`select item_id, title, price_cents from catalog_items where item_id = $1`, itemID,
	).Scan(&c.ItemID, &c.Title, &c.PriceCents)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, application.ErrItemNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Upsert insere ou atualiza; processar o mesmo evento duas vezes dá o mesmo resultado
func (s *CatalogStore) Upsert(ctx context.Context, item application.CatalogItem) error {
	_, err := s.pool.Exec(ctx,
		`insert into catalog_items (item_id, title, price_cents)
		 values ($1, $2, $3)
		 on conflict (item_id) do update
		   set title = excluded.title, price_cents = excluded.price_cents`,
		item.ItemID, item.Title, item.PriceCents)
	return err
}

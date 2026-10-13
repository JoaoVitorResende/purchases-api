// internal/infrastructure/postgres/purchase_repository.go
package postgres

import (
	"context"
	"purchases-api/internal/domain/purchase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository struct {
	pool *pgxpool.Pool
}

func NewPurchaseRepository(pool *pgxpool.Pool) *PurchaseRepository {
	return &PurchaseRepository{pool: pool}
}

// Save grava a compra e os itens na mesma transação
func (r *PurchaseRepository) Save(ctx context.Context, p *purchase.Purchase) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }() // sem efeito depois do Commit

	var id int64
	err = tx.QueryRow(ctx,
		`insert into purchases (total_cents, created_at) values ($1, $2) returning id`,
		p.TotalCents(), p.CreatedAt(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, it := range p.Items() {
		_, err = tx.Exec(ctx,
			`insert into purchase_items (purchase_id, item_id, title, unit_price_cents, quantity)
			 values ($1, $2, $3, $4, $5)`,
			id, it.ItemID(), it.Title(), it.UnitPriceCents(), it.Quantity(),
		)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

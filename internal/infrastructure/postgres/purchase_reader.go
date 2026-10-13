// internal/infrastructure/postgres/purchase_reader.go
package postgres

import (
	"context"
	application "purchases-api/internal/application"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseReader struct {
	pool *pgxpool.Pool
}

func NewPurchaseReader(pool *pgxpool.Pool) *PurchaseReader {
	return &PurchaseReader{pool: pool}
}

// FindByID lê com uma única consulta (JOIN) direto para a view, sem montar o agregado
func (r *PurchaseReader) FindByID(ctx context.Context, id int64) (*application.PurchaseView, error) {
	rows, err := r.pool.Query(ctx,
		`select p.id, p.total_cents, p.created_at,
		        i.item_id, i.title, i.unit_price_cents, i.quantity
		   from purchases p
		   join purchase_items i on i.purchase_id = p.id
		  where p.id = $1
		  order by i.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var view *application.PurchaseView
	for rows.Next() {
		var (
			purchaseID int64
			total      int64
			createdAt  time.Time
			item       application.PurchaseItemView
		)
		if err := rows.Scan(&purchaseID, &total, &createdAt,
			&item.ItemID, &item.Title, &item.UnitPriceCents, &item.Quantity); err != nil {
			return nil, err
		}
		item.SubtotalCents = item.UnitPriceCents * int64(item.Quantity)

		if view == nil {
			view = &application.PurchaseView{ID: purchaseID, TotalCents: total, CreatedAt: createdAt}
		}
		view.Items = append(view.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if view == nil {
		return nil, application.ErrPurchaseNotFound
	}
	return view, nil
}

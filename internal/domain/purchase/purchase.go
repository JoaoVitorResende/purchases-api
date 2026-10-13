// internal/domain/purchase/purchase.go
package purchase

import "time"

// Purchase é o agregado raiz: garante que a compra tem itens
// e que o total é sempre a soma dos subtotais.
type Purchase struct {
	items      []LineItem
	totalCents int64
	createdAt  time.Time
}

func New(items []LineItem, now time.Time) (*Purchase, error) {
	if len(items) == 0 {
		return nil, ErrEmptyPurchase
	}

	var total int64
	for _, it := range items {
		total += it.SubtotalCents()
	}

	return &Purchase{
		items:      append([]LineItem(nil), items...), // cópia, para ninguém alterar por fora
		totalCents: total,
		createdAt:  now,
	}, nil
}

func (p *Purchase) Items() []LineItem {
	return append([]LineItem(nil), p.items...)
}

func (p *Purchase) TotalCents() int64    { return p.totalCents }
func (p *Purchase) CreatedAt() time.Time { return p.createdAt }

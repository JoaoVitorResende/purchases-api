// internal/domain/purchase/line_item.go
package purchase

import "strings"

// LineItem é a foto do item no momento da compra (objeto de valor)
type LineItem struct {
	itemID         string
	title          string
	unitPriceCents int64
	quantity       int
}

func NewLineItem(itemID, title string, unitPriceCents int64, quantity int) (LineItem, error) {
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return LineItem{}, ErrEmptyItemID
	}
	if quantity <= 0 {
		return LineItem{}, ErrInvalidQuantity
	}
	if unitPriceCents <= 0 {
		return LineItem{}, ErrInvalidPrice
	}
	return LineItem{
		itemID:         itemID,
		title:          strings.TrimSpace(title),
		unitPriceCents: unitPriceCents,
		quantity:       quantity,
	}, nil
}

func (i LineItem) ItemID() string        { return i.itemID }
func (i LineItem) Title() string         { return i.title }
func (i LineItem) UnitPriceCents() int64 { return i.unitPriceCents }
func (i LineItem) Quantity() int         { return i.quantity }
func (i LineItem) SubtotalCents() int64  { return i.unitPriceCents * int64(i.quantity) }

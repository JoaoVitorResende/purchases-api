// internal/interfaces/httpapi/dto.go
package httpapi

import (
	application "purchases-api/internal/application"
	"time"
)

type createPurchaseItemRequest struct {
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

type createPurchaseRequest struct {
	Items []createPurchaseItemRequest `json:"items"`
}

type createPurchaseResponse struct {
	ID         int64     `json:"id"`
	TotalCents int64     `json:"totalCents"`
	CreatedAt  time.Time `json:"createdAt"`
}

type purchaseItemResponse struct {
	ItemID         string `json:"itemId"`
	Title          string `json:"title"`
	UnitPriceCents int64  `json:"unitPriceCents"`
	Quantity       int    `json:"quantity"`
	SubtotalCents  int64  `json:"subtotalCents"`
}

type purchaseResponse struct {
	ID         int64                  `json:"id"`
	TotalCents int64                  `json:"totalCents"`
	CreatedAt  time.Time              `json:"createdAt"`
	Items      []purchaseItemResponse `json:"items"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func toPurchaseResponse(v *application.PurchaseView) purchaseResponse {
	items := make([]purchaseItemResponse, 0, len(v.Items))
	for _, i := range v.Items {
		items = append(items, purchaseItemResponse{
			ItemID:         i.ItemID,
			Title:          i.Title,
			UnitPriceCents: i.UnitPriceCents,
			Quantity:       i.Quantity,
			SubtotalCents:  i.SubtotalCents,
		})
	}
	return purchaseResponse{ID: v.ID, TotalCents: v.TotalCents, CreatedAt: v.CreatedAt, Items: items}
}

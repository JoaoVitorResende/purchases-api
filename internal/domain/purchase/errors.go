// internal/domain/purchase/errors.go
package purchase

import "errors"

var (
	ErrEmptyPurchase   = errors.New("a compra precisa ter pelo menos um item")
	ErrEmptyItemID     = errors.New("o id do item e obrigatorio")
	ErrInvalidQuantity = errors.New("a quantidade deve ser maior que zero")
	ErrInvalidPrice    = errors.New("o preco deve ser maior que zero")
)

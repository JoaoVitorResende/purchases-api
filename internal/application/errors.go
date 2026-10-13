package application

import "errors"

var (
	ErrItemNotFound       = errors.New("item nao encontrado")
	ErrPurchaseNotFound   = errors.New("compra nao encontrada")
	ErrInvalidID          = errors.New("id invalido")
	ErrInvalidCatalogItem = errors.New("item de catalogo invalido")
)

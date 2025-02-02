package gann

import "errors"

var (
	errDimensionMismatch         = errors.New("dimension mismatch")
	errInvalidIndex              = errors.New("invalid index")
	errInvalidKeyVector          = errors.New("invalid key vector")
	errItemNotFoundOnGivenItemID = errors.New("item not found for give item id")
)

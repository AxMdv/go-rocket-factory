package model

import "errors"

var (
	ErrBadRequest          = errors.New("bad request")
	ErrOrderNotFound       = errors.New("order not found")
	ErrPartsNotFound       = errors.New("parts not found")
	ErrOrderStatusConflict = errors.New("conflict of order status")
)

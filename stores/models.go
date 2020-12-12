package stores

import "github.com/go-rod/rod"

// Based off: https://programming.guide/go/define-enumeration-string.html
type StockStatus int

const (
	InStock StockStatus = iota
	OutOfStock
	Unknown
)

func (s StockStatus) String() string {
	return [...]string{
		"In stock",
		"Out of stock",
		"Unknown"}[s]
}

type StockCheckResult struct {
	storeName string
	status    StockStatus
}

type Store interface {
	Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult
}

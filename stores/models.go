package stores

import (
	"fmt"

	"github.com/go-rod/rod"
)

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
	StoreName string
	Status    StockStatus
}

func (s StockCheckResult) String() string {
	return fmt.Sprintf("%s: %s", s.StoreName, s.Status.String())
}

type Store interface {
	Check(pool rod.PagePool, create func() *rod.Page) StockCheckResult
}

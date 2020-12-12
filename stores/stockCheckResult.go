package stores

// Based off: https://programming.guide/go/define-enumeration-string.html
type StockCheckResult int

const (
	InStock StockCheckResult = iota
	OutOfStock
	ErrorOccurred
)

func (s StockCheckResult) String() string {
	return [...]string{
		"In stock",
		"Out of stock",
		"Error occurred"}[s]
}

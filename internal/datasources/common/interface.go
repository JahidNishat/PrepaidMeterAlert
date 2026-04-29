package common

type DataFetcher interface {
	GetBalance(Identifier) (Balance, error)
}

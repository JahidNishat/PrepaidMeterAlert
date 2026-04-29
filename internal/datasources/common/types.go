package common

type Identifier struct {
	AccountNumber string
	MeterNumber   string
}

type Balance struct {
	Identifier
	Balance int
}

type AccountDetails struct {
	Identifier
	CustomerName    string
	CustomerContact string
	SanctionedLoad  string
	Address         string
}

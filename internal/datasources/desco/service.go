package desco

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/m4hi2/MeterAlertBot/internal/datasources"
)

const name = "DESCO"

const (
	basePath    = "https://prepaid.desco.org.bd"
	balancePath = "/api/unified/customer/getBalance"
)

const (
	paramAccountNo = "accountNo"
	paramMeterNo   = "meterNo"
)

type Service struct {
	client *datasources.Client
	name   string
}

func NewService() *Service {
	return &Service{
		client: datasources.NewClient(&datasources.ClientConfig{
			BasePath:   basePath,
			Timeout:    10 * time.Second,
			Retry:      3,
			RetryDelay: time.Second,
		}),
		name: name,
	}
}

func (s *Service) GetBalance(ctx context.Context, id datasources.Identifier) (datasources.Balance, error) {
	ctx = context.WithValue(ctx, datasources.CtxKeyDatasource, datasources.CtxDatasourceDesco)

	q := url.Values{}
	q.Set(paramAccountNo, id.AccountNumber)
	q.Set(paramMeterNo, id.MeterNumber)
	path := balancePath + "?" + q.Encode()

	var resp GetBalanceResp
	if err := s.client.Do(ctx, http.MethodGet, path, nil, nil, &resp); err != nil {
		return datasources.Balance{}, fmt.Errorf("get balance: %w", err)
	}

	if resp.Code != http.StatusOK {
		return datasources.Balance{}, fmt.Errorf("get balance: upstream code %d: %s", resp.Code, resp.Desc)
	}

	return datasources.Balance{
		Identifier: id,
		Balance:    resp.Data.Balance,
	}, nil
}

func (s *Service) Name() string {
	return s.name
}

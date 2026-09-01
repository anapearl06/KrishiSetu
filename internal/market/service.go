package market

import (
	"context"
	"time"
)

type Service interface {
	GetLatestPrice(
		ctx context.Context,
		commodity string,
	) (*MarketPrice, error)

	GetRegionalPrices(
		ctx context.Context,
		commodity string,
		state string,
		district string,
	) ([]MarketPrice, error)

	GetHistoricalPrices(
		ctx context.Context,
		commodity string,
		from time.Time,
		to time.Time,
	) ([]MarketPrice, error)

	GetPriceIntelligence(
		ctx context.Context,
		commodity string,
	) (*PriceIntelligence, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetLatestPrice(
	ctx context.Context,
	commodity string,
) (*MarketPrice, error) {
	return s.repo.GetLatestByCommodity(
		ctx,
		commodity,
	)
}

func (s *service) GetRegionalPrices(
	ctx context.Context,
	commodity string,
	state string,
	district string,
) ([]MarketPrice, error) {
	return s.repo.GetByCommodityAndLocation(
		ctx,
		commodity,
		state,
		district,
	)
}

func (s *service) GetHistoricalPrices(
	ctx context.Context,
	commodity string,
	from time.Time,
	to time.Time,
) ([]MarketPrice, error) {
	return s.repo.GetHistoricalByCommodity(
		ctx,
		commodity,
		from,
		to,
	)
}


func (s *service) GetPriceIntelligence(
	ctx context.Context,
	commodity string,
) (*PriceIntelligence, error) {
	records, err := s.repo.GetHistoricalByCommodity(
		ctx,
		commodity,
		time.Time{},
		time.Now(),
	)

	if err != nil {
		return nil, err
	}

	return CalculatePriceIntelligence(records), nil
}
package app

import (
	"fmt"
	"rate-calculator/pkg/estimator/domain"
)

type multiplier = float32

//go:generate moq -out aggregator_mock_test.go . aggregator
type aggregator interface {
	Aggregate(*domain.SegmentFare) error
}

type RateConfig struct {
	Rule func(*domain.SegmentDelta) (bool, multiplier)
	Fare domain.Fare
}

type Estimator struct {
	aggregator  aggregator
	rateConfigs []RateConfig
}

func NewEstimator(config []RateConfig, aggregator aggregator) *Estimator {
	return &Estimator{rateConfigs: config, aggregator: aggregator}
}

func (e *Estimator) Estimate(delta *domain.SegmentDelta) error {
	finalRate := &domain.SegmentFare{ID: delta.RideID, Fare: 0}

	if delta.Dirty {
		return e.sendToAggregate(finalRate)
	}

	for _, cfg := range e.rateConfigs {
		ok, mult := cfg.Rule(delta)
		if ok {
			finalRate.Fare = domain.Fare(mult * float32(cfg.Fare))
			return e.sendToAggregate(finalRate)
		}
	}

	return fmt.Errorf("unable to find a suitable rule for the rideID: %d", delta.RideID)
}

func (e *Estimator) sendToAggregate(finalRate *domain.SegmentFare) error {
	if err := e.aggregator.Aggregate(finalRate); err != nil {
		return fmt.Errorf("unable to aggregate: %w", err)
	}
	return nil
}

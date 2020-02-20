package rater

import (
	"fmt"
)

type multiplier = float32

//go:generate moq -out aggregator_mock_test.go . aggregator
type aggregator interface {
	Aggregate(*SegmentFare) error
}

type RateConfig struct {
	Rule func(*SegmentDelta) (bool, multiplier)
	Fare float32
}

type SegmentFare struct {
	ID   RideID
	Fare float32
}

type Estimator struct {
	aggregator  aggregator
	rateConfigs []RateConfig
}

func NewEstimator(config []RateConfig, aggregator aggregator) *Estimator {
	return &Estimator{rateConfigs: config, aggregator: aggregator}
}

func (e *Estimator) Estimate(delta *SegmentDelta) error {
	finalRate := &SegmentFare{ID: delta.RideID, Fare: 0}

	if delta.Dirty {
		return e.sendToAggregate(finalRate)
	}

	for _, config := range e.rateConfigs {
		ok, mult := config.Rule(delta)
		if ok {
			finalRate.Fare = mult * config.Fare
			return e.sendToAggregate(finalRate)
		}
	}

	return fmt.Errorf("unable to find a suitable rule for the rideID: %d", delta.RideID)
}

func (e *Estimator) sendToAggregate(finalRate *SegmentFare) error {
	if err := e.aggregator.Aggregate(finalRate); err != nil {
		return fmt.Errorf("unable to aggregate: %w", err)
	}
	return nil
}

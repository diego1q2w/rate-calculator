package app

import (
	"fmt"
	"rate-calculator/pkg/domain"
)

//go:generate moq -out estimator_mock_test.go . estimator
type estimator interface {
	Estimate(delta *domain.SegmentDelta) error
}

type SpeedFilter struct {
	estimator  estimator
	speedLimit float32
}

func NewSpeedFilter(estimator estimator, speedLimit float32) *SpeedFilter {
	return &SpeedFilter{estimator: estimator, speedLimit: speedLimit}
}

func (f *SpeedFilter) Filter(delta *domain.SegmentDelta) error {
	if delta.Velocity > f.speedLimit {
		delta.Dirty = true
	}

	if err := f.estimator.Estimate(delta); err != nil {
		return fmt.Errorf("unable to estimate :%w", err)
	}

	return nil
}

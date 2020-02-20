package rater

import "fmt"

//go:generate moq -out aggregator_mock_test.go . aggregator
type aggregator interface {
	Aggregate(delta *SegmentDelta) error
}

type SpeedFilter struct {
	aggregator aggregator
	speedLimit float32
}

func NewSpeedFilter(aggregator aggregator, speedLimit float32) *SpeedFilter {
	return &SpeedFilter{aggregator: aggregator, speedLimit: speedLimit}
}

func (f *SpeedFilter) Filter(delta *SegmentDelta) error {
	if delta.Velocity > f.speedLimit {
		delta.Dirty = true
	}

	if err := f.aggregator.Aggregate(delta); err != nil {
		return fmt.Errorf("unable to aggregate :%w", err)
	}

	return nil
}

package app

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"rate-calculator/pkg/estimator/domain"
	"testing"
	"time"
)

func TestEstimator(t *testing.T) {
	testCases := map[string]struct {
		aggregatorErr error
		config        []RateConfig
		deltas        []*domain.SegmentDelta
		expectedFares []*domain.SegmentFare
		expectedErr   error
	}{
		"if the agggregator has an error then it should return an error": {
			aggregatorErr: errors.New("test"),
			expectedErr:   errors.New("unable to aggregate: test"),
			deltas:        []*domain.SegmentDelta{{Dirty: true}},
			expectedFares: []*domain.SegmentFare{},
		},
		"should assign the correct fares based on the rules": {
			deltas: []*domain.SegmentDelta{
				{RideID: 1, Distance: 2, Velocity: 12, Duration: 1, Date: time.Date(1, 0, 0, 5, 0, 0, 0, time.UTC), Dirty: false},
				{RideID: 2, Distance: 2, Velocity: 12, Duration: 1, Date: time.Date(1, 0, 0, 5, 0, 1, 0, time.UTC), Dirty: false},
				{RideID: 3, Distance: 2, Velocity: 12, Duration: 1, Date: time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC), Dirty: false},
				{RideID: 4, Distance: 2, Velocity: 12, Duration: 1, Date: time.Date(1, 0, 0, 0, 1, 0, 0, time.UTC), Dirty: false},
				{RideID: 5, Distance: 2, Velocity: 10, Duration: 1, Date: time.Date(1, 0, 0, 0, 0, 1, 0, time.UTC), Dirty: false},
				{RideID: 6, Distance: 2, Velocity: 10, Duration: 1, Date: time.Date(1, 0, 0, 0, 0, 1, 0, time.UTC), Dirty: true},
			},
			config: []RateConfig{
				{Rule: func(delta *domain.SegmentDelta) (b bool, m float32) {
					start := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 5, 0, 0, 0, time.UTC) // 5:00 - 0
					end := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 0, 0, 0, 0, time.UTC)
					return delta.Velocity > 10 && inTimeSpanT(start, end, delta.Date),
						delta.Distance
				}, Fare: 1},
				{Rule: func(delta *domain.SegmentDelta) (b bool, m float32) {
					start := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 0, 0, 0, 0, time.UTC)
					end := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 5, 0, 0, 0, time.UTC)
					return delta.Velocity > 10 && inTimeSpanT(start, end, delta.Date),
						delta.Distance
				}, Fare: 2},
				{Rule: func(delta *domain.SegmentDelta) (b bool, m float32) {
					return delta.Velocity <= 10, delta.Duration
				}, Fare: 3},
			},
			expectedFares: []*domain.SegmentFare{
				{ID: 1, Fare: 4},
				{ID: 2, Fare: 2},
				{ID: 3, Fare: 2},
				{ID: 4, Fare: 4},
				{ID: 5, Fare: 3},
				{ID: 6, Fare: 0},
			},
		},
		"if no rules applicable should return an error": {
			expectedErr: errors.New("unable to find a suitable rule for the rideID: 1"),
			deltas: []*domain.SegmentDelta{
				{RideID: 1, Distance: 2, Velocity: 12, Duration: 1, Date: time.Date(2016, 1, 1, 5, 0, 0, 0, time.UTC), Dirty: false},
			},
			config: []RateConfig{
				{Rule: func(delta *domain.SegmentDelta) (b bool, m float32) {
					return false, delta.Distance
				}, Fare: 1},
			},
			expectedFares: []*domain.SegmentFare{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			segmentFares := make([]*domain.SegmentFare, 0)
			aggregator := &aggregatorMock{AggregateFunc: func(s *domain.SegmentFare) error {
				if tc.aggregatorErr == nil {
					segmentFares = append(segmentFares, s)
				}
				return tc.aggregatorErr
			}}

			estimator := NewEstimator(tc.config, aggregator)

			var estimatorErr error
			for _, d := range tc.deltas {
				if err := estimator.Estimate(d); err != nil {
					estimatorErr = err
				}
			}
			if fmt.Sprintf("%s", estimatorErr) != fmt.Sprintf("%s", tc.expectedErr) {
				t.Errorf("expected error: %s, got: %s", tc.expectedErr, estimatorErr)
			}

			assert.Equal(t, tc.expectedFares, segmentFares)
		})
	}
}

// Non inclusive start
func inTimeSpanT(start, end, check time.Time) bool {
	if start.Before(end) {
		return check.After(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return start.Before(check) || !end.Before(check)
}

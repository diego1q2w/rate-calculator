package rater

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpeedFilter(t *testing.T) {
	testCases := map[string]struct {
		delta          *SegmentDelta
		speedLimit     float32
		expectedResult *SegmentDelta
		aggregatorErr  error
		expectedErr    error
	}{
		"if the aggregator fails an error should be returned": {
			delta:         &SegmentDelta{Velocity: 2},
			aggregatorErr: errors.New("test"),
			expectedErr:   errors.New("unable to aggregate :test"),
			speedLimit:    3,
		},
		"should filter the speedy ones": {
			speedLimit:     3,
			delta:          &SegmentDelta{Velocity: 4},
			expectedResult: &SegmentDelta{Velocity: 4, Dirty: true},
		},
		"should not filter the slow ones": {
			speedLimit:     3,
			delta:          &SegmentDelta{Velocity: 2},
			expectedResult: &SegmentDelta{Velocity: 2, Dirty: false},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var delta *SegmentDelta
			aggregator := &aggregatorMock{AggregateFunc: func(d *SegmentDelta) error {
				if tc.aggregatorErr == nil {
					delta = d
				}
				return tc.aggregatorErr
			}}

			filter := NewSpeedFilter(aggregator, tc.speedLimit)

			err := filter.Filter(tc.delta)
			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tc.expectedErr) {
				t.Errorf("expected error: %s, got: %s", tc.expectedErr, err)
			}

			assert.Equal(t, tc.expectedResult, delta)
		})
	}
}

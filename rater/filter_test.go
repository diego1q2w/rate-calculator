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
		estimatorErr   error
		expectedErr    error
	}{
		"if the estimator fails an error should be returned": {
			delta:        &SegmentDelta{Velocity: 2},
			estimatorErr: errors.New("test"),
			expectedErr:  errors.New("unable to estimate :test"),
			speedLimit:   3,
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
			estimator := &estimatorMock{EstimateFunc: func(d *SegmentDelta) error {
				if tc.estimatorErr == nil {
					delta = d
				}
				return tc.estimatorErr
			}}

			filter := NewSpeedFilter(estimator, tc.speedLimit)

			err := filter.Filter(tc.delta)
			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tc.expectedErr) {
				t.Errorf("expected error: %s, got: %s", tc.expectedErr, err)
			}

			assert.Equal(t, tc.expectedResult, delta)
		})
	}
}

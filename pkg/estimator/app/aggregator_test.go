package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"rate-calculator/pkg/estimator/domain"
	fileOutput "rate-calculator/pkg/estimator/output"
	"sort"
	"testing"
	"time"
)

func TestAggregator(t *testing.T) {
	testCases := map[string]struct {
		segmentFare    []*domain.SegmentFare
		minimumFare    domain.Fare
		flagFare       domain.Fare
		outputErr      error
		expectedOutput []*domain.OutputFare
	}{
		"it should aggregate the fares correctly": {
			segmentFare: []*domain.SegmentFare{
				{ID: 1, Fare: 10},
				{ID: 1, Fare: 2},
				{ID: 1, Fare: 1},
				{ID: 2, Fare: 3},
				{ID: 3, Fare: 4},
				{ID: 4, Fare: 8},
				{ID: 4, Fare: 6},
				{ID: 7, Fare: 2},
				{ID: 9, Fare: 1},
				{ID: 9, Fare: 1},
				{ID: 10, Fare: 23},
				{ID: 12, Fare: 1},
				{ID: 12, Fare: 1},
				{ID: 12, Fare: 4},
				{ID: 12, Fare: 2},
			},
			flagFare:    2,
			minimumFare: 6,
			expectedOutput: []*domain.OutputFare{
				{ID: 1, Fare: 15},
				{ID: 2, Fare: 6},
				{ID: 3, Fare: 6},
				{ID: 4, Fare: 16},
				{ID: 7, Fare: 6},
				{ID: 9, Fare: 6},
				{ID: 10, Fare: 25},
				{ID: 12, Fare: 10},
			},
		},
		"if error creating file the process should be terminated": {
			segmentFare: []*domain.SegmentFare{
				{ID: 1, Fare: 10},
				{ID: 1, Fare: 2},
				{ID: 1, Fare: 1},
				{ID: 2, Fare: 3},
				{ID: 3, Fare: 4},
				{ID: 4, Fare: 8},
			},
			outputErr:   fileOutput.NewOpenFileError(errors.New("test")),
			flagFare:    2,
			minimumFare: 6,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var outputFares []*domain.OutputFare
			output := &outputMock{OutputFunc: func(f []*domain.OutputFare) error {
				if tc.outputErr == nil {
					sort.Slice(f, func(i, j int) bool {
						return f[i].ID < f[j].ID
					})
					outputFares = f
				}
				return tc.outputErr
			}}

			aggregator := NewAggregator(output, 4*time.Millisecond, tc.minimumFare, tc.flagFare, 200)
			for _, f := range tc.segmentFare {
				time.Sleep(2 * time.Millisecond)
				err := aggregator.Aggregate(f)
				assert.NoError(t, err)
			}
			<-aggregator.Running()

			assert.Equal(t, tc.expectedOutput, outputFares)
		})
	}
}

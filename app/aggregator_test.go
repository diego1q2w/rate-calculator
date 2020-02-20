package app

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func TestAggregator(t *testing.T) {
	testCases := map[string]struct {
		segmentFare    []*SegmentFare
		minimumFare    float32
		flagFare       float32
		expectedOutput []*OutputFare
	}{
		"it should aggregate the fares correctly": {
			segmentFare: []*SegmentFare{
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
			expectedOutput: []*OutputFare{
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
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var outputFares []*OutputFare
			output := &outputMock{OutputFunc: func(f []*OutputFare) error {
				sort.Slice(f, func(i, j int) bool {
					return f[i].ID < f[j].ID
				})
				outputFares = f
				return nil
			}}

			aggregator := NewAggregator(output, time.Millisecond*20, tc.minimumFare, tc.flagFare, 4)
			for _, f := range tc.segmentFare {
				err := aggregator.Aggregate(f)
				assert.NoError(t, err)
			}

			time.Sleep(time.Millisecond * 100)

			assert.Equal(t, tc.expectedOutput, outputFares)
		})
	}
}

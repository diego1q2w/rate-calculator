package rater

import (
	"github.com/stretchr/testify/assert"
	"github.com/umahmood/haversine"
	"testing"
	"time"
)

func TestSegmenter(t *testing.T) {
	testCases := map[string]struct {
		postitions       []*Position
		expectedSegments []*SegmentDelta
	}{
		"if not enough points nothing should happen": {
			postitions: []*Position{
				{RideID: 1, Lat: 1, Long: 1, Timestamp: 1405589100},
				{RideID: 2, Lat: 1, Long: 2, Timestamp: 1405589110},
				{RideID: 3, Lat: 1, Long: 4, Timestamp: 1405589130},
			},
			expectedSegments: []*SegmentDelta{},
		}, "should classify the segments correctly": {
			postitions: []*Position{
				{RideID: 1, Lat: 1, Long: 1, Timestamp: 3600},
				{RideID: 2, Lat: 1, Long: 2, Timestamp: 3600},
				{RideID: 1, Lat: 2, Long: 3, Timestamp: 3600 * 2},
				{RideID: 3, Lat: 1, Long: 4, Timestamp: 3600 * 2},
				{RideID: 1, Lat: 3, Long: 5, Timestamp: 3600 * 4},
				{RideID: 3, Lat: 2, Long: 6, Timestamp: 3600 * 7},
			},
			expectedSegments: []*SegmentDelta{
				{RideID: 1, Distance: 7, Duration: 1, Velocity: 7, Date: time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)},
				{RideID: 1, Distance: 13, Duration: 2, Velocity: 6.5, Date: time.Date(1970, 1, 1, 2, 0, 0, 0, time.UTC)},
				{RideID: 3, Distance: 13, Duration: 5, Velocity: 2.6, Date: time.Date(1970, 1, 1, 2, 0, 0, 0, time.UTC)},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			deltaSegments := make([]*SegmentDelta, 0)
			segmentFilter := &segmentFilterMock{
				FilterFunc: func(delta *SegmentDelta) error {
					deltaSegments = append(deltaSegments, delta)
					return nil
				},
			}

			distanceCalc := func(p1, p2 haversine.Coord) (mi, km float64) {
				d := p1.Lat + p1.Lon + p2.Lat + p2.Lon
				return 0, d
			}

			segmenter := NewSegmenter(segmentFilter, distanceCalc, 0)

			for _, p := range tc.postitions {
				if err := segmenter.Segment(p); err != nil {
					t.Fatalf("unexpected error while segmenting: %s", err)
				}
			}
			assert.Equal(t, tc.expectedSegments, deltaSegments)
		})
	}
}

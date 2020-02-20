package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/umahmood/haversine"
	"rate-calculator/pkg/domain"
	"sort"
	"testing"
	"time"
)

// It tests the whole flow
func TestAppFlow(t *testing.T) {
	var minFare float32 = 5.0                 // Minimal fare
	var flagFare float32 = 1.0                // Fare you get each ride
	var nightFare float32 = 1.0               // Fare you get during the night
	var dayFare float32 = 1.0                 // Fare you get during the day
	var idleFare float32 = 1.0                // Fare you get when you are either stop or driving to slow
	var speedLimitFare float32 = 10.0         // Speed limit for the idle fare
	var speedLimitFilter float32 = 20.0       // Speed limit for the point to be filtered
	var segmenterWorkers = 4                  // The number of workers used for the segmentation, filter and fare estimator process
	var aggregatorWorkers = 2                 // The numer of workers used to aggregate the fare estimator, ideally they would be less than the segment
	var flushInterval = time.Millisecond * 20 // The time it takes for a aggregator worker to send the acumulated data to the master aggregator

	var outputFares []*OutputFare // Final fares
	output := &outputMock{OutputFunc: func(f []*OutputFare) error {
		sort.Slice(f, func(i, j int) bool {
			return f[i].ID < f[j].ID
		})
		outputFares = f
		return nil
	}}

	distanceCalc := func(p1, p2 haversine.Coord) (mi, km float64) { // This is a dummy distance calculator just for this tests
		d := p1.Lat + p1.Lon + p2.Lat + p2.Lon
		return 0, d
	}

	// This are all the steps of the process
	aggregator := NewAggregator(output, flushInterval, minFare, flagFare, aggregatorWorkers)
	estimator := NewEstimator(getEstimatorConfig(dayFare, nightFare, idleFare, speedLimitFare), aggregator)
	filter := NewSpeedFilter(estimator, speedLimitFilter)
	segmenter := NewSegmenter(filter, distanceCalc, segmenterWorkers)
	/////

	for _, p := range getInput() {
		err := segmenter.Segment(p) // Here is where the process starts
		assert.NoError(t, err)
	}

	time.Sleep(time.Millisecond * 150)
	assert.Equal(t, expectedOutput(), outputFares)
}

func getEstimatorConfig(dayFare float32, nightFare float32, idleFare float32, speedLimitFare float32) []RateConfig {
	return []RateConfig{
		{Rule: func(delta *domain.SegmentDelta) (b bool, m multiplier) {
			start := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 5, 0, 0, 0, time.UTC) // 5:00 - 0
			end := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 0, 0, 0, 0, time.UTC)
			return delta.Velocity > speedLimitFare && inTimeSpan(start, end, delta.Date),
				delta.Distance
		}, Fare: dayFare},
		{Rule: func(delta *domain.SegmentDelta) (b bool, m multiplier) {
			start := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 0, 0, 0, 0, time.UTC)
			end := time.Date(delta.Date.Year(), delta.Date.Month(), delta.Date.Day(), 5, 0, 0, 0, time.UTC)
			return delta.Velocity > speedLimitFare && inTimeSpan(start, end, delta.Date),
				delta.Distance
		}, Fare: nightFare},
		{Rule: func(delta *domain.SegmentDelta) (b bool, m multiplier) {
			return delta.Velocity <= speedLimitFare, delta.Duration
		}, Fare: idleFare},
	}
}

func getInput() []*domain.Position {
	return []*domain.Position{
		// Ride ID 1
		{RideID: 1, Lat: 1, Long: 1, Timestamp: 3600},
		{RideID: 1, Lat: 2, Long: 3, Timestamp: 3600 * 2},
		{RideID: 1, Lat: 3, Long: 5, Timestamp: 3600 * 4},
		// Ride ID 2
		{RideID: 2, Lat: 1, Long: 2, Timestamp: 3600},
		// Ride ID 3
		{RideID: 3, Lat: 1, Long: 4, Timestamp: 3600 * 2},
		{RideID: 3, Lat: 2, Long: 6, Timestamp: 3600 * 7},
		// Ride ID 4
		{RideID: 4, Lat: 1, Long: 6, Timestamp: 3600 * 1},
		{RideID: 4, Lat: 2, Long: 5, Timestamp: 3600 * 2},
		// Ride ID 5
		{RideID: 5, Lat: 2, Long: 6, Timestamp: 3600 * 3},
		{RideID: 5, Lat: 3, Long: 5, Timestamp: 3600 * 7},
		// Ride ID 6
		{RideID: 6, Lat: 40, Long: 8, Timestamp: 3600 * 1},
		{RideID: 6, Lat: 1, Long: 6, Timestamp: 3600 * 2},
		{RideID: 6, Lat: 2, Long: 2, Timestamp: 3600 * 7},
		{RideID: 6, Lat: 3, Long: 3, Timestamp: 3600 * 10},
		// Ride ID 7
		{RideID: 7, Lat: 2, Long: 1, Timestamp: 3600},
		{RideID: 7, Lat: 20, Long: 2, Timestamp: 3600 * 1},
		{RideID: 7, Lat: 4, Long: 3, Timestamp: 3600 * 2},
		{RideID: 7, Lat: 5, Long: 4, Timestamp: 3600 * 3},
		{RideID: 7, Lat: 2, Long: 5, Timestamp: 3600 * 9},
		// Ride ID 8
		{RideID: 8, Lat: 100, Long: 2, Timestamp: 3600 * 4},
		{RideID: 8, Lat: 4, Long: 2, Timestamp: 3600 * 5},
		{RideID: 8, Lat: 5, Long: 2, Timestamp: 3600 * 9},
		{RideID: 8, Lat: 6, Long: 1, Timestamp: 3600 * 10},
		{RideID: 8, Lat: 7, Long: 1, Timestamp: 3600 * 11},
		// Ride ID 9
		{RideID: 9, Lat: 5, Long: 5, Timestamp: 3600},
		{RideID: 9, Lat: 1, Long: 5, Timestamp: 3600 * 3},
	}
}

func expectedOutput() []*OutputFare {
	return []*OutputFare{
		{ID: 1, Fare: 5},
		{ID: 3, Fare: 6},
		{ID: 4, Fare: 15},
		{ID: 5, Fare: 5},
		{ID: 6, Fare: 9},
		{ID: 7, Fare: 23},
		{ID: 8, Fare: 34},
		{ID: 9, Fare: 5},
	}
}

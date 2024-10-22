package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/umahmood/haversine"
	"rate-calculator/pkg/estimator/domain"
	"sort"
	"testing"
	"time"
)

// It tests the whole flow
func TestAppFlow(t *testing.T) {
	var minFare domain.Fare = 5.0             // Minimal fare
	var flagFare domain.Fare = 1.0            // Fare you get each ride
	var nightFare domain.Fare = 1.0           // Fare you get during the night
	var dayFare domain.Fare = 2.0             // Fare you get during the day
	var idleFare domain.Fare = 1.5            // Fare you get when you are either stop or driving to slow
	var speedLimitFare float32 = 10.0         // Speed limit for the idle fare
	var speedLimitFilter float32 = 20.0       // Speed limit for the point to be filtered
	var segmenterWorkers = 4                  // The number of workers used for the segmentation, filter and fare estimator process
	var aggregatorWorkers = 2                 // The numer of workers used to aggregate the fare estimator, ideally they would be less than the segment
	var flushInterval = time.Millisecond * 20 // The time it takes for a aggregator worker to send the acumulated data to the master aggregator
	//Time fares
	var dayStart = "05:00"
	var dayFinish = "00:00"
	var nightStart = "00:00"
	var nightFinish = "05:00"

	var outputFares []*domain.OutputFare // Final fares
	output := &outputMock{OutputFunc: func(f []*domain.OutputFare) error {
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

	dayRule := TimeRule{Start: dayStart, Finish: dayFinish, Fare: dayFare}
	nightRule := TimeRule{Start: nightStart, Finish: nightFinish, Fare: nightFare}
	speedRule := SpeedRule{Limit: speedLimitFare, Fare: idleFare}
	estimatorConfig, err := GetEstimatorConfig(dayRule, nightRule, speedRule)
	assert.NoError(t, err)

	// This are all the steps of the process
	aggregator := NewAggregator(output, flushInterval, minFare, flagFare, aggregatorWorkers)
	estimator := NewEstimator(estimatorConfig, aggregator)
	filter := NewSpeedFilter(estimator, speedLimitFilter)
	segmenter := NewSegmenter(filter, distanceCalc, segmenterWorkers)
	/////

	for _, p := range getInput() {
		err := segmenter.Segment(p) // Here is where the process starts
		assert.NoError(t, err)
	}

	<-aggregator.Running()
	assert.Equal(t, expectedOutput(), outputFares)
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

func expectedOutput() []*domain.OutputFare {
	return []*domain.OutputFare{
		{ID: 1, Fare: 5.5},
		{ID: 3, Fare: 8.5},
		{ID: 4, Fare: 15},
		{ID: 5, Fare: 7},
		{ID: 6, Fare: 13},
		{ID: 7, Fare: 26},
		{ID: 8, Fare: 65},
		{ID: 9, Fare: 5},
	}
}

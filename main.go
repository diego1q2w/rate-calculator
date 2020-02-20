package main

import (
	"context"
	"fmt"
	"github.com/umahmood/haversine"
	"os"
	"rate-calculator/pkg/api"
	"rate-calculator/pkg/app"
	"rate-calculator/pkg/domain"
	"rate-calculator/pkg/output"
	"time"
)

type multiplier = float32

const (
	// File
	csvPath        = "./paths.csv"
	outputPath     = "./output.txt"
	csvFieldLength = 4

	//Concurrency set up
	segmenterWorkers  = 15
	aggregatorWorkers = 5
	flushInterval     = time.Millisecond * 200

	//App variables
	speedLimit     = 100
	dayFare        = 0.74
	nightFare      = 1.30
	idleFare       = 11.90
	idleSpeedLimit = 10.0
	minFare        = 3.47
	flagFare       = 1.30
)

func main() {
	outputFile := output.NewFileOutput(outputPath)

	aggregator := app.NewAggregator(outputFile, flushInterval, minFare, flagFare, aggregatorWorkers)
	estimator := app.NewEstimator(getEstimatorConfig(dayFare, nightFare, idleFare, idleSpeedLimit), aggregator)
	filter := app.NewSpeedFilter(estimator, speedLimit)
	segmenter := app.NewSegmenter(filter, haversine.Distance, segmenterWorkers)

	fileReader := api.NewFileReader(segmenter, csvPath, csvFieldLength)

	fmt.Printf("Process has started, the output will appear in the file: %s. \n", outputFile)
	if err := fileReader.Process(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("The input file has been processed, due to the concurrency the process may take few moments more \n")
	<-context.Background().Done()
}

func getEstimatorConfig(dayFare, nightFare, idleFare domain.Fare, speedLimitFare float32) []app.RateConfig {
	return []app.RateConfig{
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

// Non inclusive start
func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return check.After(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return start.Before(check) || !end.Before(check)
}

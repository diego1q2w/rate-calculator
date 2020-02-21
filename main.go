package main

import (
	"fmt"
	"github.com/umahmood/haversine"
	"log"
	"os"
	"rate-calculator/pkg/estimator/api"
	"rate-calculator/pkg/estimator/app"
	"rate-calculator/pkg/estimator/output"
	"time"
)

const (
	// File
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

	//Time fares
	dayStart    = "05:00"
	dayFinish   = "00:00"
	nightStart  = "00:00"
	nightFinish = "05:00"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("One argument is required which is the directory path you wish to sync")
	}

	outputFile := output.NewFileOutput(outputPath)

	dayRule := app.TimeRule{Start: dayStart, Finish: dayFinish, Fare: dayFare}
	nightRule := app.TimeRule{Start: nightStart, Finish: nightFinish, Fare: nightFare}
	speedRule := app.SpeedRule{Limit: idleSpeedLimit, Fare: idleFare}
	estimatorConfig, err := app.GetEstimatorConfig(dayRule, nightRule, speedRule)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	aggregator := app.NewAggregator(outputFile, flushInterval, minFare, flagFare, aggregatorWorkers)
	estimator := app.NewEstimator(estimatorConfig, aggregator)
	filter := app.NewSpeedFilter(estimator, speedLimit)
	segmenter := app.NewSegmenter(filter, haversine.Distance, segmenterWorkers)

	fileReader := api.NewFileReader(segmenter, os.Args[1], csvFieldLength)

	fmt.Printf("Process has started, the output will appear in the file: %s. \n", outputFile)
	if err := fileReader.Process(); err != nil {
		log.Fatal("Error: ", err)

	}
	fmt.Printf("The input file has been processed, due to the concurrency the process may take few moments more \n")
	<-aggregator.Running()
	fmt.Printf("Process finished \n")
}

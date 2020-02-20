package app

import (
	"fmt"
	"rate-calculator/pkg/domain"
	"time"
)

//go:generate moq -out output_mock_test.go . output
type output interface {
	Output([]*domain.OutputFare) error
}

type finalFare map[domain.RideID]domain.Fare

type Aggregator struct {
	finalFare     finalFare
	flushInterval time.Duration
	flagFare      domain.Fare
	minFare       domain.Fare
	output        output
	workerCh      chan *domain.SegmentFare
}

func NewAggregator(output output, flushInterval time.Duration, minFare, flagFare domain.Fare, workers int) *Aggregator {
	a := &Aggregator{
		flushInterval: flushInterval,
		flagFare:      flagFare,
		finalFare:     make(finalFare),
		minFare:       minFare,
		output:        output,
	}

	if workers > 0 {
		ch := a.spinWorkers(workers)
		go a.masterAggregate(ch)
	}

	return a
}

// Aggregate queues the fare to be aggregated, the aggregation process is done in 2 parts
// - First is send to a "cluster" of workers where each worker would have its own count
// - Secondly and every now and then those workers are flushed into a single go routine where the final count is calculated
// and output, giving us the ability to process them without using mutex and blocking the routine
func (a *Aggregator) Aggregate(f *domain.SegmentFare) error {
	a.workerCh <- f

	return nil
}

// We aggregate indepenedntly in different go routines and then we flush every now and then it into a single routine,
// thanks to that we dont have to use Mutex while we ensure the unique increment of fares :D
func (a *Aggregator) aggregate(aggregateCh chan finalFare) {
	fFare := make(finalFare)
	flush := time.Tick(a.flushInterval)
	for {
		select {
		case fare := <-a.workerCh:
			if value, ok := fFare[fare.ID]; ok {
				fFare[fare.ID] = value + fare.Fare
			} else {
				fFare[fare.ID] = fare.Fare
			}
		case <-flush:
			aggregateCh <- fFare
			fFare = make(finalFare)
		}
	}
}

// This happens in only one gorouting and thus is syncronous ensuring an unique final aggregator
func (a *Aggregator) masterAggregate(ch chan finalFare) {
	for f := range ch {
		for rideID, fare := range f {
			if _, ok := a.finalFare[rideID]; ok {
				a.finalFare[rideID] += fare
			} else {
				a.finalFare[rideID] = fare + a.flagFare
			}
		}
		a.outputData()
	}
}

func (a *Aggregator) outputData() {
	var fareOutput = make([]*domain.OutputFare, 0)
	for rideID, fare := range a.finalFare {
		if fare < a.minFare {
			fare = a.minFare
		}
		fareOutput = append(fareOutput, &domain.OutputFare{ID: rideID, Fare: fare})
	}

	if err := a.output.Output(fareOutput); err != nil {
		fmt.Println("unable to output final fares")
	}
}

func (a *Aggregator) spinWorkers(numberOfWorkers int) chan finalFare {
	workerCh := make(chan *domain.SegmentFare)
	aggregateCh := make(chan finalFare)
	a.workerCh = workerCh
	for i := 0; i < numberOfWorkers; i++ {
		go a.aggregate(aggregateCh)
	}

	return aggregateCh
}
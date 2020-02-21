// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app

import (
	"rate-calculator/pkg/estimator/domain"
	"sync"
)

var (
	lockaggregatorMockAggregate sync.RWMutex
)

// Ensure, that aggregatorMock does implement aggregator.
// If this is not the case, regenerate this file with moq.
var _ aggregator = &aggregatorMock{}

// aggregatorMock is a mock implementation of aggregator.
//
//     func TestSomethingThatUsesaggregator(t *testing.T) {
//
//         // make and configure a mocked aggregator
//         mockedaggregator := &aggregatorMock{
//             AggregateFunc: func(in1 *domain.SegmentFare) error {
// 	               panic("mock out the Aggregate method")
//             },
//         }
//
//         // use mockedaggregator in code that requires aggregator
//         // and then make assertions.
//
//     }
type aggregatorMock struct {
	// AggregateFunc mocks the Aggregate method.
	AggregateFunc func(in1 *domain.SegmentFare) error

	// calls tracks calls to the methods.
	calls struct {
		// Aggregate holds details about calls to the Aggregate method.
		Aggregate []struct {
			// In1 is the in1 argument value.
			In1 *domain.SegmentFare
		}
	}
}

// Aggregate calls AggregateFunc.
func (mock *aggregatorMock) Aggregate(in1 *domain.SegmentFare) error {
	if mock.AggregateFunc == nil {
		panic("aggregatorMock.AggregateFunc: method is nil but aggregator.Aggregate was just called")
	}
	callInfo := struct {
		In1 *domain.SegmentFare
	}{
		In1: in1,
	}
	lockaggregatorMockAggregate.Lock()
	mock.calls.Aggregate = append(mock.calls.Aggregate, callInfo)
	lockaggregatorMockAggregate.Unlock()
	return mock.AggregateFunc(in1)
}

// AggregateCalls gets all the calls that were made to Aggregate.
// Check the length with:
//     len(mockedaggregator.AggregateCalls())
func (mock *aggregatorMock) AggregateCalls() []struct {
	In1 *domain.SegmentFare
} {
	var calls []struct {
		In1 *domain.SegmentFare
	}
	lockaggregatorMockAggregate.RLock()
	calls = mock.calls.Aggregate
	lockaggregatorMockAggregate.RUnlock()
	return calls
}
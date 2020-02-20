// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package rater

import (
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
//             AggregateFunc: func(delta *SegmentDelta) error {
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
	AggregateFunc func(delta *SegmentDelta) error

	// calls tracks calls to the methods.
	calls struct {
		// Aggregate holds details about calls to the Aggregate method.
		Aggregate []struct {
			// Delta is the delta argument value.
			Delta *SegmentDelta
		}
	}
}

// Aggregate calls AggregateFunc.
func (mock *aggregatorMock) Aggregate(delta *SegmentDelta) error {
	if mock.AggregateFunc == nil {
		panic("aggregatorMock.AggregateFunc: method is nil but aggregator.Aggregate was just called")
	}
	callInfo := struct {
		Delta *SegmentDelta
	}{
		Delta: delta,
	}
	lockaggregatorMockAggregate.Lock()
	mock.calls.Aggregate = append(mock.calls.Aggregate, callInfo)
	lockaggregatorMockAggregate.Unlock()
	return mock.AggregateFunc(delta)
}

// AggregateCalls gets all the calls that were made to Aggregate.
// Check the length with:
//     len(mockedaggregator.AggregateCalls())
func (mock *aggregatorMock) AggregateCalls() []struct {
	Delta *SegmentDelta
} {
	var calls []struct {
		Delta *SegmentDelta
	}
	lockaggregatorMockAggregate.RLock()
	calls = mock.calls.Aggregate
	lockaggregatorMockAggregate.RUnlock()
	return calls
}

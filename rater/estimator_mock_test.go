// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package rater

import (
	"sync"
)

var (
	lockestimatorMockEstimate sync.RWMutex
)

// Ensure, that estimatorMock does implement estimator.
// If this is not the case, regenerate this file with moq.
var _ estimator = &estimatorMock{}

// estimatorMock is a mock implementation of estimator.
//
//     func TestSomethingThatUsesestimator(t *testing.T) {
//
//         // make and configure a mocked estimator
//         mockedestimator := &estimatorMock{
//             EstimateFunc: func(delta *SegmentDelta) error {
// 	               panic("mock out the Estimate method")
//             },
//         }
//
//         // use mockedestimator in code that requires estimator
//         // and then make assertions.
//
//     }
type estimatorMock struct {
	// EstimateFunc mocks the Estimate method.
	EstimateFunc func(delta *SegmentDelta) error

	// calls tracks calls to the methods.
	calls struct {
		// Estimate holds details about calls to the Estimate method.
		Estimate []struct {
			// Delta is the delta argument value.
			Delta *SegmentDelta
		}
	}
}

// Estimate calls EstimateFunc.
func (mock *estimatorMock) Estimate(delta *SegmentDelta) error {
	if mock.EstimateFunc == nil {
		panic("estimatorMock.EstimateFunc: method is nil but estimator.Estimate was just called")
	}
	callInfo := struct {
		Delta *SegmentDelta
	}{
		Delta: delta,
	}
	lockestimatorMockEstimate.Lock()
	mock.calls.Estimate = append(mock.calls.Estimate, callInfo)
	lockestimatorMockEstimate.Unlock()
	return mock.EstimateFunc(delta)
}

// EstimateCalls gets all the calls that were made to Estimate.
// Check the length with:
//     len(mockedestimator.EstimateCalls())
func (mock *estimatorMock) EstimateCalls() []struct {
	Delta *SegmentDelta
} {
	var calls []struct {
		Delta *SegmentDelta
	}
	lockestimatorMockEstimate.RLock()
	calls = mock.calls.Estimate
	lockestimatorMockEstimate.RUnlock()
	return calls
}

// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package rater

import (
	"sync"
)

var (
	locksegmentFilterMockFilter sync.RWMutex
)

// Ensure, that segmentFilterMock does implement segmentFilter.
// If this is not the case, regenerate this file with moq.
var _ segmentFilter = &segmentFilterMock{}

// segmentFilterMock is a mock implementation of segmentFilter.
//
//     func TestSomethingThatUsessegmentFilter(t *testing.T) {
//
//         // make and configure a mocked segmentFilter
//         mockedsegmentFilter := &segmentFilterMock{
//             FilterFunc: func(delta *SegmentDelta) error {
// 	               panic("mock out the Filter method")
//             },
//         }
//
//         // use mockedsegmentFilter in code that requires segmentFilter
//         // and then make assertions.
//
//     }
type segmentFilterMock struct {
	// FilterFunc mocks the Filter method.
	FilterFunc func(delta *SegmentDelta) error

	// calls tracks calls to the methods.
	calls struct {
		// Filter holds details about calls to the Filter method.
		Filter []struct {
			// Delta is the delta argument value.
			Delta *SegmentDelta
		}
	}
}

// Filter calls FilterFunc.
func (mock *segmentFilterMock) Filter(delta *SegmentDelta) error {
	if mock.FilterFunc == nil {
		panic("segmentFilterMock.FilterFunc: method is nil but segmentFilter.Filter was just called")
	}
	callInfo := struct {
		Delta *SegmentDelta
	}{
		Delta: delta,
	}
	locksegmentFilterMockFilter.Lock()
	mock.calls.Filter = append(mock.calls.Filter, callInfo)
	locksegmentFilterMockFilter.Unlock()
	return mock.FilterFunc(delta)
}

// FilterCalls gets all the calls that were made to Filter.
// Check the length with:
//     len(mockedsegmentFilter.FilterCalls())
func (mock *segmentFilterMock) FilterCalls() []struct {
	Delta *SegmentDelta
} {
	var calls []struct {
		Delta *SegmentDelta
	}
	locksegmentFilterMockFilter.RLock()
	calls = mock.calls.Filter
	locksegmentFilterMockFilter.RUnlock()
	return calls
}